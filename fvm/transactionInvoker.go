package fvm

import (
	"fmt"
	"strconv"

	"github.com/onflow/cadence/runtime/interpreter"
	"github.com/onflow/cadence/runtime/stdlib"

	"github.com/onflow/cadence/runtime"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/sema"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	otelTrace "go.opentelemetry.io/otel/trace"

	"github.com/onflow/flow-go/fvm/errors"
	programsCache "github.com/onflow/flow-go/fvm/programs"
	"github.com/onflow/flow-go/fvm/state"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/trace"
)

type TransactionInvoker struct {
	logger zerolog.Logger
}

func NewTransactionInvoker(logger zerolog.Logger) *TransactionInvoker {
	return &TransactionInvoker{
		logger: logger,
	}
}

func (i *TransactionInvoker) Process(
	vm *VirtualMachine,
	ctx *Context,
	proc *TransactionProcedure,
	sth *state.StateHolder,
	programs *programsCache.Programs,
) (processErr error) {

	txIDStr := proc.ID.String()
	var span otelTrace.Span
	if ctx.Tracer != nil && proc.TraceSpan != nil {
		span = ctx.Tracer.StartSpanFromParent(proc.TraceSpan, trace.FVMExecuteTransaction)
		span.SetAttributes(
			attribute.String("transaction_id", txIDStr),
		)
		defer span.End()
	}

	var blockHeight uint64
	if ctx.BlockHeader != nil {
		blockHeight = ctx.BlockHeader.Height
	}

	nestedTxnId, err := sth.BeginNestedTransaction()
	if err != nil {
		return err
	}

	var modifiedSets programsCache.ModifiedSets
	defer func() {
		// based on the contract and frozen account updates we decide how to
		// clean up the programs for failed transactions we also do the same as
		// transaction without any deployed contracts
		programs.Cleanup(modifiedSets)

		if sth.NumNestedTransactions() > 1 {
			err := sth.RestartNestedTransaction(nestedTxnId)
			if err != nil {
				processErr = fmt.Errorf(
					"cannot restart nested transaction: %w",
					err,
				)
			}

			msg := "transaction has unexpected nested transactions"
			i.logger.Error().
				Str("txHash", txIDStr).
				Uint64("blockHeight", blockHeight).
				Msg(msg)

			proc.Err = errors.NewFVMInternalErrorf(msg)
			proc.Logs = make([]string, 0)
			proc.Events = make([]flow.Event, 0)
			proc.ServiceEvents = make([]flow.Event, 0)

		}

		var err error
		sth.RunWithAllLimitsDisabled(func() {
			err = sth.Commit(nestedTxnId)
		})
		if err != nil {
			processErr = fmt.Errorf("transaction invocation failed when merging state: %w", err)
		}
	}()

	env := NewTransactionEnvironment(*ctx, vm, sth, programs, proc.Transaction, proc.TxIndex, span)

	location := common.TransactionLocation(proc.ID)

	runtimeEnv := runtime.NewBaseInterpreterEnvironment(runtime.Config{})
	declareStandardLibraryValues(runtimeEnv, env)

	var txError error
	err = vm.Runtime.ExecuteTransaction(
		runtime.Script{
			Source:    proc.Transaction.Script,
			Arguments: proc.Transaction.Arguments,
		},
		runtime.Context{
			Interface:   env,
			Location:    location,
			Environment: runtimeEnv,
		},
	)
	if err != nil {
		var interactionLimitExceededErr *errors.LedgerInteractionLimitExceededError
		if errors.As(err, &interactionLimitExceededErr) {
			// If it is this special interaction limit error, just set it directly as the tx error
			txError = err
		} else {
			// Otherwise, do what we use to do
			txError = fmt.Errorf(
				"transaction invocation failed when executing transaction: %w",
				errors.HandleRuntimeError(err),
			)
		}
	}

	// read computationUsed from the environment. This will be used to charge fees.
	computationUsed := env.ComputationUsed()
	memoryEstimate := env.MemoryEstimate()

	// log te execution intensities here, so tha they do not contain data from storage limit checks and
	// transaction deduction, because the payer is not charged for those.
	i.logExecutionIntensities(sth, txIDStr)

	// disable the limit checks on states
	sth.DisableAllLimitEnforcements()
	// try to deduct fees even if there is an error.
	// disable the limit checks on states
	feesError := i.deductTransactionFees(env, proc, sth, computationUsed)
	if feesError != nil {
		txError = feesError
	}

	sth.EnableAllLimitEnforcements()

	// applying contract changes
	// this writes back the contract contents to accounts
	// if any error occurs we fail the tx
	// this needs to happen before checking limits, so that contract changes are committed to the state
	modifiedSets, err = env.Commit()
	if err != nil && txError == nil {
		txError = fmt.Errorf("transaction invocation failed when committing Environment: %w", err)
	}

	// if there is still no error check if all account storage limits are ok
	if txError == nil {
		// disable the computation/memory limit checks on storage checks,
		// so we don't error from computation/memory limits on this part.
		// We cannot charge the user for this part, since fee deduction already happened.
		sth.DisableAllLimitEnforcements()
		txError = NewTransactionStorageLimiter().CheckLimits(env, sth.UpdatedAddresses())
		sth.EnableAllLimitEnforcements()
	}

	// it there was any transaction error clear changes and try to deduct fees again
	if txError != nil {
		sth.DisableAllLimitEnforcements()
		defer sth.EnableAllLimitEnforcements()

		modifiedSets = programsCache.ModifiedSets{}

		// drop delta since transaction failed
		err := sth.RestartNestedTransaction(nestedTxnId)
		if err != nil {
			return fmt.Errorf(
				"cannot restart nested transaction: %w",
				err,
			)
		}

		// log transaction as failed
		i.logger.Info().
			Str("txHash", txIDStr).
			Uint64("blockHeight", blockHeight).
			Msg("transaction executed with error")

		// TODO(patrick): make env reusable on error
		// reset env
		env = NewTransactionEnvironment(*ctx, vm, sth, programs, proc.Transaction, proc.TxIndex, span)

		// try to deduct fees again, to get the fee deduction events
		feesError = i.deductTransactionFees(env, proc, sth, computationUsed)

		// if fee deduction fails just do clean up and exit
		if feesError != nil {
			i.logger.Info().
				Str("txHash", txIDStr).
				Uint64("blockHeight", blockHeight).
				Msg("transaction fee deduction executed with error")

			txError = feesError

			// drop delta
			_ = sth.RestartNestedTransaction(nestedTxnId)
		}
	}

	// if tx failed this will only contain fee deduction logs
	proc.Logs = append(proc.Logs, env.Logs()...)
	proc.ComputationUsed = proc.ComputationUsed + computationUsed
	proc.MemoryEstimate = proc.MemoryEstimate + memoryEstimate

	// if tx failed this will only contain fee deduction events
	proc.Events = append(proc.Events, env.Events()...)
	proc.ServiceEvents = append(proc.ServiceEvents, env.ServiceEvents()...)

	return txError
}

func (i *TransactionInvoker) deductTransactionFees(
	env *TransactionEnv,
	proc *TransactionProcedure,
	sth *state.StateHolder,
	computationUsed uint64) (err error) {
	if !env.ctx.TransactionFeesEnabled {
		return nil
	}

	if computationUsed > uint64(sth.TotalComputationLimit()) {
		computationUsed = uint64(sth.TotalComputationLimit())
	}

	// Hardcoded inclusion effort (of 1.0 UFix). Eventually this will be
	// dynamic.	Execution effort will be connected to computation used.
	inclusionEffort := uint64(100_000_000)
	_, err = InvokeDeductTransactionFeesContract(
		env,
		proc.Transaction.Payer,
		inclusionEffort,
		computationUsed)

	if err != nil {
		return errors.NewTransactionFeeDeductionFailedError(proc.Transaction.Payer, err)
	}
	return nil
}

var setAccountFrozenFunctionType = &sema.FunctionType{
	Parameters: []*sema.Parameter{
		{
			Label:          sema.ArgumentLabelNotRequired,
			Identifier:     "account",
			TypeAnnotation: sema.NewTypeAnnotation(&sema.AddressType{}),
		},
		{
			Label:          sema.ArgumentLabelNotRequired,
			Identifier:     "frozen",
			TypeAnnotation: sema.NewTypeAnnotation(sema.BoolType),
		},
	},
	ReturnTypeAnnotation: &sema.TypeAnnotation{
		Type: sema.VoidType,
	},
}

type StandardLibraryInterface interface {
	SetAccountFrozen(address common.Address, frozen bool) error
}

func declareStandardLibraryValues(runtimeEnv runtime.Environment, i StandardLibraryInterface) {
	// TODO return the errors instead of panicing

	runtimeEnv.Declare(stdlib.StandardLibraryValue{
		Name: "setAccountFrozen",
		Type: setAccountFrozenFunctionType,
		Kind: common.DeclarationKindFunction,
		Value: interpreter.NewUnmeteredHostFunctionValue(
			func(invocation interpreter.Invocation) interpreter.Value {
				address, ok := invocation.Arguments[0].(interpreter.AddressValue)
				if !ok {
					panic(errors.NewValueErrorf(invocation.Arguments[0].String(),
						"first argument of setAccountFrozen must be an address"))
				}

				frozen, ok := invocation.Arguments[1].(interpreter.BoolValue)
				if !ok {
					panic(errors.NewValueErrorf(invocation.Arguments[0].String(),
						"second argument of setAccountFrozen must be a boolean"))
				}

				err := i.SetAccountFrozen(common.Address(address), bool(frozen))
				if err != nil {
					panic(err)
				}

				return interpreter.VoidValue{}
			},
			setAccountFrozenFunctionType,
		),
	})
}

// logExecutionIntensities logs execution intensities of the transaction
func (i *TransactionInvoker) logExecutionIntensities(sth *state.StateHolder, txHash string) {
	if i.logger.Debug().Enabled() {
		computation := zerolog.Dict()
		for s, u := range sth.ComputationIntensities() {
			computation.Uint(strconv.FormatUint(uint64(s), 10), u)
		}
		memory := zerolog.Dict()
		for s, u := range sth.MemoryIntensities() {
			memory.Uint(strconv.FormatUint(uint64(s), 10), u)
		}
		i.logger.Info().
			Str("txHash", txHash).
			Uint64("ledgerInteractionUsed", sth.InteractionUsed()).
			Uint("computationUsed", sth.TotalComputationUsed()).
			Uint64("memoryEstimate", sth.TotalMemoryEstimate()).
			Dict("computationIntensities", computation).
			Dict("memoryIntensities", memory).
			Msg("transaction execution data")
	}
}

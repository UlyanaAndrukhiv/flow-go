package handler_test

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/onflow/cadence"

	jsoncdc "github.com/onflow/cadence/encoding/json"

	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	gethParams "github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence/runtime/common"

	"github.com/onflow/flow-go/fvm/errors"
	"github.com/onflow/flow-go/fvm/evm/emulator"
	"github.com/onflow/flow-go/fvm/evm/handler"
	"github.com/onflow/flow-go/fvm/evm/precompiles"
	"github.com/onflow/flow-go/fvm/evm/testutils"
	"github.com/onflow/flow-go/fvm/evm/types"
	"github.com/onflow/flow-go/fvm/systemcontracts"
	"github.com/onflow/flow-go/model/flow"
)

// TODO add test for fatal errors

var flowTokenAddress = common.MustBytesToAddress(systemcontracts.SystemContractsForChain(flow.Emulator).FlowToken.Address.Bytes())

func TestHandler_TransactionRun(t *testing.T) {
	t.Parallel()

	t.Run("test - transaction run (happy case)", func(t *testing.T) {
		t.Parallel()

		testutils.RunWithTestBackend(t, func(backend *testutils.TestBackend) {
			testutils.RunWithTestFlowEVMRootAddress(t, backend, func(rootAddr flow.Address) {
				testutils.RunWithEOATestAccount(t, backend, rootAddr, func(eoa *testutils.EOATestAccount) {

					bs, err := handler.NewBlockStore(backend, rootAddr)
					require.NoError(t, err)

					aa, err := handler.NewAddressAllocator(backend, rootAddr)
					require.NoError(t, err)

					result := &types.Result{
						DeployedContractAddress: types.Address(testutils.RandomAddress(t)),
						ReturnedValue:           testutils.RandomData(t),
						GasConsumed:             testutils.RandomGas(1000),
						Logs: []*gethTypes.Log{
							testutils.GetRandomLogFixture(t),
							testutils.GetRandomLogFixture(t),
						},
					}

					em := &testutils.TestEmulator{
						RunTransactionFunc: func(tx *gethTypes.Transaction) (*types.Result, error) {
							return result, nil
						},
					}
					handler := handler.NewContractHandler(flowTokenAddress, bs, aa, backend, em)

					coinbase := types.NewAddress(gethCommon.Address{})

					tx := eoa.PrepareSignAndEncodeTx(
						t,
						gethCommon.Address{},
						nil,
						nil,
						100_000,
						big.NewInt(1),
					)

					// successfully run (no-panic)
					handler.Run(tx, coinbase)

					// check gas usage
					// TODO: uncomment and investigate me
					// computationUsed, err := backend.ComputationUsed()
					// require.NoError(t, err)
					// require.Equal(t, result.GasConsumed, computationUsed)

					// check events (1 extra for block event)
					events := backend.Events()

					require.Len(t, events, 2)

					event := events[0]
					assert.Equal(t, event.Type, types.EventTypeTransactionExecuted)
					ev, err := jsoncdc.Decode(nil, event.Payload)
					require.NoError(t, err)
					cadenceEvent, ok := ev.(cadence.Event)
					require.True(t, ok)
					for j, f := range cadenceEvent.GetFields() {
						// todo add an event decoder in types.event
						if f.Identifier == "logs" {
							cadenceLogs := cadenceEvent.GetFieldValues()[j]
							encodedLogs, err := hex.DecodeString(strings.ReplaceAll(cadenceLogs.String(), "\"", ""))
							require.NoError(t, err)

							var logs []*gethTypes.Log
							err = rlp.DecodeBytes(encodedLogs, &logs)
							require.NoError(t, err)

							for i, l := range result.Logs {
								assert.Equal(t, l, logs[i])
							}
						}
					}

					// check block event
					event = events[1]
					assert.Equal(t, event.Type, types.EventTypeBlockExecuted)
					_, err = jsoncdc.Decode(nil, event.Payload)
					require.NoError(t, err)
				})
			})
		})
	})

	t.Run("test - transaction run (unhappy cases)", func(t *testing.T) {
		t.Parallel()

		testutils.RunWithTestBackend(t, func(backend *testutils.TestBackend) {
			testutils.RunWithTestFlowEVMRootAddress(t, backend, func(rootAddr flow.Address) {
				testutils.RunWithEOATestAccount(t, backend, rootAddr, func(eoa *testutils.EOATestAccount) {

					bs, err := handler.NewBlockStore(backend, rootAddr)
					require.NoError(t, err)

					aa, err := handler.NewAddressAllocator(backend, rootAddr)
					require.NoError(t, err)

					em := &testutils.TestEmulator{
						RunTransactionFunc: func(tx *gethTypes.Transaction) (*types.Result, error) {
							return &types.Result{}, types.NewEVMExecutionError(fmt.Errorf("some sort of error"))
						},
					}
					handler := handler.NewContractHandler(flowTokenAddress, bs, aa, backend, em)

					coinbase := types.NewAddress(gethCommon.Address{})

					// test RLP decoding (non fatal)
					assertPanic(t, isNotFatal, func() {
						// invalid RLP encoding
						invalidTx := "badencoding"
						handler.Run([]byte(invalidTx), coinbase)
					})

					// test gas limit (non fatal)
					assertPanic(t, isNotFatal, func() {
						gasLimit := uint64(testutils.TestComputationLimit + 1)
						tx := eoa.PrepareSignAndEncodeTx(
							t,
							gethCommon.Address{},
							nil,
							nil,
							gasLimit,
							big.NewInt(1),
						)

						handler.Run([]byte(tx), coinbase)
					})

					// tx execution failure
					tx := eoa.PrepareSignAndEncodeTx(
						t,
						gethCommon.Address{},
						nil,
						nil,
						100_000,
						big.NewInt(1),
					)

					assertPanic(t, isNotFatal, func() {
						handler.Run([]byte(tx), coinbase)
					})
				})
			})
		})
	})

	t.Run("test running transaction (with integrated emulator)", func(t *testing.T) {
		t.Parallel()

		testutils.RunWithTestBackend(t, func(backend *testutils.TestBackend) {
			testutils.RunWithTestFlowEVMRootAddress(t, backend, func(rootAddr flow.Address) {
				handler := SetupHandler(t, backend, rootAddr)

				eoa := testutils.GetTestEOAAccount(t, testutils.EOATestAccount1KeyHex)

				// deposit 1 Flow to the foa account
				addr := handler.DeployCOA()
				orgBalance := types.NewBalanceFromUFix64(types.OneFlowInUFix64)
				vault := types.NewFlowTokenVault(orgBalance)
				foa := handler.AccountByAddress(addr, true)
				foa.Deposit(vault)

				// transfer 0.1 flow to the non-foa address
				deduction := types.NewBalance(big.NewInt(1e17))
				foa.Call(eoa.Address(), nil, 400000, deduction)
				expected, err := types.SubBalance(orgBalance, deduction)
				require.NoError(t, err)
				require.Equal(t, expected, foa.Balance())

				// transfer 0.01 flow back to the foa through
				addition := types.NewBalance(big.NewInt(1e16))

				tx := eoa.PrepareSignAndEncodeTx(
					t,
					foa.Address().ToCommon(),
					nil,
					addition,
					gethParams.TxGas*10,
					big.NewInt(1e8), // high gas fee to test coinbase collection,
				)

				// setup coinbase
				foa2 := handler.DeployCOA()
				account2 := handler.AccountByAddress(foa2, true)
				require.Equal(t, types.NewBalanceFromUFix64(0), account2.Balance())

				// no panic means success here
				handler.Run(tx, account2.Address())
				expected, err = types.SubBalance(orgBalance, deduction)
				require.NoError(t, err)
				expected, err = types.AddBalance(expected, addition)
				require.NoError(t, err)
				require.Equal(t, expected, foa.Balance())

				require.NotEqual(t, types.NewBalanceFromUFix64(0), account2.Balance())
			})
		})
	})
}

func TestHandler_OpsWithoutEmulator(t *testing.T) {
	t.Parallel()

	t.Run("test last executed block call", func(t *testing.T) {
		t.Parallel()

		testutils.RunWithTestBackend(t, func(backend *testutils.TestBackend) {
			testutils.RunWithTestFlowEVMRootAddress(t, backend, func(rootAddr flow.Address) {
				handler := SetupHandler(t, backend, rootAddr)

				// test call last executed block without initialization
				b := handler.LastExecutedBlock()
				require.Equal(t, types.GenesisBlock, b)

				// do some changes
				address := testutils.RandomAddress(t)
				account := handler.AccountByAddress(address, true)
				bal := types.OneFlowBalance
				account.Deposit(types.NewFlowTokenVault(bal))

				// check if block height has been incremented
				b = handler.LastExecutedBlock()
				require.Equal(t, uint64(1), b.Height)
			})
		})
	})

}

func TestHandler_COA(t *testing.T) {
	t.Parallel()
	t.Run("test deposit/withdraw (with integrated emulator)", func(t *testing.T) {
		testutils.RunWithTestBackend(t, func(backend *testutils.TestBackend) {
			testutils.RunWithTestFlowEVMRootAddress(t, backend, func(rootAddr flow.Address) {
				handler := SetupHandler(t, backend, rootAddr)

				foa := handler.AccountByAddress(handler.DeployCOA(), true)
				require.NotNil(t, foa)

				zeroBalance := types.NewBalance(big.NewInt(0))
				require.Equal(t, zeroBalance, foa.Balance())

				balance := types.OneFlowBalance
				vault := types.NewFlowTokenVault(balance)

				foa.Deposit(vault)
				require.Equal(t, balance, foa.Balance())

				v := foa.Withdraw(balance)
				require.Equal(t, balance, v.Balance())

				require.Equal(t, zeroBalance, foa.Balance())

				events := backend.Events()
				require.Len(t, events, 6)

				// first two transactions are for COA setup

				// transaction event
				event := events[2]
				assert.Equal(t, event.Type, types.EventTypeTransactionExecuted)

				// block event
				event = events[3]
				assert.Equal(t, event.Type, types.EventTypeBlockExecuted)

				// transaction event
				event = events[4]
				assert.Equal(t, event.Type, types.EventTypeTransactionExecuted)
				_, err := jsoncdc.Decode(nil, event.Payload)
				require.NoError(t, err)
				// TODO: decode encoded tx and check for the amount and value
				// assert.Equal(t, foa.Address(), ret.Address)
				// assert.Equal(t, balance, ret.Amount)

				// block event
				event = events[5]
				assert.Equal(t, event.Type, types.EventTypeBlockExecuted)

				// check gas usage
				computationUsed, err := backend.ComputationUsed()
				require.NoError(t, err)
				require.Greater(t, computationUsed, types.DefaultDirectCallBaseGasUsage*3)

				// Withdraw with invalid balance
				assertPanic(t, types.IsWithdrawBalanceRoundingError, func() {
					// deposit some money
					foa.Deposit(vault)
					// then withdraw invalid balance
					foa.Withdraw(types.NewBalance(big.NewInt(1)))
				})
			})
		})
	})

	t.Run("test coa deployment", func(t *testing.T) {
		testutils.RunWithTestBackend(t, func(backend *testutils.TestBackend) {
			testutils.RunWithTestFlowEVMRootAddress(t, backend, func(rootAddr flow.Address) {
				h := SetupHandler(t, backend, rootAddr)

				coa := h.DeployCOA()
				acc := h.AccountByAddress(coa, true)
				require.NotEmpty(t, acc.Code())

				// make a second account with some money
				coa2 := h.DeployCOA()
				acc2 := h.AccountByAddress(coa2, true)
				acc2.Deposit(types.NewFlowTokenVault(types.MakeABalanceInFlow(100)))

				// transfer money to COA
				acc2.Transfer(
					coa,
					types.MakeABalanceInFlow(1),
				)

				// make a call to the contract
				ret := acc2.Call(
					coa,
					testutils.MakeCallData(t,
						handler.COAContractABIJSON,
						"onERC721Received",
						gethCommon.Address{1},
						gethCommon.Address{1},
						big.NewInt(0),
						[]byte{'A'},
					),
					types.GasLimit(3_000_000),
					types.EmptyBalance)

				// 0x150b7a02
				expected := types.Data([]byte{
					21, 11, 122, 2, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				})
				require.Equal(t, expected, ret)
			})
		})
	})

	t.Run("test withdraw (unhappy case)", func(t *testing.T) {
		testutils.RunWithTestBackend(t, func(backend *testutils.TestBackend) {
			testutils.RunWithTestFlowEVMRootAddress(t, backend, func(rootAddr flow.Address) {
				testutils.RunWithEOATestAccount(t, backend, rootAddr, func(eoa *testutils.EOATestAccount) {
					bs, err := handler.NewBlockStore(backend, rootAddr)
					require.NoError(t, err)

					aa, err := handler.NewAddressAllocator(backend, rootAddr)
					require.NoError(t, err)

					// Withdraw calls are only possible within FOA accounts
					assertPanic(t, types.IsAUnAuthroizedMethodCallError, func() {
						em := &testutils.TestEmulator{}

						handler := handler.NewContractHandler(flowTokenAddress, bs, aa, backend, em)

						account := handler.AccountByAddress(testutils.RandomAddress(t), false)
						account.Withdraw(types.NewBalanceFromUFix64(1))
					})

					// test insufficient total supply
					assertPanic(t, types.IsAInsufficientTotalSupplyError, func() {
						em := &testutils.TestEmulator{
							DirectCallFunc: func(call *types.DirectCall) (*types.Result, error) {
								return &types.Result{}, types.NewEVMExecutionError(fmt.Errorf("some sort of error"))
							},
						}

						handler := handler.NewContractHandler(flowTokenAddress, bs, aa, backend, em)
						account := handler.AccountByAddress(testutils.RandomAddress(t), true)

						account.Withdraw(types.NewBalanceFromUFix64(1))
					})

					// test non fatal error of emulator
					assertPanic(t, types.IsEVMExecutionError, func() {
						em := &testutils.TestEmulator{
							DirectCallFunc: func(call *types.DirectCall) (*types.Result, error) {
								return &types.Result{}, types.NewEVMExecutionError(fmt.Errorf("some sort of error"))
							},
						}

						handler := handler.NewContractHandler(flowTokenAddress, bs, aa, backend, em)
						account := handler.AccountByAddress(testutils.RandomAddress(t), true)

						account.Withdraw(types.NewBalanceFromUFix64(0))
					})

					// test fatal error of emulator
					assertPanic(t, types.IsAFatalError, func() {
						em := &testutils.TestEmulator{
							DirectCallFunc: func(call *types.DirectCall) (*types.Result, error) {
								return &types.Result{}, types.NewFatalError(fmt.Errorf("some sort of fatal error"))
							},
						}

						handler := handler.NewContractHandler(flowTokenAddress, bs, aa, backend, em)
						account := handler.AccountByAddress(testutils.RandomAddress(t), true)

						account.Withdraw(types.NewBalanceFromUFix64(0))
					})
				})
			})
		})
	})

	t.Run("test deposit (unhappy case)", func(t *testing.T) {
		t.Parallel()

		testutils.RunWithTestBackend(t, func(backend *testutils.TestBackend) {
			testutils.RunWithTestFlowEVMRootAddress(t, backend, func(rootAddr flow.Address) {
				testutils.RunWithEOATestAccount(t, backend, rootAddr, func(eoa *testutils.EOATestAccount) {
					bs, err := handler.NewBlockStore(backend, rootAddr)
					require.NoError(t, err)

					aa, err := handler.NewAddressAllocator(backend, rootAddr)
					require.NoError(t, err)

					// test non fatal error of emulator
					assertPanic(t, types.IsEVMExecutionError, func() {
						em := &testutils.TestEmulator{
							DirectCallFunc: func(call *types.DirectCall) (*types.Result, error) {
								return &types.Result{}, types.NewEVMExecutionError(fmt.Errorf("some sort of error"))
							},
						}

						handler := handler.NewContractHandler(flowTokenAddress, bs, aa, backend, em)
						account := handler.AccountByAddress(testutils.RandomAddress(t), true)

						account.Deposit(types.NewFlowTokenVault(types.NewBalanceFromUFix64(1)))
					})

					// test fatal error of emulator
					assertPanic(t, types.IsAFatalError, func() {
						em := &testutils.TestEmulator{
							DirectCallFunc: func(call *types.DirectCall) (*types.Result, error) {
								return &types.Result{}, types.NewFatalError(fmt.Errorf("some sort of fatal error"))
							},
						}

						handler := handler.NewContractHandler(flowTokenAddress, bs, aa, backend, em)
						account := handler.AccountByAddress(testutils.RandomAddress(t), true)

						account.Deposit(types.NewFlowTokenVault(types.NewBalanceFromUFix64(1)))
					})
				})
			})
		})
	})

	t.Run("test deploy/call (with integrated emulator)", func(t *testing.T) {
		t.Parallel()

		// TODO update this test with events, gas metering, etc
		testutils.RunWithTestBackend(t, func(backend *testutils.TestBackend) {
			testutils.RunWithTestFlowEVMRootAddress(t, backend, func(rootAddr flow.Address) {
				handler := SetupHandler(t, backend, rootAddr)

				foa := handler.AccountByAddress(handler.DeployCOA(), true)
				require.NotNil(t, foa)

				// deposit 10000 flow
				bal := types.MakeABalanceInFlow(10000)
				vault := types.NewFlowTokenVault(bal)
				foa.Deposit(vault)
				require.Equal(t, bal, foa.Balance())

				testContract := testutils.GetStorageTestContract(t)
				addr := foa.Deploy(testContract.ByteCode, math.MaxUint64, types.NewBalanceFromUFix64(0))
				require.NotNil(t, addr)

				num := big.NewInt(22)
				_ = foa.Call(
					addr,
					testContract.MakeCallData(t, "store", num),
					math.MaxUint64,
					types.NewBalanceFromUFix64(0))

				ret := foa.Call(
					addr,
					testContract.MakeCallData(t, "retrieve"),
					math.MaxUint64,
					types.NewBalanceFromUFix64(0))

				require.Equal(t, num, new(big.Int).SetBytes(ret))
			})
		})
	})

	t.Run("test call to cadence arch", func(t *testing.T) {
		t.Parallel()

		testutils.RunWithTestBackend(t, func(backend *testutils.TestBackend) {
			blockHeight := uint64(123)
			backend.GetCurrentBlockHeightFunc = func() (uint64, error) {
				return blockHeight, nil
			}
			testutils.RunWithTestFlowEVMRootAddress(t, backend, func(rootAddr flow.Address) {
				h := SetupHandler(t, backend, rootAddr)

				foa := h.AccountByAddress(h.DeployCOA(), true)
				require.NotNil(t, foa)

				vault := types.NewFlowTokenVault(types.MakeABalanceInFlow(10000))
				foa.Deposit(vault)

				arch := handler.MakePrecompileAddress(1)

				ret := foa.Call(arch, precompiles.FlowBlockHeightFuncSig[:], math.MaxUint64, types.NewBalanceFromUFix64(0))
				require.Equal(t, big.NewInt(int64(blockHeight)), new(big.Int).SetBytes(ret))
			})
		})
	})

	// TODO add test with test emulator for unhappy cases (emulator)
}

// returns true if error passes the checks
type checkError func(error) bool

var isNotFatal = func(err error) bool {
	return !errors.IsFailure(err)
}

func assertPanic(t *testing.T, check checkError, f func()) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("The code did not panic")
		}
		err, ok := r.(error)
		if !ok {
			t.Fatal("panic is not with an error type")
		}
		require.True(t, check(err))
	}()
	f()
}

func SetupHandler(t testing.TB, backend types.Backend, rootAddr flow.Address) *handler.ContractHandler {
	bs, err := handler.NewBlockStore(backend, rootAddr)
	require.NoError(t, err)

	aa, err := handler.NewAddressAllocator(backend, rootAddr)
	require.NoError(t, err)

	emulator := emulator.NewEmulator(backend, rootAddr)

	handler := handler.NewContractHandler(flowTokenAddress, bs, aa, backend, emulator)
	return handler
}

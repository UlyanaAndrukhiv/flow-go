package testutils

import (
	"github.com/onflow/flow-go/fvm/storage"
	"github.com/onflow/flow-go/fvm/storage/derived"
	"github.com/onflow/flow-go/fvm/storage/snapshot"
	"github.com/onflow/flow-go/fvm/storage/state"
)

// NewSimpleTransaction returns a transaction which can be used to test
// fvm evaluation.  The returned transaction should not be committed.
func NewSimpleTransaction(
	snapshot snapshot.StorageSnapshot,
) *storage.SerialTransaction {
	derivedBlockData := derived.NewEmptyDerivedBlockData()
	derivedTxnData, err := derivedBlockData.NewDerivedTransactionData(0, 0)
	if err != nil {
		panic(err)
	}

	return &storage.SerialTransaction{
		NestedTransactionPreparer: state.NewTransactionState(
			snapshot,
			state.DefaultParameters()),
		DerivedTransactionData: derivedTxnData,
	}
}

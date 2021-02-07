// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	context "context"

	delta "github.com/onflow/flow-go/engine/execution/state/delta"
	flow "github.com/onflow/flow-go/model/flow"

	messages "github.com/onflow/flow-go/model/messages"

	mock "github.com/stretchr/testify/mock"
)

// ExecutionState is an autogenerated mock type for the ExecutionState type
type ExecutionState struct {
	mock.Mock
}

// ChunkDataPackByChunkID provides a mock function with given fields: _a0, _a1
func (_m *ExecutionState) ChunkDataPackByChunkID(_a0 context.Context, _a1 flow.Identifier) (*flow.ChunkDataPack, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *flow.ChunkDataPack
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *flow.ChunkDataPack); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ChunkDataPack)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CommitDelta provides a mock function with given fields: _a0, _a1, _a2
func (_m *ExecutionState) CommitDelta(_a0 context.Context, _a1 delta.Delta, _a2 []byte) ([]byte, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, delta.Delta, []byte) []byte); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, delta.Delta, []byte) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCollection provides a mock function with given fields: identifier
func (_m *ExecutionState) GetCollection(identifier flow.Identifier) (*flow.Collection, error) {
	ret := _m.Called(identifier)

	var r0 *flow.Collection
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.Collection); ok {
		r0 = rf(identifier)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Collection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(identifier)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExecutionResultID provides a mock function with given fields: _a0, _a1
func (_m *ExecutionState) GetExecutionResultID(_a0 context.Context, _a1 flow.Identifier) (flow.Identifier, error) {
	ret := _m.Called(_a0, _a1)

	var r0 flow.Identifier
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) flow.Identifier); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Identifier)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHighestExecutedBlockID provides a mock function with given fields: _a0
func (_m *ExecutionState) GetHighestExecutedBlockID(_a0 context.Context) (uint64, flow.Identifier, error) {
	ret := _m.Called(_a0)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context) uint64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 flow.Identifier
	if rf, ok := ret.Get(1).(func(context.Context) flow.Identifier); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(flow.Identifier)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetProof provides a mock function with given fields: _a0, _a1, _a2
func (_m *ExecutionState) GetProof(_a0 context.Context, _a1 []byte, _a2 []flow.RegisterID) ([]byte, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, []byte, []flow.RegisterID) []byte); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []byte, []flow.RegisterID) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRegisters provides a mock function with given fields: _a0, _a1, _a2
func (_m *ExecutionState) GetRegisters(_a0 context.Context, _a1 []byte, _a2 []flow.RegisterID) ([][]byte, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func(context.Context, []byte, []flow.RegisterID) [][]byte); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []byte, []flow.RegisterID) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewView provides a mock function with given fields: _a0
func (_m *ExecutionState) NewView(_a0 []byte) *delta.View {
	ret := _m.Called(_a0)

	var r0 *delta.View
	if rf, ok := ret.Get(0).(func([]byte) *delta.View); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*delta.View)
		}
	}

	return r0
}

// PersistExecutionState provides a mock function with given fields: ctx, header, endState, chunkDataPacks, executionResult, events, serviceEvents, results
func (_m *ExecutionState) PersistExecutionState(ctx context.Context, header *flow.Header, endState []byte, chunkDataPacks []*flow.ChunkDataPack, executionResult *flow.ExecutionResult, events []flow.Event, serviceEvents []flow.Event, results []flow.TransactionResult) error {
	ret := _m.Called(ctx, header, endState, chunkDataPacks, executionResult, events, serviceEvents, results)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *flow.Header, []byte, []*flow.ChunkDataPack, *flow.ExecutionResult, []flow.Event, []flow.Event, []flow.TransactionResult) error); ok {
		r0 = rf(ctx, header, endState, chunkDataPacks, executionResult, events, serviceEvents, results)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RetrieveStateDelta provides a mock function with given fields: _a0, _a1
func (_m *ExecutionState) RetrieveStateDelta(_a0 context.Context, _a1 flow.Identifier) (*messages.ExecutionStateDelta, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *messages.ExecutionStateDelta
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *messages.ExecutionStateDelta); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*messages.ExecutionStateDelta)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StateCommitmentByBlockID provides a mock function with given fields: _a0, _a1
func (_m *ExecutionState) StateCommitmentByBlockID(_a0 context.Context, _a1 flow.Identifier) ([]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateHighestExecutedBlockIfHigher provides a mock function with given fields: _a0, _a1
func (_m *ExecutionState) UpdateHighestExecutedBlockIfHigher(_a0 context.Context, _a1 *flow.Header) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *flow.Header) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// BlockProducer is an autogenerated mock type for the BlockProducer type
type BlockProducer struct {
	mock.Mock
}

// MakeBlockProposal provides a mock function with given fields: view, qc, lastViewTC
func (_m *BlockProducer) MakeBlockProposal(view uint64, qc *flow.QuorumCertificate, lastViewTC *flow.TimeoutCertificate) (*flow.Header, error) {
	ret := _m.Called(view, qc, lastViewTC)

	var r0 *flow.Header
	if rf, ok := ret.Get(0).(func(uint64, *flow.QuorumCertificate, *flow.TimeoutCertificate) *flow.Header); ok {
		r0 = rf(view, qc, lastViewTC)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, *flow.QuorumCertificate, *flow.TimeoutCertificate) error); ok {
		r1 = rf(view, qc, lastViewTC)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewBlockProducer interface {
	mock.TestingT
	Cleanup(func())
}

// NewBlockProducer creates a new instance of BlockProducer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBlockProducer(t mockConstructorTestingTNewBlockProducer) *BlockProducer {
	mock := &BlockProducer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

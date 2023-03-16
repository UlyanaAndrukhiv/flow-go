// Code generated by mockery v2.21.4. DO NOT EDIT.

package mockp2p

import (
	irrecoverable "github.com/onflow/flow-go/module/irrecoverable"
	mock "github.com/stretchr/testify/mock"

	peer "github.com/libp2p/go-libp2p/core/peer"

	protocol "github.com/libp2p/go-libp2p/core/protocol"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// PubSubTracer is an autogenerated mock type for the PubSubTracer type
type PubSubTracer struct {
	mock.Mock
}

// AddPeer provides a mock function with given fields: p, proto
func (_m *PubSubTracer) AddPeer(p peer.ID, proto protocol.ID) {
	_m.Called(p, proto)
}

// DeliverMessage provides a mock function with given fields: msg
func (_m *PubSubTracer) DeliverMessage(msg *pubsub.Message) {
	_m.Called(msg)
}

// Done provides a mock function with given fields:
func (_m *PubSubTracer) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// DropRPC provides a mock function with given fields: rpc, p
func (_m *PubSubTracer) DropRPC(rpc *pubsub.RPC, p peer.ID) {
	_m.Called(rpc, p)
}

// DuplicateMessage provides a mock function with given fields: msg
func (_m *PubSubTracer) DuplicateMessage(msg *pubsub.Message) {
	_m.Called(msg)
}

// Graft provides a mock function with given fields: p, topic
func (_m *PubSubTracer) Graft(p peer.ID, topic string) {
	_m.Called(p, topic)
}

// Join provides a mock function with given fields: topic
func (_m *PubSubTracer) Join(topic string) {
	_m.Called(topic)
}

// Leave provides a mock function with given fields: topic
func (_m *PubSubTracer) Leave(topic string) {
	_m.Called(topic)
}

// Prune provides a mock function with given fields: p, topic
func (_m *PubSubTracer) Prune(p peer.ID, topic string) {
	_m.Called(p, topic)
}

// Ready provides a mock function with given fields:
func (_m *PubSubTracer) Ready() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// RecvRPC provides a mock function with given fields: rpc
func (_m *PubSubTracer) RecvRPC(rpc *pubsub.RPC) {
	_m.Called(rpc)
}

// RejectMessage provides a mock function with given fields: msg, reason
func (_m *PubSubTracer) RejectMessage(msg *pubsub.Message, reason string) {
	_m.Called(msg, reason)
}

// RemovePeer provides a mock function with given fields: p
func (_m *PubSubTracer) RemovePeer(p peer.ID) {
	_m.Called(p)
}

// SendRPC provides a mock function with given fields: rpc, p
func (_m *PubSubTracer) SendRPC(rpc *pubsub.RPC, p peer.ID) {
	_m.Called(rpc, p)
}

// Start provides a mock function with given fields: _a0
func (_m *PubSubTracer) Start(_a0 irrecoverable.SignalerContext) {
	_m.Called(_a0)
}

// ThrottlePeer provides a mock function with given fields: p
func (_m *PubSubTracer) ThrottlePeer(p peer.ID) {
	_m.Called(p)
}

// UndeliverableMessage provides a mock function with given fields: msg
func (_m *PubSubTracer) UndeliverableMessage(msg *pubsub.Message) {
	_m.Called(msg)
}

// ValidateMessage provides a mock function with given fields: msg
func (_m *PubSubTracer) ValidateMessage(msg *pubsub.Message) {
	_m.Called(msg)
}

type mockConstructorTestingTNewPubSubTracer interface {
	mock.TestingT
	Cleanup(func())
}

// NewPubSubTracer creates a new instance of PubSubTracer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPubSubTracer(t mockConstructorTestingTNewPubSubTracer) *PubSubTracer {
	mock := &PubSubTracer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

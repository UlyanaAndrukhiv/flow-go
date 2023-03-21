// Code generated by mockery v2.21.4. DO NOT EDIT.

package mockp2p

import (
	host "github.com/libp2p/go-libp2p/core/host"
	channels "github.com/onflow/flow-go/network/channels"

	irrecoverable "github.com/onflow/flow-go/module/irrecoverable"

	mock "github.com/stretchr/testify/mock"

	module "github.com/onflow/flow-go/module"

	p2p "github.com/onflow/flow-go/network/p2p"

	peer "github.com/libp2p/go-libp2p/core/peer"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	routing "github.com/libp2p/go-libp2p/core/routing"

	time "time"
)

// GossipSubBuilder is an autogenerated mock type for the GossipSubBuilder type
type GossipSubBuilder struct {
	mock.Mock
}

// Build provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) Build(_a0 irrecoverable.SignalerContext) (p2p.PubSubAdapter, p2p.PeerScoreTracer, error) {
	ret := _m.Called(_a0)

	var r0 p2p.PubSubAdapter
	var r1 p2p.PeerScoreTracer
	var r2 error
	if rf, ok := ret.Get(0).(func(irrecoverable.SignalerContext) (p2p.PubSubAdapter, p2p.PeerScoreTracer, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(irrecoverable.SignalerContext) p2p.PubSubAdapter); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.PubSubAdapter)
		}
	}

	if rf, ok := ret.Get(1).(func(irrecoverable.SignalerContext) p2p.PeerScoreTracer); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(p2p.PeerScoreTracer)
		}
	}

	if rf, ok := ret.Get(2).(func(irrecoverable.SignalerContext) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SetAppSpecificScoreParams provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) SetAppSpecificScoreParams(_a0 func(peer.ID) float64) {
	_m.Called(_a0)
}

// SetGossipSubConfigFunc provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) SetGossipSubConfigFunc(_a0 p2p.GossipSubAdapterConfigFunc) {
	_m.Called(_a0)
}

// SetGossipSubFactory provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) SetGossipSubFactory(_a0 p2p.GossipSubFactoryFunc) {
	_m.Called(_a0)
}

// SetGossipSubPeerScoring provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) SetGossipSubPeerScoring(_a0 bool) {
	_m.Called(_a0)
}

// SetGossipSubScoreTracerInterval provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) SetGossipSubScoreTracerInterval(_a0 time.Duration) {
	_m.Called(_a0)
}

// SetGossipSubTracer provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) SetGossipSubTracer(_a0 p2p.PubSubTracer) {
	_m.Called(_a0)
}

// SetHost provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) SetHost(_a0 host.Host) {
	_m.Called(_a0)
}

// SetIDProvider provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) SetIDProvider(_a0 module.IdentityProvider) {
	_m.Called(_a0)
}

// SetRoutingSystem provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) SetRoutingSystem(_a0 routing.Routing) {
	_m.Called(_a0)
}

// SetSubscriptionFilter provides a mock function with given fields: _a0
func (_m *GossipSubBuilder) SetSubscriptionFilter(_a0 pubsub.SubscriptionFilter) {
	_m.Called(_a0)
}

// SetTopicScoreParams provides a mock function with given fields: topic, topicScoreParams
func (_m *GossipSubBuilder) SetTopicScoreParams(topic channels.Topic, topicScoreParams *pubsub.TopicScoreParams) {
	_m.Called(topic, topicScoreParams)
}

type mockConstructorTestingTNewGossipSubBuilder interface {
	mock.TestingT
	Cleanup(func())
}

// NewGossipSubBuilder creates a new instance of GossipSubBuilder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGossipSubBuilder(t mockConstructorTestingTNewGossipSubBuilder) *GossipSubBuilder {
	mock := &GossipSubBuilder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
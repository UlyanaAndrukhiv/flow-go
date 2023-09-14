// Code generated by mockery v2.21.4. DO NOT EDIT.

package mockp2p

import (
	flow "github.com/onflow/flow-go/model/flow"
	irrecoverable "github.com/onflow/flow-go/module/irrecoverable"

	metrics "github.com/onflow/flow-go/module/metrics"

	mock "github.com/stretchr/testify/mock"

	module "github.com/onflow/flow-go/module"

	network "github.com/onflow/flow-go/network"

	p2p "github.com/onflow/flow-go/network/p2p"

	p2pconf "github.com/onflow/flow-go/network/p2p/p2pconf"

	zerolog "github.com/rs/zerolog"
)

// GossipSubRpcInspectorSuiteFactoryFunc is an autogenerated mock type for the GossipSubRpcInspectorSuiteFactoryFunc type
type GossipSubRpcInspectorSuiteFactoryFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5, _a6, _a7
func (_m *GossipSubRpcInspectorSuiteFactoryFunc) Execute(_a0 irrecoverable.SignalerContext, _a1 zerolog.Logger, _a2 flow.Identifier, _a3 *p2pconf.GossipSubRPCInspectorsConfig, _a4 module.GossipSubMetrics, _a5 metrics.HeroCacheMetricsFactory, _a6 network.NetworkingType, _a7 module.IdentityProvider) (p2p.GossipSubInspectorSuite, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5, _a6, _a7)

	var r0 p2p.GossipSubInspectorSuite
	var r1 error
	if rf, ok := ret.Get(0).(func(irrecoverable.SignalerContext, zerolog.Logger, flow.Identifier, *p2pconf.GossipSubRPCInspectorsConfig, module.GossipSubMetrics, metrics.HeroCacheMetricsFactory, network.NetworkingType, module.IdentityProvider) (p2p.GossipSubInspectorSuite, error)); ok {
		return rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6, _a7)
	}
	if rf, ok := ret.Get(0).(func(irrecoverable.SignalerContext, zerolog.Logger, flow.Identifier, *p2pconf.GossipSubRPCInspectorsConfig, module.GossipSubMetrics, metrics.HeroCacheMetricsFactory, network.NetworkingType, module.IdentityProvider) p2p.GossipSubInspectorSuite); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6, _a7)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.GossipSubInspectorSuite)
		}
	}

	if rf, ok := ret.Get(1).(func(irrecoverable.SignalerContext, zerolog.Logger, flow.Identifier, *p2pconf.GossipSubRPCInspectorsConfig, module.GossipSubMetrics, metrics.HeroCacheMetricsFactory, network.NetworkingType, module.IdentityProvider) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6, _a7)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGossipSubRpcInspectorSuiteFactoryFunc interface {
	mock.TestingT
	Cleanup(func())
}

// NewGossipSubRpcInspectorSuiteFactoryFunc creates a new instance of GossipSubRpcInspectorSuiteFactoryFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGossipSubRpcInspectorSuiteFactoryFunc(t mockConstructorTestingTNewGossipSubRpcInspectorSuiteFactoryFunc) *GossipSubRpcInspectorSuiteFactoryFunc {
	mock := &GossipSubRpcInspectorSuiteFactoryFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
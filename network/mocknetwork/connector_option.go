// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocknetwork

import (
	p2p "github.com/onflow/flow-go/network/p2p"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// ConnectorOption is an autogenerated mock type for the ConnectorOption type
type ConnectorOption struct {
	mock.Mock
}

// Execute provides a mock function with given fields: connector
func (_m *ConnectorOption) Execute(connector *p2p.Libp2pConnector) {
	_m.Called(connector)
}

// NewConnectorOption creates a new instance of ConnectorOption. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewConnectorOption(t testing.TB) *ConnectorOption {
	mock := &ConnectorOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
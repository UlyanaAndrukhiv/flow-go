// Code generated by mockery v2.13.0. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// AccessMetrics is an autogenerated mock type for the AccessMetrics type
type AccessMetrics struct {
	mock.Mock
}

// ConnectionAddedToPool provides a mock function with given fields:
func (_m *AccessMetrics) ConnectionAddedToPool() {
	_m.Called()
}

// ConnectionFromPoolEvicted provides a mock function with given fields:
func (_m *AccessMetrics) ConnectionFromPoolEvicted() {
	_m.Called()
}

// ConnectionFromPoolInvalidated provides a mock function with given fields:
func (_m *AccessMetrics) ConnectionFromPoolInvalidated() {
	_m.Called()
}

// ConnectionFromPoolReused provides a mock function with given fields:
func (_m *AccessMetrics) ConnectionFromPoolReused() {
	_m.Called()
}

// ConnectionFromPoolUpdated provides a mock function with given fields:
func (_m *AccessMetrics) ConnectionFromPoolUpdated() {
	_m.Called()
}

// NewConnectionEstablished provides a mock function with given fields:
func (_m *AccessMetrics) NewConnectionEstablished() {
	_m.Called()
}

// TotalConnectionsInPool provides a mock function with given fields: connectionCount, connectionPoolSize
func (_m *AccessMetrics) TotalConnectionsInPool(connectionCount uint, connectionPoolSize uint) {
	_m.Called(connectionCount, connectionPoolSize)
}

type NewAccessMetricsT interface {
	mock.TestingT
	Cleanup(func())
}

// NewAccessMetrics creates a new instance of AccessMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAccessMetrics(t NewAccessMetricsT) *AccessMetrics {
	mock := &AccessMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
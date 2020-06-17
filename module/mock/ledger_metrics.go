// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// LedgerMetrics is an autogenerated mock type for the LedgerMetrics type
type LedgerMetrics struct {
	mock.Mock
}

// ForestApproxMemorySize provides a mock function with given fields: bytes
func (_m *LedgerMetrics) ForestApproxMemorySize(bytes uint64) {
	_m.Called(bytes)
}

// ForestNumberOfTrees provides a mock function with given fields: number
func (_m *LedgerMetrics) ForestNumberOfTrees(number uint64) {
	_m.Called(number)
}

// ProofSize provides a mock function with given fields: bytes
func (_m *LedgerMetrics) ProofSize(bytes uint32) {
	_m.Called(bytes)
}

// ReadDuration provides a mock function with given fields: duration
func (_m *LedgerMetrics) ReadDuration(duration time.Duration) {
	_m.Called(duration)
}

// ReadDurationPerItem provides a mock function with given fields: duration
func (_m *LedgerMetrics) ReadDurationPerItem(duration time.Duration) {
	_m.Called(duration)
}

// ReadValuesNumber provides a mock function with given fields: number
func (_m *LedgerMetrics) ReadValuesNumber(number uint64) {
	_m.Called(number)
}

// ReadValuesSize provides a mock function with given fields: byte
func (_m *LedgerMetrics) ReadValuesSize(byte uint64) {
	_m.Called(byte)
}

// UpdateCount provides a mock function with given fields:
func (_m *LedgerMetrics) UpdateCount() {
	_m.Called()
}

// UpdateDuration provides a mock function with given fields: duration
func (_m *LedgerMetrics) UpdateDuration(duration time.Duration) {
	_m.Called(duration)
}

// UpdateDurationPerItem provides a mock function with given fields: duration
func (_m *LedgerMetrics) UpdateDurationPerItem(duration time.Duration) {
	_m.Called(duration)
}

// UpdateValuesNumber provides a mock function with given fields: number
func (_m *LedgerMetrics) UpdateValuesNumber(number uint64) {
	_m.Called(number)
}

// UpdateValuesSize provides a mock function with given fields: byte
func (_m *LedgerMetrics) UpdateValuesSize(byte uint64) {
	_m.Called(byte)
}

package irrecoverable

import (
	"context"
	"testing"
)

// MockSignalerContext is a SignalerContext which will immediately fail a test if an error is thrown.
type MockSignalerContext struct {
	context.Context
	tb testing.TB
}

var _ SignalerContext = &MockSignalerContext{}

func (m MockSignalerContext) sealed() {}

func (m MockSignalerContext) Throw(err error) {
	m.tb.Fatalf("mock signaler context received error: %v", err)
}

func NewMockSignalerContext(tb testing.TB, ctx context.Context) *MockSignalerContext {
	return &MockSignalerContext{
		Context: ctx,
		tb:      tb,
	}
}

func NewMockSignalerContextWithCancel(t *testing.T, parent context.Context) (*MockSignalerContext, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	return NewMockSignalerContext(t, ctx), cancel
}

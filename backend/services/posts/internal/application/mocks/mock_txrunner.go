package mocks

import (
	"context"
	"social-network/services/posts/internal/db/mocks"

	"github.com/stretchr/testify/mock"
)

// MockTxRunner is a simple mock usable in tests. It can be given a Querier
// instance which will be passed to the transactional function when RunTx is
// called. Callers can set the Querier field to a *mocks.MockQuerier instance.
type MockTxRunner struct {
	mock.Mock
	Queries *mocks.MockQueries
}

func (m *MockTxRunner) RunTx(ctx context.Context, fn func(*mocks.MockQueries) error) error {
	m.Called(ctx)
	if m.Queries != nil {
		return fn(m.Queries)
	}
	return nil
}

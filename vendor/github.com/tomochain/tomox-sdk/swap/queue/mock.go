package queue

import (
	"github.com/stretchr/testify/mock"
	"github.com/tomochain/tomox-sdk/types"
)

// MockQueue is a mockable queue.
type MockQueue struct {
	mock.Mock
}

func (m *MockQueue) QueueAdd(tx *types.DepositTransaction) error {
	a := m.Called(tx)
	return a.Error(0)
}

func (m *MockQueue) QueuePool() (<-chan *types.DepositTransaction, error) {
	a := m.Called()
	return a.Get(0).(<-chan *types.DepositTransaction), a.Error(1)
}

package test_helpers

import (
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
	"github.com/stretchr/testify/mock"
)

// MockQuerier is a mock of Querier interface
type MockQuerier struct {
	mock.Mock
}

func NewMockQuerier() *MockQuerier {
	return &MockQuerier{}
}

func (m *MockQuerier) Close() {
	m.Called()
}

func (m *MockQuerier) SaveDevice(device domain.Device) error {
	args := m.Called(device)
	return args.Error(0)
}

func (m *MockQuerier) GetDevices() ([]domain.Device, error) {
	args := m.Called()
	if arg := args.Get(0); arg != nil {
		return arg.([]domain.Device), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockQuerier) GetDevice(id uuid.UUID) (*domain.Device, error) {
	args := m.Called(id)
	if arg := args.Get(0); arg != nil {
		return arg.(*domain.Device), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockQuerier) UpdateDevice(device domain.Device) error {
	args := m.Called(device)
	return args.Error(0)
}

func (m *MockQuerier) SaveSignedTransaction(transaction domain.SignedTransaction) (uuid.UUID, error) {
	args := m.Called(transaction)
	if arg := args.Get(0); arg != nil {
		return arg.(uuid.UUID), args.Error(1)
	}
	return uuid.UUID{}, args.Error(1)
}

func (q *MockQuerier) GetSignedTransaction(deviceId uuid.UUID, signCounter int) (*domain.SignedTransaction, error) {
	args := q.Called(deviceId, signCounter)
	if arg := args.Get(0); arg != nil {
		return arg.(*domain.SignedTransaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockQuerier) GetSignedTransactions(deviceId uuid.UUID) ([]domain.SignedTransaction, error) {
	args := m.Called(deviceId)
	if arg := args.Get(0); arg != nil {
		return arg.([]domain.SignedTransaction), args.Error(1)
	}
	return nil, args.Error(1)
}

package test_helpers

import (
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
	"github.com/stretchr/testify/mock"
)

// mockDeviceDAO is a mock type for the DeviceDao type
type mockDeviceDAO struct {
	mock.Mock
}

// NewMockDeviceDAO creates a new instance of MockManager
func NewMockDeviceDAO() *mockDeviceDAO {
	return &mockDeviceDAO{}
}

func (m *mockDeviceDAO) CreateDevice(id uuid.UUID, label, algorithm string) (*domain.Device, error) {
	args := m.Called(id, label, algorithm)
	return args.Get(0).(*domain.Device), args.Error(1)
}

func (m *mockDeviceDAO) GetDevices() ([]domain.Device, error) {
	args := m.Called()
	return args.Get(0).([]domain.Device), args.Error(1)
}

func (m *mockDeviceDAO) GetDevice(id uuid.UUID) (*domain.Device, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Device), args.Error(1)
}

func (m *mockDeviceDAO) CreateSignedTransaction(deviceId uuid.UUID, data []byte) (*domain.SignedTransaction, error) {
	args := m.Called(deviceId, data)
	return args.Get(0).(*domain.SignedTransaction), args.Error(1)
}

func (m *mockDeviceDAO) GetSignedTransactions(deviceId uuid.UUID) ([]domain.SignedTransaction, error) {
	args := m.Called(deviceId)
	return args.Get(0).([]domain.SignedTransaction), args.Error(1)
}

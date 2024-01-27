package service

import (
	"context"
	"github.com/ildomm/ssccg/domain"
	"github.com/ildomm/ssccg/persistence"
	"github.com/ildomm/ssccg/test_helpers"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateDevice(t *testing.T) {
	ctx := context.TODO()
	mockQuerier := test_helpers.NewMockQuerier()
	sm := NewDeviceManager(mockQuerier)

	id := uuid.New()
	device := domain.Device{ID: id, Label: "Test Device", SignAlgorithm: "RSA"}

	t.Run("SuccessfulCreation", func(t *testing.T) {
		mockQuerier.On("GetDevice", id).Return(nil, nil).Once()
		mockQuerier.On("SaveDevice", mock.Anything).Return(nil).Once()
		createdDevice, err := sm.CreateDevice(ctx, id, "Test Device", "RSA")
		assert.NoError(t, err)
		assert.Equal(t, device.ID, createdDevice.ID)
		mockQuerier.AssertExpectations(t)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		mockQuerier.On("GetDevice", id).Return(&device, nil).Once()
		_, err := sm.CreateDevice(ctx, id, "Test Device", "RSA")
		assert.Equal(t, ErrDeviceExists, err)
		mockQuerier.AssertExpectations(t)
	})

	t.Run("UnsupportedAlgorithm", func(t *testing.T) {
		mockQuerier.On("GetDevice", id).Return(nil, nil).Once()
		_, err := sm.CreateDevice(ctx, id, "Test Device", "Unsupported")
		assert.Equal(t, ErrInvalidAlgorithm, err)
	})
}

func TestGetDevices(t *testing.T) {
	ctx := context.TODO()
	mockQuerier := test_helpers.NewMockQuerier()
	sm := NewDeviceManager(mockQuerier)

	devices := []domain.Device{{ID: uuid.New()}, {ID: uuid.New()}}
	mockQuerier.On("GetDevices").Return(devices, nil).Once()
	retrievedDevices, err := sm.GetDevices(ctx)
	assert.NoError(t, err)
	assert.Equal(t, devices, retrievedDevices)
	mockQuerier.AssertExpectations(t)
}

func TestGetDevice(t *testing.T) {
	ctx := context.TODO()
	mockQuerier := test_helpers.NewMockQuerier()
	sm := NewDeviceManager(mockQuerier)

	id := uuid.New()
	device := domain.Device{ID: id}
	mockQuerier.On("GetDevice", id).Return(&device, nil).Once()
	retrievedDevice, err := sm.GetDevice(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, &device, retrievedDevice)
	mockQuerier.AssertExpectations(t)
}

func TestCreateSignedTransaction(t *testing.T) {
	ctx := context.TODO()
	mockQuerier := test_helpers.NewMockQuerier()
	sm := NewDeviceManager(mockQuerier)

	deviceID := uuid.New()
	// Generate key pair
	privateKey, publicKey, err := sm.keysGenerator.GenerateKeys("RSA")
	assert.NoError(t, err)

	// Build device
	device := domain.Device{
		ID:            deviceID,
		Label:         "Test Device",
		SignAlgorithm: "RSA",
		PrivateKey:    string(privateKey),
		PublicKey:     string(publicKey),
		SignCounter:   0,
	}

	data := []byte("test data")

	t.Run("SuccessfulTransactionCreation", func(t *testing.T) {
		mockQuerier.On("GetDevice", deviceID).Return(&device, nil).Once()
		mockQuerier.On("UpdateDevice", mock.Anything).Return(nil).Once()
		mockQuerier.On("GetSignedTransaction", deviceID, mock.Anything).Return(nil, nil).Once()
		mockQuerier.On("SaveSignedTransaction", mock.Anything).Return(uuid.New(), nil).Once()
		transaction, err := sm.CreateSignedTransaction(ctx, deviceID, data)
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		mockQuerier.AssertExpectations(t)

		// TODO: implement full tests over affected resources
	})

	t.Run("DeviceNotFound", func(t *testing.T) {
		mockQuerier.On("GetDevice", deviceID).Return(nil, nil).Once()
		_, err := sm.CreateSignedTransaction(ctx, deviceID, data)
		assert.Equal(t, persistence.ErrDeviceNotFound, err)
		mockQuerier.AssertExpectations(t)
	})
}

func TestGetSignedTransactions(t *testing.T) {
	ctx := context.TODO()
	mockQuerier := test_helpers.NewMockQuerier()
	sm := NewDeviceManager(mockQuerier)

	deviceID := uuid.New()
	transactions := []domain.SignedTransaction{{ID: uuid.New()}, {ID: uuid.New()}}
	mockQuerier.On("GetSignedTransactions", deviceID).Return(transactions, nil).Once()

	retrievedTransactions, err := sm.GetSignedTransactions(ctx, deviceID)
	assert.NoError(t, err)
	assert.Equal(t, transactions, retrievedTransactions)
	mockQuerier.AssertExpectations(t)
}

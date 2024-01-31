package dao

import (
	"encoding/base64"
	"errors"
	"github.com/ildomm/ssccg/domain"
	"github.com/ildomm/ssccg/persistence"
	"github.com/ildomm/ssccg/test_helpers"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateDevice(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	sm := NewDeviceDAO(mockQuerier)

	id := uuid.New()
	device := domain.Device{ID: id, Label: "Test Device", SignAlgorithm: "RSA"}

	t.Run("SuccessfulCreation", func(t *testing.T) {
		mockQuerier.On("GetDevice", id).Return(nil, nil).Once()
		mockQuerier.On("SaveDevice", mock.Anything).Return(nil).Once()
		createdDevice, err := sm.CreateDevice(id, "Test Device", "RSA")
		assert.NoError(t, err)
		assert.Equal(t, device.ID, createdDevice.ID)
		mockQuerier.AssertExpectations(t)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		mockQuerier.On("GetDevice", id).Return(&device, nil).Once()
		_, err := sm.CreateDevice(id, "Test Device", "RSA")
		assert.Equal(t, ErrDeviceExists, err)
		mockQuerier.AssertExpectations(t)
	})

	t.Run("UnsupportedAlgorithm", func(t *testing.T) {
		mockQuerier.On("GetDevice", id).Return(nil, nil).Once()
		_, err := sm.CreateDevice(id, "Test Device", "Unsupported")
		assert.Equal(t, ErrInvalidAlgorithm, err)
	})
}

func TestGetDevices(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	sm := NewDeviceDAO(mockQuerier)

	devices := []domain.Device{{ID: uuid.New()}, {ID: uuid.New()}}
	mockQuerier.On("GetDevices").Return(devices, nil).Once()
	retrievedDevices, err := sm.GetDevices()
	assert.NoError(t, err)
	assert.Equal(t, devices, retrievedDevices)
	mockQuerier.AssertExpectations(t)
}

func TestGetDevice(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	sm := NewDeviceDAO(mockQuerier)

	id := uuid.New()
	device := domain.Device{ID: id}
	mockQuerier.On("GetDevice", id).Return(&device, nil).Once()
	retrievedDevice, err := sm.GetDevice(id)
	assert.NoError(t, err)
	assert.Equal(t, &device, retrievedDevice)
	mockQuerier.AssertExpectations(t)
}

func TestCreateSignedTransaction(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	sm := NewDeviceDAO(mockQuerier)

	deviceID := uuid.New()
	// GeneratePairs key pair
	privateKey, publicKey, err := sm.keysBuilder.Build("RSA")
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
		mockQuerier = test_helpers.NewMockQuerier()
		sm = NewDeviceDAO(mockQuerier)

		mockQuerier.On("GetDevice", deviceID).Return(&device, nil).Once()
		mockQuerier.On("UpdateDevice", mock.Anything).Return(nil).Once()
		mockQuerier.On("GetSignedTransaction", deviceID, mock.Anything).Return(nil, nil).Once()
		mockQuerier.On("SaveSignedTransaction", mock.Anything).Return(uuid.New(), nil).Once()
		transaction, err := sm.CreateSignedTransaction(deviceID, data)
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		mockQuerier.AssertExpectations(t)
	})

	t.Run("DeviceNotFound", func(t *testing.T) {
		mockQuerier = test_helpers.NewMockQuerier()
		sm = NewDeviceDAO(mockQuerier)

		mockQuerier.On("GetDevice", deviceID).Return(nil, nil).Once()
		_, err := sm.CreateSignedTransaction(deviceID, data)
		assert.Equal(t, persistence.ErrDeviceNotFound, err)
		mockQuerier.AssertExpectations(t)
	})

	t.Run("ErrorSavingTransaction", func(t *testing.T) {
		mockQuerier = test_helpers.NewMockQuerier()
		sm = NewDeviceDAO(mockQuerier)

		mockQuerier.On("GetDevice", deviceID).Return(&device, nil).Once()
		mockQuerier.On("GetSignedTransaction", deviceID, mock.Anything).Return(nil, nil).Once()
		mockQuerier.On("SaveSignedTransaction", mock.Anything).Return(uuid.Nil, errors.New("database error")).Once()

		_, err := sm.CreateSignedTransaction(deviceID, data)
		assert.Error(t, err)
		mockQuerier.AssertNotCalled(t, "UpdateDevice", mock.Anything)
	})

	t.Run("ErrorUpdatingDevice", func(t *testing.T) {
		mockQuerier = test_helpers.NewMockQuerier()
		sm = NewDeviceDAO(mockQuerier)

		mockQuerier.On("GetDevice", deviceID).Return(&device, nil).Once()
		mockQuerier.On("GetSignedTransaction", deviceID, mock.Anything).Return(nil, nil).Once()
		mockQuerier.On("SaveSignedTransaction", mock.Anything).Return(uuid.New(), nil).Once()

		// Simulating error in updating the device
		mockQuerier.On("UpdateDevice", mock.Anything).Return(errors.New("database error")).Once()

		_, err := sm.CreateSignedTransaction(deviceID, data)
		assert.Error(t, err)
		mockQuerier.AssertExpectations(t)
	})
}

func TestPreviousDeviceSignature(t *testing.T) {
	deviceId := uuid.New()
	prevSignature := "prevSignature"

	mockQuerier := test_helpers.NewMockQuerier()
	dao := NewDeviceDAO(mockQuerier)

	t.Run("NoPreviousSignature", func(t *testing.T) {
		mockQuerier.On("GetSignedTransaction", deviceId, mock.AnythingOfType("int")).Return(nil, nil).Once()

		signature, err := dao.previousDeviceSignature(deviceId, 1)
		assert.NoError(t, err)
		expected := base64.StdEncoding.EncodeToString([]byte(deviceId.String()))
		assert.Equal(t, expected, signature)
	})

	t.Run("PreviousSignatureExists", func(t *testing.T) {
		mockQuerier.On("GetSignedTransaction", deviceId, mock.AnythingOfType("int")).Return(&domain.SignedTransaction{Sign: prevSignature}, nil).Once()

		signature, err := dao.previousDeviceSignature(deviceId, 1)
		assert.NoError(t, err)
		assert.Equal(t, prevSignature, signature)
	})

	t.Run("ErrorFetchingTransaction", func(t *testing.T) {
		mockQuerier.On("GetSignedTransaction", deviceId, mock.AnythingOfType("int")).Return(nil, errors.New("database error")).Once()

		_, err := dao.previousDeviceSignature(deviceId, 1)
		assert.Error(t, err)
	})
}

func TestGetSignedTransactions(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	sm := NewDeviceDAO(mockQuerier)

	deviceID := uuid.New()
	transactions := []domain.SignedTransaction{{ID: uuid.New()}, {ID: uuid.New()}}
	mockQuerier.On("GetSignedTransactions", deviceID).Return(transactions, nil).Once()

	retrievedTransactions, err := sm.GetSignedTransactions(deviceID)
	assert.NoError(t, err)
	assert.Equal(t, transactions, retrievedTransactions)
	mockQuerier.AssertExpectations(t)
}

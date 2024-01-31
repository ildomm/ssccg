package persistence

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryQuerier(t *testing.T) {
	ctx := context.TODO()
	querier, err := NewInMemoryQuerier(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, querier)
}

func TestInMemorySaveAndGetDevice(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)
	device := domain.Device{
		ID:            uuid.New(),
		Label:         "Test Device",
		SignCounter:   0,
		SignAlgorithm: "RSA",
		PublicKey:     "public key",
		PrivateKey:    "private key",
	}

	err := querier.SaveDevice(device)
	assert.NoError(t, err)

	retrievedDevice, err := querier.GetDevice(device.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedDevice)
	assert.Equal(t, device.ID, retrievedDevice.ID)
}

func TestInMemoryGetDeviceNotFound(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)
	id := uuid.New()

	device, err := querier.GetDevice(id)
	assert.Error(t, err)
	assert.Nil(t, device)
	assert.Equal(t, ErrDeviceNotFound, err)
}

func TestInMemoryUpdateDevice(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)
	device := domain.Device{
		ID:            uuid.New(),
		Label:         "Test Device",
		SignCounter:   0,
		SignAlgorithm: "RSA",
		PublicKey:     "public key",
		PrivateKey:    "private key",
	}

	// Save first, then update
	_ = querier.SaveDevice(device)
	device.SignCounter++
	err := querier.UpdateDevice(device)
	assert.NoError(t, err)

	// Verify the update
	updatedDevice, _ := querier.GetDevice(device.ID)
	assert.Equal(t, 1, updatedDevice.SignCounter)
}

func TestInMemoryUpdateDeviceNotFound(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)
	device := domain.Device{
		ID:            uuid.New(),
		Label:         "Test Device",
		SignCounter:   0,
		SignAlgorithm: "RSA",
		PublicKey:     "public key",
		PrivateKey:    "private key",
	}

	err := querier.UpdateDevice(device)
	assert.Equal(t, ErrDeviceNotFound, err)
}

func TestInMemorySaveSignedTransactionWithDevice(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)

	// First, create and save a device
	deviceID := uuid.New()
	device := domain.Device{
		ID:            deviceID,
		Label:         "Test Device",
		SignCounter:   0,
		SignAlgorithm: "RSA",
		PublicKey:     "public key",
		PrivateKey:    "private key",
	}
	err := querier.SaveDevice(device)
	assert.NoError(t, err)

	// Then, create and save a signed transaction
	transaction := domain.SignedTransaction{
		ID:          uuid.New(),
		DeviceID:    deviceID, // Use the same device ID
		RawData:     []byte{0x01, 0x02, 0x03},
		Sign:        "sign",
		SignCounter: 0,
	}

	id, err := querier.SaveSignedTransaction(transaction)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, id)

	// Retrieving and checking the saved transaction
	signatures, err := querier.GetSignedTransactions(deviceID)
	assert.NoError(t, err)
	assert.Len(t, signatures, 1)
	assert.Equal(t, id, signatures[0].ID)
}

func TestInMemorySaveSignedTransactionWithoutDevice(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)

	// Create a transaction without saving the device first
	transaction := domain.SignedTransaction{
		ID:          uuid.New(),
		DeviceID:    uuid.New(), // Random device ID
		RawData:     []byte{0x01, 0x02, 0x03},
		Sign:        "sign",
		SignCounter: 0,
	}

	_, err := querier.SaveSignedTransaction(transaction)
	assert.Equal(t, ErrDeviceNotFound, err)
}

func TestInMemoryGetSignedTransactionsEmpty(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)
	id := uuid.New()

	signatures, err := querier.GetSignedTransactions(id)
	assert.NoError(t, err)
	assert.Empty(t, signatures)
}

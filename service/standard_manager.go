package service

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/crypto"
	"github.com/ildomm/ssccg/domain"
	"github.com/ildomm/ssccg/persistence"
)

var ErrDeviceExists = errors.New("device already exists")
var ErrInvalidAlgorithm = errors.New("invalid algorithm")

type StandardManager struct {
	querier       persistence.Querier
	keysGenerator *crypto.KeysGenerator
	SignGenerator *crypto.SignGenerator
}

func NewDeviceManager(querier persistence.Querier) *StandardManager {
	dm := StandardManager{
		querier:       querier,
		keysGenerator: crypto.NewKeysGenerator(),
		SignGenerator: crypto.NewSignGenerator(),
	}
	return &dm
}

// CreateDevice creates a new device with a new key pair
// It does check if the device already exists, return error if it does exist
// It does check if the algorithm is supported, return error if it does not
// It does generate a new key pair based on algorithm
// It does start the sign counter at 0
// It does store the device in the database
// It returns the newly created device
func (dm *StandardManager) CreateDevice(ctx context.Context, id uuid.UUID, label, algorithm string) (*domain.Device, error) {
	// Check if device exists
	existingDevice, err := dm.querier.GetDevice(id)
	if err != nil {
		return nil, err
	}
	if existingDevice != nil {
		return nil, ErrDeviceExists
	}

	// Validate algorithm
	if !dm.keysGenerator.IsValidAlgorithm(algorithm) {
		return nil, ErrInvalidAlgorithm
	}

	// Generate key pair
	privateKey, publicKey, err := dm.keysGenerator.GenerateKeys(algorithm)
	if err != nil {
		return nil, err
	}

	// Create device
	device := domain.Device{
		ID:            id,
		Label:         label,
		SignAlgorithm: algorithm,
		PrivateKey:    string(privateKey),
		PublicKey:     string(publicKey),
		SignCounter:   0,
	}

	// Store device in database
	err = dm.querier.SaveDevice(device)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// GetDevices returns all devices from the database
func (dm *StandardManager) GetDevices(ctx context.Context) ([]domain.Device, error) {
	return dm.querier.GetDevices()
}

// GetDevice returns a device from the database
func (dm *StandardManager) GetDevice(ctx context.Context, id uuid.UUID) (*domain.Device, error) {
	return dm.querier.GetDevice(id)
}

// previousDeviceSignature returns the previous device signature
// It does return the device id if no previous signature exists
// It does return the previous signature if it exists
func (dm *StandardManager) previousDeviceSignature(deviceId uuid.UUID, signCounter int) (string, error) {

	previousSignedTransaction, err := dm.querier.GetSignedTransaction(deviceId, signCounter)
	if err != nil {
		return "", err
	}

	// If no previous signed transaction exists, return the device id
	if previousSignedTransaction == nil {
		return base64.StdEncoding.EncodeToString([]byte(deviceId.String())), nil
	}

	return previousSignedTransaction.Sign, nil
}

// CreateSignedTransaction creates a new signed transaction
// It does check if the device exists, return error if it does not exist
// It does generate a new signature based on the device's algorithm
// It does increment the device's sign counter and update the device in the database
// It does persist the device sign counter with the transaction
// It does store the signed transaction in the database
// It returns the newly created signed transaction
func (dm *StandardManager) CreateSignedTransaction(ctx context.Context, deviceId uuid.UUID, data []byte) (*domain.SignedTransaction, error) {

	// Check if device exists
	device, err := dm.querier.GetDevice(deviceId)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, persistence.ErrDeviceNotFound
	}

	// Get previous signed transaction
	previousSignature, err := dm.previousDeviceSignature(deviceId, device.SignCounter)
	if err != nil {
		return nil, err
	}

	// Generate signature
	signature, err := dm.SignGenerator.Sign(device.SignAlgorithm, []byte(device.PrivateKey), data)
	if err != nil {
		return nil, err
	}

	// Increment sign counter
	device.SignCounter++
	err = dm.querier.UpdateDevice(*device)
	if err != nil {
		return nil, err
	}

	// Build signed transaction
	transaction := domain.SignedTransaction{
		ID:                 uuid.New(),
		DeviceID:           deviceId,
		RawData:            data,
		Sign:               base64.StdEncoding.EncodeToString(signature),
		SignCounter:        device.SignCounter,
		PreviousDeviceSign: previousSignature,
	}

	_, err = dm.querier.SaveSignedTransaction(transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetSignedTransactions returns all signed transactions from the database
func (dm *StandardManager) GetSignedTransactions(ctx context.Context, deviceId uuid.UUID) ([]domain.SignedTransaction, error) {
	return dm.querier.GetSignedTransactions(deviceId)
}

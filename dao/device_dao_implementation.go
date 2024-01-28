package dao

import (
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/crypto"
	"github.com/ildomm/ssccg/domain"
	"github.com/ildomm/ssccg/persistence"
	"sync"
)

var ErrDeviceExists = errors.New("device already exists")
var ErrInvalidAlgorithm = errors.New("invalid algorithm")

type deviceDao struct {
	querier     persistence.Querier
	keysBuilder *crypto.KeysBuilder
	Signer      *crypto.Signer
	lock        sync.Mutex
}

func NewDeviceDAO(querier persistence.Querier) *deviceDao {
	dm := deviceDao{
		querier:     querier,
		keysBuilder: crypto.NewKeysBuilder(),
		Signer:      crypto.NewSigner(),
	}
	return &dm
}

// CreateDevice creates a new device with a new key pair
// It does check if the device already exists, return error if it does exist
// It does check if the algorithm is supported, return error if it does not
// It does build a new key pair based on algorithm
// It does start the sign counter at 0
// It does store the device in the database
// It returns the newly created device
func (dm *deviceDao) CreateDevice(id uuid.UUID, label, algorithm string) (*domain.Device, error) {
	// Check if device exists
	existingDevice, err := dm.querier.GetDevice(id)
	if err != nil && !errors.Is(err, persistence.ErrDeviceNotFound) {
		return nil, err
	}
	if existingDevice != nil {
		return nil, ErrDeviceExists
	}

	// Validate algorithm
	if !dm.keysBuilder.IsValidAlgorithm(algorithm) {
		return nil, ErrInvalidAlgorithm
	}

	// Builds key pair
	privateKey, publicKey, err := dm.keysBuilder.Build(algorithm)
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
func (dm *deviceDao) GetDevices() ([]domain.Device, error) {
	return dm.querier.GetDevices()
}

// GetDevice returns a device from the database
func (dm *deviceDao) GetDevice(id uuid.UUID) (*domain.Device, error) {
	return dm.querier.GetDevice(id)
}

// previousDeviceSignature returns the previous device signature
// It does return the device id if no previous signature exists
// It does return the previous signature if it exists
func (dm *deviceDao) previousDeviceSignature(deviceId uuid.UUID, signCounter int) (string, error) {

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
func (dm *deviceDao) CreateSignedTransaction(deviceId uuid.UUID, data []byte) (*domain.SignedTransaction, error) {

	// Lock to prevent concurrent access
	// Doing so, we prevent the sign counter to be incremented twice wrongly
	dm.lock.Lock()
	defer dm.lock.Unlock()

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

	// Build signed transaction
	transaction := domain.SignedTransaction{
		ID:                 uuid.New(),
		DeviceID:           deviceId,
		RawData:            data,
		SignCounter:        device.SignCounter + 1,
		PreviousDeviceSign: previousSignature,
	}

	// Sign data
	signature, err := dm.Signer.Sign(device.SignAlgorithm,
		[]byte(device.PrivateKey),
		[]byte(transaction.SignedData()))
	if err != nil {
		return nil, err
	}
	transaction.Sign = base64.StdEncoding.EncodeToString(signature)

	// Store signed transaction in database
	_, err = dm.querier.SaveSignedTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// Increment sign counter
	// This is done after the transaction is stored in the database
	// to prevent the sign counter to be incremented wrongly
	// Fail safe measure
	// When using a real database, this should be done in a transaction
	device.SignCounter++
	err = dm.querier.UpdateDevice(*device)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetSignedTransactions returns all signed transactions from the database
func (dm *deviceDao) GetSignedTransactions(deviceId uuid.UUID) ([]domain.SignedTransaction, error) {
	return dm.querier.GetSignedTransactions(deviceId)
}

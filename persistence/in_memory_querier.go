package persistence

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
	"sync"
)

type InMemoryQuerier struct {
	lock            sync.Mutex
	ctx             context.Context
	devices         map[uuid.UUID]domain.Device
	signedTransacts map[uuid.UUID][]domain.SignedTransaction
}

var ErrDeviceNotFound = errors.New("device not found")

func NewInMemoryQuerier(ctx context.Context) (*InMemoryQuerier, error) {
	return &InMemoryQuerier{
		ctx:             ctx,
		devices:         make(map[uuid.UUID]domain.Device),
		signedTransacts: make(map[uuid.UUID][]domain.SignedTransaction),
	}, nil
}

////////////////////////////////// Database Querier operations /////////////////////////////////////////////////////////

func (q *InMemoryQuerier) Close() {
	// Nothing to do here for in-memory storage
}

func (q *InMemoryQuerier) SaveDevice(device domain.Device) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.devices[device.ID] = device
	return nil
}

func (q *InMemoryQuerier) GetDevices() ([]domain.Device, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	var devices []domain.Device
	for _, device := range q.devices {
		devices = append(devices, device)
	}
	return devices, nil
}

func (q *InMemoryQuerier) GetDevice(id uuid.UUID) (*domain.Device, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	device, exists := q.devices[id]
	if !exists {
		return nil, ErrDeviceNotFound
	}
	return &device, nil
}

func (q *InMemoryQuerier) UpdateDevice(device domain.Device) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	_, exists := q.devices[device.ID]
	if !exists {
		return ErrDeviceNotFound
	}
	q.devices[device.ID] = device
	return nil
}

func (q *InMemoryQuerier) SaveSignedTransaction(transaction domain.SignedTransaction) (uuid.UUID, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Check if the device exists
	_, deviceExists := q.devices[transaction.DeviceID]
	if !deviceExists {
		return uuid.Nil, ErrDeviceNotFound
	}

	transaction.ID = uuid.New() // Generate a new UUID for the transaction
	q.signedTransacts[transaction.DeviceID] = append(q.signedTransacts[transaction.DeviceID], transaction)
	return transaction.ID, nil
}

func (q *InMemoryQuerier) GetSignedTransactions(deviceId uuid.UUID) ([]domain.SignedTransaction, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.signedTransacts[deviceId], nil
}

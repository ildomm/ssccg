package persistence

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
	"sync"
)

// TODO: in-memory persistence ...

type InMemoryQuerier struct {
	lock sync.Mutex //nolint:all
	ctx  context.Context
}

func NewInMemoryQuerier(ctx context.Context) (*InMemoryQuerier, error) {
	return &InMemoryQuerier{
		ctx: ctx,
	}, nil
}

////////////////////////////////// Database Querier operations /////////////////////////////////////////////////////////

func (q *InMemoryQuerier) Close() {
	// Nothing to do here
}

func (q *InMemoryQuerier) SaveDevice(device domain.Device) error {
	//q.lock.Lock()
	//defer q.lock.Unlock()

	panic("implement me")
}

func (q *InMemoryQuerier) GetDevices() ([]domain.Device, error) {
	panic("implement me")
}

func (q *InMemoryQuerier) GetDevice(id uuid.UUID) (*domain.Device, error) {
	panic("implement me")
}

func (q *InMemoryQuerier) UpdateDevice(device domain.Device) error {
	panic("implement me")
}

func (q *InMemoryQuerier) SaveSignature(device domain.Signature) (uuid.UUID, error) {
	panic("implement me")
}

func (q *InMemoryQuerier) GetSignaturesByDevice(id uuid.UUID) ([]domain.Signature, error) {
	panic("implement me")
}

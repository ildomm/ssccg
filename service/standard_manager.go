package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
	"github.com/ildomm/ssccg/persistence"
)

type StandardManager struct {
	querier *persistence.Querier
}

func NewDeviceManager(querier *persistence.Querier) *StandardManager {
	dm := StandardManager{
		querier: querier,
	}
	return &dm
}

func (dm *StandardManager) CreateDevice(ctx context.Context, id uuid.UUID, algorithm string) (*domain.Device, error) {
	return nil, nil
}

func (dm *StandardManager) GetDevices(ctx context.Context) ([]domain.Device, error) {
	return nil, nil
}

func (dm *StandardManager) GetDevice(ctx context.Context, id uuid.UUID) (*domain.Device, error) {
	return nil, nil
}

func (dm *StandardManager) CreateSignedTransaction(ctx context.Context, deviceId uuid.UUID, data []byte) (*domain.SignedTransaction, error) {
	return nil, nil
}

func (dm *StandardManager) GetSignedTransactions(ctx context.Context, deviceId uuid.UUID) ([]domain.SignedTransaction, error) {
	return nil, nil
}

package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
)

type Manager interface {
	CreateDevice(ctx context.Context, id uuid.UUID, label, algorithm string) (*domain.Device, error)
	GetDevices(ctx context.Context) ([]domain.Device, error)
	GetDevice(ctx context.Context, id uuid.UUID) (*domain.Device, error)
	CreateSignedTransaction(ctx context.Context, deviceId uuid.UUID, data []byte) (*domain.SignedTransaction, error)
	GetSignedTransactions(ctx context.Context, deviceId uuid.UUID) ([]domain.SignedTransaction, error)
}

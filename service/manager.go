package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
)

type Manager interface {
	CreateDevice(ctx context.Context, id uuid.UUID, algorithm string) (*domain.Device, error)
	GetDevices(ctx context.Context) ([]domain.Device, error)
	GetDevice(ctx context.Context, id uuid.UUID) (*domain.Device, error)
	CreateDeviceSignature(ctx context.Context, id uuid.UUID, data []byte) (*domain.Signature, error)
	GetDeviceSignatures(ctx context.Context, id uuid.UUID) ([]domain.Signature, error)
}

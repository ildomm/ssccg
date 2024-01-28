package dao

import (
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
)

type DeviceDAO interface {
	CreateDevice(id uuid.UUID, label, algorithm string) (*domain.Device, error)
	GetDevices() ([]domain.Device, error)
	GetDevice(id uuid.UUID) (*domain.Device, error)
	CreateSignedTransaction(deviceId uuid.UUID, data []byte) (*domain.SignedTransaction, error)
	GetSignedTransactions(deviceId uuid.UUID) ([]domain.SignedTransaction, error)
}

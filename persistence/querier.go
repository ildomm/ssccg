package persistence

import (
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
)

type Querier interface {
	Close()

	SaveDevice(device domain.Device) error
	GetDevices() ([]domain.Device, error)
	GetDevice(id uuid.UUID) (*domain.Device, error)
	UpdateDevice(device domain.Device) error
	SaveSignedTransaction(transaction domain.SignedTransaction) (uuid.UUID, error)
	GetSignedTransactions(deviceId uuid.UUID) ([]domain.SignedTransaction, error)
}

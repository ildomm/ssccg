package persistence

import (
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
)

type Querier interface {
	Close()

	SaveDevice(device domain.Device) error
	GetDeviceById(id uuid.UUID) (*domain.Device, error)
	UpdateDevice(device domain.Device) error

	SaveSignature(device domain.Signature) (uuid.UUID, error)
	GetSignaturesByDeviceId(id uuid.UUID) ([]domain.Signature, error)
}

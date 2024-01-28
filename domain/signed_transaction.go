package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type SignedTransaction struct {
	ID                 uuid.UUID `db:"id"`
	DeviceID           uuid.UUID `db:"device_id"`
	RawData            []byte    `db:"raw_data"`
	Sign               string    `db:"sign"`
	PreviousDeviceSign string    `db:"previous_device_sign"`
	SignCounter        int       `db:"sign_counter"`
}

func (s *SignedTransaction) SignedData() string {
	return fmt.Sprintf("%d_%s_%s", s.SignCounter, s.RawData, s.PreviousDeviceSign)
}

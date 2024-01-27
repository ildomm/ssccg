package domain

import "github.com/google/uuid"

type SignedTransaction struct {
	ID          uuid.UUID `db:"id"`
	DeviceID    uuid.UUID `db:"device_id"`
	RawData     []byte    `db:"raw_data"`
	Sign        string    `db:"sign"`
	SignCounter int       `db:"sign_counter"`
}

package domain

import "github.com/google/uuid"

type Device struct {
	ID            uuid.UUID `db:"id"`
	Label         string    `db:"label"`
	SignCounter   int       `db:"sign_counter"`
	SignAlgorithm string    `db:"sign_algorithm"`
	PublicKey     string    `db:"public_key"`
	PrivateKey    string    `db:"private_Key"`
}

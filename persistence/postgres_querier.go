package persistence

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"

	"github.com/allisson/go-pglock/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresQuerier struct {
	ctx    context.Context
	dbURL  string
	dbConn *sqlx.DB
	lock   pglock.Lock //nolint:all
}

func NewPostgresQuerier(ctx context.Context, url string) (*PostgresQuerier, error) {

	// TODO: implement postgres initialization ...
	// TODO: implement postgres connection ...
	// TODO: implement postgres migration ...

	return &PostgresQuerier{
		ctx:    ctx,
		dbURL:  url,
		dbConn: nil,
	}, nil
}

////////////////////////////////// Database Querier operations /////////////////////////////////////////////////////////

func (q *PostgresQuerier) Close() {
	panic("implement me")
}

func (q *PostgresQuerier) SaveDevice(device domain.Device) error {
	panic("implement me")
}

func (q *PostgresQuerier) GetDevices() ([]domain.Device, error) {
	panic("implement me")
}

func (q *PostgresQuerier) GetDevice(id uuid.UUID) (*domain.Device, error) {
	panic("implement me")
}

func (q *PostgresQuerier) UpdateDevice(device domain.Device) error {
	panic("implement me")
}

func (q *PostgresQuerier) SaveSignature(device domain.Signature) (uuid.UUID, error) {
	panic("implement me")
}

func (q *PostgresQuerier) GetSignaturesByDevice(id uuid.UUID) ([]domain.Signature, error) {
	panic("implement me")
}

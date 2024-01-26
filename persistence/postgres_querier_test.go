package persistence

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewPostgresQuerier(t *testing.T) {
	ctx := context.TODO()
	querier, err := NewPostgresQuerier(ctx, "test")
	assert.NoError(t, err)
	assert.NotNil(t, querier)
}

func TestPostgresSaveDevice(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewPostgresQuerier(ctx, "test")
	device := domain.Device{ /* Initialize fields */ }

	assert.Panics(t, func() { querier.SaveDevice(device) }) //nolint:all
}

func TestIPostgresGetDeviceById(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewPostgresQuerier(ctx, "test")
	id := uuid.New()

	assert.Panics(t, func() { querier.GetDeviceById(id) }) //nolint:all
}

func TestPostgresUpdateDevice(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewPostgresQuerier(ctx, "test")
	device := domain.Device{ /* Initialize fields */ }

	assert.Panics(t, func() { querier.UpdateDevice(device) }) //nolint:all
}

func TestPostgresSaveSignature(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewPostgresQuerier(ctx, "test")
	signature := domain.Signature{ /* Initialize fields */ }

	assert.Panics(t, func() { querier.SaveSignature(signature) }) //nolint:all
}

func TestPostgresGetSignaturesByDeviceId(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewPostgresQuerier(ctx, "test")
	id := uuid.New()

	assert.Panics(t, func() { querier.GetSignaturesByDeviceId(id) }) //nolint:all
}
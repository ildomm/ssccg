package persistence

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryQuerier(t *testing.T) {
	ctx := context.TODO()
	querier, err := NewInMemoryQuerier(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, querier)
}

func TestInMemorySaveDevice(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)
	device := domain.Device{ /* Initialize fields */ }

	assert.Panics(t, func() { querier.SaveDevice(device) }) //nolint:all
}

func TestInMemoryGetDeviceById(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)
	id := uuid.New()

	assert.Panics(t, func() { querier.GetDeviceById(id) }) //nolint:all
}

func TestInMemoryUpdateDevice(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)
	device := domain.Device{ /* Initialize fields */ }

	assert.Panics(t, func() { querier.UpdateDevice(device) }) //nolint:all
}

func TestInMemorySaveSignature(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)
	signature := domain.Signature{ /* Initialize fields */ }

	assert.Panics(t, func() { querier.SaveSignature(signature) }) //nolint:all
}

func TestInMemoryGetSignaturesByDeviceId(t *testing.T) {
	ctx := context.TODO()
	querier, _ := NewInMemoryQuerier(ctx)
	id := uuid.New()

	assert.Panics(t, func() { querier.GetSignaturesByDeviceId(id) }) //nolint:all
}

package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignWithInvalidAlgorithm(t *testing.T) {
	sg := NewSigner()
	dataToBeSigned := []byte("test data")
	_, err := sg.Sign("Invalid", nil, dataToBeSigned)
	assert.Equal(t, ErrCryptoEngineNotFound, err)
}

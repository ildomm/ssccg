package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildKeysRSA(t *testing.T) {
	kg := NewKeysBuilder()
	privateKeyBytes, publicKeyBytes, err := kg.Build("RSA")
	assert.NoError(t, err)
	assert.NotNil(t, privateKeyBytes)
	assert.NotNil(t, publicKeyBytes)
}

func TestBuildKeysECC(t *testing.T) {
	kg := NewKeysBuilder()
	privateKeyBytes, publicKeyBytes, err := kg.Build("ECDSA")
	assert.NoError(t, err)
	assert.NotNil(t, privateKeyBytes)
	assert.NotNil(t, publicKeyBytes)
}

func TestBuildKeysInvalid(t *testing.T) {
	kg := NewKeysBuilder()
	_, _, err := kg.Build("Invalid")
	assert.Error(t, err)
	assert.Equal(t, ErrCryptoEngineNotFound, err)
}

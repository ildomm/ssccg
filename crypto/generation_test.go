package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestRSAGenerator tests the RSAGenerator's Generate function.
func TestRSAGenerator(t *testing.T) {
	generator := RSAGenerator{}

	keyPair, err := generator.Generate()
	assert.NoError(t, err, "RSA generator should not produce an error")
	assert.NotNil(t, keyPair, "RSA key pair should not be nil")

	// Check if the keys are non-nil
	assert.NotNil(t, keyPair.Private, "RSA private key should not be nil")
	assert.NotNil(t, keyPair.Public, "RSA public key should not be nil")
}

// TestECCGenerator tests the ECCGenerator's Generate function.
func TestECCGenerator(t *testing.T) {
	generator := ECCGenerator{}

	keyPair, err := generator.Generate()
	assert.NoError(t, err, "ECC generator should not produce an error")
	assert.NotNil(t, keyPair, "ECC key pair should not be nil")

	// Check if the keys are non-nil
	assert.NotNil(t, keyPair.Private, "ECC private key should not be nil")
	assert.NotNil(t, keyPair.Public, "ECC public key should not be nil")
}

func TestGenerateKeysRSA(t *testing.T) {
	kg := NewKeysGenerator()
	privateKeyBytes, publicKeyBytes, err := kg.GenerateKeys("RSA")
	assert.NoError(t, err)
	assert.NotNil(t, privateKeyBytes)
	assert.NotNil(t, publicKeyBytes)
}

func TestGenerateKeysECC(t *testing.T) {
	kg := NewKeysGenerator()
	privateKeyBytes, publicKeyBytes, err := kg.GenerateKeys("ECDSA")
	assert.NoError(t, err)
	assert.NotNil(t, privateKeyBytes)
	assert.NotNil(t, publicKeyBytes)
}

func TestGenerateKeysInvalid(t *testing.T) {
	kg := NewKeysGenerator()
	_, _, err := kg.GenerateKeys("Invalid")
	assert.Error(t, err)
	assert.Equal(t, ErrCryptoEngineNotFound, err)
}

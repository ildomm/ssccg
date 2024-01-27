package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestECCMarshaler tests both Encode and Decode functions of ECCMarshaler.
func TestECCMarshaler(t *testing.T) {
	// Generate a key pair for testing
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	publicKey := &privateKey.PublicKey
	keyPair := ECCKeyPair{Public: publicKey, Private: privateKey}

	marshaler := NewECCMarshaler()

	// Test Encode
	encodedPublic, encodedPrivate, err := marshaler.Marshal(keyPair)
	assert.NoError(t, err)
	assert.NotEmpty(t, encodedPublic, "Encoded public key should not be empty")
	assert.NotEmpty(t, encodedPrivate, "Encoded private key should not be empty")

	// Test Decode
	decodedKeyPair, err := marshaler.Unmarshal(encodedPrivate)
	assert.NoError(t, err)
	assert.NotNil(t, decodedKeyPair)

	// Compare the original and decoded key pairs
	assert.Equal(t, keyPair.Private.D, decodedKeyPair.Private.D)
	assert.Equal(t, keyPair.Public.X, decodedKeyPair.Public.X)
	assert.Equal(t, keyPair.Public.Y, decodedKeyPair.Public.Y)
}

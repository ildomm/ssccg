package algorithms

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestECCMarshaler tests both Encode and Decode functions of ECCMarshaler.
func TestECCMarshaler(t *testing.T) {
	// GeneratePairs a key pair for testing
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

// TestECCKeysBuilder tests the ECCKeysBuilder's GeneratePairs function.
func TestECCKeysBuilder(t *testing.T) {
	generator := ECCKeysBuilder{}

	keyPair, err := generator.Pairs()
	assert.NoError(t, err, "ECC generator should not produce an error")
	assert.NotNil(t, keyPair, "ECC key pair should not be nil")

	// Check if the keys are non-nil
	assert.NotNil(t, keyPair.Private, "ECC private key should not be nil")
	assert.NotNil(t, keyPair.Public, "ECC public key should not be nil")
}

func TestECCSignatureGeneration(t *testing.T) {
	signer := NewECCSigner()
	generator := NewECCKeysBuilder()
	keyPair, err := generator.Pairs()
	assert.NoError(t, err)

	marshaler := NewECCMarshaler()
	_, privateKeyBytes, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err)

	dataToBeSigned := []byte("test data")
	signature, err := signer.Sign(privateKeyBytes, dataToBeSigned)
	assert.NoError(t, err)
	assert.NotNil(t, signature)
}

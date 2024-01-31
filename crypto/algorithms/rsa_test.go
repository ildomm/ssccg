package algorithms

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

// generateTestRSAKeyPair generates an RSA key pair for testing.
func generateTestRSAKeyPair() (*RSAKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048) // Using a 2048-bit key for testing
	if err != nil {
		return nil, err
	}
	return &RSAKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}

// TestRSAMarshaler tests both Marshal and Unmarshal functions of RSAMarshaler.
func TestRSAMarshaler(t *testing.T) {
	keyPair, err := generateTestRSAKeyPair()
	assert.NoError(t, err)

	marshaler := NewRSAMarshaler()

	// Test Marshal
	encodedPublic, encodedPrivate, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err)
	assert.NotEmpty(t, encodedPublic, "Encoded public key should not be empty")
	assert.NotEmpty(t, encodedPrivate, "Encoded private key should not be empty")

	// Test Unmarshal
	decodedKeyPair, err := marshaler.Unmarshal(encodedPrivate)
	assert.NoError(t, err)
	assert.NotNil(t, decodedKeyPair)

	// Compare the original and decoded key pairs
	assert.Equal(t, keyPair.Private.D, decodedKeyPair.Private.D)
	assert.Equal(t, keyPair.Public.E, decodedKeyPair.Public.E)
	assert.Equal(t, keyPair.Public.N, decodedKeyPair.Public.N)
}

// TestRSAKeysBuilder tests the RSAKeysBuilder's GeneratePairs function.
func TestRSAKeysBuilder(t *testing.T) {
	generator := RSAKeysBuilder{}

	keyPair, err := generator.Pairs()
	assert.NoError(t, err, "RSA generator should not produce an error")
	assert.NotNil(t, keyPair, "RSA key pair should not be nil")

	// Check if the keys are non-nil
	assert.NotNil(t, keyPair.Private, "RSA private key should not be nil")
	assert.NotNil(t, keyPair.Public, "RSA public key should not be nil")
}

func TestRSASignatureGeneration(t *testing.T) {
	signer := NewRSASigner()
	generator := NewRSAKeysBuilder()
	keyPair, err := generator.Pairs()
	assert.NoError(t, err)

	marshaler := NewRSAMarshaler()
	_, privateKeyBytes, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err)

	dataToBeSigned := []byte("test data")
	signature, err := signer.Sign(privateKeyBytes, dataToBeSigned)
	assert.NoError(t, err)
	assert.NotNil(t, signature)
}

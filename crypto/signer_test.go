package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRSASignatureGeneration(t *testing.T) {
	signer := NewRSASigner()
	generator := NewRSAGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err)

	marshaler := NewRSAMarshaler()
	_, privateKeyBytes, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err)

	dataToBeSigned := []byte("test data")
	signature, err := signer.Sign(privateKeyBytes, dataToBeSigned)
	assert.NoError(t, err)
	assert.NotNil(t, signature)
}

func TestECCSignatureGeneration(t *testing.T) {
	signer := NewECCSigner()
	generator := NewECCGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err)

	marshaler := NewECCMarshaler()
	_, privateKeyBytes, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err)

	dataToBeSigned := []byte("test data")
	signature, err := signer.Sign(privateKeyBytes, dataToBeSigned)
	assert.NoError(t, err)
	assert.NotNil(t, signature)
}

func TestSignGeneratorWithRSA(t *testing.T) {
	sg := NewSignGenerator()
	generator := NewRSAGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err)

	marshaler := NewRSAMarshaler()
	_, privateKeyBytes, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err)

	dataToBeSigned := []byte("test data")
	signature, err := sg.Sign("RSA", privateKeyBytes, dataToBeSigned)
	assert.NoError(t, err)
	assert.NotNil(t, signature)
}

func TestSignGeneratorWithECC(t *testing.T) {
	sg := NewSignGenerator()
	generator := NewECCGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err)

	marshaler := NewECCMarshaler()
	_, privateKeyBytes, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err)

	dataToBeSigned := []byte("test data")
	signature, err := sg.Sign("ECDSA", privateKeyBytes, dataToBeSigned)
	assert.NoError(t, err)
	assert.NotNil(t, signature)
}

func TestSignGeneratorWithInvalidAlgorithm(t *testing.T) {
	sg := NewSignGenerator()
	dataToBeSigned := []byte("test data")
	_, err := sg.Sign("Invalid", nil, dataToBeSigned)
	assert.Equal(t, ErrCryptoEngineNotFound, err)
}

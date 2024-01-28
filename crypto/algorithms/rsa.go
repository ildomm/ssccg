package algorithms

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// RSAKeyPair is a DTO that holds RSA private and public keys.
type RSAKeyPair struct {
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

// RSAMarshaler can encode and decode an RSA key pair.
type RSAMarshaler struct{}

// NewRSAMarshaler creates a new RSAMarshaler.
func NewRSAMarshaler() RSAMarshaler {
	return RSAMarshaler{}
}

// Marshal takes an RSAKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func (m *RSAMarshaler) Marshal(keyPair RSAKeyPair) ([]byte, []byte, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(keyPair.Private)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(keyPair.Public)

	encodedPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PRIVATE_KEY",
		Bytes: privateKeyBytes,
	})

	encodePublic := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PUBLIC_KEY",
		Bytes: publicKeyBytes,
	})

	return encodePublic, encodedPrivate, nil
}

// Unmarshal takes an encoded RSA private key and transforms it into a rsa.PrivateKey.
func (m *RSAMarshaler) Unmarshal(privateKeyBytes []byte) (*RSAKeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}

// RSAKeysBuilder generates a RSA key pair.
type RSAKeysBuilder struct{}

func NewRSAKeysBuilder() RSAKeysBuilder {
	return RSAKeysBuilder{}
}

// Pairs generates a new RSAKeyPair.
func (g RSAKeysBuilder) Pairs() (*RSAKeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// Keys generates a new RSAKeyPair and returns the public and private key as a byte slice.
func (g RSAKeysBuilder) Keys() ([]byte, []byte, error) {
	keypair, err := g.Pairs()
	if err != nil {
		return nil, nil, err
	}

	engine := NewRSAMarshaler()
	_, privateKeyBytes, err := engine.Marshal(*keypair)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes := x509.MarshalPKCS1PublicKey(keypair.Public)
	return privateKeyBytes, publicKeyBytes, nil
}

// RSASigner signs a message using an RSA private key.
type RSASigner struct {
	marshaller RSAMarshaler
}

// NewRSASigner creates a new RSASigner.
func NewRSASigner() RSASigner {
	return RSASigner{
		marshaller: NewRSAMarshaler(),
	}
}

// Sign signs data using an RSA private key.
func (sg RSASigner) Sign(privateKeyBytes, dataToBeSigned []byte) ([]byte, error) {
	hash, err := GetHashSum(dataToBeSigned)
	if err != nil {
		return nil, err
	}
	keyPair, err := sg.marshaller.Unmarshal(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, keyPair.Private, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	err = rsa.VerifyPKCS1v15(keyPair.Public, crypto.SHA256, hash, signature)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

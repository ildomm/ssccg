package algorithms

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// ECCKeyPair is a DTO that holds ECC private and public keys.
type ECCKeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

// ECCMarshaler can encode and decode an ECC key pair.
type ECCMarshaler struct{}

// NewECCMarshaler creates a new ECCMarshaler.
func NewECCMarshaler() ECCMarshaler {
	return ECCMarshaler{}
}

// Marshal takes an ECCKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func (m ECCMarshaler) Marshal(keyPair ECCKeyPair) ([]byte, []byte, error) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(keyPair.Private)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(keyPair.Public)
	if err != nil {
		return nil, nil, err
	}

	encodedPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE_KEY",
		Bytes: privateKeyBytes,
	})

	encodedPublic := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC_KEY",
		Bytes: publicKeyBytes,
	})

	return encodedPublic, encodedPrivate, nil
}

// Unmarshal assembles an ECCKeyPair from an encoded private key.
func (m ECCMarshaler) Unmarshal(privateKeyBytes []byte) (*ECCKeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}

// ECCKeysBuilder builds an ECC key pair.
type ECCKeysBuilder struct{}

func NewECCKeysBuilder() ECCKeysBuilder {
	return ECCKeysBuilder{}
}

// Pairs builds a new ECCKeyPair.
func (g ECCKeysBuilder) Pairs() (*ECCKeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// Keys builds a new ECCKeyPair and returns the public and private keys as byte slices.
func (g ECCKeysBuilder) Keys() ([]byte, []byte, error) {
	keypair, err := g.Pairs()
	if err != nil {
		return nil, nil, err
	}
	engine := NewECCMarshaler()
	_, privateKeyBytes, err := engine.Marshal(*keypair)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(keypair.Public)
	return privateKeyBytes, publicKeyBytes, nil
}

// ECCSigner signs data using an ECC private key.
type ECCSigner struct {
	marshaller ECCMarshaler
}

// NewECCSigner creates a new ECCSigner.
func NewECCSigner() ECCSigner {
	return ECCSigner{
		marshaller: NewECCMarshaler(),
	}
}

// Sign signs data using an ECC private key.
func (sg ECCSigner) Sign(privateKeyBytes, dataToBeSigned []byte) ([]byte, error) {
	hash, err := GetHashSum(dataToBeSigned)
	if err != nil {
		return nil, err
	}
	keyPair, err := sg.marshaller.Unmarshal(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	signature, err := ecdsa.SignASN1(rand.Reader, keyPair.Private, hash[:])
	if err != nil {
		return nil, err
	}
	if !ecdsa.VerifyASN1(keyPair.Public, hash, signature) {
		return nil, errors.New("failed to verify ASN1 signature")
	}
	return signature, nil
}

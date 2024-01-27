package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
)

// RSAGenerator generates a RSA key pair.
type RSAGenerator struct{}

func NewRSAGenerator() RSAGenerator {
	return RSAGenerator{}
}

// Generate generates a new RSAKeyPair.
func (g *RSAGenerator) Generate() (*RSAKeyPair, error) {
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

// ECCGenerator generates an ECC key pair.
type ECCGenerator struct{}

func NewECCGenerator() ECCGenerator {
	return ECCGenerator{}
}

// Generate generates a new ECCKeyPair.
func (g *ECCGenerator) Generate() (*ECCKeyPair, error) {
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

var ErrCryptoEngineNotFound = errors.New("crypto engine not found")

type KeysGenerator struct {
	rsaGenerator   RSAGenerator
	ecdsaGenerator ECCGenerator
}

func NewKeysGenerator() *KeysGenerator {
	kg := KeysGenerator{
		rsaGenerator:   NewRSAGenerator(),
		ecdsaGenerator: NewECCGenerator(),
	}
	return &kg
}

func (kg *KeysGenerator) GenerateKeys(algorithm string) ([]byte, []byte, error) {

	switch algorithm {
	case "RSA":
		return kg.generateRSAKeys()
	case "ECDSA":
		return kg.generateECDSAKeys()
	default:
		return nil, nil, ErrCryptoEngineNotFound
	}
}

func (kg *KeysGenerator) generateRSAKeys() ([]byte, []byte, error) {
	keypair, err := kg.rsaGenerator.Generate()
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

func (kg *KeysGenerator) generateECDSAKeys() ([]byte, []byte, error) {

	keypair, err := kg.ecdsaGenerator.Generate()
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

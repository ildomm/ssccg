package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
)

/* Not needed

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

*/

type SignGenerator struct {
	rsaSigner   RSASigner
	ecdsaSigner ECCSigner
}

func NewSignGenerator() *SignGenerator {
	kg := SignGenerator{
		rsaSigner:   NewRSASigner(),
		ecdsaSigner: NewECCSigner(),
	}
	return &kg
}

func (sg *SignGenerator) Sign(algorithm string, privateKeyBytes, dataToBeSigned []byte) ([]byte, error) {

	switch algorithm {
	case "RSA":
		return sg.rsaSigner.Sign(privateKeyBytes, dataToBeSigned)
	case "ECDSA":
		return sg.ecdsaSigner.Sign(privateKeyBytes, dataToBeSigned)
	default:
		return nil, ErrCryptoEngineNotFound
	}
}

type RSASigner struct {
	marshaller RSAMarshaler
}

func NewRSASigner() RSASigner {
	return RSASigner{
		marshaller: NewRSAMarshaler(),
	}
}

func (sg RSASigner) Sign(privateKeyBytes, dataToBeSigned []byte) ([]byte, error) {
	hash, err := getHashSum(dataToBeSigned)
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

type ECCSigner struct {
	marshaller ECCMarshaler
}

func NewECCSigner() ECCSigner {
	return ECCSigner{
		marshaller: NewECCMarshaler(),
	}
}

func (sg ECCSigner) Sign(privateKeyBytes, dataToBeSigned []byte) ([]byte, error) {
	hash, err := getHashSum(dataToBeSigned)
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

func getHashSum(dataToBeSigned []byte) ([]byte, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write(dataToBeSigned)
	if err != nil {
		return nil, fmt.Errorf("failed to get hash sum: %w", err)
	}
	return msgHash.Sum(nil), nil
}

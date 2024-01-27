package crypto

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// TODO: implement RSA and ECDSA signing ...

// TODO: implement signers builder

// OPTION A

/*
package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type RSASigner struct {
	encodedPrivateKey []byte
	marshaller        *RSAMarshaler
}

func NewRSASigner() (Signer, error) {
	g := RSAGenerator{}
	keyPair, err := g.Generate()
	if err != nil {
		return nil, err
	}

	marshaller := &RSAMarshaler{}
	_, encodedPrivate, err := marshaller.Marshal(*keyPair)
	if err != nil {
		return nil, err
	}

	return &RSASigner{
		encodedPrivateKey: encodedPrivate,
		marshaller:        marshaller,
	}, nil
}

func (signer *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := signer.marshaller.Unmarshal(signer.encodedPrivateKey)
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256(dataToBeSigned)
	sig, err := rsa.SignPKCS1v15(nil, keyPair.Private, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}

	// TODO: Extract into a Verifier interface
	if err = rsa.VerifyPKCS1v15(keyPair.Public, crypto.SHA256, hashed[:], sig); err != nil {
		return nil, err
	}

	return sig, nil
}

type ECCSigner struct {
	encodedPrivateKey []byte
	marshaller        *ECCMarshaler
}

func NewECCSigner() (Signer, error) {
	g := ECCGenerator{}
	keyPair, err := g.Generate()
	if err != nil {
		return nil, err
	}

	marshaller := &ECCMarshaler{}
	_, encodedPrivate, err := marshaller.Encode(*keyPair)
	if err != nil {
		return nil, err
	}

	return &ECCSigner{
		encodedPrivateKey: encodedPrivate,
		marshaller:        marshaller,
	}, nil
}

func (signer *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := signer.marshaller.Decode(signer.encodedPrivateKey)
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256(dataToBeSigned)
	sig, err := ecdsa.SignASN1(rand.Reader, keyPair.Private, hashed[:])
	if err != nil {
		return nil, err
	}

	// TODO: Extract into a Verifier interface
	if valid := ecdsa.VerifyASN1(keyPair.Public, hashed[:], sig); !valid {
		return nil, errors.New("invalid ECC signature")
	}

	return sig, nil
}


*/

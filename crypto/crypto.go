package crypto

import (
	"errors"
	"github.com/ildomm/ssccg/crypto/algorithms"
)

type keysBuilder interface {
	Keys() ([]byte, []byte, error)
}

type signer interface {
	Sign(privateKeyBytes, dataToBeSigned []byte) ([]byte, error)
}

var ErrCryptoEngineNotFound = errors.New("crypto algorithm not found")

var algorithmKeyBuildersRegistry = make(map[string]keysBuilder)
var algorithmSignersRegistry = make(map[string]signer)

// RegisterAlgorithm registers a new algorithm.
func RegisterAlgorithm(name string, builder keysBuilder, signer signer) {
	algorithmKeyBuildersRegistry[name] = builder
	algorithmSignersRegistry[name] = signer
}

// init registers the cryptography algorithms.
func init() {
	RegisterAlgorithm("RSA", algorithms.NewRSAKeysBuilder(), algorithms.NewRSASigner())
	RegisterAlgorithm("ECDSA", algorithms.NewECCKeysBuilder(), algorithms.NewECCSigner())
}

// IsAlgorithmRegistered checks if a specific algorithm is registered.
func IsAlgorithmRegistered(name string) bool {
	_, kbOk := algorithmKeyBuildersRegistry[name]
	_, sOk := algorithmSignersRegistry[name]
	return kbOk && sOk
}

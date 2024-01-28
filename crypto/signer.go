package crypto

type Signer struct{}

// NewSigner creates a new Signer.
func NewSigner() *Signer {
	return &Signer{}
}

// IsValidAlgorithm checks if a specific algorithm is registered.
func (sg *Signer) IsValidAlgorithm(algorithm string) bool {
	_, exists := algorithmSignersRegistry[algorithm]
	return exists
}

// Sign signs a message using a specific algorithm.
func (sg *Signer) Sign(algorithm string, privateKeyBytes, dataToBeSigned []byte) ([]byte, error) {
	if !sg.IsValidAlgorithm(algorithm) {
		return nil, ErrCryptoEngineNotFound
	}

	return algorithmSignersRegistry[algorithm].Sign(privateKeyBytes, dataToBeSigned)
}

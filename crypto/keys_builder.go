package crypto

type KeysBuilder struct{}

// NewKeysBuilder creates a new KeysBuilder.
func NewKeysBuilder() *KeysBuilder {
	return &KeysBuilder{}
}

// IsValidAlgorithm checks if a specific algorithm is registered.
func (kg *KeysBuilder) IsValidAlgorithm(algorithm string) bool {
	return IsAlgorithmRegistered(algorithm)
}

// Build builds a new key pair using a specific algorithm.
func (kg *KeysBuilder) Build(algorithm string) ([]byte, []byte, error) {
	if !kg.IsValidAlgorithm(algorithm) {
		return nil, nil, ErrCryptoEngineNotFound
	}

	return algorithmKeyBuildersRegistry[algorithm].Keys()
}

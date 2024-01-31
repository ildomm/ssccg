package crypto_test

import (
	"github.com/ildomm/ssccg/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlgorithmRegistration(t *testing.T) {
	t.Run("RSA_Registered", func(t *testing.T) {
		assert.True(t, crypto.IsAlgorithmRegistered("RSA"))
	})

	t.Run("ECDSA_Registered", func(t *testing.T) {
		assert.True(t, crypto.IsAlgorithmRegistered("ECDSA"))
	})

	t.Run("UnregisteredAlgorithm", func(t *testing.T) {
		assert.False(t, crypto.IsAlgorithmRegistered("NonExistent"))
	})
}

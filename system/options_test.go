package system

import (
	"github.com/ildomm/ssccg/test_helpers"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestExtractServerPort tests the extractServerPort function.
func TestExtractServerPort(t *testing.T) {
	t.Run("ValidPort", func(t *testing.T) {
		os.Setenv(ListenAddressEnvVar, "8080")
		defer os.Unsetenv(ListenAddressEnvVar)

		port := ExtractServerPort()
		assert.NotNil(t, port, "Port should not be nil")
		assert.Equal(t, 8080, *port, "Port value mismatch")
	})

	t.Run("NoEnvVar", func(t *testing.T) {
		os.Unsetenv(ListenAddressEnvVar)
		port := ExtractServerPort()
		assert.Nil(t, port, "Port should be nil when environment variable is not set")
	})

	t.Run("InvalidPort", func(t *testing.T) {
		os.Setenv(ListenAddressEnvVar, "invalid")
		defer os.Unsetenv(ListenAddressEnvVar)

		buf, restoreLog := test_helpers.CaptureOutput()
		defer restoreLog()
		port := ExtractServerPort()
		assert.Nil(t, port, "Port should be nil when environment variable is not set")

		assert.Contains(t, buf.String(), "Could not parse server port address", "Expected log message not found")
	})
}

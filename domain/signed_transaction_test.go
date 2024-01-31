package domain

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestSignedData tests the SignedData method of the SignedTransaction struct.
func TestSignedData(t *testing.T) {
	transaction := SignedTransaction{
		ID:                 uuid.New(),
		DeviceID:           uuid.New(),
		RawData:            []byte("sampledata"),
		PreviousDeviceSign: "previous-signature",
		SignCounter:        5,
	}

	expected := fmt.Sprintf("%d_%s_%s", transaction.SignCounter, transaction.RawData, transaction.PreviousDeviceSign)
	result := transaction.SignedData()
	assert.Equal(t, expected, result, "SignedData method returned unexpected result")
}

package algorithms

import (
	"crypto/sha256"
	"fmt"
)

// GetHashSum returns the hash sum of the data to be signed.
func GetHashSum(dataToBeSigned []byte) ([]byte, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write(dataToBeSigned)
	if err != nil {
		return nil, fmt.Errorf("failed to get hash sum: %w", err)
	}
	return msgHash.Sum(nil), nil
}

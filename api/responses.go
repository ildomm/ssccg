package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError))) //nolint:all
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes) //nolint:all
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes) //nolint:all
}

// HealthResponse represents the response for the health check.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// CreateDeviceResponse represents the response for creating a device.
type CreateDeviceResponse struct {
	ID            uuid.UUID `db:"id"`
	Label         string    `db:"label"`
	SignAlgorithm string    `db:"sign_algorithm"`
	PublicKey     string    `db:"public_key"`
}

// CreateSignedTransactionResponse represents the response for a signed transaction.
type CreateSignedTransactionResponse struct {
	ID         uuid.UUID `json:"id"`
	Signature  string    `json:"signature"`
	SignedData string    `json:"signed_data"`
}

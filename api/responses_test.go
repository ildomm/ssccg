package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWriteInternalError tests the WriteInternalError function.
func TestWriteInternalError(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteInternalError(rr)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "expected status internal server error")
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), rr.Body.String(), "unexpected body content")
}

// TestWriteErrorResponse tests the WriteErrorResponse function.
func TestWriteErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	errors := []string{"error1", "error2"}

	WriteErrorResponse(rr, http.StatusBadRequest, errors)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "unexpected status code")

	var resp ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, errors, resp.Errors, "unexpected errors in response")
}

// TestWriteAPIResponse tests the WriteAPIResponse function.
func TestWriteAPIResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	WriteAPIResponse(rr, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, rr.Code, "unexpected status code")

	expectedBytes, err := json.MarshalIndent(Response{Data: data}, "", "  ")
	require.NoError(t, err)

	expected := string(expectedBytes)
	actual := rr.Body.String()

	assert.JSONEq(t, expected, actual, "unexpected data in response")
}

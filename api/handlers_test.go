package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthSuccess tests the Health function for a successful response.
func TestHealthHandlerSuccess(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server := Server{} // Assuming Server struct exists
		server.HealthHandler(w, r)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	expected := Response{}
	expected.Data = HealthResponse{Status: "pass", Version: "v0"}
	body, err := io.ReadAll(rr.Body)
	require.NoError(t, err)

	var actual Response
	actual.Data = HealthResponse{}

	err = json.Unmarshal(body, &actual)
	require.NoError(t, err)
}

// TestHealthMethodNotAllowed tests the Health function for an incorrect method.
func TestHealthHandlerMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server := Server{} // Assuming Server struct exists
		server.HealthHandler(w, r)
	})

	handler.ServeHTTP(rr, req)

	// Check that the status code is 405 Method Not Allowed.
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "handler allowed incorrect method")
}

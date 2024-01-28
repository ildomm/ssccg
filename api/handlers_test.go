package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/ildomm/ssccg/domain"
	"github.com/ildomm/ssccg/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

// TestListDeviceFuncSuccess tests the ListDeviceFunc for a successful response.
func TestListDeviceFuncSuccess(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/devices", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockDAO := test_helpers.NewMockDeviceDAO()
		mockDAO.On("GetDevices").Return([]domain.Device{}, nil)
		deviceHandler := NewDeviceHandler(mockDAO)
		deviceHandler.ListDeviceFunc(w, r)
	})

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "ListDeviceFunc returned wrong status code")
}

// TestCreateDeviceFuncSuccess tests the CreateDeviceFunc for a successful response using a real server.
func TestCreateDeviceFuncSuccess(t *testing.T) {
	mockDAO := test_helpers.NewMockDeviceDAO()

	// Set up mock expectations
	testDevice := &domain.Device{
		ID:            uuid.New(),
		Label:         "Test Device",
		SignAlgorithm: "RSA",
		PublicKey:     "publicKey",
		PrivateKey:    "privateKey",
	}
	mockDAO.On("CreateDevice", mock.AnythingOfType("uuid.UUID"), "Test Device", "RSA").Return(testDevice, nil)

	// Create the server and set the mock manager
	server := NewServer()
	server.WithDeviceManager(mockDAO)

	// Create the request body
	reqBody := CreateDeviceRequest{
		Algorithm: "RSA",
		Label:     "Test Device",
	}
	body, _ := json.Marshal(reqBody)

	// Create the request
	_, err := http.NewRequest(http.MethodPost, "/api/v1/devices/"+testDevice.ID.String(), bytes.NewReader(body))
	require.NoError(t, err)

	// Use httptest to create a server
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	// Execute the request
	resp, err := http.Post(testServer.URL+"/api/v1/devices/"+testDevice.ID.String(), "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "CreateDeviceFunc returned wrong status code")

	// Decode the response
	var respBody Response
	respBody.Data = CreateDeviceResponse{}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
}

// TestListSignatureFuncSuccess tests the ListSignatureFunc for a successful response.
func TestListSignatureFuncSuccess(t *testing.T) {
	mockDAO := test_helpers.NewMockDeviceDAO()

	// Set up mock expectations
	testSignatures := []domain.SignedTransaction{
		{
			ID:          uuid.New(),
			DeviceID:    uuid.New(),
			Sign:        "signature1",
			RawData:     []byte("data1"),
			SignCounter: 1,
		},
		{
			ID:          uuid.New(),
			DeviceID:    uuid.New(),
			Sign:        "signature2",
			RawData:     []byte("data2"),
			SignCounter: 2,
		},
	}
	deviceID := uuid.New()
	mockDAO.On("GetSignedTransactions", deviceID).Return(testSignatures, nil)

	// Create the server and set the mock manager
	server := NewServer()
	server.WithDeviceManager(mockDAO)

	// Use httptest to create a server
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	// Execute the request
	resp, err := http.Get(testServer.URL + "/api/v1/devices/" + deviceID.String() + "/signatures")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "ListSignatureFunc returned wrong status code")

	// Decode the response
	var respBody Response
	respBody.Data = []SignedTransactionResponse{}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
}

// TestCreateSignatureFuncSuccess tests the CreateSignatureFunc for a successful response.
func TestCreateSignatureFuncSuccess(t *testing.T) {
	mockDAO := test_helpers.NewMockDeviceDAO()

	deviceId := uuid.New()
	testSignature := &domain.SignedTransaction{
		ID:                 uuid.New(),
		DeviceID:           deviceId,
		Sign:               "signature",
		RawData:            []byte("data"),
		SignCounter:        1,
		PreviousDeviceSign: "previousSignature",
	}
	mockDAO.On("CreateSignedTransaction", deviceId, mock.AnythingOfType("[]uint8")).Return(testSignature, nil)

	// Create the server and set the mock manager
	server := NewServer()
	server.WithDeviceManager(mockDAO)

	// Create the request body
	reqBody := SignTransactionRequest{Data: "data"}
	body, _ := json.Marshal(reqBody)

	// Create the request
	url := "/api/v1/devices/" + deviceId.String() + "/signatures"
	_, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)

	// Use httptest to create a server
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	// Execute the request
	resp, err := http.Post(testServer.URL+url, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "CreateSignatureFunc returned wrong status code")

	// Decode the response
	var respBody SignedTransactionResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
}

func TestCreateSignatureFuncBadRequest(t *testing.T) {
	server := NewServer()

	// Invalid device ID
	url := "/api/v1/devices/invalid-id/signatures"
	body, _ := json.Marshal(SignTransactionRequest{Data: "data"})

	_, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)

	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	resp, err := http.Post(testServer.URL+url, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected BadRequest for invalid device ID")
}

func TestCreateSignatureFuncInternalError(t *testing.T) {
	mockDAO := test_helpers.NewMockDeviceDAO()
	mockDAO.On("CreateSignedTransaction",
		mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("[]uint8")).
		Return(nil, errors.New("internal error"))

	server := NewServer()
	server.WithDeviceManager(mockDAO)

	deviceId := uuid.New()
	reqBody := SignTransactionRequest{Data: "data"}
	body, _ := json.Marshal(reqBody)

	url := "/api/v1/devices/" + deviceId.String() + "/signatures"
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	resp, err := http.Post(testServer.URL+url, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Expected InternalServerError for simulated service error")
}

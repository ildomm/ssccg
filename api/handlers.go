package api

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ildomm/ssccg/dao"
	"github.com/ildomm/ssccg/domain"
	"net/http"
)

// HealthHandler evaluates the health of the service and writes a standardized response.
func (s *Server) HealthHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	health := HealthResponse{
		Status:  "pass",
		Version: "v1",
	}

	WriteAPIResponse(response, http.StatusOK, health)
}

// deviceHandler handles all requests related to devices.
type deviceHandler struct {
	deviceDAO dao.DeviceDAO
}

func NewDeviceHandler(deviceDAO dao.DeviceDAO) *deviceHandler {
	return &deviceHandler{
		deviceDAO: deviceDAO,
	}
}

// Transform domain.Device to api.DeviceResponse
func transformToDeviceResponse(device domain.Device) DeviceResponse {
	return DeviceResponse{
		ID:            device.ID,
		Label:         device.Label,
		SignAlgorithm: device.SignAlgorithm,
		PublicKey:     device.PublicKey,
	}
}

// ListDeviceFunc handles the request to list all devices.
func (h *deviceHandler) ListDeviceFunc(w http.ResponseWriter, r *http.Request) {
	devices, err := h.deviceDAO.GetDevices()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	deviceResponses := make([]DeviceResponse, 0, len(devices))
	for _, device := range devices {
		deviceResponses = append(deviceResponses, transformToDeviceResponse(device))
	}

	WriteAPIResponse(w, http.StatusOK, deviceResponses)
}

// CreateDeviceFunc handles the request to create a new device.
func (h *deviceHandler) CreateDeviceFunc(w http.ResponseWriter, r *http.Request) {
	var req CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"invalid request body"})
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"invalid device ID"})
		return
	}

	device, err := h.deviceDAO.CreateDevice(id, req.Label, req.Algorithm)
	if err != nil {
		if errors.Is(err, dao.ErrDeviceExists) || errors.Is(err, dao.ErrInvalidAlgorithm) {
			WriteErrorResponse(w, http.StatusBadRequest, []string{err.Error()})
		} else {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}
		return
	}

	deviceResponse := transformToDeviceResponse(*device)
	WriteAPIResponse(w, http.StatusCreated, deviceResponse)
}

// GetDeviceFunc handles the request to retrieve a specific device.
func (h *deviceHandler) GetDeviceFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"invalid device ID"})
		return
	}

	device, err := h.deviceDAO.GetDevice(id)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, []string{err.Error()})
		return
	}

	deviceResponse := transformToDeviceResponse(*device)
	WriteAPIResponse(w, http.StatusOK, deviceResponse)
}

// Transform domain.SignedTransaction to api.SignedTransactionResponse
func transformToSignedTransactionResponse(transaction domain.SignedTransaction) SignedTransactionResponse {
	return SignedTransactionResponse{
		ID:         transaction.ID,
		Signature:  transaction.Sign,
		SignedData: transaction.SignedData(),
	}
}

// CreateSignatureFunc handles the request to create a signature for a device.
func (h *deviceHandler) CreateSignatureFunc(w http.ResponseWriter, r *http.Request) {
	var req SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"invalid request body"})
		return
	}

	vars := mux.Vars(r)
	deviceId, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"invalid device ID"})
		return
	}

	signed, err := h.deviceDAO.CreateSignedTransaction(deviceId, []byte(req.Data))
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	signedResponse := transformToSignedTransactionResponse(*signed)
	WriteAPIResponse(w, http.StatusCreated, signedResponse)
}

// ListSignatureFunc handles the request to list signatures for a device.
func (h *deviceHandler) ListSignatureFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"invalid device ID"})
		return
	}

	signatures, err := h.deviceDAO.GetSignedTransactions(deviceId)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	signaturesResponses := make([]SignedTransactionResponse, 0, len(signatures))
	for _, signature := range signatures {
		signaturesResponses = append(signaturesResponses, transformToSignedTransactionResponse(signature))
	}

	WriteAPIResponse(w, http.StatusOK, signaturesResponses)
}

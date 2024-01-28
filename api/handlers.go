package api

import (
	"github.com/ildomm/ssccg/service"
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

type deviceHandler struct {
	deviceManager *service.Manager
}

func NewDeviceHandler(deviceManager *service.Manager) *deviceHandler {
	return &deviceHandler{
		deviceManager: deviceManager,
	}
}

func (h *deviceHandler) ListDeviceFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *deviceHandler) GetDeviceFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// CreateDeviceFunc creates a new device.
func (h *deviceHandler) CreateDeviceFunc(w http.ResponseWriter, r *http.Request) {
	// Extract params
	//id := mux.Vars(r)["id"]

	w.WriteHeader(http.StatusOK)
}

func (h *deviceHandler) ListSignatureFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *deviceHandler) CreateSignatureFunc(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

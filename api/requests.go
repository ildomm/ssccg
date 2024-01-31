package api

// CreateDeviceRequest represents the request body for creating a device.
type CreateDeviceRequest struct {
	Algorithm string `json:"algorithm"`
	Label     string `json:"label,omitempty"`
}

// SignTransactionRequest represents the request body for creating a signature.
type SignTransactionRequest struct {
	Data string `json:"data"`
}

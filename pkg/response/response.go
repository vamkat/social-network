package response

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Common structure for all API responses
type JSONResponse struct {
	Success   bool        `json:"success"`
	Timestamp string      `json:"timestamp"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
}

// success response helper
func Success(w http.ResponseWriter, data interface{}, message string) {
	WriteJSON(w, http.StatusOK, JSONResponse{
		Success:   true,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   message,
		Data:      data,
	})
}

// error response helper
func Error(w http.ResponseWriter, status int, message string, err error, log *zap.Logger) {
	if log != nil && err != nil {
		log.Error("API error",
			zap.String("message", message),
			zap.Error(err),
		)
	}

	WriteJSON(w, status, JSONResponse{
		Success:   false,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   message,
		Error: map[string]string{
			"details": errString(err),
		},
	})
}

// internal utility
func WriteJSON(w http.ResponseWriter, status int, payload JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// handle nil safely
func errString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

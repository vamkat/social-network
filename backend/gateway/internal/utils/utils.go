package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}
	// fmt.Println("sending this:", v)
	return json.NewEncoder(w).Encode(v)
}

func ErrorJSON(w http.ResponseWriter, code int, msg string) {
	err := WriteJSON(w, code, map[string]string{"error": msg})
	if err != nil {
		fmt.Printf("Failed to send error message: %s, code: %d, %s\n", msg, code, err)
	}
}

func B64urlEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func B64urlDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

func GenUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

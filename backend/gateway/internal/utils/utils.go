package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// Type that holds the context with value keys
type ctxKey string

// Holds the keys to values on request context.
const (
	ClaimsKey        ctxKey = "jwtClaims"
	ReqId            ctxKey = "X-Request-Id"
	ReqActionDetails ctxKey = "X-Action-Details"
	ReqTimestamp     ctxKey = "X-Timestamp"
)

// Adds value val to r context with key 'key'
func RequestWithValue[T any](r *http.Request, key ctxKey, val T) *http.Request {
	ctx := context.WithValue(r.Context(), key, val)
	return r.WithContext(ctx)
}

// Get value T from request context with key 'key'
func GetValue[T any](r *http.Request, key ctxKey) (T, bool) {
	v := r.Context().Value(key)
	if v == nil {
		fmt.Println("v is nil")
		var zero T
		return zero, false
	}
	c, ok := v.(T)
	if !ok {
		panic(1) // this should never happen, which is why I'm putting a panic here so that this mistake is obvious
	}
	return c, ok
}

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}
	fmt.Println("sending this:", v)
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

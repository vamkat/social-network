package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"slices"
	"social-network/shared/go/customtypes"

	"github.com/google/uuid"
)

// Adds value val to r context with key 'key'
func RequestWithValue(r *http.Request, key customtypes.CtxKey, val any) *http.Request {
	ctx := context.WithValue(r.Context(), key, val)
	return r.WithContext(ctx)
}

// Get value T from request context with key 'key'
func GetValue[T any](r *http.Request, key customtypes.CtxKey) (T, bool) {
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

var ErrUnmarshalFailed = errors.New("unmashal failed")

func JSON2Struct[T any](dataStruct *T, request *http.Request) (*T, error) {
	decoder := json.NewDecoder(request.Body)
	defer request.Body.Close()
	err := decoder.Decode(&dataStruct)
	if err != nil {
		return dataStruct, ErrUnmarshalFailed
	}
	return dataStruct, nil
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

var (
	ErrImageTooBig      = errors.New("image too big")
	ErrInvalidImageFile = errors.New("invalid image file: Only PNG, JPG, or GIF allowed")
	ImageTypes          = []string{"jpg, png, svg"}
)

// Parses the image file and stores it to the configured path. Returns a uuid as filename
func CheckImage(file multipart.File, header *multipart.FileHeader) (filetype string, err error) {
	if header.Size > 10*1024*1024 {
		return "", ErrImageTooBig
	}

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return "", ErrInvalidImageFile
	}
	filetype = http.DetectContentType(buf)
	file.Seek(0, io.SeekStart)

	if slices.Contains(ImageTypes, filetype) {
		return "", ErrInvalidImageFile
	}

	return filetype, nil
}

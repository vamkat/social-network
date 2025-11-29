package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/speps/go-hashids/v2"
)

var salt string = os.Getenv("ENC_KEY")

var hd = func() *hashids.HashID {
	h := hashids.NewData()
	h.Salt = salt
	h.MinLength = 12
	encoder, _ := hashids.NewWithData(h)
	return encoder
}()

type EncryptedId int64

func (e EncryptedId) MarshalJSON() ([]byte, error) {
	hash, err := hd.EncodeInt64([]int64{int64(e)})
	if err != nil {
		return nil, err
	}
	return json.Marshal(hash)
}

func (e *EncryptedId) UnmarshalJSON(data []byte) error {
	var hash string
	if err := json.Unmarshal(data, &hash); err != nil {
		return err
	}

	decoded, err := hd.DecodeInt64WithError(hash)
	if err != nil || len(decoded) == 0 {
		return err
	}

	*e = EncryptedId(decoded[0])
	return nil
}

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

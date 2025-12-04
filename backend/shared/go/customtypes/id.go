package customtypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"os"

	"github.com/lib/pq"
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

func (e EncryptedId) IsValid() bool {
	return e >= 0
}

func (e EncryptedId) Validate() error {
	if !e.IsValid() {
		return errors.Join(ErrValidation, errors.New("encryptedId must be positive"))
	}
	return nil
}

func (e EncryptedId) Int64() int64 {
	return int64(e)
}

// ------------------------------------------------------------
// Id
// ------------------------------------------------------------

func (i Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(i))
}

func (i *Id) UnmarshalJSON(data []byte) error {
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*i = Id(v)
	return nil
}

func (i Id) IsValid() bool {
	return i >= 0
}

func (i Id) Validate() error {
	if !i.IsValid() {
		return errors.Join(ErrValidation, errors.New("id must be positive"))
	}
	return nil
}

func (i Id) Int64() int64 {
	return int64(i)
}

// ------------------------------------------------------------
// Ids
// ------------------------------------------------------------

func (ids Ids) MarshalJSON() ([]byte, error) {
	return json.Marshal(ids.Int64())
}

func (ids *Ids) UnmarshalJSON(data []byte) error {
	var raw []int64
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	out := make(Ids, len(raw))
	for i, v := range raw {
		out[i] = Id(v)
	}

	*ids = out
	return nil
}

func (ids Ids) IsValid() bool {
	for _, i := range ids {
		if !i.IsValid() {
			return false
		}
	}
	return true
}

func (ids Ids) Validate() error {
	if !ids.IsValid() {
		return errors.Join(ErrValidation, errors.New("all ids must be positive"))
	}
	return nil
}

func (ids Ids) Int64() []int64 {
	out := make([]int64, len(ids))
	for i, v := range ids {
		out[i] = int64(v)
	}
	return out
}

// Value implements driver.Valuer
func (ids Ids) Value() (driver.Value, error) {
	// Convert []Id to []int64
	int64s := make([]int64, len(ids))
	for i, id := range ids {
		int64s[i] = int64(id)
	}
	return pq.Int64Array(int64s).Value()
}

// Scan implements sql.Scanner
func (ids *Ids) Scan(src interface{}) error {
	var int64Array pq.Int64Array
	if err := int64Array.Scan(src); err != nil {
		return err
	}
	// Convert []int64 to []Id
	*ids = make(Ids, len(int64Array))
	for i, v := range int64Array {
		(*ids)[i] = Id(v)
	}
	return nil
}

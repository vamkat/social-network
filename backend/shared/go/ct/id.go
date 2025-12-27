package ct

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/lib/pq"
	"github.com/speps/go-hashids/v2"
)

// Allows null values.
// When Umarshaled to JSON format the int64 value is encrypted using "github.com/speps/go-hashids/v2".
// Relies on enviromental variable "ENC_KEY" to be present.
type Id int64

// var salt string = os.Getenv("ENC_KEY")

var hd = func() *hashids.HashID {
	h := hashids.NewData()
	h.Salt = Cfgs.Salt
	h.MinLength = 12
	encoder, _ := hashids.NewWithData(h)
	return encoder
}()

func (e Id) MarshalJSON() ([]byte, error) {
	hash, err := hd.EncodeInt64([]int64{int64(e)})
	if err != nil {
		return nil, err
	}
	return json.Marshal(hash)
}

func (e *Id) UnmarshalJSON(data []byte) error {
	var hash string
	if err := json.Unmarshal(data, &hash); err != nil {
		return err
	}

	decoded, err := hd.DecodeInt64WithError(hash)
	if err != nil || len(decoded) == 0 {
		return err
	}

	*e = Id(decoded[0])
	return nil
}

func (e Id) IsValid() bool {
	return e > 0
}

func DecryptId(hash string) (Id, error) {
	decoded, err := hd.DecodeInt64WithError(hash)
	if err != nil || len(decoded) == 0 {
		return 0, err
	}
	return Id(decoded[0]), nil
}

func EncryptId(id int64) (encrypted string, err error) {
	return hd.EncodeInt64([]int64{int64(id)})
}

func (i *Id) Scan(src any) error {
	if src == nil {
		// SQL NULL reached
		*i = 0 // or whatever "invalid" means in your domain
		return nil
	}

	switch v := src.(type) {
	case int64:
		*i = Id(v)
		return nil

	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		*i = Id(n)
		return nil
	}

	return fmt.Errorf("cannot scan type %T into Id", src)
}

func (i Id) Value() (driver.Value, error) {
	if !i.IsValid() {
		return nil, nil
	}
	return i.Int64(), nil
}

func (e Id) Validate() error {
	if !e.IsValid() {
		return errors.Join(ErrValidation, errors.New("id must be positive"))
	}
	return nil
}

func (e Id) Int64() int64 {
	return int64(e)
}

// ------------------------------------------------------------
// Id
// ------------------------------------------------------------

// Validation requires for the int64 value to be above zero.
type UnsafeId int64

func (i UnsafeId) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(i))
}

func (i *UnsafeId) UnmarshalJSON(data []byte) error {
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*i = UnsafeId(v)
	return nil
}

func (i UnsafeId) IsValid() bool {
	return i > 0
}

func (i UnsafeId) Validate() error {
	if !i.IsValid() {
		return errors.Join(ErrValidation, errors.New("id must be positive"))
	}
	return nil
}

func (i *UnsafeId) Scan(src any) error {
	if src == nil {
		// SQL NULL reached
		*i = 0
		return nil
	}

	switch v := src.(type) {
	case int64:
		*i = UnsafeId(v)
		return nil

	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		*i = UnsafeId(n)
		return nil
	}

	return fmt.Errorf("cannot scan type %T into Id", src)
}

func (i UnsafeId) Value() (driver.Value, error) {
	if !i.IsValid() {
		return nil, nil
	}
	return i.Int64(), nil
}

func (i UnsafeId) Int64() int64 {
	return int64(i)
}

// ------------------------------------------------------------
// Ids
// ------------------------------------------------------------

// Validation does not allow zero len slices or indexes with null values.
type Ids []Id

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
	if len(ids) == 0 {
		return false
	}
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
func (ids *Ids) Scan(src any) error {
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

// Returns an Ids type with all unique entries from the given instance.
func (ids Ids) Unique() Ids {
	uniq := make(map[Id]struct{}, len(ids))
	cleaned := make([]Id, 0, len(ids))
	for _, id := range ids {
		if _, ok := uniq[id]; !ok {
			uniq[id] = struct{}{}
			cleaned = append(cleaned, id)
		}
	}
	return Ids(cleaned)
}

// FromInt64s converts a []int64 to Ids.
func FromInt64s(src []int64) Ids {
	out := make(Ids, len(src))
	for i, v := range src {
		out[i] = Id(v)
	}
	return out
}

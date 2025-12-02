package customtypes

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ------------------------------------------------------------
// Limit
// ------------------------------------------------------------

// Non zero type. Validation returns error if zero or above limit

func (l Limit) MarshalJSON() ([]byte, error) {
	return json.Marshal(int32(l))
}

func (l *Limit) UnmarshalJSON(data []byte) error {
	var v int32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*l = Limit(v)
	return nil
}

func (l Limit) IsValid() bool {
	return l >= 1 && l <= Limit(maxLimit)
}

func (l Limit) Validate() error {
	if !l.IsValid() {
		return errors.Join(ErrValidation, fmt.Errorf("limit must be between 1 and %d", maxLimit))
	}
	return nil
}

func (l Limit) Int32() int32 {
	return int32(l)
}

// ------------------------------------------------------------
// Offset
// ------------------------------------------------------------

// Non negative type. Validation returns error if below zero or above limit

func (o Offset) MarshalJSON() ([]byte, error) {
	return json.Marshal(int32(o))
}

func (o *Offset) UnmarshalJSON(data []byte) error {
	var v int32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*o = Offset(v)
	return nil
}

func (o Offset) IsValid() bool {
	return o >= 0
}

func (o Offset) Validate() error {
	if !o.IsValid() {
		return errors.Join(ErrValidation, errors.New("offset must be >= 0"))
	}
	return nil
}

func (o Offset) Int32() int32 {
	return int32(o)
}

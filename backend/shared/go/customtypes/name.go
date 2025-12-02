package customtypes

import (
	"encoding/json"
	"errors"
)

// ------------------------------------------------------------
// Name
// ------------------------------------------------------------

// General type for names and surnames. Name type is not nullable. All smaller than len 2 values return error

func (n Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(n))
}

func (n *Name) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*n = Name(s)
	return nil
}

func (n Name) IsValid() bool {
	// empty check for now. Add regex for name
	return len(n) > 1
}

func (n Name) Validate() error {
	if !n.IsValid() {
		return errors.Join(ErrValidation, errors.New("name must be at least 2 characters"))
	}
	return nil
}

func (n Name) String() string {
	return string(n)
}

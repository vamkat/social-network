package customtypes

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ------------------------------------------------------------
// About
// ------------------------------------------------------------

// Can be used for bio or descritpion. About is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.

func (a About) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(a))
}

func (a *About) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*a = About(s)
	return nil
}

func (a About) IsValid() bool {
	if len(a) == 0 {
		return true
	}
	if len(a) < aboutCharsMin || len(a) > aboutCharsMax {
		return false
	}
	for _, r := range a {
		if r < 32 { // control characters
			return false
		}
	}
	return true
}

func (a About) Validate() error {
	if !a.IsValid() {
		return errors.Join(ErrValidation,
			fmt.Errorf("about must be %d–%d chars and contain no control characters",
				aboutCharsMin,
				aboutCharsMax,
			))
	}
	return nil
}

func (a About) String() string {
	return string(a)
}

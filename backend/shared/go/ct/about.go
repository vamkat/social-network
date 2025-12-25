package ct

//CT stands for Custom Types

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Can be used for bio or descritpion.
//
// Usage:
//
//	var bioCt ct
//	var bioStr string
//	bioCt = ct.About("about me")
//	bioStr = bio.String()
type About string

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
		return false
	}
	if len(a) < aboutCharsMin || len(a) > aboutCharsMax {
		return false
	}

	return controlCharsFree(a.String())
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

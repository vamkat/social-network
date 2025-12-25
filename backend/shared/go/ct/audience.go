package ct

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ------------------------------------------------------------
// Audience
// ------------------------------------------------------------

// Can be used for post, comment, event body.
type Audience string

func (au Audience) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(au))
}

func (au *Audience) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*au = Audience(s)
	return nil
}

func (au Audience) IsValid() bool {
	if au == "" {
		return false
	}
	for _, permittedValue := range permittedAudienceValues {
		if strings.EqualFold(au.String(), permittedValue) {
			return true
		}
	}
	return controlCharsFree(au.String())
}

func (au Audience) Validate() error {
	if !au.IsValid() {
		return errors.Join(ErrValidation,
			fmt.Errorf("audience must be one of the following: %v",
				permittedAudienceValues,
			))
	}
	return nil
}

func (au Audience) String() string {
	return string(au)
}

package customtypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ------------------------------------------------------------
// Title (group/chat title)
// ------------------------------------------------------------

// Refers to title of content not mr, mrs. Title is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.

func (t Title) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

func (t *Title) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = Title(s)
	return nil
}

func (t Title) IsValid() bool {
	if t == "" {
		return true
	}

	s := strings.TrimSpace(string(t))
	if len(s) < minTitleChars || len(s) > maxTitleChars {
		return false
	}

	// No control chars
	for _, r := range s {
		if r < 32 {
			return false
		}
	}

	return true
}

func (t Title) Validate() error {
	if !t.IsValid() {
		return errors.Join(ErrValidation,
			fmt.Errorf("invalid title: must be %d-%d chars, no control characters, and not an honorific",
				minTitleChars,
				maxTitleChars,
			))
	}
	return nil
}

func (t Title) String() string {
	return string(t)
}

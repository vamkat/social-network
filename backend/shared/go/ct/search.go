package ct

import (
	"encoding/json"
	"errors"
	"regexp"
)

// ------------------------------------------------------------
// Search
// ------------------------------------------------------------

// SearchTerm represents a validated search query term
type SearchTerm string

func (s SearchTerm) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *SearchTerm) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = SearchTerm(str)
	return nil
}

// IsValid checks if the search term meets minimum validation rules.
func (s SearchTerm) IsValid() bool {
	// Basic length check
	if len(s) < 2 {
		return false
	}

	// Optional: enforce allowed characters (letters, numbers, spaces, hyphens)
	// Adjust regex as needed.
	re := regexp.MustCompile(`^[A-Za-z0-9\s\-]+$`)
	return re.MatchString(string(s)) && controlCharsFree(s.String())
}

// Validate returns a descriptive error if the value is invalid.
func (s SearchTerm) Validate() error {
	if len(s) < 1 {
		return errors.Join(
			ErrValidation,
			errors.New("search term must be at least 1 character"),
		)
	}

	// Same regex as IsValid()
	re := regexp.MustCompile(`^[A-Za-z0-9\s\-]+$`)
	if !re.MatchString(string(s)) {
		return errors.Join(
			ErrValidation,
			errors.New("search term contains invalid characters"),
		)
	}

	return nil
}

func (s SearchTerm) String() string {
	return string(s)
}

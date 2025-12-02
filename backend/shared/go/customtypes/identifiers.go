package customtypes

import (
	"encoding/json"
	"errors"
)

// ------------------------------------------------------------
// Email
// ------------------------------------------------------------

func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

func (e *Email) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*e = Email(s)
	return nil
}

func (e Email) IsValid() bool {
	return emailRegex.MatchString(string(e))
}

func (e Email) Validate() error {
	if !e.IsValid() {
		return errors.Join(ErrValidation, errors.New("invalid email format"))
	}
	return nil
}

func (e Email) String() string {
	return string(e)
}

// ------------------------------------------------------------
// Username
// ------------------------------------------------------------

func (u Username) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(u))
}

func (u *Username) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*u = Username(s)
	return nil
}

func (u Username) IsValid() bool {
	if u == "" {
		return true
	}
	return usernameRegex.MatchString(string(u))
}

func (u Username) Validate() error {
	if !u.IsValid() {
		return errors.Join(ErrValidation, errors.New("invalid username format"))
	}
	return nil
}

func (u Username) String() string {
	return string(u)
}

// ------------------------------------------------------------
// Identifier (username or email)
// ------------------------------------------------------------

func (i Identifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(i))
}

func (i *Identifier) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*i = Identifier(s)
	return nil
}

func (i Identifier) IsValid() bool {
	if i == "" {
		return true
	}
	s := string(i)
	return usernameRegex.MatchString(s) || emailRegex.MatchString(s)
}

func (i Identifier) Validate() error {
	if !i.IsValid() {
		return errors.Join(ErrValidation, errors.New("identifier must be a valid username or email"))
	}
	return nil
}

func (i Identifier) String() string {
	return string(i)
}

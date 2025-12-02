package customtypes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
)

// ------------------------------------------------------------
// Password
// (Hash on Unmarshal; store hashed value only)
// ------------------------------------------------------------

func (p Password) MarshalJSON() ([]byte, error) {
	// No encoder required – return placeholder or omit
	return json.Marshal("********")
}

func (p *Password) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) == 0 {
		return errors.Join(ErrValidation, errors.New("password is required"))
	}

	secret := os.Getenv("PASSWORD_SECRET")
	if secret == "" {
		return errors.Join(ErrValidation, errors.New("missing env var PASSWORD_SECRET"))
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(raw))
	hashed := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	*p = Password(hashed)
	return nil
}

func (p Password) IsValid() bool {
	// After hashing, length always valid; check before hashing instead? Up to you.
	return len(p) > 0
}

func (p Password) Validate() error {
	if !p.IsValid() {
		return errors.Join(ErrValidation, errors.New("invalid password"))
	}
	return nil
}

func (p Password) String() string {
	return string(p)
}

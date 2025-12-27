package ct

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

// ------------------------------------------------------------
// Password
// ------------------------------------------------------------

// Password is not nullable. The length is checked and error is returned during json unmarshall and validation methods.
type Password string

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

	*p = Password(raw)
	return nil
}

func (p Password) Hash() (Password, error) {
	fmt.Println(Cfgs.PassSecret)
	// TODO: Change this to configs when in place
	// secret := func() string {
	// 	s := Cfgs.PassSecret
	// 	if s == "" {
	// 		s = os.Getenv("PASSWORD_SECRET")
	// 	}
	// 	return s
	// }()
	secret := Cfgs.PassSecret

	if secret == "" {
		return "", errors.Join(ErrValidation, errors.New("missing env var PASSWORD_SECRET"))
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(p))
	p = Password(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	return p, nil
}

// one symbol, one capital letter, one number min 8 chars max 64 chars
var (
	uppercase = regexp.MustCompile(`[A-Z]`)
	lowercase = regexp.MustCompile(`[a-z]`)
	digit     = regexp.MustCompile(`[0-9]`)
	symbol    = regexp.MustCompile(`[^A-Za-z0-9]`)
)

// Validates raw password for one symbol, one capital letter, one number, min 8 chars, max 64 chars
func (p Password) IsValid() bool {
	s := string(p)

	if len(s) < 8 || len(s) > 64 {
		return false
	}
	if !uppercase.MatchString(s) {
		return false
	}
	if !lowercase.MatchString(s) {
		return false
	}
	if !digit.MatchString(s) {
		return false
	}
	if !symbol.MatchString(s) {
		return false
	}
	return controlCharsFree(p.String())
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

// ------------------------------------------------------------
// HashedPassword
// ------------------------------------------------------------

// HashedPassword is not nullable. The length is checked and error is returned during json unmarshall and validation methods.
type HashedPassword string

func (hp HashedPassword) MarshalJSON() ([]byte, error) {
	// No encoder required – return placeholder or omit
	return json.Marshal("********")
}

func (hp *HashedPassword) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) == 0 {
		return errors.Join(ErrValidation, errors.New("password is required"))
	}

	*hp = HashedPassword(raw)
	return nil
}

func (hp HashedPassword) IsValid() bool {
	return hp != "" && controlCharsFree(hp.String())
}

func (hp HashedPassword) Validate() error {
	if !hp.IsValid() {
		return errors.Join(ErrValidation, errors.New("invalid hashed password"))
	}
	return nil
}

func (hp HashedPassword) String() string {
	return string(hp)
}

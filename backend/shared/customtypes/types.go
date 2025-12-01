package customtypes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/speps/go-hashids/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Validator interface {
	Validate() error
}

var ErrValidation error = errors.New("type validation error")

// ------------------------------------------------------------
// EncryptedId
// ------------------------------------------------------------

// Encrytped id is not nullable. Validation always returns error if not len>0.
type EncryptedId int64

var salt string = os.Getenv("ENC_KEY")

var hd = func() *hashids.HashID {
	h := hashids.NewData()
	h.Salt = salt
	h.MinLength = 12
	encoder, _ := hashids.NewWithData(h)
	return encoder
}()

func (e EncryptedId) MarshalJSON() ([]byte, error) {
	hash, err := hd.EncodeInt64([]int64{int64(e)})
	if err != nil {
		return nil, err
	}
	return json.Marshal(hash)
}

func (e *EncryptedId) UnmarshalJSON(data []byte) error {
	var hash string
	if err := json.Unmarshal(data, &hash); err != nil {
		return err
	}

	decoded, err := hd.DecodeInt64WithError(hash)
	if err != nil || len(decoded) == 0 {
		return err
	}

	*e = EncryptedId(decoded[0])
	return nil
}

func (e EncryptedId) IsValid() bool {
	return e > 0
}

func (e EncryptedId) Validate() error {
	if !e.IsValid() {
		return errors.Join(ErrValidation, errors.New("encryptedId must be positive"))
	}
	return nil
}

// ------------------------------------------------------------
// Id
// ------------------------------------------------------------

// Id is not nullable. Validation always returns error if not len>0.
type Id int64

func (i Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(i))
}

func (i *Id) UnmarshalJSON(data []byte) error {
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*i = Id(v)
	return nil
}

func (i Id) IsValid() bool {
	return i > 0
}

func (i Id) Validate() error {
	if !i.IsValid() {
		return errors.Join(ErrValidation, errors.New("id must be positive"))
	}
	return nil
}

// ------------------------------------------------------------
// Name
// ------------------------------------------------------------

// General type for names and surnames. Name is a nullable field. If needs to be present tag with `validate:"required"`
type Name string

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
	return len(n) >= 0
}

func (n Name) Validate() error {
	if !n.IsValid() {
		return errors.Join(ErrValidation, errors.New("name must be at least 2 characters"))
	}
	return nil
}

// ------------------------------------------------------------
// Username
// ------------------------------------------------------------
// Username is a nullable field. If needs to be present tag with `validate:"required"`
type Username string

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

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

// ------------------------------------------------------------
// Email
// ------------------------------------------------------------

type Email string

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

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
	if e == "" {
		return true
	}
	return emailRegex.MatchString(string(e))
}

func (e Email) Validate() error {
	if !e.IsValid() {
		return errors.Join(ErrValidation, errors.New("invalid email format"))
	}
	return nil
}

// ------------------------------------------------------------
// Limit
// ------------------------------------------------------------

type Limit int32

var maxLimit = 500

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

// ------------------------------------------------------------
// Offset
// ------------------------------------------------------------

type Offset int32

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

// ------------------------------------------------------------
// Password
// (Hash on Unmarshal; store hashed value only)
// ------------------------------------------------------------

// Password is not nullable. The length is checked and error is returned during json unmarshall.
// Password unmarshaling also depends on "PASSWORD_SECRET" env variable to be present.
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

// ------------------------------------------------------------
// DateOfBirth
// ------------------------------------------------------------

// DateOfBirth is a nullable field. If needs to be present tag with `validate:"required"`
type DateOfBirth time.Time

const (
	dobLayout     = "2006-01-02" // JSON date format
	minAgeInYears = 13
)

func (d DateOfBirth) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return json.Marshal(t.Format(dobLayout))
}

func (d *DateOfBirth) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	t, err := time.Parse(dobLayout, s)
	if err != nil {
		return err
	}

	*d = DateOfBirth(t)
	return nil
}

func (d DateOfBirth) IsValid() bool {
	t := time.Time(d)
	// if t.IsZero() {
	// 	return false
	// }

	now := time.Now().UTC()

	// cannot be in the future
	if t.After(now) {
		return false
	}

	// must be minimum age
	age := now.Year() - t.Year()
	if now.YearDay() < t.YearDay() {
		age--
	}

	return age >= minAgeInYears
}

func (d DateOfBirth) Validate() error {
	if !d.IsValid() {
		return errors.Join(ErrValidation, errors.New("invalid date of birth"))
	}
	return nil
}

// Helper to get time.Time if needed
func (d DateOfBirth) Time() time.Time {
	return time.Time(d)
}

func (d DateOfBirth) ToProto() *timestamppb.Timestamp {
	return timestamppb.New(time.Time(d))
}

// ------------------------------------------------------------
// Identifier (username or email)
// ------------------------------------------------------------

// Represents user name or email. Identifier is a nullable field. If needs to be present tag with `validate:"required"`
type Identifier string

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

// ------------------------------------------------------------
// About
// ------------------------------------------------------------

// Can be used for bio or descritpion. About is a nullable field. If needs to be present tag with `validate:"required"`
type About string

var (
	aboutCharsMin = 3
	aboutCharsMax = 300
)

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

// ------------------------------------------------------------
// Title (group/chat title)
// ------------------------------------------------------------

// Refers to title of content not mr, mrs. Title is a nullable field. If needs to be present tag with `validate:"required"`
type Title string

var (
	minTitleChars = 1
	maxTitleChars = 50
)

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

// ValidateStruct iterates over exported struct fields and validates them.
// - If a field implements Validator, its Validate() method is called.
// - If a field has `validate:"required"` tag, zero values are flagged as errors.
// - Optional fields (no tag) are ignored if empty.
// Example:
//
//	type RegisterRequest struct {
//	    Username  customtypes.Username `json:"username,omitempty" validate:"required"`
//	    FirstName customtypes.Name     `json:"first_name,omitempty" validate:"required"`
//	    LastName  customtypes.Name     `json:"last_name,omitempty"` // optional
//	    About     customtypes.About    `json:"about,omitempty"`     // optional
//	    Email     customtypes.Email    `json:"email,omitempty" validate:"required"`
//	}
func ValidateStruct(s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()

	var allErrors []string

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !fieldVal.CanInterface() {
			continue
		}

		val := fieldVal.Interface()
		validator, ok := val.(Validator)

		required := fieldType.Tag.Get("validate") == "required"

		if required {
			if isZeroValue(fieldVal) {
				allErrors = append(allErrors, fmt.Sprintf("%s: required field missing", fieldType.Name))
				continue
			}
		}
		// Call Validate() if the field implements Validator
		if ok {
			if err := validator.Validate(); err != nil {
				allErrors = append(allErrors, fmt.Sprintf("%s: %v", fieldType.Name, err))
			}
		}
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("validation errors: %v", allErrors)
	}
	return nil
}

// isZeroValue returns true if the reflect.Value is its type's zero value
func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

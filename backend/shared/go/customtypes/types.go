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

// Encrypted id is nullable. If `validation:"nullable` tag is present zero values don't return error.
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
	return e >= 0
}

func (e EncryptedId) Validate() error {
	if !e.IsValid() {
		return errors.Join(ErrValidation, errors.New("encryptedId must be positive"))
	}
	return nil
}

func (e EncryptedId) Int64() int64 {
	return int64(e)
}

// ------------------------------------------------------------
// Id
// ------------------------------------------------------------

// Id is nullable. If `validation:"nullable"` tag is present zero values don't return error.
// Negative values return error.
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
	return i >= 0
}

func (i Id) Validate() error {
	if !i.IsValid() {
		return errors.Join(ErrValidation, errors.New("id must be positive"))
	}
	return nil
}

func (i Id) Int64() int64 {
	return int64(i)
}

// ------------------------------------------------------------
// Name
// ------------------------------------------------------------

// General type for names and surnames. Name type is not nullable. All smaller than len 2 values return error
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

// ------------------------------------------------------------
// Username
// ------------------------------------------------------------

// Username is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
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

func (u Username) String() string {
	return string(u)
}

// ------------------------------------------------------------
// Email
// ------------------------------------------------------------

// Not nullable.
// Error upon validation is returned if string doesn't match email format or is empty.
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
// Limit
// ------------------------------------------------------------

// Non zero type. Validation returns error if zero or above limit
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

func (l Limit) Int32() int32 {
	return int32(l)
}

// ------------------------------------------------------------
// Offset
// ------------------------------------------------------------

// Non negative type. Validation returns error if below zero or above limit
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

func (o Offset) Int32() int32 {
	return int32(o)
}

// ------------------------------------------------------------
// Password
// (Hash on Unmarshal; store hashed value only)
// ------------------------------------------------------------

// Password is not nullable. The length is checked and error is returned during json unmarshall.
// Password unmarshaling depends on "PASSWORD_SECRET" env variable to be present.
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

func (p Password) String() string {
	return string(p)
}

// ------------------------------------------------------------
// DateOfBirth
// ------------------------------------------------------------

// TODO: Make not nullable
// DateOfBirth is non nullable. If value is the zero time instant, January 1, year 1, 00:00:00 UTC validation returns error.
type DateOfBirth time.Time

const (
	dobLayout     = "2006-01-02" // JSON date format
	minAgeInYears = 13
	maxAgeInYears = 120
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
	if t.IsZero() {
		return false
	}

	now := time.Now().UTC()

	// cannot be in the future
	if t.After(now) {
		return false
	}

	// compute age
	age := now.Year() - t.Year()
	if now.YearDay() < t.YearDay() {
		age--
	}

	// must be at least minAge and not older than maxAge
	if age < minAgeInYears {
		return false
	}

	if age > maxAgeInYears {
		return false
	}

	return true
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

// Helper to parse time.Time value to proto *timestamppb.Timestamp
func (d DateOfBirth) ToProto() *timestamppb.Timestamp {
	return timestamppb.New(time.Time(d))
}

// ------------------------------------------------------------
// EventDate
// ------------------------------------------------------------

// TODO: Make not nullable
// DateOfBirth is non nullable. If value is the zero time instant, January 1, year 1, 00:00:00 UTC validation returns error.
type EventDate time.Time

const (
	edLayout       = "2006-01-02" // JSON date format
	maxMonthsAhead = 6
)

func (ed EventDate) MarshalJSON() ([]byte, error) {
	t := time.Time(ed)
	return json.Marshal(t.Format(edLayout))
}

func (ed *EventDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	t, err := time.Parse(edLayout, s)
	if err != nil {
		return err
	}

	*ed = EventDate(t)
	return nil
}

func (ed EventDate) IsValid() bool {
	t := time.Time(ed)
	if t.IsZero() {
		return false
	}

	now := time.Now().UTC()

	// Normalize to the same location and remove time-of-day if needed
	t = t.In(now.Location())

	// Must be today or later
	if t.Before(now) {
		return false
	}

	// Must not be more than N months ahead
	limit := now.AddDate(0, maxMonthsAhead, 0)
	if t.After(limit) {
		return false
	}

	return true

}

func (ed EventDate) Validate() error {
	if !ed.IsValid() {
		return errors.Join(ErrValidation, errors.New("invalid event date"))
	}
	return nil
}

// Helper to get time.Time if needed
func (ed EventDate) Time() time.Time {
	return time.Time(ed)
}

// Helper to parse time.Time value to proto *timestamppb.Timestamp
func (ed EventDate) ToProto() *timestamppb.Timestamp {
	return timestamppb.New(time.Time(ed))
}

// ------------------------------------------------------------
// Identifier (username or email)
// ------------------------------------------------------------

// Represents user name or email. Identifier is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
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

func (i Identifier) String() string {
	return string(i)
}

// ------------------------------------------------------------
// About
// ------------------------------------------------------------

// Can be used for bio or descritpion. About is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
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

func (a About) String() string {
	return string(a)
}

// ------------------------------------------------------------
// Search
// ------------------------------------------------------------

// SearchTerm represents a validated search query term. Not nullable value
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
	return re.MatchString(string(s))
}

// Validate returns a descriptive error if the value is invalid.
func (s SearchTerm) Validate() error {
	if len(s) < 2 {
		return errors.Join(
			ErrValidation,
			errors.New("search term must be at least 2 characters"),
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

// ------------------------------------------------------------
// Title (group/chat title)
// ------------------------------------------------------------

// Refers to title of content not mr, mrs. Title is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
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

func (t Title) String() string {
	return string(t)
}

// ------------------------------------------------------------
// PostBody
// ------------------------------------------------------------

// Can be used for post body. PostBody is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type PostBody string

var (
	postBodyCharsMin = 3
	postBodyCharsMax = 500
)

func (b PostBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(b))
}

func (b *PostBody) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*b = PostBody(s)
	return nil
}

func (b PostBody) IsValid() bool {
	if len(b) == 0 {
		return true
	}
	if len(b) < postBodyCharsMin || len(b) > postBodyCharsMax {
		return false
	}
	for _, r := range b {
		if r < 32 { // control characters
			return false
		}
	}
	return true
}

func (b PostBody) Validate() error {
	if !b.IsValid() {
		return errors.Join(ErrValidation,
			fmt.Errorf("post body must be %d–%d chars and contain no control characters",
				postBodyCharsMin,
				postBodyCharsMax,
			))
	}
	return nil
}

func (b PostBody) String() string {
	return string(b)
}

// ------------------------------------------------------------
// CommentBody
// ------------------------------------------------------------

// Can be used for comment body. CommentBody is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type CommentBody string

var (
	commentBodyCharsMin = 3
	commentBodyCharsMax = 400
)

func (c CommentBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}

func (c *CommentBody) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*c = CommentBody(s)
	return nil
}

func (c CommentBody) IsValid() bool {
	if len(c) == 0 {
		return true
	}
	if len(c) < commentBodyCharsMin || len(c) > commentBodyCharsMax {
		return false
	}
	for _, r := range c {
		if r < 32 { // control characters
			return false
		}
	}
	return true
}

func (c CommentBody) Validate() error {
	if !c.IsValid() {
		return errors.Join(ErrValidation,
			fmt.Errorf("comment body must be %d–%d chars and contain no control characters",
				commentBodyCharsMin,
				commentBodyCharsMax,
			))
	}
	return nil
}

func (c CommentBody) String() string {
	return string(c)
}

// ------------------------------------------------------------
// EventBody
// ------------------------------------------------------------

// Can be used for event body. EventBody is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type EventBody string

var (
	eventBodyCharsMin = 3
	eventBodyCharsMax = 400
)

func (eb EventBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(eb))
}

func (eb *EventBody) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*eb = EventBody(s)
	return nil
}

func (eb EventBody) IsValid() bool {
	if len(eb) == 0 {
		return true
	}
	if len(eb) < eventBodyCharsMin || len(eb) > eventBodyCharsMax {
		return false
	}
	for _, r := range eb {
		if r < 32 { // control characters
			return false
		}
	}
	return true
}

func (eb EventBody) Validate() error {
	if !eb.IsValid() {
		return errors.Join(ErrValidation,
			fmt.Errorf("event body must be %d–%d chars and contain no control characters",
				eventBodyCharsMin,
				eventBodyCharsMax,
			))
	}
	return nil
}

func (eb EventBody) String() string {
	return string(eb)
}

// ------------------------------------------------------------
// Audience
// ------------------------------------------------------------

// Can be used for post, comment, event body. Audience is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type Audience string

var permittedAudienceValues = []string{"everyone", "group", "followers", "selected"}

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
	for _, permittedValue := range permittedAudienceValues {
		if strings.EqualFold(au.String(), permittedValue) {
			return true
		}
	}
	return false
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

// ValidateStruct iterates over exported struct fields and validates them.
// - If a field implements Validator, its Validate() method is called.
// - If a field does not have `validate:"nullable"` tag, zero values are flagged as errors.
// - Nullable fields if empty return nil error.
// Example:
//
//	type RegisterRequest struct {
//	    Username  customtypes.Username `json:"username,omitempty" validate:"nullable"` // optional
//	    FirstName customtypes.Name     `json:"first_name,omitempty" validate:"nullable"` // optional
//	    LastName  customtypes.Name     `json:"last_name"` // required
//	    About     customtypes.About    `json:"about"`     // required
//	    Email     customtypes.Email    `json:"email,omitempty" validate:"nullable"` // optional
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

		nullable := fieldType.Tag.Get("validate") == "nullable"
		isPrimitive := fieldVal.Type().PkgPath() == "" // exclude primitives
		_, zeroOk := allowedZeroVal[fieldVal.Type().Name()]

		if !nullable && !isPrimitive && !zeroOk {
			if isZeroValue(fieldVal) {
				allErrors = append(allErrors, fmt.Sprintf("%s: required field missing", fieldType.Name))
				continue
			}
		}

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

// Excluded types from nul check
var allowedZeroVal = map[string]struct{}{
	"Offset": {},
}

// isZeroValue returns true if the reflect.Value is its type's zero value
func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

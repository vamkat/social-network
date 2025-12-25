package ct

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const genDateTimeLayout = time.RFC3339

// ------------------------------------------------------------
// DateOfBirth
// ------------------------------------------------------------

// DateOfBirth is non nullable. If value is the zero time instant, January 1, year 1, 00:00:00 UTC validation returns error.
// It is marshaled and unmarshaled in "2006-01-02" format.
type DateOfBirth time.Time

const dobLayout = "2006-01-02"

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
	if age < dobMinAgeInYears {
		return false
	}

	if age > dobMaxAgeInYears {
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
// If 'd' is the zero time instant returns nil
func (d DateOfBirth) ToProto() *timestamppb.Timestamp {
	if d.Time().IsZero() {
		return nil
	}
	return timestamppb.New(time.Time(d))
}

func ParseDateOfBirth(s string) (DateOfBirth, error) {
	if s == "" {
		return DateOfBirth{}, errors.New("date_of_birth is required")
	}

	t, err := time.Parse(dobLayout, s)
	if err != nil {
		return DateOfBirth{}, fmt.Errorf("invalid date_of_birth format: %w", err)
	}

	dob := DateOfBirth(t)
	if err := dob.Validate(); err != nil {
		return DateOfBirth{}, err
	}

	return dob, nil
}

// ------------------------------------------------------------
// EventDateTime
// ------------------------------------------------------------

// It formats a time.Time value to genDateTimeLayout format.
// It Umarshals to time.Time type but Marshals to time.RFC3339 format.
//
// Null values are not allowed. If value is the zero time instant, January 1, year 1, 00:00:00 UTC validation returns error.
//
// Usage convert to proto type '*timestamppb.Timestamp':
//
//	return &pb.Event{
//			EventDateTime: resp.CreatedAt.ToProto(),
//	}, nil
type EventDateTime time.Time

func (edt EventDateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(edt)
	return json.Marshal(t.UTC().Format(genDateTimeLayout))
}

func (edt *EventDateTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	t, err := time.Parse(genDateTimeLayout, s)
	if err != nil {
		return err
	}

	*edt = EventDateTime(t)
	return nil
}

func (edt EventDateTime) IsValid() bool {
	t := time.Time(edt)
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
	limit := now.AddDate(0, eventDateMaxMonthsAhead, 0)

	return !t.After(limit)
}

func (edt EventDateTime) Validate() error {
	if !edt.IsValid() {
		return errors.Join(ErrValidation, errors.New("invalid event date"))
	}
	return nil
}

// Helper to get time.Time if needed
func (edt EventDateTime) Time() time.Time {
	return time.Time(edt)
}

// Helper to parse time.Time value to proto *timestamppb.Timestamp
func (edt EventDateTime) ToProto() *timestamppb.Timestamp {
	if edt.Time().IsZero() {
		return nil
	}
	return timestamppb.New(time.Time(edt))
}

// ------------------------------------------------------------
// Generic Date Time
// ------------------------------------------------------------

// GenDateTime (Generic) allows null values.
// It Umarshals to time.Time type but Marshals to time.RFC3339 format.
//
// Usage convert to proto type '*timestamppb.Timestamp':
//
//	return &pb.Conversation{
//			CreatedAt: resp.CreatedAt.ToProto(),
//	}, nil
type GenDateTime time.Time

// Marshal to RFC3339
func (g GenDateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(g)
	return json.Marshal(t.UTC().Format(genDateTimeLayout))
}

// Unmarshal from RFC3339 string
func (g *GenDateTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == "" {
		*g = GenDateTime(time.Time{})
		return nil
	}

	t, err := time.Parse(genDateTimeLayout, s)
	if err != nil {
		return err
	}

	*g = GenDateTime(t.UTC())
	return nil
}

func (g GenDateTime) IsValid() bool {
	t := time.Time(g)
	return !t.IsZero()
}

func (g GenDateTime) Validate() error {
	if !g.IsValid() {
		return errors.Join(ErrValidation, errors.New("invalid event date"))
	}
	return nil
}

// Scan implements the sql.Scanner interface
func (g *GenDateTime) Scan(src any) error {
	if src == nil {
		*g = GenDateTime(time.Time{})
		return nil
	}

	switch t := src.(type) {
	case time.Time:
		*g = GenDateTime(t) // store exactly as DB returns it
		return nil
	case []byte:
		parsed, err := time.Parse(time.RFC3339Nano, string(t))
		if err != nil {
			return err
		}
		*g = GenDateTime(parsed)
		return nil
	case string:
		parsed, err := time.Parse(time.RFC3339Nano, t)
		if err != nil {
			return err
		}
		*g = GenDateTime(parsed)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into GenDateTime", src)
	}
}

// Value implements the driver.Valuer interface
func (g GenDateTime) Value() (driver.Value, error) {
	if !g.IsValid() {
		return nil, nil // SQL NULL for invalid timestamps
	}
	return time.Time(g), nil // store exactly as is
}

// Helper to get time.Time if needed
func (g GenDateTime) Time() time.Time {
	return time.Time(g)
}

// Helper to parse time.Time value to proto *timestamppb.Timestamp
func (g GenDateTime) ToProto() *timestamppb.Timestamp {
	if g.Time().IsZero() {
		return nil
	}
	return timestamppb.New(time.Time(g))
}

package customtypes

import (
	"encoding/json"
	"errors"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ------------------------------------------------------------
// DateOfBirth
// ------------------------------------------------------------

// TODO: Make not nullable
// DateOfBirth is non nullable. If value is the zero time instant, January 1, year 1, 00:00:00 UTC validation returns error.

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
func (d DateOfBirth) ToProto() *timestamppb.Timestamp {
	return timestamppb.New(time.Time(d))
}

// ------------------------------------------------------------
// EventDate
// ------------------------------------------------------------

// TODO: Make not nullable
// DateOfBirth is non nullable. If value is the zero time instant, January 1, year 1, 00:00:00 UTC validation returns error.

func (ed EventDate) MarshalJSON() ([]byte, error) {
	t := time.Time(ed)
	return json.Marshal(t.Format(eventDateLayout))
}

func (ed *EventDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	t, err := time.Parse(eventDateLayout, s)
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
	limit := now.AddDate(0, eventDateMaxMonthsAhead, 0)
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

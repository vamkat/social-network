package ct

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// ------------------------------------------------------------
// PostBody
// ------------------------------------------------------------

// Can be used for post body.
type PostBody string

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
		return false
	}
	if len(b) < postBodyCharsMin || len(b) > postBodyCharsMax {
		return false
	}
	return controlCharsFree(b.String())
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

// Can be used for comment body
type CommentBody string

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
		return false
	}
	if len(c) < commentBodyCharsMin || len(c) > commentBodyCharsMax {
		return false
	}

	return controlCharsFree(c.String())
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

// Can be used for event body.
type EventBody string

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
		return false
	}
	if len(eb) < eventBodyCharsMin || len(eb) > eventBodyCharsMax {
		return false
	}
	return controlCharsFree(eb.String())
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

type MsgBody string

func (m MsgBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(m))
}

func (m *MsgBody) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*m = MsgBody(s)
	return nil
}

func (m MsgBody) IsValid() bool {
	if len(m) == 0 {
		return false
	}
	if len(m) < commentBodyCharsMin || len(m) > commentBodyCharsMax {
		return false
	}

	return controlCharsFree(m.String())
}

func (m MsgBody) Validate() error {
	if !m.IsValid() {
		return errors.Join(ErrValidation,
			fmt.Errorf("message body must be %d–%d chars and contain no control characters",
				commentBodyCharsMin,
				commentBodyCharsMax,
			))
	}
	return nil
}

func (i *MsgBody) Scan(src any) error {
	if src == nil {
		// SQL NULL reached
		*i = "" // or whatever "invalid" means in your domain
		return nil
	}

	switch v := src.(type) {
	case string:
		*i = MsgBody(v)
		return nil

	case []byte:
		*i = MsgBody(string(v))
		return nil
	}

	return fmt.Errorf("cannot scan type %T into MsgBody", src)
}

func (i MsgBody) Value() (driver.Value, error) {
	if !i.IsValid() {
		return nil, nil
	}
	return i.String(), nil
}

func (m MsgBody) String() string {
	return string(m)
}

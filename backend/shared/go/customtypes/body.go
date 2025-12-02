package customtypes

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ------------------------------------------------------------
// PostBody
// ------------------------------------------------------------

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

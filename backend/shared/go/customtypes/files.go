package customtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// =======================
// FileVisibility
// =======================

// Describes file visibility
// private, public
type FileVisibility string

const (
	Private FileVisibility = "private"
	Public  FileVisibility = "public"
)

func (v FileVisibility) String() string {
	return string(v)
}

func (v FileVisibility) isValid() bool {
	switch v {
	case Private, Public:
		return true
	default:
		return false
	}
}

func (v FileVisibility) Validate() error {
	if !v.isValid() {
		return fmt.Errorf("invalid FileVisibility: %q", v)
	}
	return nil
}

func (v FileVisibility) MarshalJSON() ([]byte, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(string(v))
}

func (v *FileVisibility) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val := FileVisibility(s)
	if !val.isValid() {
		return fmt.Errorf("invalid FileVisibility: %q", s)
	}

	*v = val
	return nil
}

func (v *FileVisibility) Scan(src any) error {
	if src == nil {
		*v = ""
		return nil
	}

	switch s := src.(type) {
	case string:
		val := FileVisibility(s)
		if !val.isValid() {
			return fmt.Errorf("invalid FileVisibility: %q", s)
		}
		*v = val
		return nil
	case []byte:
		val := FileVisibility(string(s))
		if !val.isValid() {
			return fmt.Errorf("invalid FileVisibility: %q", s)
		}
		*v = val
		return nil
	default:
		return fmt.Errorf("cannot scan FileVisibility from %T", src)
	}
}

func (v FileVisibility) Value() (driver.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return string(v), nil
}

func (v FileVisibility) SetExp() time.Duration {
	switch v {
	case Private:
		return time.Duration(3 * time.Minute)
	case Public:
		return time.Duration(6 * time.Hour)
	}
	return time.Duration(0)
}

// =======================
// ImgVariant
// =======================

// Describes the type of image
// thumb, small, medium, large
type ImgVariant string

const (
	Thumbnail ImgVariant = "thumb"
	Small     ImgVariant = "small"
	Medium    ImgVariant = "medium"
	Large     ImgVariant = "large"
)

func (v ImgVariant) String() string {
	return string(v)
}

func (v ImgVariant) IsValid() bool {
	switch v {
	case Thumbnail, Small, Medium, Large:
		return true
	default:
		return false
	}
}

func (v ImgVariant) Validate() error {
	if !v.IsValid() {
		return fmt.Errorf("invalid ImgVariant: %q", v)
	}
	return nil
}

func (v ImgVariant) MarshalJSON() ([]byte, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(string(v))
}

func (v *ImgVariant) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val := ImgVariant(s)
	if !val.IsValid() {
		return fmt.Errorf("invalid ImgVariant: %q", s)
	}

	*v = val
	return nil
}

func (v *ImgVariant) Scan(src any) error {
	if src == nil {
		*v = ""
		return nil
	}

	switch s := src.(type) {
	case string:
		val := ImgVariant(s)
		if !val.IsValid() {
			return fmt.Errorf("invalid ImgVariant: %q", s)
		}
		*v = val
		return nil
	case []byte:
		val := ImgVariant(string(s))
		if !val.IsValid() {
			return fmt.Errorf("invalid ImgVariant: %q", s)
		}
		*v = val
		return nil
	default:
		return fmt.Errorf("cannot scan ImgVariant from %T", src)
	}
}

func (v ImgVariant) Value() (driver.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return string(v), nil
}

// =======================
// UploadStatus
// =======================

// Describes the file upload status:
// pending, complete, failed
type UploadStatus string

const (
	Pending  UploadStatus = "pending"
	Complete UploadStatus = "complete"
	Failed   UploadStatus = "failed"
)

func (v UploadStatus) String() string {
	return string(v)
}

func (v UploadStatus) isValid() bool {
	switch v {
	case Pending, Complete, Failed:
		return true
	default:
		return false
	}
}

func (v UploadStatus) Validate() error {
	if !v.isValid() {
		return fmt.Errorf("invalid ImgVariant: %q", v)
	}
	return nil
}

func (v UploadStatus) MarshalJSON() ([]byte, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(string(v))
}

func (v *UploadStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val := UploadStatus(s)
	if !val.isValid() {
		return fmt.Errorf("invalid ImgVariant: %q", s)
	}

	*v = val
	return nil
}

func (v *UploadStatus) Scan(src any) error {
	if src == nil {
		*v = ""
		return nil
	}

	switch s := src.(type) {
	case string:
		val := UploadStatus(s)
		if !val.isValid() {
			return fmt.Errorf("invalid ImgStatus: %q", s)
		}
		*v = val
		return nil
	case []byte:
		val := UploadStatus(string(s))
		if !val.isValid() {
			return fmt.Errorf("invalid ImgStatus: %q", s)
		}
		*v = val
		return nil
	default:
		return fmt.Errorf("cannot scan ImgStatus from %T", src)
	}
}

func (v UploadStatus) Value() (driver.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return string(v), nil
}

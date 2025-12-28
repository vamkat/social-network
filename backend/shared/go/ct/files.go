package ct

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

// Sets 3 minutes expiration for private and 6 hours exp for public
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

// Describes the type of file
// original, thumb, small, medium, large
type FileVariant string

const (
	ImgThumbnail FileVariant = "thumb"
	ImgSmall     FileVariant = "small"
	ImgMedium    FileVariant = "medium"
	ImgLarge     FileVariant = "large"
	Original     FileVariant = "original"
)

func (v FileVariant) String() string {
	return string(v)
}

func (v FileVariant) isValid() bool {
	switch v {
	case ImgThumbnail, ImgSmall, ImgMedium, ImgLarge, Original:
		return true
	default:
		return false
	}
}

func (v FileVariant) Validate() error {
	if !v.isValid() {
		return fmt.Errorf("invalid ImgVariant: %q", v)
	}
	return nil
}

func (v FileVariant) MarshalJSON() ([]byte, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(string(v))
}

func (v *FileVariant) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val := FileVariant(s)
	if !val.isValid() {
		return fmt.Errorf("invalid ImgVariant: %q", s)
	}

	*v = val
	return nil
}

func (v *FileVariant) Scan(src any) error {
	if src == nil {
		*v = ""
		return nil
	}

	switch s := src.(type) {
	case string:
		val := FileVariant(s)
		if !val.isValid() {
			return fmt.Errorf("invalid ImgVariant: %q", s)
		}
		*v = val
		return nil
	case []byte:
		val := FileVariant(string(s))
		if !val.isValid() {
			return fmt.Errorf("invalid ImgVariant: %q", s)
		}
		*v = val
		return nil
	default:
		return fmt.Errorf("cannot scan ImgVariant from %T", src)
	}
}

func (v FileVariant) Value() (driver.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return string(v), nil
}

// =======================
// UploadStatus
// =======================

// Describes the file upload status:
// pending, processing, complete, failed
type UploadStatus string

const (
	Pending    UploadStatus = "pending"
	Processing UploadStatus = "processing"
	Complete   UploadStatus = "complete"
	Failed     UploadStatus = "failed"
)

func (v UploadStatus) String() string {
	return string(v)
}

func (v UploadStatus) isValid() bool {
	switch v {
	case Pending, Complete, Failed, Processing:
		return true
	default:
		return false
	}
}

func (v UploadStatus) Validate() error {
	if !v.isValid() {
		return fmt.Errorf("invalid UploadStatus: %q", v)
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
		return fmt.Errorf("invalid UploadStatus: %q", s)
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
			return fmt.Errorf("invalid UploadStatus: %q", s)
		}
		*v = val
		return nil
	case []byte:
		val := UploadStatus(string(s))
		if !val.isValid() {
			return fmt.Errorf("invalid UploadStatus: %q", s)
		}
		*v = val
		return nil
	default:
		return fmt.Errorf("cannot scan UploadStatus from %T", src)
	}
}

func (v UploadStatus) Value() (driver.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return string(v), nil
}

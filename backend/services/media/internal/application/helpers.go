package application

import (
	"fmt"
	"social-network/shared/go/ct"
	"time"
)

// validateUploadRequest validates all inputs required to create an image upload.
// It ensures the request metadata is well-formed, allowed by configuration,
// and safe to process.
//
// Validation rules:
//   - Filename must be non-empty.
//   - MimeType must be provided and allowed by file service constraints.
//   - SizeBytes must be greater than zero.
//   - Visibility must be a valid enum value.
//   - Expiration must be between 1 minute and 24 hours.
//   - At least one file variant must be provided.
//   - Each variant must be valid.
//   - The ct.Original variant is not allowed, as it is created implicitly.
//
// Returns a descriptive error on validation failure, or nil if the request is valid.
func (m *MediaService) validateUploadRequest(
	req UploadImageReq,
	exp time.Duration,
	variants []ct.FileVariant,
) error {

	if req.Filename == "" {
		return fmt.Errorf("upload image: invalid filename %q", req.Filename)
	}

	if req.MimeType == "" {
		return fmt.Errorf("upload image: missing mime type")
	}

	if !m.Cfgs.FileService.FileConstraints.AllowedMIMEs[req.MimeType] {
		return fmt.Errorf("upload image: mime type %q not allowed", req.MimeType)
	}

	if req.SizeBytes < 1 || req.SizeBytes > m.Cfgs.FileService.FileConstraints.MaxImageUpload {
		return fmt.Errorf("upload image: invalid size %d", req.SizeBytes)
	}

	if !req.Visibility.IsValid() {
		return fmt.Errorf("upload image: invalid visibility %v", req.Visibility)
	}

	if exp < time.Minute || exp > 24*time.Hour {
		return fmt.Errorf("upload image: invalid expiration %v", exp)
	}

	for _, v := range variants {
		if !v.IsValid() {
			return fmt.Errorf("invalid variant %v", v)
		} else if v == ct.Original {
			return fmt.Errorf("cannot accept customtypes.Original as variant to create.")
		}
	}

	return nil
}

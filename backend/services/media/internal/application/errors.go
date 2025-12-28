package application

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/services/media/internal/db/dbservice"
	ct "social-network/shared/go/ct"
	"time"
)

var (
	ErrReqValidation = errors.New("request validation error")               // invalid arguments
	ErrNotValidated  = errors.New("file not yet validated")                 // means that validation is pending
	ErrFailed        = errors.New("file has permanently failed validation") // means that file validation has failed permanently
	ErrNotFound      = errors.New("not found")                              // Usually equivalent to sql.ErrNoRows
	ErrInternal      = errors.New("internal error")
)

type MediaError struct {
	Kind error  // Classification: ErrNotFound, ErrInternal, etc.
	Err  error  // Cause: wrapped original error.
	Msg  string // Context: Func, args etc.
}

func (e *MediaError) Error() string {
	switch {
	case e.Msg != "" && e.Err != nil:
		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Msg, e.Err)
	case e.Msg != "":
		return fmt.Sprintf("%s: %s", e.Kind, e.Msg)
	case e.Err != nil:
		return fmt.Sprintf("%s: %v", e.Kind, e.Err)
	default:
		return e.Kind.Error()
	}
}

func (e *MediaError) Unwrap() error {
	return e.Err
}

// TODO: Check cases of nil kind
func Wrap(kind error, err error, msg ...string) error {
	if err == nil {
		return nil
	}

	// If it's already a MediaError, just add context
	var me *MediaError
	if errors.As(err, &me) && kind == nil {
		if len(msg) > 0 {
			return &MediaError{
				Kind: me.Kind, // preserve classification
				Msg:  msg[0],
				Err:  err,
			}
		}
		return err
	}

	// Fresh classification
	e := &MediaError{
		Kind: kind,
		Err:  err,
	}
	if len(msg) > 0 {
		e.Msg = msg[0]
	}
	return e
}

// Maps a file status to application errors and returns error.
// Caller decides if adding extra info about the file
func validateFileStatus(fm dbservice.File) error {
	errMsg := fmt.Sprintf(
		"file id %v file name %v status %v",
		fm.Id,
		fm.Filename,
		fm.Status,
	)

	if fm.Status == ct.Complete {
		return nil
	}

	if err := fm.Status.Validate(); err != nil {
		return Wrap(ErrFailed, err, errMsg)
	}

	if fm.Status == ct.Failed {
		return Wrap(ErrFailed, nil, errMsg)
	}

	if fm.Status == ct.Pending || fm.Status == ct.Processing {
		return Wrap(ErrNotValidated, nil, errMsg)
	}

	return nil
}

// Maps an db error to application custom error types
// sql.NoRows == ErrNotFound all other errors ErrInternal.
func mapDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return Wrap(ErrNotFound, err)
	}

	return Wrap(ErrInternal, err)
}

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
		return fmt.Errorf("upload image: missing mime type for file %v", req.Filename)
	}

	if !m.Cfgs.FileService.FileConstraints.AllowedMIMEs[req.MimeType] {
		return fmt.Errorf("upload image: mime type %q not allowed  for file %v", req.MimeType, req.Filename)
	}

	if req.SizeBytes < 1 || req.SizeBytes > m.Cfgs.FileService.FileConstraints.MaxImageUpload {
		return fmt.Errorf("upload image: invalid size %d for file %v", req.SizeBytes, req.Filename)
	}

	if err := req.Visibility.Validate(); err != nil {
		return fmt.Errorf("upload image: invalid visibility %v for file %v", req.Visibility, req.Filename)
	}

	if exp < time.Minute || exp > 24*time.Hour {
		return fmt.Errorf("upload image: invalid expiration %v for file %v", exp, req.Filename)
	}

	for _, v := range variants {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("invalid variant %v for file %v", v, req.Filename)
		}
		if v == ct.Original {
			return fmt.Errorf("original is not a variant %v for file %v", v, req.Filename)
		}
	}

	return nil
}

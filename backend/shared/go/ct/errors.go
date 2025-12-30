package ct

import (
	"errors"
	"fmt"
)

// Error represents a custom error type that includes classification, cause, and context.
// It implements the error interface and supports error wrapping and classification.
type Error struct {
	kind error  // Classification: ErrNotFound, ErrInternal, etc. Enusured to never be nil // TODO change to use grpc error codes
	err  error  // Cause: wrapped original error. //TODO change to be original error, that started the chain
	msg  string // Context: Func, args etc. //TODO change so that error chain survives
	// TODO add public field, which if it exists we trust to send it to user, if not then we decide what we do
}

var ErrUnknownClass = errors.New("unknown classification error") // kind is nil

// Returns a string with all available fields of err
func (e *Error) Error() string {
	e.kind = preventNilKind(e.kind)
	switch {
	case e.msg != "" && e.err != nil:
		return fmt.Sprintf("%s: %s: %v", e.kind, e.msg, e.err)
	case e.msg != "":
		return fmt.Sprintf("%s: %s", e.kind, e.msg)
	case e.err != nil:
		return fmt.Sprintf("%s: %v", e.kind, e.err)
	default:
		return e.kind.Error()
	}
}

// Returns a string containing only the kind field of Error.
func (e *Error) Public() string {
	return preventNilKind(e.kind).Error()
}

func preventNilKind(k error) error {
	if k != nil {
		return k
	}
	return ErrUnknownClass
}

// Method for errors.Is parsing. Returns `MediaError.Kind`.
func (e *Error) Is(target error) bool {
	return e.kind == target
}

// Method for error.As parsing. Returns the `MediaError.Err`.
func (e *Error) Unwrap() error {
	return e.err
}

// Wrap creates a MediaError that classifies and optionally wraps an existing error.
//
// Usage:
//   - kind: the classification of the error (e.g., ErrFailed, ErrNotFound). If nil, ErrUnknownClass is used.
//   - err: the underlying error to wrap; if nil, Wrap returns nil.
//   - msg: optional context message describing where or why the error occurred.
//
// Behavior:
//   - If `err` is already a MediaError and `kind` is nil, it preserves the original Kind and optionally adds a new message.
//   - Otherwise, it creates a new MediaError with the specified Kind, Err, and message.
//   - The resulting MediaError supports errors.Is (matches Kind) and errors.As (type assertion) and preserves the wrapped cause.
//   - If kind is nil and the err is not media error or lacks kind then kind is set to ErrUnknownClass.
//
// It is recommended to only use nil kind if the underlying error is of type MediaError and its kind is not nil.
func Wrap(kind error, err error, msg ...string) error {
	if err == nil {
		return nil
	}

	// If it's already a MediaError, just add context
	var me *Error
	if errors.As(err, &me) && kind == nil {
		if len(msg) > 0 {
			return &Error{
				kind: me.kind, // preserve classification
				msg:  msg[0],
				err:  err,
			}
		}
		return err
	}

	if kind == nil {
		kind = ErrUnknownClass
	}

	e := &Error{
		kind: kind,
		err:  err,
	}
	if len(msg) > 0 {
		e.msg = msg[0]
	}
	return e
}

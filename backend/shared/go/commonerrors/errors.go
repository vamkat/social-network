package commonerrors

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error represents a custom error type that includes classification, cause, and context.
// It implements the error interface and supports error wrapping and classification.
type Error struct {
	Code      error  // Classification: ErrNotFound, ErrInternal, etc. Enusured to never be nil
	Err       error  // Cause: wrapped original error.
	Msg       string // Context: Func, args etc.
	PublicMsg string // A message that will be displayed to clients.
	Source    error  // The original most underlying error.
}

// Returns a string with all available fields of err
func (e *Error) Error() string {
	e.Code = preventNilKind(e.Code)
	switch {
	case e.Msg != "" && e.Err != nil:
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Msg, e.Err)
	case e.Msg != "":
		return fmt.Sprintf("%s: %s", e.Code, e.Msg)
	case e.Err != nil:
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	default:
		return e.Code.Error()
	}
}

func preventNilKind(k error) error {
	if k != nil {
		return k
	}
	return ErrUnknown
}

// Method for errors.Is parsing. Returns `MediaError.Kind`.
func (e *Error) Is(target error) bool {
	return e.Code == target
}

// Method for error.As parsing. Returns the `MediaError.Err`.
func (e *Error) Unwrap() error {
	return e.Err
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
// It is recommended to only use nil kind if the underlying error is of type Error and its kind is not nil.
func Wrap(kind error, err error, msg ...string) *Error {
	if err == nil {
		return nil
	}

	var ce *Error
	if errors.As(err, &ce) {
		// Wrapping an existing custom error
		e := &Error{
			Code:      ce.Code,
			Err:       err,
			Source:    ce.Source,    // retain original source
			PublicMsg: ce.PublicMsg, // retain public message by default
		}

		if kind != nil {
			e.Code = kind
		}
		if len(msg) > 0 {
			e.Msg = msg[0]
		}

		return e
	}

	if kind == nil {
		kind = ErrUnknown
	}

	e := &Error{
		Code:   kind,
		Err:    err,
		Source: err, // ORIGINAL ROOT ERROR
	}
	if len(msg) > 0 {
		e.Msg = msg[0]
	}
	return e
}

// Add a Public Message to be displayed on APIs and other public endpoints
//
// Usage:
//
//	 return Wrap(ErrUnauthorized, err, "token expired").
//		WithPublic("Authentication required")
func (e *Error) WithPublic(msg string) *Error {
	e.PublicMsg = msg
	return e
}

// Helper mapper from error to grpc code.
func ToGRPCCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	// TODO: Check this
	// Propagate gRPC status errors
	// if st, ok := status.FromError(err); ok {
	// 	return st.Code()
	// }

	// Handle context errors
	if errors.Is(err, context.DeadlineExceeded) {
		return codes.DeadlineExceeded
	}
	if errors.Is(err, context.Canceled) {
		return codes.Canceled
	}

	// Handle your domain error
	var e *Error
	if errors.As(err, &e) {
		if code, ok := errorToGRPC[e.Code]; ok {
			return code
		}
	}

	// 4. Fallback
	return codes.Unknown
}

func GRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	// TODO: Check this
	// Propagate gRPC status errors
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	// Handle context errors
	if errors.Is(err, context.DeadlineExceeded) {
		return status.Errorf(codes.DeadlineExceeded, "deadline exceeded")
	}
	if errors.Is(err, context.Canceled) {
		return status.Errorf(codes.Canceled, "request canceled")
	}

	// Handle domain error
	var e *Error
	if errors.As(err, &e) {
		if code, ok := errorToGRPC[e.Code]; ok {
			return status.Errorf(code, "service error: %v", e.PublicMsg)
		}
	}
	return status.Errorf(codes.Unknown, "unknown error")
}

package ct

// import (
// 	"errors"
// 	"fmt"
// )

// // Error represents a custom error type that includes classification, cause, and context.
// // It implements the error interface and supports error wrapping and classification.
// type Error struct {
// 	Kind      error  // Classification: ErrNotFound, ErrInternal, etc. Enusured to never be nil
// 	Err       error  // Cause: wrapped original error.
// 	Msg       string // Context: Func, args etc.
// 	PublicMsg string // A message that will be displayed to clients.
// 	Source    error  // The original most underlying error.
// }

// var ErrUnknownClass = errors.New("unknown classification error") // kind is nil

// // Returns a string with all available fields of err
// func (e *Error) Error() string {
// 	e.Kind = preventNilKind(e.Kind)
// 	switch {
// 	case e.Msg != "" && e.Err != nil:
// 		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Msg, e.Err)
// 	case e.Msg != "":
// 		return fmt.Sprintf("%s: %s", e.Kind, e.Msg)
// 	case e.Err != nil:
// 		return fmt.Sprintf("%s: %v", e.Kind, e.Err)
// 	default:
// 		return e.Kind.Error()
// 	}
// }

// // Returns a string containing only the kind field of Error.
// func (e *Error) Public() string {
// 	return preventNilKind(e.Kind).Error()
// }

// func preventNilKind(k error) error {
// 	if k != nil {
// 		return k
// 	}
// 	return ErrUnknownClass
// }

// // Method for errors.Is parsing. Returns `MediaError.Kind`.
// func (e *Error) Is(target error) bool {
// 	return e.Kind == target
// }

// // Method for error.As parsing. Returns the `MediaError.Err`.
// func (e *Error) Unwrap() error {
// 	return e.Err
// }

// // Wrap creates a MediaError that classifies and optionally wraps an existing error.
// //
// // Usage:
// //   - kind: the classification of the error (e.g., ErrFailed, ErrNotFound). If nil, ErrUnknownClass is used.
// //   - err: the underlying error to wrap; if nil, Wrap returns nil.
// //   - msg: optional context message describing where or why the error occurred.
// //
// // Behavior:
// //   - If `err` is already a MediaError and `kind` is nil, it preserves the original Kind and optionally adds a new message.
// //   - Otherwise, it creates a new MediaError with the specified Kind, Err, and message.
// //   - The resulting MediaError supports errors.Is (matches Kind) and errors.As (type assertion) and preserves the wrapped cause.
// //   - If kind is nil and the err is not media error or lacks kind then kind is set to ErrUnknownClass.
// //
// // It is recommended to only use nil kind if the underlying error is of type Error and its kind is not nil.
// func Wrap(kind error, err error, msg ...string) *Error {
// 	if err == nil {
// 		return nil
// 	}

// 	var ce *Error
// 	if errors.As(err, &ce) {
// 		// Wrapping an existing custom error
// 		e := &Error{
// 			Kind:      ce.Kind,
// 			Err:       err,
// 			Source:    ce.Source,    // retain original source
// 			PublicMsg: ce.PublicMsg, // retain public message by default
// 		}

// 		if kind != nil {
// 			e.Kind = kind
// 		}
// 		if len(msg) > 0 {
// 			e.Msg = msg[0]
// 		}

// 		return e
// 	}

// 	if kind == nil {
// 		kind = ErrUnknownClass
// 	}

// 	e := &Error{
// 		Kind:   kind,
// 		Err:    err,
// 		Source: err, // ORIGINAL ROOT ERROR
// 	}
// 	if len(msg) > 0 {
// 		e.Msg = msg[0]
// 	}
// 	return e
// }

// // Add a Public Message to be displayed on APIs and other public endpoints
// //
// // Usage:
// //
// //	 return Wrap(ErrUnauthorized, err, "token expired").
// //		WithPublic("Authentication required")
// func (e *Error) WithPublic(msg string) *Error {
// 	e.PublicMsg = msg
// 	return e
// }

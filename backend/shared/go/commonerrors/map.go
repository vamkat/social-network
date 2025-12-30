package commonerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
)

var (
	// ErrOK indicates successful completion.
	// This error should generally not be returned; use nil instead.
	ErrOK = errors.New("OK")

	// ErrCanceled indicates the operation was canceled, typically by the caller.
	// Use when the client explicitly cancels the request.
	ErrCanceled = errors.New("CANCELED")

	// ErrUnknown indicates an unknown error.
	// Use when no other error code is appropriate.
	ErrUnknown = errors.New("UNKNOWN")

	// ErrInvalidArgument indicates the client specified an invalid argument.
	// Use when arguments are malformed or invalid regardless of system state.
	ErrInvalidArgument = errors.New("INVALID_ARGUMENT")

	// ErrDeadlineExceeded indicates the operation timed out before completion.
	// Use when the deadline expires before the operation completes.
	ErrDeadlineExceeded = errors.New("DEADLINE_EXCEEDED")

	// ErrNotFound indicates a requested entity was not found.
	// Use when a resource does not exist.
	ErrNotFound = errors.New("NOT_FOUND")

	// ErrAlreadyExists indicates an attempt to create an entity that already exists.
	// Use for idempotent create operations.
	ErrAlreadyExists = errors.New("ALREADY_EXISTS")

	// ErrPermissionDenied indicates the caller does not have permission.
	// Use when the caller is authenticated but lacks authorization.
	ErrPermissionDenied = errors.New("PERMISSION_DENIED")

	// ErrResourceExhausted indicates resource limits have been exceeded.
	// Use for rate limits, quotas, or out-of-capacity errors.
	ErrResourceExhausted = errors.New("RESOURCE_EXHAUSTED")

	// ErrFailedPrecondition indicates the system is in a state
	// that prevents the operation from executing.
	// Use when the operation is rejected due to system state, not argument value.
	ErrFailedPrecondition = errors.New("FAILED_PRECONDITION")

	// ErrAborted indicates the operation was aborted, typically due to concurrency issues.
	// Use for transaction aborts or concurrency conflicts.
	ErrAborted = errors.New("ABORTED")

	// ErrOutOfRange indicates an operation was attempted past the valid range.
	// Use when values are outside allowable bounds but otherwise well-formed.
	ErrOutOfRange = errors.New("OUT_OF_RANGE")

	// ErrUnimplemented indicates the operation is not implemented or supported.
	// Use when the API method or functionality is not available.
	ErrUnimplemented = errors.New("UNIMPLEMENTED")

	// ErrInternal indicates an internal server error.
	// Use when invariants are broken or unexpected conditions occur.
	ErrInternal = errors.New("INTERNAL")

	// ErrUnavailable indicates the service is currently unavailable.
	// Use for transient failures where retrying may succeed.
	ErrUnavailable = errors.New("UNAVAILABLE")

	// ErrDataLoss indicates unrecoverable data loss or corruption.
	// Use when data integrity cannot be guaranteed.
	ErrDataLoss = errors.New("DATA_LOSS")

	// ErrUnauthenticated indicates the caller is not authenticated.
	// Use when authentication credentials are missing or invalid.
	ErrUnauthenticated = errors.New("UNAUTHENTICATED")
)

var errorToGRPC = map[error]codes.Code{
	ErrOK:                 codes.OK,
	ErrCanceled:           codes.Canceled,
	ErrUnknown:            codes.Unknown,
	ErrInvalidArgument:    codes.InvalidArgument,
	ErrDeadlineExceeded:   codes.DeadlineExceeded,
	ErrNotFound:           codes.NotFound,
	ErrAlreadyExists:      codes.AlreadyExists,
	ErrPermissionDenied:   codes.PermissionDenied,
	ErrResourceExhausted:  codes.ResourceExhausted,
	ErrFailedPrecondition: codes.FailedPrecondition,
	ErrAborted:            codes.Aborted,
	ErrOutOfRange:         codes.OutOfRange,
	ErrUnimplemented:      codes.Unimplemented,
	ErrInternal:           codes.Internal,
	ErrUnavailable:        codes.Unavailable,
	ErrDataLoss:           codes.DataLoss,
	ErrUnauthenticated:    codes.Unauthenticated,
}

package commonerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name:     "nil kind defaults to ErrUnknownClass",
			err:      &Error{Code: nil, Msg: "test"},
			expected: "unknown classification error: test",
		},
		{
			name:     "kind, msg, and wrapped error",
			err:      &Error{Code: ErrNotFound, Msg: "user not found", Err: errors.New("db error")},
			expected: "NOT_FOUND: user not found: db error",
		},
		{
			name:     "kind and msg only",
			err:      &Error{Code: ErrInternal, Msg: "internal error"},
			expected: "INTERNAL: internal error",
		},
		{
			name:     "kind and wrapped error only",
			err:      &Error{Code: ErrInvalidArgument, Err: errors.New("invalid input")},
			expected: "INVALID_ARGUMENT: invalid input",
		},
		{
			name:     "kind only",
			err:      &Error{Code: ErrPermissionDenied},
			expected: "PERMISSION_DENIED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestError_Is(t *testing.T) {
	err := &Error{Code: ErrNotFound}

	assert.True(t, errors.Is(err, ErrNotFound))
	assert.False(t, errors.Is(err, ErrInternal))
}

func TestError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := &Error{Err: underlying}

	assert.Equal(t, underlying, errors.Unwrap(err))
}

func TestWrap(t *testing.T) {
	underlying := errors.New("underlying")

	t.Run("wrap nil returns nil", func(t *testing.T) {
		assert.Nil(t, Wrap(ErrInternal, nil))
	})

	t.Run("wrap with kind and msg", func(t *testing.T) {
		err := Wrap(ErrNotFound, underlying, "not found")
		require.NotNil(t, err)
		assert.Equal(t, ErrNotFound, err.Code)
		assert.Equal(t, underlying, err.Err)
		assert.Equal(t, underlying, err.Source)
		assert.Equal(t, "not found", err.Msg)
		assert.Equal(t, "", err.PublicMsg)
	})

	t.Run("wrap existing Error with new kind and msg", func(t *testing.T) {
		existing := &Error{Code: ErrInternal, Err: underlying, Source: underlying, PublicMsg: "public"}
		wrapped := Wrap(ErrNotFound, existing, "updated")
		require.NotNil(t, wrapped)
		assert.Equal(t, ErrNotFound, wrapped.Code)
		assert.Equal(t, existing, wrapped.Err)
		assert.Equal(t, underlying, wrapped.Source)
		assert.Equal(t, "updated", wrapped.Msg)
		assert.Equal(t, "public", wrapped.PublicMsg)
	})

	t.Run("wrap existing Error preserving kind when new kind nil", func(t *testing.T) {
		existing := &Error{Code: ErrInternal, Err: underlying, Source: underlying}
		wrapped := Wrap(nil, existing, "preserved")
		require.NotNil(t, wrapped)
		assert.Equal(t, ErrInternal, wrapped.Code)
		assert.Equal(t, existing, wrapped.Err)
		assert.Equal(t, underlying, wrapped.Source)
		assert.Equal(t, "preserved", wrapped.Msg)
	})

	t.Run("wrap with nil kind defaults to ErrUnknownClass", func(t *testing.T) {
		err := Wrap(nil, underlying)
		require.NotNil(t, err)
		assert.Equal(t, ErrUnknown, err.Code)
	})
}

func TestWithPublic(t *testing.T) {
	err := &Error{Code: ErrInternal}
	result := err.WithPublic("public message")

	assert.Equal(t, err, result) // returns same instance
	assert.Equal(t, "public message", err.PublicMsg)
}

func TestToGRPCCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected codes.Code
	}{
		{name: "nil error", err: nil, expected: codes.OK},
		{name: "ErrOK", err: &Error{Code: ErrOK}, expected: codes.OK},
		{name: "ErrNotFound", err: &Error{Code: ErrNotFound}, expected: codes.NotFound},
		{name: "ErrInternal", err: &Error{Code: ErrInternal}, expected: codes.Internal},
		{name: "ErrInvalidArgument", err: &Error{Code: ErrInvalidArgument}, expected: codes.InvalidArgument},
		{name: "ErrAlreadyExists", err: &Error{Code: ErrAlreadyExists}, expected: codes.AlreadyExists},
		{name: "ErrPermissionDenied", err: &Error{Code: ErrPermissionDenied}, expected: codes.PermissionDenied},
		{name: "ErrUnauthenticated", err: &Error{Code: ErrUnauthenticated}, expected: codes.Unauthenticated},
		{name: "unknown kind", err: &Error{Code: errors.New("custom")}, expected: codes.Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ToGRPCCode(tt.err))
		})
	}
}

func TestIntegration(t *testing.T) {
	underlying := errors.New("db connection failed")
	customErr := Wrap(ErrInternal, underlying, "failed to connect")

	// Test errors.Is works
	assert.True(t, errors.Is(customErr, ErrInternal))
	assert.False(t, errors.Is(customErr, ErrNotFound))

	// Test errors.As works
	var target *Error
	assert.True(t, errors.As(customErr, &target))
	assert.Equal(t, ErrInternal, target.Code)

	// Test unwrapping chain
	assert.Equal(t, underlying, errors.Unwrap(customErr))
}

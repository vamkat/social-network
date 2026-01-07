package commonerrors

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name:     "nil kind defaults to ErrUnknownClass",
			err:      &Error{class: nil, input: "test"},
			expected: ": test",
		},
		{
			name:     "kind, msg, and wrapped error",
			err:      &Error{class: ErrNotFound, input: "user not found", err: errors.New("db error")},
			expected: "not found: user not found: db error",
		},
		{
			name:     "kind and msg only",
			err:      &Error{class: ErrInternal, input: "internal error"},
			expected: "internal error: internal error",
		},
		{
			name:     "kind and wrapped error only",
			err:      &Error{class: ErrInvalidArgument, err: errors.New("invalid input")},
			expected: "invalid argument: invalid input",
		},
		{
			name:     "kind only",
			err:      &Error{class: ErrPermissionDenied},
			expected: "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestError_Is(t *testing.T) {
	err := &Error{class: ErrNotFound}

	assert.True(t, errors.Is(err, ErrNotFound))
	assert.False(t, errors.Is(err, ErrInternal))
}

func TestError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := &Error{err: underlying}

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
		assert.Equal(t, ErrNotFound, err.class)
		assert.Equal(t, underlying, err.err)
		assert.Equal(t, "not found", err.input)
		assert.Equal(t, "", err.publicMsg)
	})

	t.Run("wrap existing Error with new kind and msg", func(t *testing.T) {
		existing := &Error{class: ErrInternal, err: underlying, publicMsg: "public"}
		wrapped := Wrap(ErrNotFound, existing, "updated")
		require.NotNil(t, wrapped)
		assert.Equal(t, ErrNotFound, wrapped.class)
		assert.Equal(t, existing, wrapped.err)
		assert.Equal(t, "updated", wrapped.input)
		assert.Equal(t, "public", wrapped.publicMsg)
	})

	t.Run("wrap existing Error preserving kind when new kind nil", func(t *testing.T) {
		existing := &Error{class: ErrInternal, err: underlying}
		wrapped := Wrap(nil, existing, "preserved")
		require.NotNil(t, wrapped)
		assert.Equal(t, ErrInternal, wrapped.class)
		assert.Equal(t, existing, wrapped.err)
		assert.Equal(t, "preserved", wrapped.input)
	})

	t.Run("wrap with nil kind defaults to ErrUnknownClass", func(t *testing.T) {
		err := Wrap(nil, underlying)
		require.NotNil(t, err)
		assert.Equal(t, ErrUnknown, err.class)
	})
}

func TestWithPublic(t *testing.T) {
	err := &Error{class: ErrInternal}
	result := err.WithPublic("public message")

	assert.Equal(t, err, result) // returns same instance
	assert.Equal(t, "public message", err.publicMsg)
}

func TestToGRPCCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected codes.Code
	}{
		{name: "nil error", err: nil, expected: codes.OK},
		{name: "ErrOK", err: &Error{class: ErrOK}, expected: codes.OK},
		{name: "ErrNotFound", err: &Error{class: ErrNotFound}, expected: codes.NotFound},
		{name: "ErrInternal", err: &Error{class: ErrInternal}, expected: codes.Internal},
		{name: "ErrInvalidArgument", err: &Error{class: ErrInvalidArgument}, expected: codes.InvalidArgument},
		{name: "ErrAlreadyExists", err: &Error{class: ErrAlreadyExists}, expected: codes.AlreadyExists},
		{name: "ErrPermissionDenied", err: &Error{class: ErrPermissionDenied}, expected: codes.PermissionDenied},
		{name: "ErrUnauthenticated", err: &Error{class: ErrUnauthenticated}, expected: codes.Unauthenticated},
		{name: "unknown kind", err: &Error{class: errors.New("custom")}, expected: codes.Unknown},
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
	assert.Equal(t, ErrInternal, target.class)

	// Test unwrapping chain
	assert.Equal(t, underlying, errors.Unwrap(customErr))
}

func TestMultiLayerWrap_ErrorString(t *testing.T) {
	root := root()
	e1 := err1(root)
	e2 := Err2(e1)
	e3 := err3(e2)

	out := e3.Error()
	// t.Fatal(out)

	assert.Contains(t, out, "level 3")
	assert.Contains(t, out, "level 2")
	assert.Contains(t, out, "level 1")
	assert.Contains(t, out, "sql: no rows")

	// stack appears exactly once (entire block)
	assert.Equal(t, 1, strings.Count(out, e1.(*Error).stack))
}

func root() error {
	return errors.New("sql: no rows")
}
func err1(err error) error {
	return New(ErrNotFound, err, "level 1")
}
func Err2(err error) error {
	return Wrap(nil, err, "level 2")
}
func err3(err error) error {
	return Wrap(ErrInternal, err, "level 3")
}

func TestAs_ReturnsOutermostError(t *testing.T) {
	root := errors.New("io failure")
	e1 := New(ErrUnavailable, root, "dial")
	e2 := Wrap(nil, e1, "retry")

	var ce *Error
	require.True(t, errors.As(e2, &ce))
	assert.Equal(t, ErrUnavailable, ce.class)
	assert.Equal(t, e2, ce)
}

func TestGetSource_MultiWrap(t *testing.T) {
	root := errors.New("disk full")
	err := Wrap(
		ErrInternal,
		Wrap(nil,
			New(ErrUnavailable, root, "storage"),
			"service",
		),
		"handler",
	)

	assert.Equal(t, "disk full", GetSource(err))
}

func TestGRPCStatus_DoesNotMutateError(t *testing.T) {
	root := errors.New("token expired")
	err := Wrap(ErrUnauthenticated, root, "auth").
		WithPublic("authentication required")

	_ = GRPCStatus(err)

	// publicMsg must remain unchanged
	var ce *Error
	require.True(t, errors.As(err, &ce))
	assert.Equal(t, "authentication required", ce.publicMsg)
}

// Ensure most outer code prevails over nested codes
func TestGRPCStatus_MultipleWrapedCodes(t *testing.T) {
	root := New(ErrUnknown, errors.New("token expired"))
	err := Wrap(ErrUnauthenticated, root, "auth")

	out := GRPCStatus(err)
	require.NotNil(t, out)

	st, ok := status.FromError(out)
	require.True(t, ok, "expected gRPC status error")

	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestGRPCStatus_DefaultPublicMessage(t *testing.T) {
	err := Wrap(ErrInternal, errors.New("panic"), "handler")

	st, ok := status.FromError(GRPCStatus(err))
	require.True(t, ok)

	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "missing error message")
}

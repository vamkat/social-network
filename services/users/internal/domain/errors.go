package userservice

import "errors"

var (
	ErrUsernameConflict  = errors.New("username already exists")
	ErrEmailConflict     = errors.New("email already exists")
	ErrProfilePrivate    = errors.New("no permission to view profile")
	ErrNotAuthorized     = errors.New("user not authorized")
	ErrInvalidDateFormat = errors.New("invalid date format: expected YYYY-MM-DD")
	ErrWrongCredentials  = errors.New("invalid identifier or password")
)

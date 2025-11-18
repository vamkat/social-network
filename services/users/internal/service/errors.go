package userservice

import "errors"

var (
	ErrUsernameConflict = errors.New("username already exists")
	ErrEmailConflict    = errors.New("email already exists")
	ErrProfilePrivate   = errors.New("no permission to view profile")
	ErrNotAuthorized    = errors.New("user not authorized")
)

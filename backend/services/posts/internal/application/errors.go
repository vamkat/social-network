package application

import "errors"

var (
	ErrNoAudienceSelected = errors.New("no post audience selected for private post")
	ErrNotFound           = errors.New("no action was performed because no entity exists fitting the given criteria")
	ErrNoGroupIdGiven     = errors.New("group id required but not provided")
	ErrNotAllowed         = errors.New("user is not allowed to see this content")
)

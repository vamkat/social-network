package application

import "errors"

var (
	ErrNoAudienceSelected = errors.New("no post audience selected for private post")
)

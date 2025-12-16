package customtypes

import (
	"errors"
	"regexp"
)

type Validator interface {
	Validate() error
}

type Configs struct {
	PassSecret string
}

var Cfgs Configs

const (
	aboutCharsMin           = 3
	aboutCharsMax           = 300
	postBodyCharsMin        = 3
	postBodyCharsMax        = 500
	commentBodyCharsMin     = 3
	commentBodyCharsMax     = 400
	eventBodyCharsMin       = 3
	eventBodyCharsMax       = 400
	dobMinAgeInYears        = 13
	dobMaxAgeInYears        = 120
	eventDateMaxMonthsAhead = 6
	maxLimit                = 500
	minTitleChars           = 1
	maxTitleChars           = 50
)

var permittedAudienceValues = []string{"everyone", "group", "followers", "selected"}
var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)
var nameRegex = regexp.MustCompile(`^\p{L}+([\p{L}'\- ]*\p{L})?$`)

// Excluded types from nul check
var alwaysAllowZero = map[string]struct{}{
	"Offset": {},
}

var ErrValidation error = errors.New("type validation error")

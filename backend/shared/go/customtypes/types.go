package customtypes

import (
	"errors"
	"regexp"
	"time"
)

type Validator interface {
	Validate() error
}

type About string
type Audience string
type PostBody string
type CommentBody string
type EventBody string
type DateOfBirth time.Time
type EventDate time.Time
type EncryptedId int64
type Id int64
type Email string
type Username string
type Identifier string
type Name string
type Limit int32
type Offset int32
type Password string
type SearchTerm string
type Title string

const (
	aboutCharsMin           = 3
	aboutCharsMax           = 300
	postBodyCharsMin        = 3
	postBodyCharsMax        = 500
	commentBodyCharsMin     = 3
	commentBodyCharsMax     = 400
	eventBodyCharsMin       = 3
	eventBodyCharsMax       = 400
	dobLayout               = "2006-01-02" // JSON date format
	eventDateLayout         = "2006-01-02" // JSON date format
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

// Excluded types from nul check
var allowedZeroVal = map[string]struct{}{
	"Offset": {},
}

var ErrValidation error = errors.New("type validation error")

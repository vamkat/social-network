package customtypes

import (
	"errors"
	"regexp"
	"time"
)

type Validator interface {
	Validate() error
}

// new types
type About string          // Can be used for bio or descritpion. About is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type Audience string       // Can be used for post, comment, event body. Audience is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type PostBody string       // Can be used for post body. PostBody is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type CommentBody string    // Can be used for comment body. CommentBody is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type EventBody string      // Can be used for event body. EventBody is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type DateOfBirth time.Time // DateOfBirth is non nullable. If value is the zero time instant, January 1, year 1, 00:00:00 UTC validation returns error.
type EventDate time.Time   // Non nullable. If value is the zero time instant, January 1, year 1, 00:00:00 UTC validation returns error.
type EncryptedId int64     // Encrypted id is nullable. If `validation:"nullable` tag is present zero values don't return error. When json decoded string is unhashed to int64. The reverse is applied on Json encode.
type Id int64              // Id is nullable. If `validation:"nullable"` tag is present zero values don't return error. Negative values return error.
type Ids []Id              // Ids is nullable. If `validation:"nullable"` tag is present zero values don't return error. Contents are allowed to be nullable with tag elements=nullable
type Email string          // Not nullable. Error upon validation is returned if string doesn't match email format or is empty.
type Username string       // Username is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type Identifier string     // Represents user name or email. Identifier is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.
type Name string           // General type for names and surnames. Name type is not nullable. All smaller than len 2 values return error
type Limit int32           // Non zero type. Validation returns error if zero or above limit
type Offset int32          // Non negative type. Validation returns error if below zero or above limit
type Password string       // Password is not nullable. The length is checked and error is returned during json unmarshall. Password unmarshaling depends on "PASSWORD_SECRET" env variable to be present.
type SearchTerm string     // SearchTerm represents a validated search query term. Not nullable value
type Title string          // Refers to title of content not mr, mrs. Title is a nullable field. If `validation:"nullable"` tag is present zero values don't return error.

const (
	aboutCharsMin           = 3
	aboutCharsMax           = 300
	postBodyCharsMin        = 3
	postBodyCharsMax        = 500
	commentBodyCharsMin     = 3
	commentBodyCharsMax     = 400
	eventBodyCharsMin       = 3
	eventBodyCharsMax       = 400
	dobLayout               = "2006-01-02"
	eventDateLayout         = "2006-01-02"
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

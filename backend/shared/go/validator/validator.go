package validator

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

// parse a regexp pattern for sanity checking the format of an email address
// local part: alphanumeric plus all RFC5322 allowed special characters
// domain part: starts and ends with alphanumeric, allows hyphens, limit 63 chars
// requires at least one dot in domain
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$")

// regex pattern for username allows only alphanumeric and .-_
var UsernameRX = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// returns true if the fielderrors map is empty
func (v *Validator) Valid() bool {
	return (len(v.FieldErrors)) == 0 && len(v.NonFieldErrors) == 0
}

// adds an error message to the map
func (v *Validator) AddFieldError(key, message string) {
	//if map not initialized, we need to initialize it
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}

}

// adds error message to the map if validation check is not "ok"
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// returns true if value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// returns true if email is valid according to RFC5322 (more lax than our regex)
func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// returns true if value containts no more than n chars
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// returns true if value is in list of permitted ints
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// returns true if value is in list of permitted strings
func PermittedString(value string, permittedValues []string) bool {
	for _, permittedValue := range permittedValues {
		if strings.EqualFold(value, permittedValue) {
			return true
		}
	}
	return false
}

// returns true if slice of strings is not empty
func SliceNotEmpty(slice []string) bool {
	return len(slice) > 0
}

// returns true if value contains at least n chars
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// returns true if value matches provided compiled regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

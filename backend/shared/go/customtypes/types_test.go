package customtypes_test

import (
	"encoding/json"
	"os"
	"reflect"
	"social-network/shared/go/customtypes"

	"strings"
	"testing"
	"time"
)

// Utility: mustSetEnv
func mustSetEnv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("failed to set env %s: %v", key, err)
	}
}

// ------------------------------------------------------------
// Id
// ------------------------------------------------------------
func TestIdJSON(t *testing.T) {
	mustSetEnv(t, "ENC_KEY", "test-salt")

	id := customtypes.Id(123)
	b, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded customtypes.Id
	err = json.Unmarshal(b, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded != id {
		t.Fatalf("expected %d, got %d", id, decoded)
	}
}

func TestIdValidate(t *testing.T) {
	if err := customtypes.Id(-5).Validate(); err == nil {
		t.Fatal("expected validation error for negative Id")
	}
}

// ------------------------------------------------------------
// Id
// ------------------------------------------------------------
func TestIdValidation(t *testing.T) {
	if err := customtypes.Id(-1).Validate(); err == nil {
		t.Fatal("expected error for invalid id")
	}
	if err := customtypes.Id(5).Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

// ------------------------------------------------------------
// Name
// ------------------------------------------------------------
func TestNameValidation(t *testing.T) {
	if err := customtypes.Name("A").Validate(); err == nil {
		t.Fatal("expected name length error")
	}
	if err := customtypes.Name("John").Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

// ------------------------------------------------------------
// Username
// ------------------------------------------------------------
func TestUsernameValidation(t *testing.T) {
	if err := customtypes.Username("ab").Validate(); err == nil {
		t.Fatal("should fail: too short")
	}
	if err := customtypes.Username("valid_user123").Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

// ------------------------------------------------------------
// Email
// ------------------------------------------------------------
func TestEmailValidation(t *testing.T) {
	if err := customtypes.Email("not-an-email").Validate(); err == nil {
		t.Fatal("expected invalid email error")
	}
	if err := customtypes.Email("test@example.com").Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

// ------------------------------------------------------------
// Limit
// ------------------------------------------------------------
func TestLimitValidation(t *testing.T) {
	if err := customtypes.Limit(0).Validate(); err == nil {
		t.Fatal("expected error")
	}
	if err := customtypes.Limit(500).Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if err := customtypes.Limit(501).Validate(); err == nil {
		t.Fatal("expected upper bound error")
	}
}

// ------------------------------------------------------------
// Offset
// ------------------------------------------------------------
func TestOffsetValidation(t *testing.T) {
	if err := customtypes.Offset(-1).Validate(); err == nil {
		t.Fatal("expected error for negative offset")
	}
	if err := customtypes.Offset(10).Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

// ------------------------------------------------------------
// Password
// ------------------------------------------------------------
func TestPasswordJSON(t *testing.T) {
	mustSetEnv(t, "PASSWORD_SECRET", "supersecret")

	// raw password JSON
	body := []byte(`"mySecretPass"`)

	var p customtypes.Password
	if err := json.Unmarshal(body, &p); err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	// marshal must return "********"
	out, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(out) != `"********"` {
		t.Fatalf("expected masked password, got %s", out)
	}
}

func TestPasswordValidation(t *testing.T) {
	mustSetEnv(t, "PASSWORD_SECRET", "supersecret")

	var p customtypes.Password
	_ = json.Unmarshal([]byte(`"Password!123"`), &p)

	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

// ------------------------------------------------------------
// DateOfBirth
// ------------------------------------------------------------
func TestDOBValidation(t *testing.T) {
	now := time.Now().UTC()
	under13 := now.AddDate(-10, 0, 0)
	valid := now.AddDate(-20, 0, 0)
	future := now.AddDate(1, 0, 0)

	d := customtypes.DateOfBirth(valid)
	if err := d.Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	d = customtypes.DateOfBirth(under13)
	if err := d.Validate(); err == nil {
		t.Fatal("expected min-age error")
	}

	d = customtypes.DateOfBirth(future)
	if err := d.Validate(); err == nil {
		t.Fatal("expected future-date error")
	}
}

// ------------------------------------------------------------
// Identifier
// ------------------------------------------------------------
func TestIdentifierValidation(t *testing.T) {
	if err := customtypes.Identifier("bad@format@x").Validate(); err == nil {
		t.Fatal("expected invalid identifier")
	}
	if err := customtypes.Identifier("validUser_123").Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if err := customtypes.Identifier("email@test.com").Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

// ------------------------------------------------------------
// About
// ------------------------------------------------------------
func TestAboutValidation(t *testing.T) {
	if err := customtypes.About("ok!").Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if err := customtypes.About("\x01bad").Validate(); err == nil {
		t.Fatal("expected control char error")
	}
	if err := customtypes.About("ab").Validate(); err == nil {
		t.Fatal("expected min length error")
	}
}

// ------------------------------------------------------------
// Title
// ------------------------------------------------------------
func TestTitleValidation(t *testing.T) {
	if err := customtypes.Title("A title").Validate(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if err := customtypes.Title(" ").Validate(); err == nil {
		t.Fatal("expected trimmed length error")
	}
}

// ------------------------------------------------------------
// ValidateStruct
// ------------------------------------------------------------
func TestValidateStruct(t *testing.T) {
	type TestReq struct {
		Name     customtypes.Name     `validate:"nullable"`
		Email    customtypes.Email    // required
		About    customtypes.About    `validate:"nullable"`
		Username customtypes.Username `validate:"nullable"`
	}

	ok := TestReq{
		Name:     "John Doe",
		Email:    "valid@example.com",
		About:    "This is ok",
		Username: "user_1",
	}

	if err := customtypes.ValidateStruct(ok); err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	// Missing required Email
	bad := TestReq{}
	err := customtypes.ValidateStruct(bad)
	if err == nil {
		t.Fatal("expected missing required field error")
	}

	if !contains(err.Error(), "Email: required field missing") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// helpers
func contains(haystack, needle string) bool {
	return reflect.ValueOf(haystack).String() != "" &&
		len(haystack) >= len(needle) &&
		(len(needle) == 0 || (len(haystack) >= len(needle) && (index(haystack, needle) != -1)))
}

func index(s, sep string) int {
	return len([]rune(s[:])) - len([]rune(stringsAfter(s, sep)))
}

func stringsAfter(s, sep string) string {
	if sep == "" {
		return s
	}
	i := len([]rune(s)) - len([]rune(sep))
	if i < 0 {
		return s
	}
	return s[i:]
}

func TestValidateStruct_BoolAndOffsetExempt(t *testing.T) {
	type TestStruct struct {
		Flag   bool               `validate:""` // bool = false should NOT trigger required
		Number customtypes.Offset `validate:""` // Offset = 0 should NOT trigger required
		Name   customtypes.Name   `validate:""` // string = "" SHOULD trigger required
	}
	s := TestStruct{
		Flag:   false, // should NOT fail
		Number: 0,     // should NOT fail
		Name:   "",    // should fail
	}

	err := customtypes.ValidateStruct(s)
	if err == nil {
		t.Fatalf("expected validation error but got none")
	}

	// We expect ONLY Name to fail
	msg := err.Error()
	if !strings.Contains(msg, "Name: required field missing") {
		t.Fatalf("expected missing name error, got: %v", msg)
	}

	// Verify bool and Offset are NOT included in errors
	if strings.Contains(msg, "Flag") {
		t.Fatalf("bool=false should not produce error, got: %v", msg)
	}
	if strings.Contains(msg, "Number") {
		t.Fatalf("Offset=0 should not produce error, got: %v", msg)
	}
}

func TestValidateStruct_SliceOfCustomTypes(t *testing.T) {
	type TestStruct struct {
		// Slice of custom types - nullable
		NullableIDs customtypes.Ids `validate:"nullable"`

		// Required field to satisfy other validations
		Email customtypes.Email
	}

	type TestStructRequired struct {
		// Slice of custom types - not nullable
		RequiredIDs customtypes.Ids

		// Required field to satisfy other validations
		Email customtypes.Email
	}

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid required IDs",
			input: TestStructRequired{
				RequiredIDs: customtypes.Ids{1, 2},
				Email:       "test@example.com",
			},
			wantError: false,
		},
		{
			name: "nil required IDs - should fail",
			input: TestStructRequired{
				RequiredIDs: nil,
				Email:       "test@example.com",
			},
			wantError: true,
			errorMsg:  "RequiredIDs: required field missing",
		},
		{
			name: "empty required IDs - should fail",
			input: TestStructRequired{
				RequiredIDs: customtypes.Ids{},
				Email:       "test@example.com",
			},
			wantError: true,
			errorMsg:  "RequiredIDs: required field missing",
		},
		{
			name: "required IDs with zero element - should fail",
			input: TestStructRequired{
				RequiredIDs: customtypes.Ids{1, 0},
				Email:       "test@example.com",
			},
			wantError: true,
			errorMsg:  "RequiredIDs[1]: required element missing",
		},
		{
			name: "required IDs with negative element - should fail on validation",
			input: TestStructRequired{
				RequiredIDs: customtypes.Ids{1, -1},
				Email:       "test@example.com",
			},
			wantError: true,
			errorMsg:  "RequiredIDs[1]:",
		},
		{
			name: "nullable IDs with nil - should pass",
			input: TestStruct{
				NullableIDs: nil,
				Email:       "test@example.com",
			},
			wantError: false,
		},
		{
			name: "nullable IDs with empty slice - should pass",
			input: TestStruct{
				NullableIDs: customtypes.Ids{},
				Email:       "test@example.com",
			},
			wantError: false,
		},
		{
			name: "nullable IDs with zero element - should fail (elements still validated)",
			input: TestStruct{
				NullableIDs: customtypes.Ids{1, 0},
				Email:       "test@example.com",
			},
			wantError: true,
			errorMsg:  "NullableIDs[1]: required element missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := customtypes.ValidateStruct(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error on %v but got none", tt)
					return
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

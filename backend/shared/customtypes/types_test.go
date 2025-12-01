package customtypes

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncryptedIdJSONAndValidation(t *testing.T) {
	os.Setenv("ENC_KEY", "test_salt")

	e := EncryptedId(12345)
	data, err := json.Marshal(e)
	assert.NoError(t, err)

	var decoded EncryptedId
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, e, decoded)

	assert.NoError(t, e.Validate())

	invalid := EncryptedId(0)
	assert.Error(t, invalid.Validate())
}

func TestIdJSONAndValidation(t *testing.T) {
	i := Id(42)
	data, err := json.Marshal(i)
	assert.NoError(t, err)

	var decoded Id
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, i, decoded)

	assert.NoError(t, i.Validate())
	assert.Error(t, Id(0).Validate())
}

func TestNameJSONAndValidation(t *testing.T) {
	n := Name("Alice")
	data, err := json.Marshal(n)
	assert.NoError(t, err)

	var decoded Name
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, n, decoded)

	assert.NoError(t, n.Validate())
}

func TestUsernameValidation(t *testing.T) {
	valid := Username("user_123")
	assert.True(t, valid.IsValid())
	assert.NoError(t, valid.Validate())

	invalid := Username("ab") // too short
	assert.False(t, invalid.IsValid())
	assert.Error(t, invalid.Validate())
}

func TestEmailValidation(t *testing.T) {
	valid := Email("test@example.com")
	assert.True(t, valid.IsValid())
	assert.NoError(t, valid.Validate())

	invalid := Email("invalid-email")
	assert.False(t, invalid.IsValid())
	assert.Error(t, invalid.Validate())
}

func TestLimitValidation(t *testing.T) {
	assert.NoError(t, Limit(1).Validate())
	assert.NoError(t, Limit(500).Validate())
	assert.Error(t, Limit(0).Validate())
	assert.Error(t, Limit(501).Validate())
}

func TestOffsetValidation(t *testing.T) {
	assert.NoError(t, Offset(0).Validate())
	assert.NoError(t, Offset(100).Validate())
	assert.Error(t, Offset(-1).Validate())
}

func TestPasswordJSONAndValidation(t *testing.T) {
	os.Setenv("PASSWORD_SECRET", "secret_key")
	raw := "mypassword"
	var p Password
	data, _ := json.Marshal(raw)
	err := p.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.NotEqual(t, raw, string(p)) // Should be hashed
	assert.NoError(t, p.Validate())

	// Missing secret
	os.Unsetenv("PASSWORD_SECRET")
	var p2 Password
	err = p2.UnmarshalJSON(data)
	assert.Error(t, err)
}

func TestDateOfBirthValidation(t *testing.T) {
	now := time.Now()
	dob := DateOfBirth(now.AddDate(-15, 0, 0))
	assert.NoError(t, dob.Validate())

	tooYoung := DateOfBirth(now.AddDate(-10, 0, 0))
	assert.Error(t, tooYoung.Validate())

	future := DateOfBirth(now.AddDate(1, 0, 0))
	assert.Error(t, future.Validate())
}

func TestIdentifierValidation(t *testing.T) {
	assert.NoError(t, Identifier("user_123").Validate())
	assert.NoError(t, Identifier("test@example.com").Validate())
	assert.NoError(t, Identifier("").Validate())
	assert.Error(t, Identifier("invalid*id").Validate())
}

func TestAboutValidation(t *testing.T) {
	valid := About("This is a bio")
	assert.NoError(t, valid.Validate())

	tooShort := About("ab")
	assert.Error(t, tooShort.Validate())

	controlChar := About("Hello\x01World")
	assert.Error(t, controlChar.Validate())
}

func TestTitleValidation(t *testing.T) {
	valid := Title("Group Chat")
	assert.NoError(t, valid.Validate())

	tooShort := Title("")
	assert.NoError(t, tooShort.Validate()) // Nullable allowed

	tooLong := Title("This title is definitely way too long for validation purposes")
	assert.Error(t, tooLong.Validate())
}

func TestValidateStruct(t *testing.T) {
	type RegisterRequest struct {
		Username  Username `validate:"required"`
		FirstName Name     `validate:"required"`
		LastName  Name
		About     About
		Email     Email `validate:"required"`
	}

	req := RegisterRequest{
		Username:  "user_123",
		FirstName: "Alice",
		Email:     "test@example.com",
	}

	err := ValidateStruct(req)
	assert.NoError(t, err)

	invalidReq := RegisterRequest{
		Username:  "ab",
		FirstName: "",
		Email:     "bad-email",
	}
	err = ValidateStruct(invalidReq)
	assert.Error(t, err)
}

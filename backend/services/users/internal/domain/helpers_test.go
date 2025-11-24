package userservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {
	password := "mypassword123"

	hash, err := hashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestHashPassword_Consistency(t *testing.T) {
	password := "mypassword123"

	hash1, _ := hashPassword(password)
	hash2, _ := hashPassword(password)

	// Two hashes of the same password should be different (bcrypt uses random salt)
	assert.NotEqual(t, hash1, hash2)
}

func TestCheckPassword_Correct(t *testing.T) {
	password := "mypassword123"
	hash, _ := hashPassword(password)

	result := checkPassword(hash, password)

	assert.True(t, result)
}

func TestCheckPassword_Incorrect(t *testing.T) {
	password := "mypassword123"
	wrongPassword := "wrongpassword"
	hash, _ := hashPassword(password)

	result := checkPassword(hash, wrongPassword)

	assert.False(t, result)
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	password := "mypassword123"
	hash, _ := hashPassword(password)

	result := checkPassword(hash, "")

	assert.False(t, result)
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	result := checkPassword("invalid_hash", "password")

	assert.False(t, result)
}

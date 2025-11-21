package userservice

import (
	"context"
	"errors"
	"testing"

	"social-network/services/users/internal/db/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser_InvalidDateFormat(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	req := RegisterUserRequest{
		Username:    "testuser",
		FirstName:   "Test",
		LastName:    "User",
		DateOfBirth: "invalid-date",
		Avatar:      "avatar.jpg",
		About:       "Test about me",
		Public:      true,
		Email:       "test@example.com",
		Password:    "password123",
	}

	ctx := context.Background()
	_, err := service.RegisterUser(ctx, req)

	assert.Equal(t, ErrInvalidDateFormat, err)
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	t.Skip("Requires transaction support with real pool - use integration tests")
}

func TestUpdateUserPassword_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	oldPassword := "oldpassword"
	newPassword := "newpassword"
	hashedOld, _ := hashPassword(oldPassword)

	req := UpdatePasswordRequest{
		UserId:      1,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	ctx := context.Background()

	mockDB.On("GetUserPassword", ctx, int64(1)).Return(hashedOld, nil)
	mockDB.On("UpdateUserPassword", ctx, mock.MatchedBy(func(arg sqlc.UpdateUserPasswordParams) bool {
		return arg.UserID == 1
	})).Return(nil)

	err := service.UpdateUserPassword(ctx, req)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUpdateUserPassword_WrongOldPassword(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	correctPassword := "correctpassword"
	hashedCorrect, _ := hashPassword(correctPassword)

	req := UpdatePasswordRequest{
		UserId:      1,
		OldPassword: "wrongpassword",
		NewPassword: "newpassword",
	}

	ctx := context.Background()

	mockDB.On("GetUserPassword", ctx, int64(1)).Return(hashedCorrect, nil)

	err := service.UpdateUserPassword(ctx, req)

	assert.Equal(t, ErrNotAuthorized, err)
	mockDB.AssertExpectations(t)
}

func TestUpdateUserEmail_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	req := UpdateEmailRequest{
		UserId: 1,
		Email:  "newemail@example.com",
	}

	ctx := context.Background()

	mockDB.On("UpdateUserEmail", ctx, sqlc.UpdateUserEmailParams{
		UserID: 1,
		Email:  "newemail@example.com",
	}).Return(nil)

	err := service.UpdateUserEmail(ctx, req)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUpdateUserEmail_Error(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	req := UpdateEmailRequest{
		UserId: 1,
		Email:  "duplicate@example.com",
	}

	ctx := context.Background()

	expectedErr := errors.New("email already exists")
	mockDB.On("UpdateUserEmail", ctx, sqlc.UpdateUserEmailParams{
		UserID: 1,
		Email:  "duplicate@example.com",
	}).Return(expectedErr)

	err := service.UpdateUserEmail(ctx, req)

	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}

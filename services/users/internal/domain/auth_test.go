package userservice

import (
	"context"
	"errors"
	"testing"

	"social-network/services/users/internal/db/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
		UserId:      "550e8400-e29b-41d4-a716-446655440000",
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	ctx := context.Background()

	userUUID, _ := stringToUUID(req.UserId)
	mockDB.On("GetUserPassword", ctx, userUUID).Return(hashedOld, nil)
	mockDB.On("UpdateUserPassword", ctx, mock.MatchedBy(func(arg sqlc.UpdateUserPasswordParams) bool {
		return arg.Pub == userUUID
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
		UserId:      "550e8400-e29b-41d4-a716-446655440000",
		OldPassword: "wrongpassword",
		NewPassword: "newpassword",
	}

	ctx := context.Background()

	userUUID, _ := stringToUUID(req.UserId)
	mockDB.On("GetUserPassword", ctx, userUUID).Return(hashedCorrect, nil)

	err := service.UpdateUserPassword(ctx, req)

	assert.Equal(t, ErrNotAuthorized, err)
	mockDB.AssertExpectations(t)
}

func TestUpdateUserEmail_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	userUUID, _ := stringToUUID("550e8400-e29b-41d4-a716-446655440000")
	req := UpdateEmailRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440000",
		Email:  "newemail@example.com",
	}

	ctx := context.Background()

	mockDB.On("UpdateUserEmail", ctx, sqlc.UpdateUserEmailParams{
		Pub:   userUUID,
		Email: "newemail@example.com",
	}).Return(nil)

	err := service.UpdateUserEmail(ctx, req)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUpdateUserEmail_Error(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	userUUID, _ := stringToUUID("550e8400-e29b-41d4-a716-446655440000")
	req := UpdateEmailRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440000",
		Email:  "duplicate@example.com",
	}

	ctx := context.Background()

	expectedErr := errors.New("email already exists")
	mockDB.On("UpdateUserEmail", ctx, sqlc.UpdateUserEmailParams{
		Pub:   userUUID,
		Email: "duplicate@example.com",
	}).Return(expectedErr)

	err := service.UpdateUserEmail(ctx, req)

	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}

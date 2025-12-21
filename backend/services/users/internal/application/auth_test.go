package application

// func TestLoginUser_InvalidCredentials(t *testing.T) {
// 	t.Skip("Requires transaction support with real pool - use integration tests")
// }

// func TestUpdateUserPassword_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	oldPassword := "OldPass123!"
// 	newPassword := "NewPass456!"
// 	storedPassword := "OldPass123!"

// 	req := models.UpdatePasswordRequest{
// 		UserId:      ct.Id(1),
// 		OldPassword: ct.HashedPassword(oldPassword),
// 		NewPassword: ct.HashedPassword(newPassword),
// 	}

// 	ctx := context.Background()

// 	mockDB.On("GetUserPassword", ctx, int64(1)).Return(string(storedPassword), nil)
// 	mockDB.On("UpdateUserPassword", ctx, mock.MatchedBy(func(arg sqlc.UpdateUserPasswordParams) bool {
// 		return arg.UserID == 1
// 	})).Return(nil)

// 	err := service.UpdateUserPassword(ctx, req)

// 	assert.NoError(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestUpdateUserPassword_WrongOldPassword(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	correctPassword := "CorrectPass123!"
// 	hashedCorrect, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

// 	req := models.UpdatePasswordRequest{
// 		UserId:      ct.Id(1),
// 		OldPassword: ct.HashedPassword("WrongPass456!"),
// 		NewPassword: ct.HashedPassword("NewPass789!"),
// 	}

// 	ctx := context.Background()

// 	mockDB.On("GetUserPassword", ctx, int64(1)).Return(string(hashedCorrect), nil)

// 	err := service.UpdateUserPassword(ctx, req)

// 	assert.Equal(t, ErrNotAuthorized, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestUpdateUserEmail_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	req := models.UpdateEmailRequest{
// 		UserId: ct.Id(1),
// 		Email:  ct.Email("newemail@example.com"),
// 	}

// 	ctx := context.Background()

// 	mockDB.On("UpdateUserEmail", ctx, sqlc.UpdateUserEmailParams{
// 		UserID: 1,
// 		Email:  "newemail@example.com",
// 	}).Return(nil)

// 	err := service.UpdateUserEmail(ctx, req)

// 	assert.NoError(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestUpdateUserEmail_Error(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	req := models.UpdateEmailRequest{
// 		UserId: ct.Id(1),
// 		Email:  ct.Email("duplicate@example.com"),
// 	}

// 	ctx := context.Background()

// 	expectedErr := errors.New("email already exists")
// 	mockDB.On("UpdateUserEmail", ctx, sqlc.UpdateUserEmailParams{
// 		UserID: 1,
// 		Email:  "duplicate@example.com",
// 	}).Return(expectedErr)

// 	err := service.UpdateUserEmail(ctx, req)

// 	assert.Equal(t, expectedErr, err)
// 	mockDB.AssertExpectations(t)
// }

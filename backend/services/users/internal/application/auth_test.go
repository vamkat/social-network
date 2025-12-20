package application

// func TestLoginUser_InvalidCredentials(t *testing.T) {
// 	t.Skip("Requires transaction support with real pool - use integration tests")
// }

// func TestApplication_Workflows_Compact(t *testing.T) {
// 	ctx := context.Background()
// 	db, cleanup := setupTestDB(t)
// 	defer cleanup()

// 	app := NewApplication(db, nil, nil)

// 	// Define all operations as a table: [operation name, function args..., expected error]
// 	workflow := []struct {
// 		name string
// 		fn   func() error
// 	}{
// 		{"Create Alice", func() error { _, err := app.CreateGroup(ctx, "alice@example.com", "Alice"); return err }},
// 	}

// 	// Execute sequentially
// 	for i, op := range workflow {
// 		if err := op.fn(); err != nil {
// 			t.Fatalf("[%03d] %s failed: %v", i+1, op.name, err)
// 		}
// 	}
// }

// func TestUpdateUserPassword_Success(t *testing.T) {
// 	ctx := context.Background()

// 	mockDB := &MockDatabase{&sqlc.Queries{}}
// 	service := NewApplication(mockDB, nil, nil)

// 	oldPassword := "OldPass123!"
// 	newPassword := "NewPass456!"
// 	storedPassword := "OldPass123!"

// 	req := models.UpdatePasswordRequest{
// 		UserId:      ct.Id(1),
// 		OldPassword: ct.HashedPassword(oldPassword),
// 		NewPassword: ct.HashedPassword(newPassword),
// 	}

// 	mockQueries.
// 		On("GetUserPassword", ctx, int64(1)).
// 		Return(storedPassword, nil)

// 	mockQueries.
// 		On("UpdateUserPassword", ctx, mock.MatchedBy(func(arg sqlc.UpdateUserPasswordParams) bool {
// 			return arg.UserID == 1
// 		})).
// 		Return(nil)

// 	err := service.UpdateUserPassword(ctx, req)

// 	assert.NoError(t, err)
// 	mockQueries.AssertExpectations(t)
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

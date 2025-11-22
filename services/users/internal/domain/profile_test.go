package userservice

import (
	"context"
	"errors"
	"testing"
	"time"

	"social-network/services/users/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetBasicUserInfo_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID, _ := stringToUUID(userID)

	expectedRow := sqlc.GetUserBasicRow{
		PublicID: userUUID,
		Username: "testuser",
		Avatar:   "avatar.jpg",
	}

	mockDB.On("GetUserBasic", ctx, userUUID).Return(expectedRow, nil)

	user, err := service.GetBasicUserInfo(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.UserId)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "avatar.jpg", user.Avatar)
	mockDB.AssertExpectations(t)
}

func TestGetBasicUserInfo_NotFound(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440999"
	userUUID, _ := stringToUUID(userID)

	mockDB.On("GetUserBasic", ctx, userUUID).Return(nil, errors.New("user not found"))

	_, err := service.GetBasicUserInfo(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockDB.AssertExpectations(t)
}

func TestGetUserProfile_Public_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	requesterID := "550e8400-e29b-41d4-a716-446655440001"
	userUUID, _ := stringToUUID(userID)
	requesterUUID, _ := stringToUUID(requesterID)

	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	dobDate := pgtype.Date{
		Time:  dob,
		Valid: true,
	}

	expectedRow := sqlc.GetUserProfileRow{
		PublicID:      userUUID,
		Username:      "testuser",
		FirstName:     "Test",
		LastName:      "User",
		DateOfBirth:   dobDate,
		Avatar:        "avatar.jpg",
		AboutMe:       "About me",
		ProfilePublic: true,
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	req := UserProfileRequest{
		UserId:      userID,
		RequesterId: requesterID,
	}

	mockDB.On("GetUserProfile", ctx, userUUID).Return(expectedRow, nil)
	mockDB.On("IsFollowing", ctx, mock.MatchedBy(func(arg sqlc.IsFollowingParams) bool {
		return arg.Pub == requesterUUID && arg.Pub_2 == userUUID
	})).Return(false, nil)
	mockDB.On("GetFollowerCount", ctx, userUUID).Return(int64(10), nil)
	mockDB.On("GetFollowingCount", ctx, userUUID).Return(int64(5), nil)
	mockDB.On("UserGroupCountsPerRole", ctx, userUUID).Return(sqlc.UserGroupCountsPerRoleRow{
		TotalMemberships: 3,
		OwnerCount:       1,
	}, nil)

	profile, err := service.GetUserProfile(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, userID, profile.UserId)
	assert.Equal(t, "testuser", profile.Username)
	assert.Equal(t, int64(10), profile.FollowersCount)
	assert.Equal(t, int64(5), profile.FollowingCount)
	mockDB.AssertExpectations(t)
}

func TestGetUserProfile_Private_NotFollowing(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	requesterID := "550e8400-e29b-41d4-a716-446655440001"
	userUUID, _ := stringToUUID(userID)
	requesterUUID, _ := stringToUUID(requesterID)

	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	dobDate := pgtype.Date{
		Time:  dob,
		Valid: true,
	}

	expectedRow := sqlc.GetUserProfileRow{
		PublicID:      userUUID,
		Username:      "testuser",
		FirstName:     "Test",
		LastName:      "User",
		DateOfBirth:   dobDate,
		Avatar:        "avatar.jpg",
		AboutMe:       "About me",
		ProfilePublic: false,
	}

	req := UserProfileRequest{
		UserId:      userID,
		RequesterId: requesterID,
	}

	mockDB.On("GetUserProfile", ctx, userUUID).Return(expectedRow, nil)
	mockDB.On("IsFollowing", ctx, mock.MatchedBy(func(arg sqlc.IsFollowingParams) bool {
		return arg.Pub == requesterUUID && arg.Pub_2 == userUUID
	})).Return(false, nil)

	_, err := service.GetUserProfile(ctx, req)

	assert.Equal(t, ErrProfilePrivate, err)
	mockDB.AssertExpectations(t)
}

func TestGetUserProfile_Private_IsFollowing(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	requesterID := "550e8400-e29b-41d4-a716-446655440001"
	userUUID, _ := stringToUUID(userID)
	requesterUUID, _ := stringToUUID(requesterID)

	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	dobDate := pgtype.Date{
		Time:  dob,
		Valid: true,
	}

	expectedRow := sqlc.GetUserProfileRow{
		PublicID:      userUUID,
		Username:      "testuser",
		FirstName:     "Test",
		LastName:      "User",
		DateOfBirth:   dobDate,
		Avatar:        "avatar.jpg",
		AboutMe:       "About me",
		ProfilePublic: false,
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	req := UserProfileRequest{
		UserId:      userID,
		RequesterId: requesterID,
	}

	mockDB.On("GetUserProfile", ctx, userUUID).Return(expectedRow, nil)
	mockDB.On("IsFollowing", ctx, mock.MatchedBy(func(arg sqlc.IsFollowingParams) bool {
		return arg.Pub == requesterUUID && arg.Pub_2 == userUUID
	})).Return(true, nil)
	mockDB.On("GetFollowerCount", ctx, userUUID).Return(int64(10), nil)
	mockDB.On("GetFollowingCount", ctx, userUUID).Return(int64(5), nil)
	mockDB.On("UserGroupCountsPerRole", ctx, userUUID).Return(sqlc.UserGroupCountsPerRoleRow{
		TotalMemberships: 3,
		OwnerCount:       1,
	}, nil)

	profile, err := service.GetUserProfile(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, userID, profile.UserId)
	mockDB.AssertExpectations(t)
}

func TestSearchUsers_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	searchReq := UserSearchReq{
		SearchTerm: "test",
		Limit:      10,
	}

	uuid1, _ := stringToUUID("550e8400-e29b-41d4-a716-446655440000")
	uuid2, _ := stringToUUID("550e8400-e29b-41d4-a716-446655440001")

	expectedRows := []sqlc.SearchUsersRow{
		{
			PublicID:      uuid1,
			Username:      "testuser1",
			Avatar:        "avatar1.jpg",
			ProfilePublic: true,
		},
		{
			PublicID:      uuid2,
			Username:      "testuser2",
			Avatar:        "avatar2.jpg",
			ProfilePublic: true,
		},
	}

	mockDB.On("SearchUsers", ctx, sqlc.SearchUsersParams{
		Username: "test",
		Limit:    10,
	}).Return(expectedRows, nil)

	users, err := service.SearchUsers(ctx, searchReq)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "testuser1", users[0].Username)
	assert.Equal(t, "testuser2", users[1].Username)
	mockDB.AssertExpectations(t)
}

func TestSearchUsers_NoResults(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	searchReq := UserSearchReq{
		SearchTerm: "nonexistent",
		Limit:      10,
	}

	mockDB.On("SearchUsers", ctx, sqlc.SearchUsersParams{
		Username: "nonexistent",
		Limit:    10,
	}).Return([]sqlc.SearchUsersRow{}, nil)

	users, err := service.SearchUsers(ctx, searchReq)

	assert.NoError(t, err)
	assert.Len(t, users, 0)
	mockDB.AssertExpectations(t)
}

func TestUpdateUserProfile_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID, _ := stringToUUID(userID)

	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	dobDate := pgtype.Date{
		Time:  dob,
		Valid: true,
	}

	req := UpdateProfileRequest{
		UserId:      userID,
		Username:    "newusername",
		FirstName:   "NewFirst",
		LastName:    "NewLast",
		DateOfBirth: dob,
		Avatar:      "newavatar.jpg",
		About:       "New about",
	}

	expectedUser := sqlc.User{
		PublicID:      userUUID,
		Username:      "newusername",
		FirstName:     "NewFirst",
		LastName:      "NewLast",
		DateOfBirth:   dobDate,
		Avatar:        "newavatar.jpg",
		AboutMe:       "New about",
		ProfilePublic: true,
	}

	mockDB.On("UpdateUserProfile", ctx, mock.MatchedBy(func(arg sqlc.UpdateUserProfileParams) bool {
		return arg.Pub == userUUID && arg.Username == "newusername"
	})).Return(expectedUser, nil)

	profile, err := service.UpdateUserProfile(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "newusername", profile.Username)
	assert.Equal(t, "NewFirst", profile.FirstName)
	mockDB.AssertExpectations(t)
}

func TestUpdateProfilePrivacy_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID, _ := stringToUUID(userID)

	req := UpdateProfilePrivacyRequest{
		UserId: userID,
		Public: false,
	}

	mockDB.On("UpdateProfilePrivacy", ctx, sqlc.UpdateProfilePrivacyParams{
		Pub:           userUUID,
		ProfilePublic: false,
	}).Return(nil)

	err := service.UpdateProfilePrivacy(ctx, req)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUpdateProfilePrivacy_Error(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440999"
	userUUID, _ := stringToUUID(userID)

	req := UpdateProfilePrivacyRequest{
		UserId: userID,
		Public: false,
	}

	mockDB.On("UpdateProfilePrivacy", ctx, sqlc.UpdateProfilePrivacyParams{
		Pub:           userUUID,
		ProfilePublic: false,
	}).Return(errors.New("user not found"))

	err := service.UpdateProfilePrivacy(ctx, req)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

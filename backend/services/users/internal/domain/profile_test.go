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
	userID := int64(1)

	expectedRow := sqlc.GetUserBasicRow{
		ID:       userID,
		Username: "testuser",
		Avatar:   "avatar.jpg",
	}

	mockDB.On("GetUserBasic", ctx, userID).Return(expectedRow, nil)

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
	userID := int64(999)

	mockDB.On("GetUserBasic", ctx, userID).Return(nil, errors.New("user not found"))

	_, err := service.GetBasicUserInfo(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockDB.AssertExpectations(t)
}

func TestGetUserProfile_Public_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := int64(1)
	requesterID := int64(2)

	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	dobDate := pgtype.Date{
		Time:  dob,
		Valid: true,
	}

	expectedRow := sqlc.GetUserProfileRow{
		ID:            userID,
		Username:      "testuser",
		FirstName:     "Test",
		LastName:      "User",
		DateOfBirth:   dobDate,
		Avatar:        "avatar.jpg",
		AboutMe:       "About me",
		ProfilePublic: true,
	}

	req := UserProfileRequest{
		UserId:      userID,
		RequesterId: requesterID,
	}

	mockDB.On("GetUserProfile", ctx, userID).Return(expectedRow, nil)
	mockDB.On("IsFollowing", ctx, sqlc.IsFollowingParams{
		FollowerID:  requesterID,
		FollowingID: userID,
	}).Return(true, nil)
	mockDB.On("GetFollowerCount", ctx, userID).Return(int64(10), nil)
	mockDB.On("GetFollowingCount", ctx, userID).Return(int64(5), nil)
	mockDB.On("UserGroupCountsPerRole", ctx, userID).Return(sqlc.UserGroupCountsPerRoleRow{
		OwnerCount:       0,
		MemberOnlyCount:  0,
		TotalMemberships: 0,
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
	userID := int64(1)
	requesterID := int64(2)

	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	dobDate := pgtype.Date{
		Time:  dob,
		Valid: true,
	}

	expectedRow := sqlc.GetUserProfileRow{
		ID:            userID,
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

	mockDB.On("GetUserProfile", ctx, userID).Return(expectedRow, nil)
	mockDB.On("IsFollowing", ctx, sqlc.IsFollowingParams{
		FollowerID:  requesterID,
		FollowingID: userID,
	}).Return(false, nil)

	_, err := service.GetUserProfile(ctx, req)

	assert.Equal(t, ErrProfilePrivate, err)
	mockDB.AssertExpectations(t)
}

func TestGetUserProfile_Private_IsFollowing(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := int64(1)
	requesterID := int64(2)

	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	dobDate := pgtype.Date{
		Time:  dob,
		Valid: true,
	}

	expectedRow := sqlc.GetUserProfileRow{
		ID:            userID,
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

	mockDB.On("GetUserProfile", ctx, userID).Return(expectedRow, nil)
	mockDB.On("IsFollowing", ctx, sqlc.IsFollowingParams{
		FollowerID:  requesterID,
		FollowingID: userID,
	}).Return(true, nil)
	mockDB.On("GetFollowerCount", ctx, userID).Return(int64(10), nil)
	mockDB.On("GetFollowingCount", ctx, userID).Return(int64(5), nil)
	mockDB.On("UserGroupCountsPerRole", ctx, userID).Return(sqlc.UserGroupCountsPerRoleRow{
		OwnerCount:       0,
		MemberOnlyCount:  0,
		TotalMemberships: 0,
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

	expectedRows := []sqlc.SearchUsersRow{
		{
			ID:            1,
			Username:      "testuser1",
			Avatar:        "avatar1.jpg",
			ProfilePublic: true,
		},
		{
			ID:            2,
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
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	dobDate := pgtype.Date{
		Time:  dob,
		Valid: true,
	}

	req := UpdateProfileRequest{
		UserId:      1,
		Username:    "newusername",
		FirstName:   "NewFirst",
		LastName:    "NewLast",
		DateOfBirth: dob,
		Avatar:      "newavatar.jpg",
		About:       "New about",
	}

	expectedUser := sqlc.User{
		ID:            1,
		Username:      "newusername",
		FirstName:     "NewFirst",
		LastName:      "NewLast",
		DateOfBirth:   dobDate,
		Avatar:        "newavatar.jpg",
		AboutMe:       "New about",
		ProfilePublic: true,
	}

	mockDB.On("UpdateUserProfile", ctx, mock.MatchedBy(func(arg sqlc.UpdateUserProfileParams) bool {
		return arg.ID == 1 && arg.Username == "newusername"
	})).Return(expectedUser, nil)

	profile, err := service.UpdateUserProfile(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "newusername", profile.Username)
	assert.Equal(t, "NewFirst", profile.FirstName)
	mockDB.AssertExpectations(t)
}

// DateOfBirth is now a time.Time in the request model; invalid string parsing
// tests are no longer applicable at the domain layer. If needed, validate
// parsing at the API layer. Omit invalid-format test here.

func TestUpdateProfilePrivacy_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()

	req := UpdateProfilePrivacyRequest{
		UserId: 1,
		Public: false,
	}

	mockDB.On("UpdateProfilePrivacy", ctx, sqlc.UpdateProfilePrivacyParams{
		ID:            1,
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

	req := UpdateProfilePrivacyRequest{
		UserId: 999,
		Public: false,
	}

	mockDB.On("UpdateProfilePrivacy", ctx, sqlc.UpdateProfilePrivacyParams{
		ID:            999,
		ProfilePublic: false,
	}).Return(errors.New("user not found"))

	err := service.UpdateProfilePrivacy(ctx, req)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

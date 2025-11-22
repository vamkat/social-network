package userservice

import (
	"context"
	"errors"
	"testing"

	"social-network/services/users/internal/db/sqlc"

	"github.com/stretchr/testify/assert"
)

func TestGetFollowersPaginated_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID, _ := stringToUUID(userID)

	uuid1, _ := stringToUUID("550e8400-e29b-41d4-a716-446655440001")
	uuid2, _ := stringToUUID("550e8400-e29b-41d4-a716-446655440002")

	req := GetFollowersReq{
		FollowingID: userID,
		Limit:       10,
		Offset:      0,
	}

	expectedRows := []sqlc.GetFollowersRow{
		{
			PublicID:      uuid1,
			Username:      "follower1",
			Avatar:        "avatar1.jpg",
			ProfilePublic: true,
		},
		{
			PublicID:      uuid2,
			Username:      "follower2",
			Avatar:        "avatar2.jpg",
			ProfilePublic: true,
		},
	}

	mockDB.On("GetFollowers", ctx, sqlc.GetFollowersParams{
		Pub:    userUUID,
		Limit:  10,
		Offset: 0,
	}).Return(expectedRows, nil)

	users, err := service.GetFollowersPaginated(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "follower1", users[0].Username)
	assert.Equal(t, "follower2", users[1].Username)
	mockDB.AssertExpectations(t)
}

func TestGetFollowersPaginated_Empty(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440999"
	userUUID, _ := stringToUUID(userID)

	req := GetFollowersReq{
		FollowingID: userID,
		Limit:       10,
		Offset:      0,
	}

	mockDB.On("GetFollowers", ctx, sqlc.GetFollowersParams{
		Pub:    userUUID,
		Limit:  10,
		Offset: 0,
	}).Return([]sqlc.GetFollowersRow{}, nil)

	users, err := service.GetFollowersPaginated(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, users, 0)
	mockDB.AssertExpectations(t)
}

func TestGetFollowingPaginated_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID, _ := stringToUUID(userID)

	uuid1, _ := stringToUUID("550e8400-e29b-41d4-a716-446655440001")
	uuid2, _ := stringToUUID("550e8400-e29b-41d4-a716-446655440002")

	req := GetFollowingReq{
		FollowerID: userID,
		Limit:      10,
		Offset:     0,
	}

	expectedRows := []sqlc.GetFollowingRow{
		{
			PublicID:      uuid1,
			Username:      "following1",
			Avatar:        "avatar1.jpg",
			ProfilePublic: true,
		},
		{
			PublicID:      uuid2,
			Username:      "following2",
			Avatar:        "avatar2.jpg",
			ProfilePublic: false,
		},
	}

	mockDB.On("GetFollowing", ctx, sqlc.GetFollowingParams{
		Pub:    userUUID,
		Limit:  10,
		Offset: 0,
	}).Return(expectedRows, nil)

	users, err := service.GetFollowingPaginated(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "following1", users[0].Username)
	assert.Equal(t, "following2", users[1].Username)
	mockDB.AssertExpectations(t)
}

func TestGetFollowingPaginated_Empty(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440999"
	userUUID, _ := stringToUUID(userID)

	req := GetFollowingReq{
		FollowerID: userID,
		Limit:      10,
		Offset:     0,
	}

	mockDB.On("GetFollowing", ctx, sqlc.GetFollowingParams{
		Pub:    userUUID,
		Limit:  10,
		Offset: 0,
	}).Return([]sqlc.GetFollowingRow{}, nil)

	users, err := service.GetFollowingPaginated(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, users, 0)
	mockDB.AssertExpectations(t)
}

func TestFollowUser_Immediate(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	followerID := "550e8400-e29b-41d4-a716-446655440001"
	targetID := "550e8400-e29b-41d4-a716-446655440002"
	followerUUID, _ := stringToUUID(followerID)
	targetUUID, _ := stringToUUID(targetID)

	req := FollowUserReq{
		FollowerId:   followerID,
		TargetUserId: targetID,
	}

	mockDB.On("FollowUser", ctx, sqlc.FollowUserParams{
		Pub:   followerUUID,
		Pub_2: targetUUID,
	}).Return("accepted", nil)

	pending, err := service.FollowUser(ctx, req)

	assert.NoError(t, err)
	assert.False(t, pending)
	mockDB.AssertExpectations(t)
}

func TestFollowUser_Pending(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	followerID := "550e8400-e29b-41d4-a716-446655440001"
	targetID := "550e8400-e29b-41d4-a716-446655440002"
	followerUUID, _ := stringToUUID(followerID)
	targetUUID, _ := stringToUUID(targetID)

	req := FollowUserReq{
		FollowerId:   followerID,
		TargetUserId: targetID,
	}

	mockDB.On("FollowUser", ctx, sqlc.FollowUserParams{
		Pub:   followerUUID,
		Pub_2: targetUUID,
	}).Return("requested", nil)

	pending, err := service.FollowUser(ctx, req)

	assert.NoError(t, err)
	assert.True(t, pending)
	mockDB.AssertExpectations(t)
}

func TestFollowUser_Error(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	followerID := "550e8400-e29b-41d4-a716-446655440001"
	targetID := "550e8400-e29b-41d4-a716-446655440002"
	followerUUID, _ := stringToUUID(followerID)
	targetUUID, _ := stringToUUID(targetID)

	req := FollowUserReq{
		FollowerId:   followerID,
		TargetUserId: targetID,
	}

	mockDB.On("FollowUser", ctx, sqlc.FollowUserParams{
		Pub:   followerUUID,
		Pub_2: targetUUID,
	}).Return("", errors.New("database error"))

	_, err := service.FollowUser(ctx, req)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestUnFollowUser_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	followerID := "550e8400-e29b-41d4-a716-446655440001"
	targetID := "550e8400-e29b-41d4-a716-446655440002"
	followerUUID, _ := stringToUUID(followerID)
	targetUUID, _ := stringToUUID(targetID)

	req := FollowUserReq{
		FollowerId:   followerID,
		TargetUserId: targetID,
	}

	mockDB.On("UnfollowUser", ctx, sqlc.UnfollowUserParams{
		Pub:   followerUUID,
		Pub_2: targetUUID,
	}).Return(nil)

	err := service.UnFollowUser(ctx, req)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUnFollowUser_Error(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	followerID := "550e8400-e29b-41d4-a716-446655440001"
	targetID := "550e8400-e29b-41d4-a716-446655440999"
	followerUUID, _ := stringToUUID(followerID)
	targetUUID, _ := stringToUUID(targetID)

	req := FollowUserReq{
		FollowerId:   followerID,
		TargetUserId: targetID,
	}

	mockDB.On("UnfollowUser", ctx, sqlc.UnfollowUserParams{
		Pub:   followerUUID,
		Pub_2: targetUUID,
	}).Return(errors.New("not following"))

	err := service.UnFollowUser(ctx, req)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestHandleFollowRequest_Accept(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440001"
	requesterID := "550e8400-e29b-41d4-a716-446655440002"
	userUUID, _ := stringToUUID(userID)
	requesterUUID, _ := stringToUUID(requesterID)

	req := HandleFollowRequestReq{
		UserId:      userID,
		RequesterId: requesterID,
		Accept:      true,
	}

	mockDB.On("AcceptFollowRequest", ctx, sqlc.AcceptFollowRequestParams{
		Pub:   requesterUUID,
		Pub_2: userUUID,
	}).Return(nil)

	err := service.HandleFollowRequest(ctx, req)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestHandleFollowRequest_Reject(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440001"
	requesterID := "550e8400-e29b-41d4-a716-446655440002"
	userUUID, _ := stringToUUID(userID)
	requesterUUID, _ := stringToUUID(requesterID)

	req := HandleFollowRequestReq{
		UserId:      userID,
		RequesterId: requesterID,
		Accept:      false,
	}

	mockDB.On("RejectFollowRequest", ctx, sqlc.RejectFollowRequestParams{
		Pub:   requesterUUID,
		Pub_2: userUUID,
	}).Return(nil)

	err := service.HandleFollowRequest(ctx, req)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestHandleFollowRequest_AcceptError(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440001"
	requesterID := "550e8400-e29b-41d4-a716-446655440002"
	userUUID, _ := stringToUUID(userID)
	requesterUUID, _ := stringToUUID(requesterID)

	req := HandleFollowRequestReq{
		UserId:      userID,
		RequesterId: requesterID,
		Accept:      true,
	}

	mockDB.On("AcceptFollowRequest", ctx, sqlc.AcceptFollowRequestParams{
		Pub:   requesterUUID,
		Pub_2: userUUID,
	}).Return(errors.New("request not found"))

	err := service.HandleFollowRequest(ctx, req)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

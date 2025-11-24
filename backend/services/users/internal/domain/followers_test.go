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
	req := GetFollowersReq{
		FollowingID: 1,
		Limit:       10,
		Offset:      0,
	}

	expectedRows := []sqlc.GetFollowersRow{
		{
			ID:            2,
			Username:      "follower1",
			Avatar:        "avatar1.jpg",
			ProfilePublic: true,
		},
		{
			ID:            3,
			Username:      "follower2",
			Avatar:        "avatar2.jpg",
			ProfilePublic: true,
		},
	}

	mockDB.On("GetFollowers", ctx, sqlc.GetFollowersParams{
		FollowingID: 1,
		Limit:       10,
		Offset:      0,
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
	req := GetFollowersReq{
		FollowingID: 999,
		Limit:       10,
		Offset:      0,
	}

	mockDB.On("GetFollowers", ctx, sqlc.GetFollowersParams{
		FollowingID: 999,
		Limit:       10,
		Offset:      0,
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
	req := GetFollowingReq{
		FollowerID: 1,
		Limit:      10,
		Offset:     0,
	}

	expectedRows := []sqlc.GetFollowingRow{
		{
			ID:            2,
			Username:      "following1",
			Avatar:        "avatar1.jpg",
			ProfilePublic: true,
		},
		{
			ID:            3,
			Username:      "following2",
			Avatar:        "avatar2.jpg",
			ProfilePublic: false,
		},
	}

	mockDB.On("GetFollowing", ctx, sqlc.GetFollowingParams{
		FollowerID: 1,
		Limit:      10,
		Offset:     0,
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
	req := GetFollowingReq{
		FollowerID: 999,
		Limit:      10,
		Offset:     0,
	}

	mockDB.On("GetFollowing", ctx, sqlc.GetFollowingParams{
		FollowerID: 999,
		Limit:      10,
		Offset:     0,
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
	req := FollowUserReq{
		FollowerId:   1,
		TargetUserId: 2,
	}

	mockDB.On("FollowUser", ctx, sqlc.FollowUserParams{
		PFollower: 1,
		PTarget:   2,
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
	req := FollowUserReq{
		FollowerId:   1,
		TargetUserId: 2,
	}

	mockDB.On("FollowUser", ctx, sqlc.FollowUserParams{
		PFollower: 1,
		PTarget:   2,
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
	req := FollowUserReq{
		FollowerId:   1,
		TargetUserId: 2,
	}

	mockDB.On("FollowUser", ctx, sqlc.FollowUserParams{
		PFollower: 1,
		PTarget:   2,
	}).Return("", errors.New("database error"))

	_, err := service.FollowUser(ctx, req)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestUnFollowUser_Success(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	req := FollowUserReq{
		FollowerId:   1,
		TargetUserId: 2,
	}

	mockDB.On("UnfollowUser", ctx, sqlc.UnfollowUserParams{
		FollowerID:  1,
		FollowingID: 2,
	}).Return(nil)

	err := service.UnFollowUser(ctx, req)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUnFollowUser_Error(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	req := FollowUserReq{
		FollowerId:   1,
		TargetUserId: 999,
	}

	mockDB.On("UnfollowUser", ctx, sqlc.UnfollowUserParams{
		FollowerID:  1,
		FollowingID: 999,
	}).Return(errors.New("not following"))

	err := service.UnFollowUser(ctx, req)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestHandleFollowRequest_Accept(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	req := HandleFollowRequestReq{
		UserId:      1,
		RequesterId: 2,
		Accept:      true,
	}

	mockDB.On("AcceptFollowRequest", ctx, sqlc.AcceptFollowRequestParams{
		RequesterID: 2,
		TargetID:    1,
	}).Return(nil)

	err := service.HandleFollowRequest(ctx, req)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestHandleFollowRequest_Reject(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	req := HandleFollowRequestReq{
		UserId:      1,
		RequesterId: 2,
		Accept:      false,
	}

	mockDB.On("RejectFollowRequest", ctx, sqlc.RejectFollowRequestParams{
		RequesterID: 2,
		TargetID:    1,
	}).Return(nil)

	err := service.HandleFollowRequest(ctx, req)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestHandleFollowRequest_AcceptError(t *testing.T) {
	mockDB := new(MockQuerier)
	service := NewUserService(mockDB, nil)

	ctx := context.Background()
	req := HandleFollowRequestReq{
		UserId:      1,
		RequesterId: 2,
		Accept:      true,
	}

	mockDB.On("AcceptFollowRequest", ctx, sqlc.AcceptFollowRequestParams{
		RequesterID: 2,
		TargetID:    1,
	}).Return(errors.New("request not found"))

	err := service.HandleFollowRequest(ctx, req)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

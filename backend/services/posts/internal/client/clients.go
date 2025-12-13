package client

import (
	"context"
	cm "social-network/shared/gen-go/common"
	userpb "social-network/shared/gen-go/users"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Holds connections to clients
type Clients struct {
	UserClient userpb.UserServiceClient
}

func (c *Clients) IsFollowing(ctx context.Context, userId, targetUserId int64) (bool, error) {
	resp, err := c.UserClient.IsFollowing(ctx, &userpb.FollowUserRequest{
		FollowerId:   userId,
		TargetUserId: targetUserId,
	})
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

func (c *Clients) IsGroupMember(ctx context.Context, userId, groupId int64) (bool, error) {
	resp, err := c.UserClient.IsGroupMember(ctx, &userpb.GeneralGroupRequest{
		GroupId: groupId,
		UserId:  userId,
	})
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

func (c *Clients) GetBatchBasicUserInfo(ctx context.Context, userIds []int64) (*cm.ListUsers, error) {
	req := &cm.Int64Arr{
		Values: userIds,
	}
	resp, err := c.UserClient.GetBatchBasicUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Clients) GetBasicUserInfo(ctx context.Context, userId int64) (*cm.User, error) {
	req := &wrapperspb.Int64Value{Value: userId}

	resp, err := c.UserClient.GetBasicUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Clients) GetFollowingIds(ctx context.Context, userId int64) ([]int64, error) {
	req := &wrapperspb.Int64Value{Value: userId}

	resp, err := c.UserClient.GetFollowingIds(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}

// func (c *Clients) GetFollowerIds(ctx context.Context, userId int64) ([]int64, error) {
// 	//need to make this in users
// 	return nil, nil
// }

// func (c *Clients) GetUserGroupIds(ctx context.Context, userId int64) ([]int64, error) {
// 	//need to make this in users
// 	return nil, nil
// }

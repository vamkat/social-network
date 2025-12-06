package client

import (
	"context"
	userpb "social-network/shared/gen-go/users"
)

// Holds connections to clients
type Clients struct {
	UserClient userpb.UserServiceClient
}

func (c *Clients) IsFollowing(ctx context.Context, userId, targetUserId int64) (bool, error) {
	resp, err := c.UserClient.IsFollowing(ctx, &userpb.IsFollowingRequest{
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

func (c *Clients) GetBatchBasicUserInfo(ctx context.Context, userIds []int64) (*userpb.ListUsers, error) {
	req := &userpb.Int64Arr{
		Values: userIds,
	}
	resp, err := c.UserClient.GetBatchUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// func (c *Clients) GetFollowingIds(ctx context.Context, userId int64) ([]int64, error) {
// 	req := &wrapperspb.Int64Value{Value: userId}

// 	// Call the gRPC method
// 	resp, err := c.UserClient.GetFollowingIds(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Assuming Int64Arr has a field called Ids []int64
// 	return resp.Values, nil
// }

// func (c *Clients) GetFollowerIds(ctx context.Context, userId int64) ([]int64, error) {
// 	//need to make this in users
// 	return nil, nil
// }

// func (c *Clients) GetUserGroupIds(ctx context.Context, userId int64) ([]int64, error) {
// 	//need to make this in users
// 	return nil, nil
// }

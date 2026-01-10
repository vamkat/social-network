package client

import (
	"context"
	cm "social-network/shared/gen-go/common"
	"social-network/shared/gen-go/media"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

// Converts a slice of ct.Ids representing users to models.User slice.
// DEPRECATED
func (c *Clients) UserIdsToUsers(ctx context.Context,
	ids ct.Ids) (userInfo []md.User, err error) {
	req := &cm.UserIds{Values: ids.Int64()}
	resp, err := c.UserClient.GetBatchBasicUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, u := range resp.Users {
		userInfo = append(userInfo, md.User{
			UserId:    ct.Id(u.UserId),
			Username:  ct.Username(u.Username),
			AvatarId:  ct.Id(u.Avatar),
			AvatarURL: u.AvatarUrl,
		})
	}
	return userInfo, nil
}

// Function to inject to retrive users
func (c *Clients) GetBatchBasicUserInfo(ctx context.Context, req *cm.UserIds) (*cm.ListUsers, error) {
	resp, err := c.UserClient.GetBatchBasicUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// function to inject to retrieve media
func (c *Clients) GetImages(ctx context.Context, req *media.GetImagesRequest) (*media.GetImagesResponse, error) {
	resp, err := c.MediaClient.GetImages(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Clients) IsGroupMember(ctx context.Context,
	groupId ct.Id, userId ct.Id) (bool, error) {
	resp, err := c.UserClient.IsGroupMember(ctx, &users.GeneralGroupRequest{
		GroupId: groupId.Int64(),
		UserId:  userId.Int64(),
	})
	return resp.GetValue(), err
}

func (c *Clients) AreConnected(ctx context.Context, userA, userB ct.Id) (bool, error) {
	resp, err := c.UserClient.AreFollowingEachOther(ctx, &users.FollowUserRequest{
		FollowerId:   userA.Int64(),
		TargetUserId: userB.Int64(),
	})
	if err != nil {
		return false, err
	}

	connected := resp.FollowerFollowsTarget || resp.TargetFollowsFollower
	return connected, nil
}

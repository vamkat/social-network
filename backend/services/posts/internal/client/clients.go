package client

import (
	"context"
	cm "social-network/shared/gen-go/common"
	"social-network/shared/gen-go/media"
	mediapb "social-network/shared/gen-go/media"
	userpb "social-network/shared/gen-go/users"
	ct "social-network/shared/go/ct"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Holds connections to clients
type Clients struct {
	UserClient  userpb.UserServiceClient
	MediaClient mediapb.MediaServiceClient
}

func NewClients(UserClient userpb.UserServiceClient, MediaClient mediapb.MediaServiceClient) *Clients {
	return &Clients{
		UserClient:  UserClient,
		MediaClient: MediaClient,
	}
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

func (c *Clients) GetBatchBasicUserInfo(ctx context.Context, userIds ct.Ids) (*cm.ListUsers, error) {
	req := &cm.UserIds{
		Values: userIds.Int64(),
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

func (c *Clients) GetImages(ctx context.Context, imageIds ct.Ids, variant media.FileVariant) (map[int64]string, []int64, error) {
	req := &mediapb.GetImagesRequest{
		ImgIds:  &mediapb.ImageIds{ImgIds: imageIds.Int64()},
		Variant: variant,
	}
	resp, err := c.MediaClient.GetImages(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	var imagesToDelete []int64
	for _, failedImage := range resp.FailedIds {
		if failedImage.GetStatus() == 4 || failedImage.GetStatus() == 0 {
			imagesToDelete = append(imagesToDelete, failedImage.FileId)
		}
	}
	return resp.DownloadUrls, imagesToDelete, nil
}

func (c *Clients) GetImage(ctx context.Context, imageId int64, variant media.FileVariant) (string, error) {
	req := &mediapb.GetImageRequest{
		ImageId: imageId,
		Variant: variant,
	}
	resp, err := c.MediaClient.GetImage(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.DownloadUrl, nil
}

// func (c *Clients) GetFollowerIds(ctx context.Context, userId int64) ([]int64, error) {
// 	//need to make this in users
// 	return nil, nil
// }

// func (c *Clients) GetUserGroupIds(ctx context.Context, userId int64) ([]int64, error) {
// 	//need to make this in users
// 	return nil, nil
// }

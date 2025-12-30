package client

import (
	"context"
	"fmt"
	cm "social-network/shared/gen-go/common"
	mediapb "social-network/shared/gen-go/media"
	"social-network/shared/gen-go/notifications"
	notifpb "social-network/shared/gen-go/notifications"
	userpb "social-network/shared/gen-go/users"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Holds connections to clients
type Clients struct {
	UserClient   userpb.UserServiceClient
	MediaClient  mediapb.MediaServiceClient
	NotifsClient notifications.NotificationServiceClient
}

func NewClients(UserClient userpb.UserServiceClient, MediaClient mediapb.MediaServiceClient, NotifsClient notifpb.NotificationServiceClient) *Clients {
	return &Clients{
		UserClient:   UserClient,
		MediaClient:  MediaClient,
		NotifsClient: NotifsClient,
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

func (c *Clients) GetBatchBasicUserInfo(ctx context.Context, req *cm.UserIds) (*cm.ListUsers, error) {
	resp, err := c.UserClient.GetBatchBasicUserInfo(ctx, req)
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

func (c *Clients) CreateNewEvent(ctx context.Context, userId, groupId, eventId int64, groupName, eventTitle string) error {
	req := &notifpb.CreateNewEventRequest{
		UserId:     userId,
		GroupId:    groupId,
		EventId:    eventId,
		GroupName:  groupName,
		EventTitle: eventTitle,
	}
	if c.NotifsClient == nil {
		return fmt.Errorf("NotifsClient is nil")
	}
	_, err := c.NotifsClient.CreateNewEvent(ctx, req)
	return err
}

// for comments too?
func (c *Clients) CreatePostLike(ctx context.Context, userId, likerUserId, postId int64, likerUsername string) error {
	req := &notifpb.CreatePostLikeRequest{
		UserId:        userId,
		LikerUserId:   likerUserId,
		PostId:        postId,
		LikerUsername: likerUsername,
		Aggregate:     true,
	}
	if c.NotifsClient == nil {
		return fmt.Errorf("NotifsClient is nil")
	}
	_, err := c.NotifsClient.CreatePostLike(ctx, req)
	return err
}

func (c *Clients) CreatePostComment(ctx context.Context, userId, commenterId, postId int64, commenterUsername, commentContent string) error {
	req := &notifpb.CreatePostCommentRequest{
		UserId:            userId,
		CommenterUserId:   commenterId,
		PostId:            postId,
		CommenterUsername: commenterUsername,
		CommentContent:    commentContent,
		Aggregate:         true,
	}
	if c.NotifsClient == nil {
		return fmt.Errorf("NotifsClient is nil")
	}
	_, err := c.NotifsClient.CreatePostComment(ctx, req)
	return err
}

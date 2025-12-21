package client

import (
	"context"
	chatpb "social-network/shared/gen-go/chat"
	mediapb "social-network/shared/gen-go/media"
	"social-network/shared/gen-go/notifications"
)

// Holds connections to clients
type Clients struct {
	ChatClient   chatpb.ChatServiceClient
	NotifsClient notifications.NotificationServiceClient
	MediaClient  mediapb.MediaServiceClient
}

func NewClients(chatClient chatpb.ChatServiceClient, notifClient notifications.NotificationServiceClient, mediaClient mediapb.MediaServiceClient) *Clients {
	c := &Clients{
		ChatClient:   chatClient,
		NotifsClient: notifClient,
		MediaClient:  mediaClient,
	}
	return c
}

func (c *Clients) GetImages(ctx context.Context, imageIds []int64) (map[int64]string, []int64, error) {
	req := &mediapb.GetImagesRequest{
		ImgIds:  &mediapb.ImageIds{ImgIds: imageIds},
		Variant: 1,
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

func (c *Clients) GetImage(ctx context.Context, imageId int64) (string, error) {
	req := &mediapb.GetImageRequest{
		ImageId: imageId,
		Variant: 1,
	}
	resp, err := c.MediaClient.GetImage(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.DownloadUrl, nil
}

// // on successful follow (public profile or accept follow request)
// func (c *Clients) CreatePrivateConversation(ctx context.Context, userId1, userId2 int64) error {
// 	_, err := c.ChatClient.CreatePrivateConversation(ctx, &chatpb.CreatePrivateConvParams{
// 		UserA: userId1,
// 		UserB: userId2,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // when group is created there's only the owner
// func (c *Clients) CreateGroupConversation(ctx context.Context, groupId, ownerId int64) error {
// 	_, err := c.ChatClient.CreateGroupConversation(ctx, &chatpb.CreateGroupConvParams{
// 		GroupId: groupId,
// 		UserIds: []int64{ownerId},
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *Clients) AddMembersToGroupConversation(ctx context.Context, groupId int64, userIds []int64) error {
// 	_, err := c.ChatClient.AddMembersToGroupConversation(ctx, &chatpb.AddMembersToGroupConversationParams{
// 		GroupId: groupId,
// 		UserIds: userIds,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *Clients) DeleteConversationByExactMembers(ctx context.Context, userIds []int64) error {
// 	_, err := c.ChatClient.DeleteConversationByExactMembers(ctx, &chatpb.UserIds{
// 		UserIds: userIds,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

//remove members from group conversation?
//delete group conversation on group delete?

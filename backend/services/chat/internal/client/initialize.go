package client

import (
	"context"
	cm "social-network/shared/gen-go/common"
	"social-network/shared/gen-go/media"
	notificationspb "social-network/shared/gen-go/notifications"
	userpb "social-network/shared/gen-go/users"
	rds "social-network/shared/go/redis"
)

type Clients struct {
	UserClient         userpb.UserServiceClient
	NotificationClient notificationspb.NotificationServiceClient
	MediaClient        media.MediaServiceClient
	RedisClient        *rds.RedisClient
	// MediaRetriever     *retrievemedia.MediaRetriever
}

type RetriveUsers func(ctx context.Context, userIds *cm.UserIds) (*cm.ListUsers, error)

func NewClients(
	userClient userpb.UserServiceClient,
	notificationClient notificationspb.NotificationServiceClient,
	mediaClient media.MediaServiceClient,
	redis *rds.RedisClient) Clients {

	return Clients{
		UserClient:         userClient,
		NotificationClient: notificationClient,
		MediaClient:        mediaClient,
		RedisClient:        redis,
	}
}

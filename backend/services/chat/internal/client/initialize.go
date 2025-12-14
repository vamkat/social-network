package client

import (
	notificationspb "social-network/shared/gen-go/notifications"
	userpb "social-network/shared/gen-go/users"
)

type Clients struct {
	UserClient         userpb.UserServiceClient
	NotificationClient notificationspb.NotificationServiceClient
}

func NewClients(userClient userpb.UserServiceClient, notificationClient notificationspb.NotificationServiceClient) Clients {
	return Clients{
		UserClient:         userClient,
		NotificationClient: notificationClient,
	}
}

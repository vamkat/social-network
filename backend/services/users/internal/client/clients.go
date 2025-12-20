package client

import (
	chatpb "social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/notifications"
)

// Holds connections to clients
type Clients struct {
	ChatClient   chatpb.ChatServiceClient
	NotifsClient notifications.NotificationServiceClient
}

func NewClients(chatClient chatpb.ChatServiceClient, notifClient notifications.NotificationServiceClient) *Clients {
	c := &Clients{
		ChatClient:   chatClient,
		NotifsClient: notifClient,
	}
	return c
}

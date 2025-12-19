package application

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/notifications"
)

type Database interface {
	TxQueries(context.Context) (*sqlc.Queries, func(context.Context) error, func(context.Context) error, error)
	Queries() *sqlc.Queries
}
type Application struct {
	db            Database
	ChatService   chat.ChatServiceClient
	NotifsService notifications.NotificationServiceClient
}

// NewApplication constructs a new UserService
func NewApplication(db Database, chatService chat.ChatServiceClient, notifService notifications.NotificationServiceClient) *Application {
	return &Application{
		db:            db,
		ChatService:   chatService,
		NotifsService: notifService,
	}
}

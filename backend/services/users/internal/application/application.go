package application

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/notifications"
)

type TxRunner interface {
	RunTx(ctx context.Context, fn func(sqlc.Querier) error) error
}

type Application struct {
	db           sqlc.Querier
	txRunner     TxRunner
	chatService  *chat.ChatServiceClient
	notifService *notifications.NotificationServiceClient
}

// NewApplication constructs a new UserService
func NewApplication[T TxRunner](db sqlc.Querier, txRunner T, chatService *chat.ChatServiceClient, notifService *notifications.NotificationServiceClient) *Application {
	return &Application{
		db:           db,
		txRunner:     txRunner,
		chatService:  chatService,
		notifService: notifService,
	}
}

// func NewApplicationWithMocks(db sqlc.Querier, clients ClientsInterface) *Application {
// 	return &Application{
// 		db:      db,
// 		clients: clients,
// 	}
// }
// func NewApplicationWithMocksTx(db sqlc.Querier, clients ClientsInterface, txRunner TxRunner) *Application {
// 	return &Application{
// 		db:       db,
// 		clients:  clients,
// 		txRunner: txRunner,
// 	}
// }

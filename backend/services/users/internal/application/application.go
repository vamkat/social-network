package application

import (
	"social-network/services/users/internal/client"
	"social-network/services/users/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	db       sqlc.Querier
	txRunner TxRunner
	clients  ClientsInterface
}

// NewApplication constructs a new UserService
func NewApplication(db sqlc.Querier, pool *pgxpool.Pool, clients *client.Clients) *Application {
	var txRunner TxRunner
	if pool != nil {
		queries, ok := db.(*sqlc.Queries)
		if !ok {
			panic("db must be *sqlc.Queries for transaction support")
		}
		txRunner = NewPgxTxRunner(pool, queries)
	}

	return &Application{
		db:       db,
		txRunner: txRunner,
		clients:  clients,
	}
}

// ClientsInterface defines the methods that Application needs from clients.
type ClientsInterface interface {
	// CreateGroupConversation(ctx context.Context, groupId int64, ownerId int64) error
	// CreatePrivateConversation(ctx context.Context, userId1, userId2 int64) error
	// AddMembersToGroupConversation(ctx context.Context, groupId int64, userIds []int64) error
	// DeleteConversationByExactMembers(ctx context.Context, userIds []int64) error
}

func NewApplicationWithMocks(db sqlc.Querier, clients ClientsInterface) *Application {
	return &Application{
		db:      db,
		clients: clients,
	}
}
func NewApplicationWithMocksTx(db sqlc.Querier, clients ClientsInterface, txRunner TxRunner) *Application {
	return &Application{
		db:       db,
		clients:  clients,
		txRunner: txRunner,
	}
}

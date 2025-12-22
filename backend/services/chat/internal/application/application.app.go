package application

import (
	"context"
	"fmt"
	"social-network/services/chat/internal/client"
	"social-network/services/chat/internal/db/dbservice"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
	postgresql "social-network/shared/go/postgre"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxRunner defines the interface for running database transactions
type TxRunner interface {
	RunTx(ctx context.Context, fn func(*dbservice.Queries) error) error
}

// Holds logic for requests and calls
type ChatService struct {
	Clients Clients
	// Clients  *client.Clients
	Queries  dbservice.Querier
	txRunner TxRunner
}

type Clients interface {
	// Calls user client to convert a slice of ct.Ids representing users to a
	// map[ct.Id]models.User.
	UserIdsToMap(ctx context.Context,
		ids ct.Ids) (map[ct.Id]md.User, error)

	// Converts a slice of ct.Ids representing users to models.User slice.
	UserIdsToUsers(ctx context.Context,
		ids ct.Ids) (userInfo []md.User, err error)
}

func NewChatService(pool *pgxpool.Pool, clients *client.Clients, queries dbservice.Querier) (*ChatService, error) {
	var txRunner TxRunner
	var err error
	if pool != nil {
		queries, ok := queries.(*dbservice.Queries)
		if !ok {
			panic("db must be *sqlc.Queries for transaction support")
		}
		txRunner, err = postgresql.NewPgxTxRunner(pool, queries)
		if err != nil {
			return nil, fmt.Errorf("failed to create pgxTxRunner %w", err)
		}
	}
	return &ChatService{
		Clients:  clients,
		Queries:  queries,
		txRunner: txRunner,
	}, nil
}

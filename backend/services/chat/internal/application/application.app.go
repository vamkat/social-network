package application

import (
	"context"
	"fmt"
	"social-network/services/chat/internal/client"
	"social-network/services/chat/internal/db/dbservice"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
	postgresql "social-network/shared/go/postgre"
	"social-network/shared/go/retrieveusers"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxRunner defines the interface for running database transactions
type TxRunner interface {
	RunTx(ctx context.Context, fn func(*dbservice.Queries) error) error
}

// Holds logic for requests and calls
type ChatService struct {
	Clients      Clients
	RetriveUsers *retrieveusers.UserRetriever
	Queries      dbservice.Querier
	txRunner     TxRunner
}

type Clients interface {
	// Converts a slice of ct.Ids representing users to models.User slice.
	UserIdsToUsers(ctx context.Context,
		ids ct.Ids) (userInfo []md.User, err error)
}

func NewChatService(pool *pgxpool.Pool, clients *client.Clients, queries dbservice.Querier, userRetriever *retrieveusers.UserRetriever) (*ChatService, error) {
	var txRunner TxRunner
	var err error
	if pool != nil {
		queries, ok := queries.(*dbservice.Queries)
		if !ok {
			panic("db must be *db.Queries for transaction support")
		}
		txRunner, err = postgresql.NewPgxTxRunner(pool, queries)
		if err != nil {
			return nil, fmt.Errorf("failed to create pgxTxRunner %w", err)
		}
	}

	return &ChatService{
		Clients:      clients,
		Queries:      queries,
		txRunner:     txRunner,
		RetriveUsers: userRetriever,
	}, nil
}

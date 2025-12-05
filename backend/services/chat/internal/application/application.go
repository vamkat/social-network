package application

import (
	"social-network/services/chat/internal/client"
	"social-network/services/chat/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Holds logic for requests and calls
type ChatService struct {
	Pool     *pgxpool.Pool
	Clients  *client.Clients
	Queries  sqlc.Querier
	txRunner TxRunner
}

func NewChatService(pool *pgxpool.Pool, clients *client.Clients, queries sqlc.Querier) *ChatService {
	var txRunner TxRunner
	if pool != nil {
		queries, ok := queries.(*sqlc.Queries)
		if !ok {
			panic("db must be *sqlc.Queries for transaction support")
		}
		txRunner = NewPgxTxRunner(pool, queries)
	}
	return &ChatService{
		Pool:     pool,
		Clients:  clients,
		Queries:  queries,
		txRunner: txRunner,
	}
}

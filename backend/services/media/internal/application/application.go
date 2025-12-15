package application

import (
	"social-network/services/media/internal/client"
	"social-network/services/media/internal/configs"
	"social-network/services/media/internal/db/dbservice"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Holds logic for requests and calls
type MediaService struct {
	Pool     *pgxpool.Pool
	Clients  Clients
	Queries  dbservice.Querier
	txRunner TxRunner
	Cfgs     configs.Config
}

func NewMediaService(pool *pgxpool.Pool, clients *client.Clients,
	queries dbservice.Querier, cfgs configs.Config) *MediaService {
	var txRunner TxRunner
	if pool != nil {
		queries, ok := queries.(*dbservice.Queries)
		if !ok {
			panic("db must be *dbservice.Queries for transaction support")
		}
		txRunner = NewPgxTxRunner(pool, queries)
	}
	return &MediaService{
		Pool:     pool,
		Clients:  clients,
		Queries:  queries,
		txRunner: txRunner,
		Cfgs:     cfgs,
	}
}

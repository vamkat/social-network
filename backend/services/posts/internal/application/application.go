package application

import (
	"social-network/services/posts/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	db       sqlc.Querier
	txRunner TxRunner
}

// NewApplication constructs a new Application with transaction support
func NewApplication(db sqlc.Querier, pool *pgxpool.Pool) *Application {
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
	}
}

// NewApplicationWithTxRunner allows injecting a custom transaction runner (for testing)
func NewApplicationWithTxRunner(db sqlc.Querier, txRunner TxRunner) *Application {
	return &Application{
		db:       db,
		txRunner: txRunner,
	}
}

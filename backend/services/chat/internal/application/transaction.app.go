package application

import (
	"context"
	"social-network/services/chat/internal/db/dbservice"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxRunner defines the interface for running database transactions
type TxRunner interface {
	RunTx(ctx context.Context, fn func(dbservice.Querier) error) error
}

// PgxTxRunner is the production implementation using pgxpool
type PgxTxRunner struct {
	pool *pgxpool.Pool
	db   *dbservice.Queries
}

// NewPgxTxRunner creates a new transaction runner
func NewPgxTxRunner(pool *pgxpool.Pool, db *dbservice.Queries) *PgxTxRunner {
	return &PgxTxRunner{
		pool: pool,
		db:   db,
	}
}

// RunTx wraps a function that contains two or more queries inside a database transaction
func (r *PgxTxRunner) RunTx(ctx context.Context, fn func(dbservice.Querier) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.db.WithTx(tx)

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

package application

import (
	"context"
	"social-network/services/posts/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxRunner defines the interface for running database transactions
type TxRunner interface {
	RunTx(ctx context.Context, fn func(*sqlc.Queries) error) error
}

// PgxTxRunner is the production implementation using pgxpool
type PgxTxRunner struct {
	pool *pgxpool.Pool
	db   *sqlc.Queries
}

// NewPgxTxRunner creates a new transaction runner
func NewPgxTxRunner(pool *pgxpool.Pool, db *sqlc.Queries) *PgxTxRunner {
	return &PgxTxRunner{
		pool: pool,
		db:   db,
	}
}

// RunTx runs a function inside a database transaction
func (r *PgxTxRunner) RunTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	// start tx
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// create queries with transaction
	qtx := r.db.WithTx(tx)

	// run the function
	if err := fn(qtx); err != nil {
		return err
	}

	// commit transaction
	return tx.Commit(ctx)
}

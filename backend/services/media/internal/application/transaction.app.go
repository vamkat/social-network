package application

import (
	"context"
	"social-network/services/media/internal/db/dbservice"

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

// RunTx runs a function inside a database transaction
// The function receives a sqlc.Querier interface, not *sqlc.Queries
func (r *PgxTxRunner) RunTx(ctx context.Context, fn func(dbservice.Querier) error) error {
	// start tx
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// create queries with transaction - returns *sqlc.Queries
	qtx := r.db.WithTx(tx)

	// run the function, passing qtx as sqlc.Querier interface
	if err := fn(qtx); err != nil {
		return err
	}

	// commit transaction
	return tx.Commit(ctx)
}

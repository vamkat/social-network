package postgresql

import (
	"context"
	"errors"
	ce "social-network/shared/go/commonerrors"
	tele "social-network/shared/go/telemetry"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBTX is the minimal interface required by sqlc-generated queries.
//
// It is implemented by both *pgxpool.Pool and pgx.Tx, allowing the same
// queries to run inside or outside a transaction.
// type DBTX interface {
// 	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
// 	Query(context.Context, string, ...any) (pgx.Rows, error)
// 	QueryRow(context.Context, string, ...any) pgx.Row
// }

type HasWithTx[T any] interface {
	WithTx(pgx.Tx) T
}

// PgxTxRunner is the production implementation using pgxpool.
type PgxTxRunner[T HasWithTx[T]] struct {
	pool *pgxpool.Pool
	db   T
}

var ErrNilPassed = errors.New("Passed nil argument")

// NewPgxTxRunner creates a new transaction runner.
func NewPgxTxRunner[T HasWithTx[T]](pool *pgxpool.Pool, db T) (*PgxTxRunner[T], error) {
	if pool == nil {
		return nil, ErrNilPassed
	}
	return &PgxTxRunner[T]{
		pool: pool,
		db:   db,
	}, nil
}

// RunTx runs a function inside a database transaction.
// The function receives a sqlc.Querier interface, not *sqlc.Queries.
func (r *PgxTxRunner[T]) RunTx(ctx context.Context, fn func(T) error) error {
	// start tx.
	tele.Info(ctx, "starting transaction")
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		tele.Error(ctx, "failed to begin transaction @1", "error", err)
		return ce.Wrap(ce.ErrInternal, err, "run tx error")
	}
	defer tx.Rollback(ctx)

	// create queries with transaction - returns *sqlc.Queries.
	qtx := r.db.WithTx(tx)

	// run the function, passing qtx as sqlc.Querier interface.
	if err := fn(qtx); err != nil {
		//TODO add no rows check to avoid logging non error
		tele.Error(ctx, "querier @1", "error", err)
		return err
	}

	// commit transaction.
	tele.Info(ctx, "committing transaction")
	err = tx.Commit(ctx)
	if err != nil {
		return ce.Wrap(ce.ErrInternal, err, "transaction commit error")
	}
	return nil
}

// NewPool creates a pgx connection pool.
//
// Arguments:
//   - ctx: context used to initialize the pool.
//   - address: PostgreSQL connection string.
//
// Returns:
//   - DBTX: pool exposed as a DBTX interface.
//   - *pgxpool.Pool: concrete pool for transaction creation and shutdown.
//   - error: non-nil if the pool cannot be created.
//
// Usage:
//
//	dbtx, pool, err := NewPool(ctx, dsn)
func NewPool(ctx context.Context, address string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, address)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

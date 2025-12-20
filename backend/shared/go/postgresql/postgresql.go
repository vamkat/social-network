package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Package postgresql provides a thin, generic abstraction over pgx/pgxpool
// to support sqlc-style query structs with optional transactional execution.
//

// Database defines the public contract exposed to application layers.
//
// T is typically a sqlc-generated Queries type.
//
// TxQueries:
//   - Starts a new database transaction.
//   - Returns a queries instance bound to that transaction.
//   - Returns commit and rollback functions.
//   - The caller must call exactly one of commit or rollback.
//   - Rollback becomes intempodent after commit is called.
//
// Queries:
//   - Returns a non-transactional queries instance bound to the connection pool.
type Database[T any] interface {
	TxQueries(context.Context) (T, func(context.Context) error, func(context.Context) error, error)
	Queries() T
}

// Postgre is the concrete PostgreSQL implementation of Database.
//
// Q is the sqlc-generated Queries type.
//
// Fields:
//   - queriesTx: helper that can bind queries to a pgx.Tx.
//   - queries: non-transactional queries bound to the pool.
//   - pool: pgx connection pool used for queries and transactions.
type Postgre[Q any] struct {
	QueriesTx   HasWithTX[Q]
	QueriesNorm Q
	pool        *pgxpool.Pool
}

var ErrAlreadyFinished = errors.New("this transaction is already completed")

// TxQueries starts a new database transaction and returns transaction-scoped queries.
//
// Returns:
//   - T: queries instance bound to the started transaction.
//   - commit: function that commits the transaction.
//   - rollback: function that rolls back the transaction.
//   - error: non-nil if the transaction could not be started.
//
// Usage:
//
//		q, commit, rollback, err := db.TxQueries(ctx)
//		if err != nil {
//		    return err
//		}
//		defer rollback(ctx)
//
//		err = q.InsertUser(ctx, params)
//		if err != nil {
//		    return err
//		}
//
//	 	commit(ctx)
//
// Behavior:
//   - Calling either commit or rollback more than once returns ErrAlreadyFinished.
//   - It is suggested to always defer rollback, as it becomes intempodent after commit is called
func (p Postgre[T]) TxQueries(ctx context.Context) (T, func(context.Context) error, func(context.Context) error, error) {
	tx, err := p.newTx(ctx)
	if err != nil {
		return *new(T), nil, nil, err
	}

	finished := false

	commit := func(ctx context.Context) error {
		if finished {
			return ErrAlreadyFinished
		}
		finished = true
		return tx.Commit(ctx)
	}

	rollback := func(ctx context.Context) error {
		if finished {
			return ErrAlreadyFinished
		}
		finished = true
		return tx.Rollback(ctx)
	}

	return p.QueriesTx.WithTx(tx), commit, rollback, nil
}

// Queries returns a non-transactional queries instance.
//
// Usage:
//
//	user, err := db.Queries().GetUser(ctx, id)
//
// Notes:
//   - Queries are executed directly against the pool.
//   - No transaction is created.
func (p Postgre[T]) Queries() T {
	return p.QueriesNorm
}

// NewPostgre initializes a PostgreSQL-backed Database implementation.
//
// Type parameters:
//   - Q: sqlc-generated Queries type.
//   - DB: underlying database interface expected by the sqlc constructor
//     (usually postgresql.DBTX).
//
// Arguments:
//   - ctx: context used to create the connection pool.
//   - address: PostgreSQL connection string.
//   - constructor: sqlc-generated constructor (e.g. sqlc.New).
//   - txQueries: value implementing HasWithTX to support transactions.
//
// Returns:
//   - postgre[Q]: initialized database wrapper.
//   - func(): function that closes the connection pool.
//   - error: non-nil if initialization fails.
//
// Usage:
//
//	db, closeFunc, err := postgresql.NewPostgre(
//	    ctx,
//	    os.Getenv("DATABASE_URL"),
//	    sqlc.New,
//	    &sqlc.Queries{},
//	)
//	defer closeFn()
func NewPostgre[Q any, DB any](
	ctx context.Context,
	address string,
	constuctor func(DB) Q,
	txQueries HasWithTX[Q],
) (Postgre[Q], func(), error) {
	dbtx, pool, err := NewPool(ctx, address)
	if err != nil {
		return Postgre[Q]{}, nil, err
	}

	newPool, ok := dbtx.(DB)
	if !ok {
		return Postgre[Q]{}, nil, fmt.Errorf("DB type %T does not implement expected interface", dbtx)
	}

	queries := constuctor(newPool)

	p := Postgre[Q]{
		QueriesTx:   txQueries,
		QueriesNorm: queries,
		pool:        pool,
	}

	return p, pool.Close, nil
}

// HasWithTX is implemented by sqlc-generated Queries structs that support
// rebinding to a pgx transaction.
type HasWithTX[Q any] interface {
	WithTx(tx pgx.Tx) Q
}

// DBTX is the minimal interface required by sqlc-generated queries.
//
// It is implemented by both *pgxpool.Pool and pgx.Tx, allowing the same
// queries to run inside or outside a transaction.
type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
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
func NewPool(ctx context.Context, address string) (DBTX, *pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, address)
	if err != nil {
		return nil, nil, err
	}
	return pool, pool, nil
}

// newTx starts a new pgx transaction from the pool.
//
// Arguments:
//   - ctx: context used to begin the transaction.
//
// Returns:
//   - pgx.Tx: started transaction.
//   - error: non-nil if the transaction cannot be created.
func (p *Postgre[T]) newTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// FakeDB is used for controlling db responses for testing purposes
type FakeDB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	// We don't need other methods since sqlc only uses DBTX methods
}

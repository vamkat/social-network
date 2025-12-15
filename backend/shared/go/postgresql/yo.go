package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WithTx[T any] interface {
	WithTx(pgx.Tx) T
}

// type TxRunner[T any] interface {
// 	RunTx(ctx context.Context, fn func(T) error) error
// }

type PgxTxRunner[T any] struct {
	pool *pgxpool.Pool
	db   WithTx[T]
}

func NewPgxTxRunner[T any](pool *pgxpool.Pool, db WithTx[T]) *PgxTxRunner[T] {
	return &PgxTxRunner[T]{
		pool: pool,
		db:   db,
	}
}

func (r *PgxTxRunner[T]) RunTx(ctx context.Context, fn func(T) error) error {
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

// var txRunner TxRunner
// 	if pool != nil {
// 		queries, ok := db.(*sqlc.Queries)
// 		if !ok {
// 			panic("db must be *sqlc.Queries for transaction support")
// 		}
// 		txRunner = NewPgxTxRunner(pool, queries)
// 	}
// return txRunner

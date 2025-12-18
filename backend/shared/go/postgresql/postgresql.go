package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database[T any] interface {
	TxQueries(context.Context) (T, pgx.Tx, error)
	Queries() T
}

type postgre[Q any] struct {
	queriesTx HasWithTX[Q]
	queries   Q
	pool      *pgxpool.Pool
}

type HasWithTX[Q any] interface {
	WithTx(tx pgx.Tx) Q
}

type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

func NewPostgre[Q any, DB any](ctx context.Context, address string, constuctor func(DB) Q, txQueries HasWithTX[Q]) (Database[Q], error) {
	dbtx, pool, err := NewPool(ctx, address)
	if err != nil {
		return nil, err
	}
	newPool, ok := dbtx.(DB)
	if !ok {
		return nil, fmt.Errorf("DB type %T does not implement expected interface", dbtx)
	}
	queries := constuctor(newPool)
	p := postgre[Q]{
		queriesTx: txQueries,
		queries:   queries,
		pool:      pool,
	}
	return p, nil
}

func NewPool(ctx context.Context, address string) (DBTX, *pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, address)
	if err != nil {
		return nil, nil, err
	}
	return pool, pool, nil
}

func (p *postgre[T]) newTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (p postgre[T]) TxQueries(ctx context.Context) (T, pgx.Tx, error) {
	tx, err := p.newTx(ctx)
	if err != nil {
		return *new(T), nil, err
	}
	return p.queriesTx.WithTx(tx), tx, nil
}

func (p postgre[T]) Queries() T {
	return p.queries
}

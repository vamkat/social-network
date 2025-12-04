package application

import (
	"social-network/services/posts/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	db   sqlc.Querier  // interface, can be *sqlc.Queries or mock
	pool *pgxpool.Pool // needed to start transactions
}

// NewApplication constructs a new UserService
func NewApplication(db sqlc.Querier, pool *pgxpool.Pool) *Application {
	return &Application{
		db:   db,
		pool: pool,
	}
}

package userservice

import (
	"social-network/services/users/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db   sqlc.Querier  // interface, can be *sqlc.Queries or mock
	pool *pgxpool.Pool // needed to start transactions
}

// NewUserService constructs a new UserService
func NewUserService(db sqlc.Querier, pool *pgxpool.Pool) *UserService {
	return &UserService{
		db:   db,
		pool: pool,
	}
}

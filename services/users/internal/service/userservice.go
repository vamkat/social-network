package userservice

import (
	"context"
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

// kept temporarily for grpc
// -----------------------------------------------------------------------------------
type BasicUserInfo struct {
	UserName      string
	Avatar        string
	PublicProfile bool
}

func GetBasicUserInfo(ctx context.Context, userID int64) (resp BasicUserInfo, err error) {
	return BasicUserInfo{UserName: "Mitsos", Avatar: "M", PublicProfile: true}, nil
}

//-----------------------------------------------------------------------------------

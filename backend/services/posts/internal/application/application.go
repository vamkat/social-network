package application

import (
	"context"
	"social-network/services/posts/internal/client"
	"social-network/services/posts/internal/db/sqlc"
	cm "social-network/shared/gen-go/common"
	userhydrate "social-network/shared/go/hydrateusers"
	"social-network/shared/go/models"
	redis_connector "social-network/shared/go/redis"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	db       sqlc.Querier
	txRunner TxRunner
	clients  ClientsInterface
	hydrator Hydrator
}

type UserHydrator struct {
	clients UsersBatchClient
	cache   RedisCache
	ttl     time.Duration
}

// UsersBatchClient abstracts the single RPC used by the hydrator to fetch basic user info.
type UsersBatchClient interface {
	GetBatchBasicUserInfo(ctx context.Context, userIds []int64) (*cm.ListUsers, error)
}

// RedisCache defines the minimal Redis operations used by the hydrator.
type RedisCache interface {
	GetObj(ctx context.Context, key string, dest any) error
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
}

// Hydrator defines the subset of behavior used by application for user hydration.
type Hydrator interface {
	GetUsers(ctx context.Context, userIDs []int64) (map[int64]models.User, error)
	// HydrateUsers(ctx context.Context, items []models.HasUser) error
	// HydrateUserSlice(ctx context.Context, users []models.User) error
}

// ClientsInterface defines the methods that Application needs from clients.
type ClientsInterface interface {
	IsFollowing(ctx context.Context, userId, targetUserId int64) (bool, error)
	IsGroupMember(ctx context.Context, userId, groupId int64) (bool, error)
	GetFollowingIds(ctx context.Context, userId int64) ([]int64, error)
}

// NewApplication constructs a new Application with transaction support
func NewApplication(db sqlc.Querier, pool *pgxpool.Pool, clients *client.Clients) *Application {
	var txRunner TxRunner
	if pool != nil {
		queries, ok := db.(*sqlc.Queries)
		if !ok {
			panic("db must be *sqlc.Queries for transaction support")
		}
		txRunner = NewPgxTxRunner(pool, queries)
	}

	cache := redis_connector.NewRedisClient("redis:6379", "", 0)

	return &Application{
		db:       db,
		txRunner: txRunner,
		clients:  clients,
		hydrator: userhydrate.NewUserHydrator(clients, cache, 3*time.Minute),
	}
}

func NewApplicationWithMocks(db sqlc.Querier, clients ClientsInterface) *Application {
	return &Application{
		db:      db,
		clients: clients,
	}
}
func NewApplicationWithMocksTx(db sqlc.Querier, clients ClientsInterface, txRunner TxRunner) *Application {
	return &Application{
		db:       db,
		clients:  clients,
		txRunner: txRunner,
	}
}

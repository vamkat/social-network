package application

import (
	"context"
	"fmt"
	"social-network/services/posts/internal/client"
	"social-network/services/posts/internal/db/sqlc"
	cm "social-network/shared/gen-go/common"
	"social-network/shared/go/models"
	postgresql "social-network/shared/go/postgre"
	redis_connector "social-network/shared/go/redis"
	ur "social-network/shared/go/retrieveusers"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxRunner defines the interface for running database transactions
type TxRunner interface {
	RunTx(ctx context.Context, fn func(*sqlc.Queries) error) error
}

type Application struct {
	db            *sqlc.Queries
	txRunner      TxRunner
	clients       ClientsInterface
	userRetriever UserRetriever
}

// UsersBatchClient abstracts the single RPC used by the hydrator to fetch basic user info.
type UsersBatchClient interface {
	GetBatchBasicUserInfo(ctx context.Context, userIds []int64) (*cm.ListUsers, error)
	GetImages(ctx context.Context, imageIds []int64) (map[int64]string, []int64, error)
}

// RedisCache defines the minimal Redis operations used by the hydrator.
type RedisCache interface {
	GetObj(ctx context.Context, key string, dest any) error
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
}

// Hydrator defines the subset of behavior used by application for user hydration.
type UserRetriever interface {
	GetUsers(ctx context.Context, userIDs []int64) (map[int64]models.User, error)
	GetImages(ctx context.Context, imageIds []int64) (map[int64]string, []int64, error)
}

// ClientsInterface defines the methods that Application needs from clients.
type ClientsInterface interface {
	IsFollowing(ctx context.Context, userId, targetUserId int64) (bool, error)
	IsGroupMember(ctx context.Context, userId, groupId int64) (bool, error)
	GetFollowingIds(ctx context.Context, userId int64) ([]int64, error)
}

// NewApplication constructs a new Application with transaction support
func NewApplication(db *sqlc.Queries, pool *pgxpool.Pool, clients *client.Clients) (*Application, error) {
	var txRunner TxRunner
	var err error
	if pool != nil {
		txRunner, err = postgresql.NewPgxTxRunner(pool, db)
		if err != nil {
			return nil, fmt.Errorf("failed to create pgxTxRunner %w", err)
		}
	}

	cache := redis_connector.NewRedisClient("redis:6379", "", 0)

	return &Application{
		db:            db,
		txRunner:      txRunner,
		clients:       clients,
		userRetriever: ur.NewUserRetriever(clients, cache, 3*time.Minute),
	}, nil
}

func NewApplicationWithMocks(db *sqlc.Queries, clients ClientsInterface) *Application {
	return &Application{
		db:      db,
		clients: clients,
	}
}
func NewApplicationWithMocksTx(db *sqlc.Queries, clients ClientsInterface, txRunner TxRunner) *Application {
	return &Application{
		db:       db,
		clients:  clients,
		txRunner: txRunner,
	}
}

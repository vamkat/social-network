package application

import (
	"context"
	"fmt"
	"social-network/services/posts/internal/client"
	ds "social-network/services/posts/internal/db/dbservice"
	cm "social-network/shared/gen-go/common"
	"social-network/shared/gen-go/media"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	postgresql "social-network/shared/go/postgre"
	rds "social-network/shared/go/redis"
	ur "social-network/shared/go/retrieveusers"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxRunner defines the interface for running database transactions
type TxRunner interface {
	RunTx(ctx context.Context, fn func(*ds.Queries) error) error
}

type Application struct {
	db            *ds.Queries
	txRunner      TxRunner
	clients       ClientsInterface
	userRetriever UserRetriever
}

// UsersBatchClient abstracts the single RPC used by the hydrator to fetch basic user info.
type UsersBatchClient interface {
	GetBatchBasicUserInfo(ctx context.Context, userIds ct.Ids) (*cm.ListUsers, error)
	GetImages(ctx context.Context, imageIds ct.Ids, variant media.FileVariant) (map[int64]string, []int64, error)
}

// RedisCache defines the minimal Redis operations used by the hydrator.
type RedisCache interface {
	GetObj(ctx context.Context, key string, dest any) error
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
}

// UserRetriever defines the subset of behavior used by application for user hydration.
type UserRetriever interface {
	GetUsers(ctx context.Context, userIDs ct.Ids) (map[ct.Id]models.User, error)
	GetImages(ctx context.Context, imageIds ct.Ids, variant media.FileVariant) (map[int64]string, []int64, error)
}

// ClientsInterface defines the methods that Application needs from clients.
type ClientsInterface interface {
	IsFollowing(ctx context.Context, userId, targetUserId int64) (bool, error)
	IsGroupMember(ctx context.Context, userId, groupId int64) (bool, error)
	GetFollowingIds(ctx context.Context, userId int64) ([]int64, error)
	GetImage(ctx context.Context, imageId int64, variant media.FileVariant) (string, error)
	GetImages(ctx context.Context, imageIds ct.Ids, variant media.FileVariant) (map[int64]string, []int64, error)
}

// NewApplication constructs a new Application with transaction support
func NewApplication(db *ds.Queries, pool *pgxpool.Pool, clients *client.Clients, redisConnector *rds.RedisClient) (*Application, error) {
	var txRunner TxRunner
	var err error
	if pool != nil {
		txRunner, err = postgresql.NewPgxTxRunner(pool, db)
		if err != nil {
			return nil, fmt.Errorf("failed to create pgxTxRunner %w", err)
		}
	}

	return &Application{
		db:            db,
		txRunner:      txRunner,
		clients:       clients,
		userRetriever: ur.NewUserRetriever(clients, redisConnector, 3*time.Minute),
	}, nil
}

func NewApplicationWithMocks(db *ds.Queries, clients ClientsInterface) *Application {
	return &Application{
		db:      db,
		clients: clients,
	}
}
func NewApplicationWithMocksTx(db *ds.Queries, clients ClientsInterface, txRunner TxRunner) *Application {
	return &Application{
		db:       db,
		clients:  clients,
		txRunner: txRunner,
	}
}

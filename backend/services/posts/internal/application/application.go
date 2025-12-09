package application

import (
	"context"
	"social-network/services/posts/internal/client"
	"social-network/services/posts/internal/db/sqlc"
	redis_connector "social-network/shared/go/redis"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	db       sqlc.Querier
	txRunner TxRunner
	clients  ClientsInterface
	hydrator *UserHydrator
}

type UserHydrator struct {
	clients *client.Clients
	cache   *redis_connector.RedisClient
	ttl     time.Duration
}

// ClientsInterface defines the methods that Application needs from clients.
type ClientsInterface interface {
	IsFollowing(ctx context.Context, userId, targetUserId int64) (bool, error)
	IsGroupMember(ctx context.Context, userId, groupId int64) (bool, error)
	GetFollowingIds(ctx context.Context, userId int64) ([]int64, error)
}

func NewUserHydrator(clients *client.Clients, cache *redis_connector.RedisClient, ttl time.Duration) *UserHydrator {
	return &UserHydrator{
		clients: clients,
		cache:   cache,
		ttl:     ttl,
	}
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

	cache := redis_connector.NewRedisClient("localhost:6379", "", 0)

	return &Application{
		db:       db,
		txRunner: txRunner,
		clients:  clients,
		hydrator: NewUserHydrator(clients, cache, 3*time.Minute),
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
		//hydrator: NewUserHydrator(clients.(*client.Clients)), // <- pass real type
	}
}

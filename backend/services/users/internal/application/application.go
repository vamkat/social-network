package application

import (
	"context"
	"social-network/services/users/internal/client"
	ds "social-network/services/users/internal/db/dbservice"
	"social-network/shared/go/retrievemedia"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxRunner defines the interface for running database transactions
type TxRunner interface {
	RunTx(ctx context.Context, fn func(*ds.Queries) error) error
}

type Application struct {
	db             ds.Querier
	txRunner       TxRunner
	clients        ClientsInterface
	mediaRetriever *retrievemedia.MediaRetriever
}

// NewApplication constructs a new UserService
func NewApplication(db ds.Querier, txRunner TxRunner, pool *pgxpool.Pool, clients *client.Clients) *Application {
	mediaRetriever := retrievemedia.NewMediaRetriever(clients.MediaClient, clients.RedisClient, 3*time.Minute)
	return &Application{
		db:             db,
		txRunner:       txRunner,
		clients:        clients,
		mediaRetriever: mediaRetriever,
	}
}

// ClientsInterface defines the methods that Application needs from clients.
type ClientsInterface interface {
	GetImages(ctx context.Context, imageIds []int64) (map[int64]string, []int64, error)
	GetImage(ctx context.Context, imageId int64) (string, error)
	GetObj(ctx context.Context, key string, dest any) error
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
	// CreateGroupConversation(ctx context.Context, groupId int64, ownerId int64) error
	// CreatePrivateConversation(ctx context.Context, userId1, userId2 int64) error
	// AddMembersToGroupConversation(ctx context.Context, groupId int64, userIds []int64) error
	// DeleteConversationByExactMembers(ctx context.Context, userIds []int64) error
}

func NewApplicationWithMocks(db ds.Querier, clients ClientsInterface) *Application {
	return &Application{
		db:      db,
		clients: clients,
	}
}
func NewApplicationWithMocksTx(db ds.Querier, clients ClientsInterface, txRunner TxRunner) *Application {
	return &Application{
		db:       db,
		clients:  clients,
		txRunner: txRunner,
	}
}

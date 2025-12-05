package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"social-network/services/chat/internal/client"
	"social-network/services/chat/internal/db/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Holds logic for requests and calls
type ChatService struct {
	pool    *pgxpool.Pool
	clients *client.Clients
	db      sqlc.Querier
}

func Run(ctx context.Context) (*ChatService, error) {
	var pool *pgxpool.Pool
	var err error

	connStr := os.Getenv("DATABASE_URL")

	for i := range 10 {
		pool, err = pgxpool.New(ctx, connStr)
		if err == nil {
			break
		}
		log.Printf("DB not ready yet (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to chat database")

	return &ChatService{
		pool:    pool,
		clients: client.InitClients(),
		db:      sqlc.New(pool),
	}, nil
}

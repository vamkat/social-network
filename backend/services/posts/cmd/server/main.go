package main

import (
	"context"
	"log"
	"os"
	"social-network/services/posts/internal/application"
	"social-network/services/posts/internal/db/sqlc"

	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	// cfg := db.LoadConfigFromEnv()

	var pool *pgxpool.Pool
	var err error

	for i := range 10 {
		connStr := os.Getenv("DATABASE_URL")
		pool, err = pgxpool.New(ctx, connStr)
		if err == nil {
			break
		}
		log.Printf("DB not ready yet (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to users database")

	log.Println("Service ready!")

	queries := sqlc.New(pool)
	postsService := application.NewApplication(queries, pool)
	_ = postsService

	// server := server.NewPostsServer(postsService)
	// server.RunGRPCServer()
}

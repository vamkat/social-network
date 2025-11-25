package main

import (
	"context"
	"log"
	"os"
	"time"

	"social-network/services/users/internal/db/sqlc"
	userservice "social-network/services/users/internal/domain"
	"social-network/services/users/internal/server"

	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	userService := userservice.NewUserService(queries, pool)

	service := server.NewUsersServer(userService)
	service.RunGRPCServer()
}

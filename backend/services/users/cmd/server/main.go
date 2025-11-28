package main

import (
	"context"
	"log"
	"os"
	"time"

	userservice "social-network/services/users/internal/application"
	"social-network/services/users/internal/db/sqlc"
	"social-network/services/users/internal/server"

	// Seems to be a relic from when migrations were run from here
	// _ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
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

	server := server.NewUsersServer(userService)
	server.InitClients()
	server.RunGRPCServer()
}

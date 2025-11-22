package main

import (
	"context"
	"log"

	"social-network/services/users/internal/db/sqlc"
	userservice "social-network/services/users/internal/domain"
	"social-network/services/users/internal/server"
	"social-network/shared/db"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	cfg := db.LoadConfigFromEnv()

	pool, err := db.ConnectOrCreateDB(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to users database")

	if err := db.RunMigrations(cfg, "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Service ready!")

	queries := sqlc.New(pool)
	userService := userservice.NewUserService(queries, pool)

	service := server.NewUsersServer(userService)
	service.RunGRPCServer()
}

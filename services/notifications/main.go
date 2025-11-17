package main

import (
	"context"
	"log"

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
}

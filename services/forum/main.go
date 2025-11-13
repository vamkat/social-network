package main

import (
	"context"
	"database/sql"
	"log"

	"social-network/shared/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	cfg := db.LoadConfigFromEnv()

	// Connect to database
	pool, err := db.ConnectOrCreateDB(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	defer pool.Close()

	log.Println("✅ Connected to forum database")

	// Run migrations
	sqlDB, err := sql.Open("postgres", cfg.ConnString())
	if err != nil {
		log.Fatalf("Failed to open DB for migrations: %v", err)
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("✅ Forum migrations applied successfully!")
}

package main

import (
	"context"
	"log"
	"os"
	"social-network/services/notifications/internal/application"
	"social-network/services/notifications/internal/server"
	"time"

	"social-network/services/notifications/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
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
		log.Fatalf("Failed to connect DB: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to notifications database")

	// Create the database queries instance
	queries := sqlc.New(pool)

	// Create the application service
	app := application.NewApplication(queries)

	// Initialize default notification types
	if err := app.CreateDefaultNotificationTypes(ctx); err != nil {
		log.Printf("Warning: failed to create default notification types: %v", err)
	}

	// Create and initialize the server
	s := server.NewNotificationsServer(app)
	s.InitClients()
	s.RunGRPCServer()
}

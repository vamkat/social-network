package entry

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"social-network/services/notifications/internal/application"
	"social-network/services/notifications/internal/db/sqlc"
	"social-network/services/notifications/internal/handler"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run() error {
	pool, err := connectToDb(context.Background())
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to notifications database")

	clients := InitClients()
	app := application.NewApplication(sqlc.New(pool), pool, clients)

	// Initialize default notification types
	if err := app.CreateDefaultNotificationTypes(context.Background()); err != nil {
		log.Printf("Warning: failed to create default notification types: %v", err)
	}

	service := handler.NewNotificationsHandler(app)

	log.Println("Running gRpc service...")
	grpc, err := handler.RunGRPCServer(service)
	if err != nil {
		log.Fatalf("couldn't start gRpc Server: %s", err.Error())
	}

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	grpc.GracefulStop()
	log.Println("Server stopped")
	return nil

}

func connectToDb(ctx context.Context) (pool *pgxpool.Pool, err error) {
	connStr := os.Getenv("DATABASE_URL")
	for i := range 10 {
		pool, err = pgxpool.New(ctx, connStr)
		if err == nil {
			break
		}
		log.Printf("DB not ready yet (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return pool, err
}
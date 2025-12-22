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
	"social-network/shared/gen-go/notifications"
	contextkeys "social-network/shared/go/context-keys"
	"social-network/shared/go/gorpc"
	postgresql "social-network/shared/go/postgre"
	"syscall"
)

func Run() error {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal() //check if ok
	pool, err := postgresql.NewPool(ctx, "postgres://postgres:secret@notifications-db:5432/social_notifications?sslmode=disable")
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
	startServerFunc, endServerFunc, err := gorpc.CreateGRpcServer[notifications.NotificationServiceServer](notifications.RegisterNotificationServiceServer, service, ":50051", contextkeys.CommonKeys())
	if err != nil {
		return err
	}
	defer endServerFunc()

	go func() {
		err := startServerFunc()
		if err != nil {
			log.Fatal("server failed to start")
		}
		fmt.Println("server finished")
	}()

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	endServerFunc()
	log.Println("Server stopped")
	return nil

}

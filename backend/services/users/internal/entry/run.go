package entry

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"social-network/services/users/internal/application"
	"social-network/services/users/internal/db/sqlc"
	"social-network/services/users/internal/handler"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/notifications"
	"social-network/shared/gen-go/users"
	contextkeys "social-network/shared/go/context-keys"
	"social-network/shared/go/gorpc"
	"social-network/shared/go/postgresql"
	"syscall"
)

//TODO add logs as things are getting initialized

func Run() error {
	ctx := context.Background()
	// DATABASE
	pg, closeFunc, err := postgresql.NewPostgre(ctx, os.Getenv("DATABASE_URL"), sqlc.New, &sqlc.Queries{})
	if err != nil {
		log.Fatal("failed to create postgress object:", err.Error())
	}
	defer closeFunc()
	log.Println("Connected to users-db database")

	// CLIENT SERVICES
	chatService, err := gorpc.GetGRpcClient(
		chat.NewChatServiceClient,
		"chat:50051",
		contextkeys.CommonKeys())
	if err != nil {
		log.Fatal("failed to create chat client")
	}

	notifService, err := gorpc.GetGRpcClient(
		notifications.NewNotificationServiceClient,
		"notifications:50051",
		contextkeys.CommonKeys(),
	)
	if err != nil {
		log.Fatal("failed to create notifications client")
	}

	// APPLICATION
	app := application.NewApplication(pg, chatService, notifService)
	service := *handler.NewUsersHanlder(app)

	port := ":50051"

	//
	//
	//
	// SERVER
	startServerFunc, stopServerFunc, err := gorpc.CreateGRpcServer[users.UserServiceServer](
		users.RegisterUserServiceServer,
		&service,
		port,
		contextkeys.CommonKeys(),
	)
	if err != nil {
		log.Fatalf("couldn't create gRpc Server: %s", err.Error())
	}

	go func() {
		err := startServerFunc()
		if err != nil {
			log.Fatal("server failed to start")
		}
		fmt.Println("server finished")
	}()

	//
	//
	//
	// SHUTDOWN

	log.Printf("gRPC server listening on %s", port)

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	stopServerFunc()
	log.Println("Server stopped")
	return nil
}

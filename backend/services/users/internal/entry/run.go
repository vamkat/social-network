package entry

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"social-network/services/users/internal/application"
	"social-network/services/users/internal/client"
	"social-network/services/users/internal/db/sqlc"
	"social-network/services/users/internal/handler"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/notifications"
	"social-network/shared/gen-go/users"
	contextkeys "social-network/shared/go/context-keys"
	"social-network/shared/go/gorpc"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//TODO add logs as things are getting initialized

func Run() error {

	// DATABASE
	pool, err := connectToDb(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to users-db database")

	// CLIENT SERVICES
	chatClient, err := gorpc.GetGRpcClient(chat.NewChatServiceClient, "chat:50051", contextkeys.CommonKeys())
	if err != nil {
		log.Fatal("failed to create chat client")
	}
	notificationsClient, err := gorpc.GetGRpcClient(
		notifications.NewNotificationServiceClient,
		"notifications:50051",
		contextkeys.CommonKeys(),
	)
	if err != nil {
		log.Fatal("failed to create chat client")
	}

	// APPLICATION
	clients := client.NewClients(chatClient, notificationsClient)
	app := application.NewApplication(sqlc.New(pool), pool, clients)
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

func connectToDb(ctx context.Context, address string) (pool *pgxpool.Pool, err error) {
	for i := range 10 {
		pool, err = pgxpool.New(ctx, address)
		if err == nil {
			break
		}
		log.Printf("DB not ready yet (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return pool, err
}

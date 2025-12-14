package entry

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"social-network/services/chat/internal/application"
	"social-network/services/chat/internal/client"
	"social-network/services/chat/internal/db/dbservice"
	"social-network/services/chat/internal/handler"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/notifications"
	"social-network/shared/gen-go/users"
	configutil "social-network/shared/go/configs"
	contextkeys "social-network/shared/go/context-keys"
	"social-network/shared/go/gorpc"

	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type configs struct {
	DatabaseConn        string `env:"DATABASE_URL"`
	GrpcServerPort      string `env:"GRPC_SERVER_PORT"`
	NotificationsAdress string `env:"NOTIFICATIONS_ADDRESS"`
	UsersAdress         string `env:"USERS_ADDRESS"`
}

var cfgs configs

func init() {
	cfgs = configs{
		DatabaseConn:        "postgres://postgres:secret@chat-db:5432/social_chat?sslmode=disable",
		GrpcServerPort:      ":50051",
		NotificationsAdress: "notifications:50051",
		UsersAdress:         "users:50051",
	}
	configutil.LoadConfigs(&cfgs)
}

func Run() error {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

	pool, err := connectToDb(ctx, cfgs.DatabaseConn)
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	fmt.Println("Conneted to DB")

	notifClient, err := gorpc.GetGRpcClient(notifications.NewNotificationServiceClient, cfgs.NotificationsAdress, contextkeys.CommonKeys())
	if err != nil {
		log.Fatal("failed to create notification client: ", err)
	}
	userClient, err := gorpc.GetGRpcClient(users.NewUserServiceClient, cfgs.UsersAdress, contextkeys.CommonKeys())
	if err != nil {
		log.Fatal("failed to create user client: ", err)
	}

	clients := client.NewClients(userClient, notifClient)
	app := application.NewChatService(
		pool,
		&clients,
		dbservice.New(pool),
	)

	handler := handler.NewChatHandler(app)

	startServerFunc, stopServerFunc, err := gorpc.CreateGRpcServer[chat.ChatServiceServer](chat.RegisterChatServiceServer, handler, cfgs.GrpcServerPort, contextkeys.CommonKeys())
	if err != nil {
		log.Fatal("failed to create server:", err.Error())
	}

	go func() {
		fmt.Println("Starting grpc server at port: ", cfgs.GrpcServerPort)
		err := startServerFunc()
		if err != nil {
			log.Fatal("server failed to start")
		}
		fmt.Println("server finished")
	}()

	log.Printf("gRPC server listening on %s", cfgs.GrpcServerPort)

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

package entry

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"social-network/services/notifications/internal/application"
	"social-network/services/notifications/internal/client"
	"social-network/services/notifications/internal/db/sqlc"
	"social-network/services/notifications/internal/handler"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/notifications"
	"social-network/shared/gen-go/posts"
	"social-network/shared/gen-go/users"
	configutil "social-network/shared/go/configs"
	"social-network/shared/go/ct"

	"social-network/shared/go/gorpc"
	postgresql "social-network/shared/go/postgre"
	"syscall"
)

func Run() error {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal() //check if ok

	cfgs := getConfigs()

	pool, err := postgresql.NewPool(ctx, "postgres://postgres:secret@notifications-db:5432/social_notifications?sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to notifications database")

	postsService, err := gorpc.GetGRpcClient(
		posts.NewPostsServiceClient,
		cfgs.PostsGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		log.Fatalf("failed to connect to posts service: %v", err)
	}

	chatService, err := gorpc.GetGRpcClient(
		chat.NewChatServiceClient,
		cfgs.ChatGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		log.Fatalf("failed to connect to chat service: %v", err)
	}

	usersService, err := gorpc.GetGRpcClient(
		users.NewUserServiceClient,
		cfgs.UsersGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		log.Fatalf("failed to connect to media service: %v", err)
	}

	clients := client.NewClients(usersService, chatService, postsService)
	app := application.NewApplication(sqlc.New(pool), clients)

	// Initialize default notification types
	if err := app.CreateDefaultNotificationTypes(context.Background()); err != nil {
		log.Printf("Warning: failed to create default notification types: %v", err)
	}

	service := handler.NewNotificationsHandler(app)

	log.Println("Running gRpc service...")
	startServerFunc, endServerFunc, err := gorpc.CreateGRpcServer[notifications.NotificationServiceServer](
		notifications.RegisterNotificationServiceServer,
		service,
		cfgs.GrpcServerPort,
		ct.CommonKeys())
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

type configs struct {
	RedisAddr     string `env:"REDIS_ADDR"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`

	UsersGRPCAddr  string `env:"USERS_GRPC_ADDR"`
	PostsGRPCAddr  string `env:"POSTS_GRPC_ADDR"`
	ChatGRPCAddr   string `env:"CHAT_GRPC_ADDR"`
	GrpcServerPort string `env:"GRPC_SERVER_PORT"`
}

func getConfigs() configs { // sensible defaults
	cfgs := configs{
		RedisAddr:     "redis:6379",
		RedisPassword: "",
		RedisDB:       0,
		UsersGRPCAddr: "users:50051",
		PostsGRPCAddr: "posts:50051",
		ChatGRPCAddr:  "chat:50051",
	}

	// load environment variables if present
	if err := configutil.LoadConfigs(&cfgs); err != nil {
		log.Fatalf("failed to load env variables into config struct: %v", err)
	}

	return cfgs
}

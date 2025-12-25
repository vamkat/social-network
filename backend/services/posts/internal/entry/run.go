package entry

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"social-network/services/posts/internal/application"
	"social-network/services/posts/internal/client"
	ds "social-network/services/posts/internal/db/dbservice"
	"social-network/services/posts/internal/handler"
	"social-network/shared/gen-go/media"
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
	defer stopSignal()

	cfgs := getConfigs()

	dbUrl := os.Getenv("DATABASE_URL")
	pool, err := postgresql.NewPool(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to posts-db database")

	UsersService, err := gorpc.GetGRpcClient(
		users.NewUserServiceClient,
		cfgs.UsersGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		log.Fatalf("failed to connect to users service: %v", err)
	}

	MediaService, err := gorpc.GetGRpcClient(
		media.NewMediaServiceClient,
		cfgs.MediaGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		log.Fatalf("failed to connect to media service: %v", err)
	}

	clients := client.NewClients(UsersService, MediaService)
	app, err := application.NewApplication(ds.New(pool), pool, clients)
	if err != nil {
		return fmt.Errorf("failed to create posts application: %v", err)
	}

	service := handler.NewPostsHandler(app)

	log.Println("Running gRpc service...")
	startServerFunc, endServerFunc, err := gorpc.CreateGRpcServer[posts.PostsServiceServer](posts.RegisterPostsServiceServer, service, ":50051", ct.CommonKeys())
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

	UsersGRPCAddr string `env:"USERS_GRPC_ADDR"`
	PostsGRPCAddr string `env:"POSTS_GRPC_ADDR"`
	ChatGRPCAddr  string `env:"CHAT_GRPC_ADDR"`
	MediaGRPCAddr string `env:"MEDIA_GRPC_ADDR"`

	HTTPAddr        string `env:"HTTP_ADDR"`
	ShutdownTimeout int    `env:"SHUTDOWN_TIMEOUT_SECONDS"`
}

func getConfigs() configs { // sensible defaults
	cfgs := configs{
		RedisAddr:       "redis:6379",
		RedisPassword:   "",
		RedisDB:         0,
		UsersGRPCAddr:   "users:50051",
		PostsGRPCAddr:   "posts:50051",
		ChatGRPCAddr:    "chat:50051",
		MediaGRPCAddr:   "media:50051",
		HTTPAddr:        "0.0.0.0:8081",
		ShutdownTimeout: 5,
	}

	// load environment variables if present
	if err := configutil.LoadConfigs(&cfgs); err != nil {
		log.Fatalf("failed to load env variables into config struct: %v", err)
	}

	return cfgs
}

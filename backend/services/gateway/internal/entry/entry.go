package entry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"social-network/services/gateway/internal/handlers"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/posts"
	"social-network/shared/gen-go/users"
	configutil "social-network/shared/go/configs"
	contextkeys "social-network/shared/go/context-keys"
	"social-network/shared/go/gorpc"
	redis_connector "social-network/shared/go/redis"
	"syscall"
	"time"
)

type configs struct {
	RedisAddr     string `env:"REDIS_ADDR"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`

	UsersGRPCAddr string `env:"USERS_GRPC_ADDR"`
	PostsGRPCAddr string `env:"POSTS_GRPC_ADDR"`
	ChatGRPCAddr  string `env:"CHAT_GRPC_ADDR"`

	HTTPAddr        string `env:"HTTP_ADDR"`
	ShutdownTimeout int    `env:"SHUTDOWN_TIMEOUT_SECONDS"`
}

var cfg configs

func init() {
	// sensible defaults
	cfg = configs{
		RedisAddr:       "redis:6379",
		RedisPassword:   "",
		RedisDB:         0,
		UsersGRPCAddr:   "users:50051",
		PostsGRPCAddr:   "posts:50051",
		ChatGRPCAddr:    "chat:50051",
		HTTPAddr:        "0.0.0.0:8081",
		ShutdownTimeout: 5,
	}

	// load environment variables if present
	if err := configutil.LoadConfigs(&cfg); err != nil {
		log.Fatalf("failed to load env variables into config struct: %v", err)
	}
}

func Start() {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// Cache
	CacheService := redis_connector.NewRedisClient(
		cfg.RedisAddr,
		cfg.RedisPassword,
		cfg.RedisDB,
	)
	if err := CacheService.TestRedisConnection(); err != nil {
		log.Fatalf("connection test failed, ERROR: %v", err)
	}
	fmt.Println("Cache service connection started correctly")

	//
	//
	//
	// GRPC CLIENTS
	var err error
	UsersService, err := gorpc.GetGRpcClient(
		users.NewUserServiceClient,
		cfg.UsersGRPCAddr,
		contextkeys.CommonKeys(),
	)
	if err != nil {
		log.Fatalf("failed to connect to users service: %v", err)
	}

	PostsService, err := gorpc.GetGRpcClient(
		posts.NewPostsServiceClient,
		cfg.PostsGRPCAddr,
		contextkeys.CommonKeys(),
	)
	if err != nil {
		log.Fatalf("failed to connect to posts service: %v", err)
	}

	ChatService, err := gorpc.GetGRpcClient(
		chat.NewChatServiceClient,
		cfg.ChatGRPCAddr,
		contextkeys.CommonKeys(),
	)
	if err != nil {
		log.Fatalf("failed to connect to chat service: %v", err)
	}

	//
	//
	//
	// HANDLER
	apiMux := handlers.NewHandlers(
		"gateway",
		CacheService,
		UsersService,
		PostsService,
		ChatService,
	)
	if err != nil {
		log.Fatal("Can't create handlers, ERROR:", err)
	}

	//
	//
	//
	// SERVER
	server := &http.Server{
		Handler:     apiMux,
		Addr:        cfg.HTTPAddr,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	srvErr := make(chan error, 1)
	go func() {
		log.Printf("Starting server on http://%s\n", server.Addr)
		srvErr <- server.ListenAndServe()
	}()

	//
	//
	//
	// SHUTDOWN
	select {
	case err = <-srvErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	case <-ctx.Done():
		stopSignal()
	}

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.ShutdownTimeout)*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Graceful server Shutdown Failed: %v", err)
	}

	log.Println("Server stopped")
}

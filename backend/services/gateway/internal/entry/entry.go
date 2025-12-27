package entry

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"social-network/services/gateway/internal/handlers"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/media"
	"social-network/shared/gen-go/posts"
	"social-network/shared/gen-go/users"
	configutil "social-network/shared/go/configs"
	"social-network/shared/go/ct"
	"social-network/shared/go/gorpc"
	redis_connector "social-network/shared/go/redis"

	// tele "social-network/shared/go/telemetry"
	"syscall"
	"time"
)

func Run() {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	cfgs := getConfigs()

	// close := tele.InitTelemetry(ctx, "gateway", ct.CommonKeys(), cfgs.EnableDebugLogs, cfgs.SimplePrint)
	// defer close()

	// tele.Info(ctx, "initialized telemetry")

	// Cache
	CacheService := redis_connector.NewRedisClient(
		cfgs.RedisAddr,
		cfgs.RedisPassword,
		cfgs.RedisDB,
	)
	if err := CacheService.TestRedisConnection(); err != nil {
		// tele.Fatalf("connection test failed, ERROR: %v", err)
	}
	// tele.Info(ctx, "Cache service connection started correctly")

	//
	//
	//
	// GRPC CLIENTS
	var err error
	UsersService, err := gorpc.GetGRpcClient(
		users.NewUserServiceClient,
		cfgs.UsersGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		// tele.Fatalf("failed to connect to users service: %v", err)
	}

	PostsService, err := gorpc.GetGRpcClient(
		posts.NewPostsServiceClient,
		cfgs.PostsGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		// tele.Fatalf("failed to connect to posts service: %v", err)
	}

	ChatService, err := gorpc.GetGRpcClient(
		chat.NewChatServiceClient,
		cfgs.ChatGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		// tele.Fatalf("failed to connect to chat service: %v", err)
	}

	MediaService, err := gorpc.GetGRpcClient(
		media.NewMediaServiceClient,
		cfgs.MediaGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		// tele.Fatalf("failed to connect to media service: %v", err)
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
		MediaService,
	)

	//
	//
	//
	// SERVER
	server := &http.Server{
		Handler:     apiMux,
		Addr:        cfgs.HTTPAddr,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	srvErr := make(chan error, 1)
	go func() {
		// tele.Info(ctx, "Starting server", "address", server.Addr)
		srvErr <- server.ListenAndServe()
	}()

	//
	//
	//
	// SHUTDOWN
	select {
	case err = <-srvErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			// tele.Fatalf("Failed to listen and serve: %v", err)
		}
	case <-ctx.Done():
		stopSignal()
	}

	// tele.Info(ctx, "Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfgs.ShutdownTimeout)*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		// tele.Fatalf("Graceful server Shutdown Failed: %v", err)
	}

	// tele.Info(ctx, "Server stopped")
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
	EnableDebugLogs bool   `env:"ENABLE_DEBUG_LOGS"`
	SimplePrint     bool   `env:"ENABLE_SIMPLE_PRINT"`
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
		EnableDebugLogs: true,
		SimplePrint:     true,
	}

	// load environment variables if present
	if err := configutil.LoadConfigs(&cfgs); err != nil {
		// tele.Fatalf("failed to load env variables into config struct: %v", err)
	}

	return cfgs
}

package entry

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"social-network/services/posts/internal/application"
	"social-network/services/posts/internal/client"
	ds "social-network/services/posts/internal/db/dbservice"
	"social-network/services/posts/internal/handler"
	"social-network/shared/gen-go/media"
	"social-network/shared/gen-go/notifications"
	"social-network/shared/gen-go/posts"
	"social-network/shared/gen-go/users"
	configutil "social-network/shared/go/configs"
	"social-network/shared/go/ct"
	rds "social-network/shared/go/redis"
	tele "social-network/shared/go/telemetry"

	"social-network/shared/go/gorpc"
	postgresql "social-network/shared/go/postgre"
	"syscall"
)

func Run() error {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

	cfgs := getConfigs()

	//
	//
	//
	// TELEMETRY
	closeTelemetry, err := tele.InitTelemetry(ctx, "posts", "PST", cfgs.TelemetryCollectorAddress, ct.CommonKeys(), cfgs.EnableDebugLogs, cfgs.SimplePrint)
	if err != nil {
		tele.Fatalf("failed to init telemetry: %s", err.Error())
	}
	defer closeTelemetry()

	tele.Info(ctx, "initialized telemetry")

	//
	//
	//
	// DATABASE
	dbUrl := os.Getenv("DATABASE_URL")
	pool, err := postgresql.NewPool(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	tele.Info(ctx, "Connected to posts-db database")

	//
	//
	//
	// GRPC CLIENTS
	UsersService, err := gorpc.GetGRpcClient(
		users.NewUserServiceClient,
		cfgs.UsersGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		tele.Fatalf("failed to connect to users service: %v", err)
	}

	MediaService, err := gorpc.GetGRpcClient(
		media.NewMediaServiceClient,
		cfgs.MediaGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		tele.Fatalf("failed to connect to media service: %v", err)
	}

	NotifService, err := gorpc.GetGRpcClient(
		notifications.NewNotificationServiceClient,
		cfgs.NotifGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		tele.Fatalf("failed to connect to notifications service: %v", err)
	}
	if NotifService == nil {
		tele.Fatal("NotifService is nil after initialization")
	}

	//
	//
	//
	// REDIS
	redisConnector := rds.NewRedisClient(cfgs.RedisAddr, cfgs.RedisPassword, cfgs.RedisDB)

	clients := client.NewClients(UsersService, MediaService, NotifService)

	app, err := application.NewApplication(ds.New(pool), pool, clients, redisConnector)
	if err != nil {
		return fmt.Errorf("failed to create posts application: %v", err)
	}

	service := handler.NewPostsHandler(app)
	tele.Info(ctx, "Running gRpc service...")

	//
	//
	//
	// GRPC SERVER
	startServerFunc, endServerFunc, err := gorpc.CreateGRpcServer[posts.PostsServiceServer](
		posts.RegisterPostsServiceServer,
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
			tele.Fatal("server failed to start")
		}
		tele.Info(ctx, "server finished")
	}()

	//
	//
	//
	// SHUTDOWN
	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	tele.Info(ctx, "Shutting down server...")
	endServerFunc()
	tele.Info(ctx, "Server stopped")
	return nil

}

type configs struct {
	RedisAddr     string `env:"REDIS_ADDR"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`

	UsersGRPCAddr  string `env:"USERS_GRPC_ADDR"`
	PostsGRPCAddr  string `env:"POSTS_GRPC_ADDR"`
	ChatGRPCAddr   string `env:"CHAT_GRPC_ADDR"`
	MediaGRPCAddr  string `env:"MEDIA_GRPC_ADDR"`
	NotifGRPCAddr  string `env:"NOTIFICATIONS_ADDRESS"`
	GrpcServerPort string `env:"GRPC_SERVER_PORT"`

	HTTPAddr        string `env:"HTTP_ADDR"`
	ShutdownTimeout int    `env:"SHUTDOWN_TIMEOUT_SECONDS"`

	EnableDebugLogs bool `env:"ENABLE_DEBUG_LOGS"`
	SimplePrint     bool `env:"ENABLE_SIMPLE_PRINT"`

	OtelResourceAttributes    string `env:"OTEL_RESOURCE_ATTRIBUTES"`
	TelemetryCollectorAddress string `env:"TELEMETRY_COLLECTOR_ADDR"`
	PassSecret                string `env:"PASSWORD_SECRET"`
	EncrytpionKey             string `env:"ENC_KEY"`
}

func getConfigs() configs { // sensible defaults
	cfgs := configs{
		RedisAddr:     "redis:6379",
		RedisPassword: "",
		RedisDB:       0,

		UsersGRPCAddr: "users:50051",
		PostsGRPCAddr: "posts:50051",
		ChatGRPCAddr:  "chat:50051",
		MediaGRPCAddr: "media:50051",
		NotifGRPCAddr: "notifications:50051",

		HTTPAddr:        "0.0.0.0:8081",
		ShutdownTimeout: 5,
		GrpcServerPort:  ":50051",

		EnableDebugLogs: true,
		SimplePrint:     true,

		OtelResourceAttributes:    "service.name=posts,service.namespace=social-network,deployment.environment=dev",
		TelemetryCollectorAddress: "alloy:4317",
		PassSecret:                "a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28=",
		EncrytpionKey:             "a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28=",
	}

	// load environment variables if present
	_, err := configutil.LoadConfigs(&cfgs)
	if err != nil {
		tele.Fatalf("failed to load env variables into config struct: %v", err)
	}

	return cfgs
}

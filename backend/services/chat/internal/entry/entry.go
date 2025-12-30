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
	"social-network/shared/gen-go/media"
	"social-network/shared/gen-go/notifications"
	"social-network/shared/gen-go/users"
	configutil "social-network/shared/go/configs"
	"social-network/shared/go/ct"
	"social-network/shared/go/gorpc"
	postgresql "social-network/shared/go/postgre"
	rds "social-network/shared/go/redis"
	"social-network/shared/go/retrievemedia"
	"social-network/shared/go/retrieveusers"
	"time"

	"syscall"
)

type configs struct {
	RedisAddr           string `env:"REDIS_ADDR"`
	RedisPassword       string `env:"REDIS_PASSWORD"`
	RedisDB             int    `env:"REDIS_DB"`
	DatabaseConn        string `env:"DATABASE_URL"`
	GrpcServerPort      string `env:"GRPC_SERVER_PORT"`
	NotificationsAdress string `env:"NOTIFICATIONS_ADDRESS"`
	UsersAdress         string `env:"USERS_ADDRESS"`
	MediaGRPCAddr       string `env:"MEDIA_GRPC_ADDR"`
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
	defer stopSignal() //TODO check if this is ok

	clients := initClients()

	retriveUsers := retrieveusers.NewUserRetriever(
		clients.GetBatchBasicUserInfo,
		clients.RedisClient,
		retrievemedia.NewMediaRetriever(
			clients.MediaClient, clients.RedisClient, 3*time.Minute),
		3*time.Minute,
	)

	pool, err := postgresql.NewPool(ctx, cfgs.DatabaseConn)
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	fmt.Println("Conneted to DB")

	app, err := application.NewChatService(
		pool,
		&clients,
		dbservice.New(pool),
		retriveUsers,
	)
	if err != nil {
		log.Fatal("failed to create chat service application: ", err)
	}

	handler := handler.NewChatHandler(app)

	startServerFunc, stopServerFunc, err := gorpc.CreateGRpcServer[chat.ChatServiceServer](
		chat.RegisterChatServiceServer,
		handler,
		cfgs.GrpcServerPort,
		ct.CommonKeys(),
	)
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

func initClients() client.Clients {
	notifClient, err := gorpc.GetGRpcClient(
		notifications.NewNotificationServiceClient, cfgs.NotificationsAdress, ct.CommonKeys())
	if err != nil {
		log.Fatal("failed to create notification client: ", err)
	}

	userClient, err := gorpc.GetGRpcClient(
		users.NewUserServiceClient, cfgs.UsersAdress, ct.CommonKeys(),
	)
	if err != nil {
		log.Fatal("failed to create user client: ", err)
	}

	mediaClient, err := gorpc.GetGRpcClient(
		media.NewMediaServiceClient,
		cfgs.MediaGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		log.Fatal("failed to create media client: ", err)
	}

	redisClient := rds.NewRedisClient(
		cfgs.RedisAddr, cfgs.RedisPassword, cfgs.RedisDB,
	)

	return client.NewClients(
		userClient,
		notifClient,
		mediaClient,
		redisClient,
	)
}

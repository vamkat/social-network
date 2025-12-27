package entry

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"social-network/services/users/internal/application"
	"social-network/services/users/internal/client"
	ds "social-network/services/users/internal/db/dbservice"
	"social-network/services/users/internal/handler"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/media"
	"social-network/shared/gen-go/notifications"
	"social-network/shared/gen-go/users"
	configutil "social-network/shared/go/configs"
	"social-network/shared/go/ct"
	tele "social-network/shared/go/telemetry"

	"social-network/shared/go/gorpc"
	postgresql "social-network/shared/go/postgre"
	"syscall"
)

//TODO add logs as things are getting initialized

func Run() error {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal() //TODO REMOVE, stop signal should be called when appropriat, during shutdown of server or somethig

	cfgs := getConfigs()

	//
	//
	// CLIENT SERVICES
	chatClient, err := gorpc.GetGRpcClient(
		chat.NewChatServiceClient,
		cfgs.ChatGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		log.Fatal("failed to create chat client")
	}
	mediaClient, err := gorpc.GetGRpcClient(
		media.NewMediaServiceClient,
		cfgs.MediaGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		log.Fatal("failed to create media client")
	}
	notificationsClient, err := gorpc.GetGRpcClient(
		notifications.NewNotificationServiceClient,
		cfgs.NotificationsGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		log.Fatal("failed to create chat client")
	}

	//
	//
	// DATABASE
	pool, err := postgresql.NewPool(ctx, cfgs.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to users-db database")

	//
	//
	// APPLICATION
	clients := client.NewClients(
		chatClient,
		notificationsClient,
		mediaClient,
	)
	pgxTxRunner, err := postgresql.NewPgxTxRunner(pool, ds.New(pool))
	if err != nil {
		log.Fatal("failed to create pgxTxRunner")
	}
	app := application.NewApplication(ds.New(pool), pgxTxRunner, pool, clients)
	service := *handler.NewUsersHanlder(app)

	// port := ":50051"

	//
	//
	// SERVER
	startServerFunc, stopServerFunc, err := gorpc.CreateGRpcServer[users.UserServiceServer](
		users.RegisterUserServiceServer,
		&service,
		cfgs.GrpcServerPort,
		ct.CommonKeys(),
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

type configs struct {
	DatabaseURL           string `env:"DATABASE_URL"`
	ChatGRPCAddr          string `env:"CHAT_GRPC_ADDR"`
	MediaGRPCAddr         string `env:"MEDIA_GRPC_ADDR"`
	NotificationsGRPCAddr string `env:"NOTIFICATIONS_GRPC_ADDR"`
	ShutdownTimeout       int    `env:"SHUTDOWN_TIMEOUT_SECONDS"`
	GrpcServerPort        string `env:"GRPC_SERVER_PORT"`
}

func getConfigs() configs {
	cfgs := configs{
		DatabaseURL:           "postgres://postgres:secret@users-db:5432/social_users?sslmode=disable",
		ChatGRPCAddr:          "chat:50051",
		MediaGRPCAddr:         "media:50051",
		NotificationsGRPCAddr: "notifications:50051",
		ShutdownTimeout:       5,
	}

	_, err := configutil.LoadConfigs(&cfgs)
	if err != nil {
		tele.Fatalf("failed to load env variables into config struct: %v", err)
	}

	return cfgs
}

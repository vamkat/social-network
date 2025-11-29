package main

import (
	"context"
	"log"
	"net"
	"os"

	"social-network/services/notifications/internal/application"
	"social-network/services/notifications/internal/db/sqlc"
	"social-network/services/notifications/internal/server"

	"social-network/shared/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	notificationpb "social-network/shared/gen/notifications"
)

func main() {
	ctx := context.Background()
	cfg := db.LoadConfigFromEnv()

	pool, err := db.ConnectOrCreateDB(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)
	notificationService := application.NewService(queries)
	grpcServer := server.NewNotificationServer(notificationService)

	// Set up gRPC server
	grpcServerImpl := grpc.NewServer()
	reflection.Register(grpcServerImpl)

	// Register the notification service
	notificationpb.RegisterNotificationServiceServer(grpcServerImpl, grpcServer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting notification service on port %s", port)
	if err := grpcServerImpl.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
package handler

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"social-network/services/notifications/internal/application"
	pb "social-network/shared/gen-go/notifications"
)

// Holds Client conns, services and handler funcs
type Server struct {
	pb.UnimplementedNotificationServiceServer
	Application *application.Application
}

func NewNotificationsHandler(app *application.Application) *Server {
	return &Server{
		Application: app,
	}
}

// RunGRPCServer starts the gRPC server and returns the server instance for graceful shutdown
func RunGRPCServer(s *Server) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", ":50051", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterNotificationServiceServer(grpcServer, s)

	log.Printf("gRPC server listening on %s", ":50051")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	return grpcServer, nil
}
package server

import (
	"log"
	"net"

	"social-network/services/notifications/internal/application"
	pb "social-network/shared/gen-go/notifications"
	usersPb "social-network/shared/gen-go/users"

	"google.golang.org/grpc"
)

// Holds Client conns, services and handler funcs
type Server struct {
	pb.UnimplementedNotificationServiceServer
	Clients     Clients
	Port        string
	Application *application.Application
}

// Holds connections to clients
type Clients struct {
	// define here all grpc connections
	// Example client definitions (uncomment when respective service clients are generated):
	UsersClient usersPb.UserServiceClient
	// PostsClient postsPb.PostsServiceClient
}

// RunGRPCServer starts the gRPC server and blocks
func (s *Server) RunGRPCServer() {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterNotificationServiceServer(grpcServer, s)

	log.Printf("gRPC server listening on %s", s.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

// NewNotificationsServer creates a new notification server
func NewNotificationsServer(app *application.Application) *Server {
	return &Server{
		Port:        ":50051", // Default port for notifications service (internal)
		Application: app,
		Clients:     Clients{},
	}
}

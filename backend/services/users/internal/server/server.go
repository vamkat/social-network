package server

import (
	"log"
	"net"

	us "social-network/services/users/internal/application"
	pb "social-network/shared/gen-go/users"

	"google.golang.org/grpc"
)

// Holds Client conns, services and handler funcs
type Server struct {
	pb.UnimplementedUserServiceServer
	Clients Clients
	Port    string
	Service *us.UserService
}

// Holds connections to clients
type Clients struct {
}

// RunGRPCServer starts the gRPC server and blocks
func (s *Server) RunGRPCServer() {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &Server{})

	log.Printf("gRPC server listening on %s", s.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func NewUsersServer(service *us.UserService) *Server {
	return &Server{
		Port:    ":50051",
		Clients: Clients{},
		Service: service,
	}
}

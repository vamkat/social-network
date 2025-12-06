package handler

import (
	"log"
	"net"

	"social-network/services/users/internal/application"
	pb "social-network/shared/gen-go/users"

	"google.golang.org/grpc"
)

// Holds Client conns, services and handler funcs
type UsersHandler struct {
	pb.UnimplementedUserServiceServer
	Port        string
	Application *application.Application
}

// RunGRPCServer starts the gRPC server and blocks
func (s *UsersHandler) RunGRPCServer() {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, s)

	services := grpcServer.GetServiceInfo()
	log.Printf("Registered services: %v", services)

	log.Printf("gRPC server listening on %s", s.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func NewUsersHanlder(service *application.Application) *UsersHandler {
	return &UsersHandler{
		Port:        ":50051",
		Application: service,
	}
}

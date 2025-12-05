package handler

import (
	"log"
	"net"

	"social-network/services/posts/internal/application"
	//pb "social-network/shared/gen-go/users"

	"google.golang.org/grpc"
)

// Holds Client conns, services and handler funcs
type PostsHandler struct {
	//pb.UnimplementedPostsServiceServer
	Port        string
	Application *application.Application
}

// RunGRPCServer starts the gRPC server and blocks
func (s *PostsHandler) RunGRPCServer() {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	grpcServer := grpc.NewServer()

	//pb.RegisterUserServiceServer(grpcServer, s)

	log.Printf("gRPC server listening on %s", s.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func NewPostsHandler(service *application.Application) *PostsHandler {
	return &PostsHandler{
		Port:        ":50051",
		Application: service,
	}
}

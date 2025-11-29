package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "social-network/shared/gen-go/users"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Server struct placeholder
type Server struct {
	pb.UnimplementedUserServiceServer
	Clients Clients
	Port    string
	// Service *us.UserService
}

type Clients struct {
	UserClient pb.UserServiceClient
}

// RunGRPCServer starts the gRPC server and blocks
func (s *Server) RunGRPCServer() {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	grpcServer := grpc.NewServer()

	// pb.RegisterUserServiceServer(grpcServer, &Server{})

	// ================================================
	// TEST GRPC CONN
	// ================================================

	// Wait for other servers to run first
	time.Sleep(time.Second)
	s.InitClients()

	// Example call to UserClient
	resp, err := s.Clients.UserClient.GetBasicUserInfo(
		context.Background(),
		&wrapperspb.Int64Value{Value: 1234},
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.String())

	// -------------------------------------------------

	log.Printf("gRPC server listening on %s", s.Port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func NewPostsServer() *Server {
	return &Server{
		Port:    ":50051",
		Clients: Clients{},
	}
}

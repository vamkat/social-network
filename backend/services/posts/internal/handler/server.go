package handler

import (
	"log"
	"net"
	"social-network/services/posts/internal/application"
	pb "social-network/shared/gen-go/posts"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/gorpc"

	"google.golang.org/grpc"
)

// Holds Client conns, services and handler funcs
type PostsHandler struct {
	pb.UnimplementedPostsServiceServer
	Port        string
	Application *application.Application
}

func NewPostsHandler(service *application.Application) *PostsHandler {
	return &PostsHandler{
		Port:        ":50051",
		Application: service,
	}
}

// RunGRPCServer starts the gRPC server and blocks
func RunGRPCServer(s *PostsHandler) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	customUnaryInterceptor, err := gorpc.UnaryServerInterceptorWithContextKeys([]gorpc.StringableKey{ct.UserId, ct.ReqID, ct.TraceId}...)
	if err != nil {
		return nil, err
	}
	customStreamInterceptor, err := gorpc.StreamServerInterceptorWithContextKeys([]gorpc.StringableKey{ct.UserId, ct.ReqID, ct.TraceId}...)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(customUnaryInterceptor),
		grpc.StreamInterceptor(customStreamInterceptor),
	)

	pb.RegisterPostsServiceServer(grpcServer, s)

	log.Printf("gRPC server listening on %s", s.Port)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	return grpcServer, nil
}

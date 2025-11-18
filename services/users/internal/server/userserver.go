package server

import (
	"context"
	"log"
	"net"
	userservice "social-network/services/users/internal/service"

	commonpb "social-network/shared/gen/common"
	pb "social-network/shared/gen/users"

	"google.golang.org/grpc"
)

// UserServer struct placeholder
type UserServer struct {
	// Add DB or other dependencies here later
}

// RunGRPCServer starts the gRPC server and blocks
func RunGRPCServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	// TODO: Register services here, e.g.,
	// pb.RegisterUserServiceServer(grpcServer, &UserServer{})

	log.Printf("gRPC server listening on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func (s *UserServer) GetBasicUserInfo(ctx context.Context, req *commonpb.UserId) (*pb.BasicUserInfo, error) {
	u, err := userservice.GetBasicUserInfo(ctx, req.Id)
	return &pb.BasicUserInfo{
		UserName:      u.UserName,
		Avatar:        u.Avatar,
		PublicProfile: u.PublicProfile,
	}, err
}

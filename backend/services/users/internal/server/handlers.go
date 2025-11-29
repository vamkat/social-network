/*
Expose methods via gRpc
*/

package server

import (
	"context"
	pb "social-network/shared/gen-go/users"
)

func (s *Server) GetBasicUserInfo(ctx context.Context, req *pb.UserBasicInfoRequest) (*pb.UserBasicInfoResponse, error) {
	u, err := s.Service.GetBasicUserInfo(ctx, req.Id)
	return &pb.UserBasicInfoResponse{
		UserName: u.Username,
		Avatar:   u.Avatar,
	}, err
}

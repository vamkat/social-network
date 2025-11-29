/*
Expose methods via gRpc
*/

package server

import (
	"context"
	pb "social-network/shared/gen-go/users"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Server) GetBasicUserInfo(ctx context.Context, req *wrapperspb.Int64Value) (*pb.User, error) {
	u, err := s.Service.GetBasicUserInfo(ctx, req.GetValue())
	return &pb.User{
		UserId:   u.UserId,
		Username: u.Username,
		Avatar:   u.Avatar,
	}, err
}

/*
Expose methods via gRpc
*/

package server

import (
	"context"
	commonpb "social-network/shared/gen/common"
	pb "social-network/shared/gen/users"
)

func (s *Server) GetBasicUserInfo(ctx context.Context, req *commonpb.UserId) (*pb.BasicUserInfo, error) {
	u, err := s.Service.GetBasicUserInfo(ctx, req.Id)
	return &pb.BasicUserInfo{
		UserName: u.Username,
		Avatar:   u.Avatar,
	}, err
}

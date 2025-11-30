/*
Expose methods via gRpc
*/

package server

import (
	"context"
	"fmt"
	"social-network/services/users/internal/application"
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

func (s *Server) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfileResponse, error) {
	fmt.Println("GetUserProfile gRPC method called")
	userProfileRequest := application.UserProfileRequest{
		UserId:      req.GetUserId(),
		RequesterId: req.GetRequesterId(),
	}

	profile, err := s.Service.GetUserProfile(ctx, userProfileRequest)
	if err != nil {
		fmt.Println("Error in GetUserProfile:", err)
		return nil, err
	}

	return &pb.UserProfileResponse{
		UserId:   profile.UserId,
		Username: profile.Username,
	}, nil
}

func (s *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.User, error) {
	fmt.Println("RegisterUser gRPC method called")
	user, err := s.Service.RegisterUser(ctx, application.RegisterUserRequest{
		Username:    req.GetUsername(),
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		Email:       req.GetEmail(),
		Password:    req.GetPassword(),
		DateOfBirth: req.GetDateOfBirth().AsTime(),
		Avatar:      req.GetAvatar(),
		About:       req.GetAbout(),
		Public:      req.GetPublic(),
	})
	if err != nil {
		fmt.Println("Error in RegisterUser:", err)
		return nil, err
	}

	return &pb.User{
		UserId:   user.UserId,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.User, error) {
	fmt.Println("LoginUser gRPC method called")
	user, err := s.Service.LoginUser(ctx, application.LoginRequest{
		Identifier: req.GetIdentifier(),
		Password:   req.GetPassword(),
	})
	if err != nil {
		fmt.Println("Error in LoginUser:", err)
		return nil, err
	}

	return &pb.User{
		UserId:   user.UserId,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}

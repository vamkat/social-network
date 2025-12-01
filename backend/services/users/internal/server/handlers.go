/*
Expose methods via gRpc
*/

package server

import (
	"context"
	"fmt"
	"runtime"
	"social-network/services/users/internal/application"
	pb "social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// AUTH
func (s *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.User, error) {
	fmt.Println("RegisterUser gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	user, err := s.Application.RegisterUser(ctx, application.RegisterUserRequest{
		Username:    ct.Username(req.GetUsername()),
		FirstName:   ct.Name(req.GetFirstName()),
		LastName:    ct.Name(req.GetLastName()),
		Email:       ct.Email(req.GetEmail()),
		Password:    ct.Password(req.GetPassword()),
		DateOfBirth: ct.DateOfBirth(req.GetDateOfBirth().AsTime()),
		Avatar:      req.GetAvatar(),
		About:       ct.About(req.GetAbout()),
		Public:      req.GetPublic(),
	})
	if err != nil {
		fmt.Println("Error in RegisterUser:", err)
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.User{
		UserId:   user.UserId.Int64(),
		Username: user.Username.String(),
		Avatar:   user.Avatar,
	}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.User, error) {
	fmt.Println("LoginUser gRPC method called")

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "LoginUser: request is nil")
	}

	Identifier := req.GetIdentifier()
	if err := invalidString("ident", Identifier); err != nil {
		return nil, err
	}

	Password := req.GetPassword()
	if err := invalidString("pass", Password); err != nil {
		return nil, err
	}

	user, err := s.Application.LoginUser(ctx, application.LoginRequest{
		Identifier: ct.Identifier(Identifier),
		Password:   ct.Password(Password),
	})
	if err != nil {
		fmt.Println("Error in LoginUser:", err)
		return nil, status.Errorf(codes.Internal, "LoginUser: failed to login user: %v", err)
	}

	return &pb.User{
		UserId:   user.UserId.Int64(),
		Username: user.Username.String(),
		Avatar:   user.Avatar,
	}, nil
}

func (s *Server) UpdateUserPassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "UpdateUserPassword: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	newPassword := req.GetNewPassword()
	if err := invalidString("newPassword", newPassword); err != nil {
		return nil, err
	}

	err := s.Application.UpdateUserPassword(ctx, application.UpdatePasswordRequest{
		UserId:      ct.Id(userId),
		NewPassword: ct.Password(newPassword),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UpdateUserPassword: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateUserEmail(ctx context.Context, req *pb.UpdateEmailRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "UpdateUserEmail: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("user id", userId); err != nil {
		return nil, err
	}

	newEmail := req.GetEmail()
	if err := invalidString("newEmail", newEmail); err != nil {
		return nil, err
	}

	err := s.Application.UpdateUserEmail(ctx, application.UpdateEmailRequest{
		UserId: ct.Id(userId),
		Email:  ct.Email(newEmail),
	})
	if err != nil {
		fmt.Println("Error in UpdateUserEmail:", err)
		return nil, status.Errorf(codes.Internal, "UpdateUserEmail: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// FOLLOW
func (s *Server) GetFollowersPaginated(ctx context.Context, req *pb.Pagination) (*pb.ListUsers, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetFollowersPaginated: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	limit := req.GetLimit()
	offset := req.GetOffset()
	if err := checkLimOff(limit, offset); err != nil {
		return nil, err
	}

	pag := application.Pagination{
		UserId: ct.Id(userId),
		Limit:  ct.Limit(limit),
		Offset: ct.Offset(offset),
	}

	resp, err := s.Application.GetFollowersPaginated(ctx, pag)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetFollowersPaginated: %v", err)
	}
	return usersToPB(resp), nil
}

func (s *Server) GetFollowingPaginated(ctx context.Context, req *pb.Pagination) (*pb.ListUsers, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetFollowingPaginated: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	limit := req.GetLimit()
	offset := req.GetOffset()
	if err := checkLimOff(limit, offset); err != nil {
		return nil, err
	}

	pag := application.Pagination{
		UserId: ct.Id(userId),
		Limit:  ct.Limit(limit),
		Offset: ct.Offset(offset),
	}

	resp, err := s.Application.GetFollowingPaginated(ctx, pag)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetFollowingPaginated: %v", err)
	}
	return usersToPB(resp), nil
}

func (s *Server) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "FollowUser: request is nil")
	}

	followerId := req.GetFollowerId()
	if err := invalidId("followerId", followerId); err != nil {
		return nil, err
	}

	targetUserId := req.GetTargetUserId()
	if err := invalidId("targetUserId", targetUserId); err != nil {
		return nil, err
	}

	resp, err := s.Application.FollowUser(ctx, application.FollowUserReq{
		FollowerId:   ct.Id(followerId),
		TargetUserId: ct.Id(targetUserId),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "FollowUser: %v", err)
	}

	return &pb.FollowUserResponse{
		IsPending:         resp.IsPending,
		ViewerIsFollowing: resp.ViewerIsFollowing,
	}, nil
}

func (s *Server) UnFollowUser(ctx context.Context, req *pb.FollowUserRequest) (*wrapperspb.BoolValue, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "UnFollowUser: request is nil")
	}

	followerId := req.GetFollowerId()
	if err := invalidId("followerId", followerId); err != nil {
		return nil, err
	}

	targetUserId := req.GetTargetUserId()
	if err := invalidId("targetUserId", targetUserId); err != nil {
		return nil, err
	}

	resp, err := s.Application.UnFollowUser(ctx, application.FollowUserReq{
		FollowerId:   ct.Id(followerId),
		TargetUserId: ct.Id(targetUserId),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UnFollowUser: %v", err)
	}

	return wrapperspb.Bool(resp), nil
}

func (s *Server) HandleFollowRequest(ctx context.Context, req *pb.HandleFollowRequestRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "HandleFollowRequest: request is nil")
	}

	userID := req.GetUserId()
	if err := invalidId("userId", userID); err != nil {
		return nil, err
	}

	RequesterId := req.GetRequesterId()
	if err := invalidId("requesterId", RequesterId); err != nil {
		return nil, err
	}

	acc := req.GetAccept()

	err := s.Application.HandleFollowRequest(ctx, application.HandleFollowRequestReq{
		UserId:      ct.Id(userID),
		RequesterId: ct.Id(RequesterId),
		Accept:      acc,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "HandleFollowRequest: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetFollowingIds(ctx context.Context, req *wrapperspb.Int64Value) (*pb.Int64Arr, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetFollowingIds: request is nil")
	}
	userId := req.GetValue()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	resp, err := s.Application.GetFollowingIds(ctx, ct.Id(userId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetFollowingIds: %v", err)
	}
	return &pb.Int64Arr{Values: resp}, nil
}

func (s *Server) GetFollowSuggestions(ctx context.Context, req *wrapperspb.Int64Value) (*pb.ListUsers, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetFollowSuggestions: request is nil")
	}

	userId := req.GetValue()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	resp, err := s.Application.GetFollowSuggestions(ctx, ct.Id(userId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetFollowSuggestions: %v", err)
	}

	return usersToPB(resp), nil
}

// GROUPS
func (s *Server) GetAllGroupsPaginated(ctx context.Context, req *pb.Pagination) (*pb.GroupArr, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetAllGroupsPaginated: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	limit := req.GetLimit()
	offset := req.GetOffset()
	if err := checkLimOff(limit, offset); err != nil {
		return nil, err
	}

	pag := application.Pagination{
		UserId: ct.Id(userId),
		Limit:  ct.Limit(limit),
		Offset: ct.Offset(offset),
	}

	resp, err := s.Application.GetAllGroupsPaginated(ctx, pag)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetAllGroupsPaginated: %v", err)
	}
	return groupsToPb(resp), nil
}

func (s *Server) GetUserGroupsPaginated(ctx context.Context, req *pb.Pagination) (*pb.GroupArr, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetUserGroupsPaginated: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	limit := req.GetLimit()
	offset := req.GetOffset()
	if err := checkLimOff(limit, offset); err != nil {
		return nil, err
	}

	pag := application.Pagination{
		UserId: ct.Id(userId),
		Limit:  ct.Limit(limit),
		Offset: ct.Offset(offset),
	}

	resp, err := s.Application.GetUserGroupsPaginated(ctx, pag)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetUserGroupsPaginated: %v", err)
	}
	return groupsToPb(resp), nil
}

func (s *Server) GetGroupInfo(ctx context.Context, req *pb.GeneralGroupRequest) (*pb.Group, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetGroupInfo: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	groupId := req.GetGroupId()
	if err := invalidId("groupId", groupId); err != nil {
		return nil, err
	}

	resp, err := s.Application.GetGroupInfo(ctx, application.GeneralGroupReq{
		UserId:  ct.Id(userId),
		GroupId: ct.Id(groupId),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetGroupInfo: %v", err)
	}

	return &pb.Group{
		GroupId:          resp.GroupId.Int64(),
		GroupOwnerId:     resp.GroupOwnerId.Int64(),
		GroupTitle:       resp.GroupTitle.String(),
		GroupDescription: resp.GroupDescription.String(),
		GroupImage:       resp.GroupDescription.String(),
		MembersCount:     resp.MembersCount,
		IsMember:         resp.IsMember,
		IsOwner:          resp.IsOwner,
		IsPending:        resp.IsPending,
	}, nil
}

func (s *Server) GetGroupMembers(ctx context.Context, req *pb.GroupMembersRequest) (*pb.GroupUserArr, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetGroupMembers: request is nil")
	}
	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	groupId := req.GetGroupId()
	if err := invalidId("groupId", groupId); err != nil {
		return nil, err
	}

	limit := req.Limit
	offset := req.Offset
	if err := checkLimOff(limit, offset); err != nil {
		return nil, err
	}

	resp, err := s.Application.GetGroupMembers(ctx, application.GroupMembersReq{
		UserId:  ct.Id(userId),
		GroupId: ct.Id(groupId),
		Limit:   ct.Limit(limit),
		Offset:  ct.Offset(offset),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetGroupMembers: %v", err)
	}
	return groupUsersToPB(resp), nil
}

func (s *Server) SearchGroups(ctx context.Context, req *pb.GroupSearchRequest) (*pb.GroupArr, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "SearchGroups: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}
	search := req.SearchTerm
	if err := invalidString("search", search); err != nil {
		return nil, err
	}

	limit := req.Limit
	offset := req.Offset
	if err := checkLimOff(limit, offset); err != nil {
		return nil, err
	}

	resp, err := s.Application.SearchGroups(ctx, application.GroupSearchReq{
		UserId:     ct.Id(userId),
		SearchTerm: ct.SearchTerm(search),
		Limit:      ct.Limit(limit),
		Offset:     ct.Offset(offset),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "SearchGroups: %v", err)
	}

	return groupsToPb(resp), nil
}

func (s *Server) InviteToGroup(ctx context.Context, req *pb.InviteToGroupRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "InviteToGroup: request is nil")
	}

	inviterId := req.InviterId
	if err := invalidId("inviterId", inviterId); err != nil {
		return nil, err
	}

	invitedId := req.InvitedId
	if err := invalidId("invitedId", invitedId); err != nil {
		return nil, err
	}

	groupId := req.GroupId
	if err := invalidId("groupId", groupId); err != nil {
		return nil, err
	}

	err := s.Application.InviteToGroup(ctx, application.InviteToGroupReq{
		InviterId: ct.Id(inviterId),
		InvitedId: ct.Id(invitedId),
		GroupId:   ct.Id(groupId),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "InviteToGroup: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) RequestJoinGroupOrCancel(ctx context.Context, req *pb.GroupJoinRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "RequestJoinGroupOrCancel: request is nil")
	}

	groupId := req.GroupId
	if err := invalidId("groupId", groupId); err != nil {
		return nil, err
	}

	requesterId := req.RequesterId
	if err := invalidId("RequesterId", requesterId); err != nil {
		return nil, err
	}

	err := s.Application.RequestJoinGroupOrCancel(ctx, application.GroupJoinRequest{
		GroupId:     ct.Id(groupId),
		RequesterId: ct.Id(requesterId),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "RequestJoinGroupOrCancel: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) RespondToGroupInvite(ctx context.Context, req *pb.HandleGroupInviteRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "RespondToGroupInvite: request is nil")
	}

	groupId := req.GroupId
	if err := invalidId("groupId", groupId); err != nil {
		return nil, err
	}

	InvitedId := req.InvitedId
	if err := invalidId("InvitedId", InvitedId); err != nil {
		return nil, err
	}

	acc := req.Accepted

	err := s.Application.RespondToGroupInvite(ctx, application.HandleGroupInviteRequest{
		GroupId:   ct.Id(groupId),
		InvitedId: ct.Id(InvitedId),
		Accepted:  acc,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "RespondToGroupInvite: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) HandleGroupJoinRequest(ctx context.Context, req *pb.HandleJoinRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "HandleGroupJoinRequest: request is nil")
	}

	groupId := req.GroupId
	if err := invalidId("groupId", groupId); err != nil {
		return nil, err
	}

	RequesterId := req.RequesterId
	if err := invalidId("RequesterId", RequesterId); err != nil {
		return nil, err
	}

	ownerId := req.OwnerId
	if err := invalidId("ownerId", ownerId); err != nil {
		return nil, err
	}

	acc := req.Accepted

	err := s.Application.HandleGroupJoinRequest(ctx, application.HandleJoinRequest{
		GroupId:     ct.Id(groupId),
		RequesterId: ct.Id(RequesterId),
		OwnerId:     ct.Id(ownerId),
		Accepted:    acc,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "HandleGroupJoinRequest: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) LeaveGroup(ctx context.Context, req *pb.GeneralGroupRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "LeaveGroup: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	groupId := req.GetGroupId()
	if err := invalidId("groupId", groupId); err != nil {
		return nil, err
	}

	err := s.Application.LeaveGroup(ctx, application.GeneralGroupReq{
		UserId:  ct.Id(userId),
		GroupId: ct.Id(groupId),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "LeaveGroup: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*wrapperspb.Int64Value, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "CreateGroup: request is nil")
	}

	OwnerId := req.OwnerId
	if err := invalidId("owner", OwnerId); err != nil {
		return nil, err
	}

	GroupTitle := req.GroupTitle
	if err := invalidString("GroupTitle", GroupTitle); err != nil {
		return nil, err
	}

	GroupDescription := req.GroupDescription
	if err := invalidString("GroupDescription", GroupDescription); err != nil {
		return nil, err
	}

	GroupImage := req.GroupImage
	if err := invalidString("GroupImage", GroupImage); err != nil {
		return nil, err
	}

	resp, err := s.Application.CreateGroup(ctx, application.CreateGroupRequest{
		OwnerId:          ct.Id(OwnerId),
		GroupTitle:       ct.Title(GroupTitle),
		GroupDescription: ct.About(GroupDescription),
		GroupImage:       GroupImage,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateGroup: %v", err)
	}
	return wrapperspb.Int64(int64(resp)), nil
}

// PROFILE
func (s *Server) GetBasicUserInfo(ctx context.Context, req *wrapperspb.Int64Value) (*pb.User, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userId := req.GetValue()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	u, err := s.Application.GetBasicUserInfo(ctx, ct.Id(req.GetValue()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetBasicUserInfo: %v", err)
	}

	return &pb.User{
		UserId:   u.UserId.Int64(),
		Username: u.Username.String(),
		Avatar:   u.Avatar,
	}, nil
}

func (s *Server) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfileResponse, error) {
	fmt.Println("GetUserProfile gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}
	RequesterId := req.GetRequesterId()
	if err := invalidId("RequesterId", RequesterId); err != nil {
		return nil, err
	}

	userProfileRequest := application.UserProfileRequest{
		UserId:      ct.Id(req.GetUserId()),
		RequesterId: ct.Id(req.GetRequesterId()),
	}

	profile, err := s.Application.GetUserProfile(ctx, userProfileRequest)
	if err != nil {
		fmt.Println("Error in GetUserProfile:", err)
		return nil, status.Errorf(codes.Internal, "GetUserProfile: %v", err)
	}

	return &pb.UserProfileResponse{
		UserId:   profile.UserId.Int64(),
		Username: profile.Username.String(),
	}, nil
}

func (s *Server) SearchUsers(ctx context.Context, req *pb.UserSearchRequest) (*pb.ListUsers, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "SearchUsers: request is nil")
	}

	SearchTerm := req.SearchTerm
	if err := invalidString("SearchTerm", SearchTerm); err != nil {
		return nil, err
	}

	limit := req.Limit
	if err := checkLimOff(limit, 1); err != nil {
		return nil, err
	}

	resp, err := s.Application.SearchUsers(ctx, application.UserSearchReq{
		SearchTerm: ct.SearchTerm(SearchTerm),
		Limit:      ct.Limit(limit),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "SearchUsers: %v", err)
	}
	return usersToPB(resp), nil
}

func (s *Server) UpdateUserProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "UpdateUserProfile: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	resp, err := s.Application.UpdateUserProfile(ctx, application.UpdateProfileRequest{
		UserId:      ct.Id(userId),
		Username:    ct.Username(req.GetUsername()),
		FirstName:   ct.Name(req.GetFirstName()),
		LastName:    ct.Name(req.GetLastName()),
		DateOfBirth: ct.DateOfBirth(req.GetDateOfBirth().AsTime()),
		Avatar:      req.GetAvatar(),
		About:       ct.About(req.GetAbout()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UpdateUserProfile: %v", err)
	}

	dob := timestamppb.New(resp.DateOfBirth.Time())
	if resp.DateOfBirth.Time().IsZero() {
		dob = nil
	}

	return &pb.UserProfileResponse{
		UserId:      resp.UserId.Int64(),
		Username:    resp.Username.String(),
		FirstName:   resp.FirstName.String(),
		LastName:    resp.LastName.String(),
		DateOfBirth: dob,
		Avatar:      resp.Avatar,
		About:       resp.About.String(),
		Public:      resp.Public,
	}, nil
}

func (s *Server) UpdateProfilePrivacy(ctx context.Context, req *pb.UpdateProfilePrivacyRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "UpdateProfilePrivacy: request is nil")
	}

	userId := req.GetUserId()
	if err := invalidId("userId", userId); err != nil {
		return nil, err
	}

	public := req.Public

	err := s.Application.UpdateProfilePrivacy(ctx, application.UpdateProfilePrivacyRequest{
		UserId: ct.Id(userId),
		Public: public,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UpdateProfilePrivacy: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// CONVERTORS
func usersToPB(dbUsers []application.User) *pb.ListUsers {
	pbUsers := make([]*pb.User, 0, len(dbUsers))

	for _, u := range dbUsers {
		pbUsers = append(pbUsers, &pb.User{
			UserId:   u.UserId.Int64(),
			Username: u.Username.String(),
			Avatar:   u.Avatar,
		})
	}

	return &pb.ListUsers{Users: pbUsers}
}

func groupsToPb(groups []application.Group) *pb.GroupArr {
	pbGroups := make([]*pb.Group, 0, len(groups))
	for _, g := range groups {
		pbGroups = append(pbGroups, &pb.Group{
			GroupId:          g.GroupId.Int64(),
			GroupOwnerId:     g.GroupOwnerId.Int64(),
			GroupTitle:       g.GroupTitle.String(),
			GroupDescription: g.GroupDescription.String(),
			GroupImage:       g.GroupImage,
			MembersCount:     g.MembersCount,
			IsMember:         g.IsMember,
			IsOwner:          g.IsOwner,
			IsPending:        g.IsPending,
		})
	}

	return &pb.GroupArr{
		GroupArr: pbGroups,
	}
}

func groupUsersToPB(users []application.GroupUser) *pb.GroupUserArr {
	out := &pb.GroupUserArr{
		GroupUserArr: make([]*pb.GroupUser, 0, len(users)),
	}

	for _, u := range users {
		out.GroupUserArr = append(out.GroupUserArr, &pb.GroupUser{
			UserId:    u.UserId.Int64(),
			Username:  u.Username.String(),
			Avatar:    u.Avatar,
			GroupRole: u.GroupRole,
		})
	}

	return out
}

func invalidId(varName string, value int64) error {
	if value <= 0 {
		pc, _, _, ok := runtime.Caller(1)
		funcName := "unknown"
		if ok {
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				funcName = fn.Name()
			}
		}

		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("[%s] variable: %v, value: %v must be larger than zero", funcName, varName, value),
		)
	}
	return nil
}

func invalidString(varName string, value string) error {
	if value == "" {
		pc, _, _, ok := runtime.Caller(1)
		funcName := "unknown"
		if ok {
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				funcName = fn.Name()
			}
		}

		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("[%s] variable: %v, value: %v must be non empty", funcName, varName, value),
		)
	}
	return nil

}

func checkLimOff(limit, offset int32) error {
	pc, _, _, ok := runtime.Caller(1)
	funcName := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			funcName = fn.Name()
		}
	}

	var maxLimit int32 = 100
	if limit > maxLimit {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("[%s] limit value: %v must be less than %v", funcName, limit, maxLimit),
		)
	}

	if offset < 0 {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("[%s] offset value: %v must be larger than 0", funcName, offset),
		)
	}
	return nil
}

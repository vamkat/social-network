/*
Expose methods via gRpc
*/

package handler

import (
	"context"
	"fmt"
	cm "social-network/shared/gen-go/common"
	pb "social-network/shared/gen-go/posts"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// POSTS

func (s *PostsHandler) GetPostById(ctx context.Context, req *pb.GenericReq) (*pb.Post, error) {
	fmt.Println("GetPostById gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	post, err := s.Application.GetPostById(ctx, models.GenericReq{
		RequesterId: ct.Id(req.RequesterId),
		EntityId:    ct.Id(req.EntityId),
	})
	if err != nil {
		fmt.Println("Error in GetPostById:", err)
		return nil, status.Errorf(codes.Internal, "failed to get post: %v", err)
	}
	return &pb.Post{
		PostId:   int64(post.PostId),
		PostBody: string(post.Body),
		User: &cm.User{
			UserId:   post.User.UserId.Int64(),
			Username: post.User.Username.String(),
			Avatar:   post.User.AvatarId.Int64(),
		},
		GroupId:         int64(post.GroupId),
		Audience:        post.Audience.String(),
		CommentsCount:   int32(post.CommentsCount),
		ReactionsCount:  int32(post.ReactionsCount),
		LastCommentedAt: post.LastCommentedAt.ToProto(),
		CreatedAt:       post.CreatedAt.ToProto(),
		UpdatedAt:       post.UpdatedAt.ToProto(),
		LikedByUser:     post.LikedByUser,
		Image:           int64(post.Image),
	}, nil
}

func (s *PostsHandler) CreatePost(ctx context.Context, req *pb.CreatePostReq) (*emptypb.Empty, error) {
	fmt.Println("CreatePost gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	err := s.Application.CreatePost(ctx, models.CreatePostReq{
		CreatorId:   ct.Id(req.CreatorId),
		Body:        ct.PostBody(req.Body),
		GroupId:     ct.Id(req.GroupId),
		Audience:    ct.Audience(req.Audience),
		AudienceIds: ct.FromInt64s(req.AudienceIds.Values),
		Image:       ct.Id(req.Image),
	})
	if err != nil {
		fmt.Println("Error in CreatePost:", err)
		return nil, status.Errorf(codes.Internal, "failed to create post: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) DeletePost(ctx context.Context, req *pb.GenericReq) (*emptypb.Empty, error) {
	fmt.Println("DeletePost gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	err := s.Application.DeletePost(ctx, models.GenericReq{
		RequesterId: ct.Id(req.RequesterId),
		EntityId:    ct.Id(req.EntityId),
	})
	if err != nil {
		fmt.Println("Error in DeletePost:", err)
		return nil, status.Errorf(codes.Internal, "failed to delete post: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) EditPost(ctx context.Context, req *pb.EditPostReq) (*emptypb.Empty, error) {
	fmt.Println("EditPost gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	err := s.Application.EditPost(ctx, models.EditPostReq{
		RequesterId: ct.Id(req.RequesterId),
		PostId:      ct.Id(req.PostId),
		NewBody:     ct.PostBody(req.Body),
		Image:       ct.Id(req.Image),
		Audience:    ct.Audience(req.Audience),
		AudienceIds: ct.FromInt64s(req.AudienceIds.Values),
	})
	if err != nil {
		fmt.Println("Error in EditPost:", err)
		return nil, status.Errorf(codes.Internal, "failed to edit post: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) GetMostPopularPostInGroup(ctx context.Context, req *pb.SimpleIdReq) (*pb.Post, error) {
	fmt.Println("GetMostPopularPostInGroup gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	post, err := s.Application.GetMostPopularPostInGroup(ctx, models.SimpleIdReq{
		Id: ct.Id(req.Id),
	})
	if err != nil {
		fmt.Println("Error in GetMostPopularPostInGroup:", err)
		return nil, status.Errorf(codes.Internal, "failed to get post: %v", err)
	}
	return &pb.Post{
		PostId:   int64(post.PostId),
		PostBody: string(post.Body),
		User: &cm.User{
			UserId:   post.User.UserId.Int64(),
			Username: post.User.Username.String(),
			Avatar:   post.User.AvatarId.Int64(),
		},
		GroupId:         int64(post.GroupId),
		Audience:        post.Audience.String(),
		CommentsCount:   int32(post.CommentsCount),
		ReactionsCount:  int32(post.ReactionsCount),
		LastCommentedAt: post.LastCommentedAt.ToProto(),
		CreatedAt:       post.CreatedAt.ToProto(),
		UpdatedAt:       post.UpdatedAt.ToProto(),
		LikedByUser:     post.LikedByUser,
		Image:           int64(post.Image),
	}, nil
}

func (s *PostsHandler) GetPersonalizedFeed(ctx context.Context, req *pb.GetPersonalizedFeedReq) (*pb.ListPosts, error) {
	fmt.Println("GetPersonalizedFeed gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	posts, err := s.Application.GetPersonalizedFeed(ctx, models.GetPersonalizedFeedReq{
		RequesterId: ct.Id(req.RequesterId),
		Limit:       ct.Limit(req.Limit),
		Offset:      ct.Offset(req.Offset),
	})
	if err != nil {
		fmt.Println("Error in GetPersonalizedFeed:", err)
		return nil, status.Errorf(codes.Internal, "failed to get personalized feed: %v", err)
	}
	pbPosts := make([]*pb.Post, 0, len(posts))
	for _, p := range posts {
		pbPosts = append(pbPosts, &pb.Post{
			PostId:   int64(p.PostId),
			PostBody: string(p.Body),
			User: &cm.User{
				UserId:   p.User.UserId.Int64(),
				Username: p.User.Username.String(),
				Avatar:   p.User.AvatarId.Int64(),
			},
			GroupId:         int64(p.GroupId),
			Audience:        p.Audience.String(),
			CommentsCount:   int32(p.CommentsCount),
			ReactionsCount:  int32(p.ReactionsCount),
			LastCommentedAt: p.LastCommentedAt.ToProto(),
			CreatedAt:       p.CreatedAt.ToProto(),
			UpdatedAt:       p.UpdatedAt.ToProto(),
			LikedByUser:     p.LikedByUser,
			Image:           int64(p.Image),
		})
	}
	return &pb.ListPosts{Posts: pbPosts}, nil
}

func (s *PostsHandler) GetPublicFeed(ctx context.Context, req *pb.GenericPaginatedReq) (*pb.ListPosts, error) {
	fmt.Println("GetPublicFeed gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	posts, err := s.Application.GetPublicFeed(ctx, models.GenericPaginatedReq{
		RequesterId: ct.Id(req.RequesterId),
		Limit:       ct.Limit(req.Limit),
		Offset:      ct.Offset(req.Offset),
	})
	if err != nil {
		fmt.Println("Error in GetPublicFeed:", err)
		return nil, status.Errorf(codes.Internal, "failed to get public feed: %v", err)
	}
	pbPosts := make([]*pb.Post, 0, len(posts))
	for _, p := range posts {
		pbPosts = append(pbPosts, &pb.Post{
			PostId:   int64(p.PostId),
			PostBody: string(p.Body),
			User: &cm.User{
				UserId:   p.User.UserId.Int64(),
				Username: p.User.Username.String(),
				Avatar:   p.User.AvatarId.Int64(),
			},
			GroupId:         int64(p.GroupId),
			Audience:        p.Audience.String(),
			CommentsCount:   int32(p.CommentsCount),
			ReactionsCount:  int32(p.ReactionsCount),
			LastCommentedAt: p.LastCommentedAt.ToProto(),
			CreatedAt:       p.CreatedAt.ToProto(),
			UpdatedAt:       p.UpdatedAt.ToProto(),
			LikedByUser:     p.LikedByUser,
			Image:           int64(p.Image),
		})
	}
	return &pb.ListPosts{Posts: pbPosts}, nil
}

func (s *PostsHandler) GetUserPostsPaginated(ctx context.Context, req *pb.GetUserPostsReq) (*pb.ListPosts, error) {
	fmt.Println("GetUserPostsPaginated gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	posts, err := s.Application.GetUserPostsPaginated(ctx, models.GetUserPostsReq{
		CreatorId:   ct.Id(req.CreatorId),
		RequesterId: ct.Id(req.RequesterId),
		Limit:       ct.Limit(req.Limit),
		Offset:      ct.Offset(req.Offset),
	})
	if err != nil {
		fmt.Println("Error in GetUserPostsPaginated:", err)
		return nil, status.Errorf(codes.Internal, "failed to get user posts: %v", err)
	}
	pbPosts := make([]*pb.Post, 0, len(posts))
	for _, p := range posts {
		pbPosts = append(pbPosts, &pb.Post{
			PostId:   int64(p.PostId),
			PostBody: string(p.Body),
			User: &cm.User{
				UserId:   p.User.UserId.Int64(),
				Username: p.User.Username.String(),
				Avatar:   p.User.AvatarId.Int64(),
			},
			GroupId:         int64(p.GroupId),
			Audience:        p.Audience.String(),
			CommentsCount:   int32(p.CommentsCount),
			ReactionsCount:  int32(p.ReactionsCount),
			LastCommentedAt: p.LastCommentedAt.ToProto(),
			CreatedAt:       p.CreatedAt.ToProto(),
			UpdatedAt:       p.UpdatedAt.ToProto(),
			LikedByUser:     p.LikedByUser,
			Image:           int64(p.Image),
		})
	}
	return &pb.ListPosts{Posts: pbPosts}, nil
}

func (s *PostsHandler) GetGroupPostsPaginated(ctx context.Context, req *pb.GetGroupPostsReq) (*pb.ListPosts, error) {
	fmt.Println("GetGroupPostsPaginated gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	posts, err := s.Application.GetGroupPostsPaginated(ctx, models.GetGroupPostsReq{
		GroupId:     ct.Id(req.GroupId),
		RequesterId: ct.Id(req.RequesterId),
		Limit:       ct.Limit(req.Limit),
		Offset:      ct.Offset(req.Offset),
	})
	if err != nil {
		fmt.Println("Error in GetGroupPostsPaginated:", err)
		return nil, status.Errorf(codes.Internal, "failed to get group posts: %v", err)
	}
	pbPosts := make([]*pb.Post, 0, len(posts))
	for _, p := range posts {
		pbPosts = append(pbPosts, &pb.Post{
			PostId:   int64(p.PostId),
			PostBody: string(p.Body),
			User: &cm.User{
				UserId:   p.User.UserId.Int64(),
				Username: p.User.Username.String(),
				Avatar:   p.User.AvatarId.Int64(),
			},
			GroupId:         int64(p.GroupId),
			Audience:        p.Audience.String(),
			CommentsCount:   int32(p.CommentsCount),
			ReactionsCount:  int32(p.ReactionsCount),
			LastCommentedAt: p.LastCommentedAt.ToProto(),
			CreatedAt:       p.CreatedAt.ToProto(),
			UpdatedAt:       p.UpdatedAt.ToProto(),
			LikedByUser:     p.LikedByUser,
			Image:           int64(p.Image),
		})
	}
	return &pb.ListPosts{Posts: pbPosts}, nil
}

func (s *PostsHandler) CreateComment(ctx context.Context, req *pb.CreateCommentReq) (*emptypb.Empty, error) {
	fmt.Println("CreateComment gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	err := s.Application.CreateComment(ctx, models.CreateCommentReq{
		CreatorId: ct.Id(req.CreatorId),
		ParentId:  ct.Id(req.ParentId),
		Body:      ct.CommentBody(req.Body),
		Image:     ct.Id(req.Image),
	})
	if err != nil {
		fmt.Println("Error in CreateComment:", err)
		return nil, status.Errorf(codes.Internal, "failed to create comment: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) EditComment(ctx context.Context, req *pb.EditCommentReq) (*emptypb.Empty, error) {
	fmt.Println("EditComment gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	err := s.Application.EditComment(ctx, models.EditCommentReq{
		CreatorId: ct.Id(req.CreatorId),
		CommentId: ct.Id(req.CommentId),
		Body:      ct.CommentBody(req.Body),
		Image:     ct.Id(req.Image),
	})
	if err != nil {
		fmt.Println("Error in EditComment:", err)
		return nil, status.Errorf(codes.Internal, "failed to edit comment: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) DeleteComment(ctx context.Context, req *pb.GenericReq) (*emptypb.Empty, error) {
	fmt.Println("DeleteComment gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	err := s.Application.DeleteComment(ctx, models.GenericReq{
		RequesterId: ct.Id(req.RequesterId),
		EntityId:    ct.Id(req.EntityId),
	})
	if err != nil {
		fmt.Println("Error in DeleteComment:", err)
		return nil, status.Errorf(codes.Internal, "failed to delete comment: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) GetCommentsByParentId(ctx context.Context, req *pb.EntityIdPaginatedReq) (*pb.ListComments, error) {
	fmt.Println("GetCommentsByParentId gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	comments, err := s.Application.GetCommentsByParentId(ctx, models.EntityIdPaginatedReq{
		RequesterId: ct.Id(req.RequesterId),
		EntityId:    ct.Id(req.EntityId),
		Limit:       ct.Limit(req.Limit),
		Offset:      ct.Offset(req.Offset),
	})
	if err != nil {
		fmt.Println("Error in GetCommentsByParentId:", err)
		return nil, status.Errorf(codes.Internal, "failed to get comments: %v", err)
	}
	pbComments := make([]*pb.Comment, 0, len(comments))
	for _, c := range comments {
		pbComments = append(pbComments, &pb.Comment{
			CommentId: int64(c.CommentId),
			ParentId:  int64(c.ParentId),
			Body:      string(c.Body),
			User: &cm.User{
				UserId:   c.User.UserId.Int64(),
				Username: c.User.Username.String(),
				Avatar:   c.User.AvatarId.Int64(),
			},
			ReactionsCount: int32(c.ReactionsCount),
			CreatedAt:      c.CreatedAt.ToProto(),
			UpdatedAt:      c.UpdatedAt.ToProto(),
			LikedByUser:    c.LikedByUser,
			Image:          int64(c.Image),
		})
	}
	return &pb.ListComments{Comments: pbComments}, nil
}

func (s *PostsHandler) CreateEvent(ctx context.Context, req *pb.CreateEventReq) (*emptypb.Empty, error) {
	fmt.Println("CreateEvent gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	err := s.Application.CreateEvent(ctx, models.CreateEventReq{
		Title:     ct.Title(req.Title),
		Body:      ct.EventBody(req.Body),
		CreatorId: ct.Id(req.CreatorId),
		GroupId:   ct.Id(req.GroupId),
		Image:     ct.Id(req.Image),
		EventDate: ct.EventDateTime(req.EventDate.AsTime()),
	})
	if err != nil {
		fmt.Println("Error in CreateEvent:", err)
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) DeleteEvent(ctx context.Context, req *pb.GenericReq) (*emptypb.Empty, error) {
	fmt.Println("DeleteEvent gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	err := s.Application.DeleteEvent(ctx, models.GenericReq{
		RequesterId: ct.Id(req.RequesterId),
		EntityId:    ct.Id(req.EntityId),
	})
	if err != nil {
		fmt.Println("Error in DeleteEvent:", err)
		return nil, status.Errorf(codes.Internal, "failed to delete event: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) EditEvent(ctx context.Context, req *pb.EditEventReq) (*emptypb.Empty, error) {
	fmt.Println("EditEvent gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	err := s.Application.EditEvent(ctx, models.EditEventReq{
		EventId:     ct.Id(req.EventId),
		RequesterId: ct.Id(req.RequesterId),
		Title:       ct.Title(req.Title),
		Body:        ct.EventBody(req.Body),
		Image:       ct.Id(req.Image),
		EventDate:   ct.EventDateTime(req.EventDate.AsTime()),
	})
	if err != nil {
		fmt.Println("Error in EditEvent:", err)
		return nil, status.Errorf(codes.Internal, "failed to edit event: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) GetEventsByGroupId(ctx context.Context, req *pb.EntityIdPaginatedReq) (*pb.ListEvents, error) {
	fmt.Println("GetEventsByGroupId gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	events, err := s.Application.GetEventsByGroupId(ctx, models.EntityIdPaginatedReq{
		RequesterId: ct.Id(req.RequesterId),
		EntityId:    ct.Id(req.EntityId),
		Limit:       ct.Limit(req.Limit),
		Offset:      ct.Offset(req.Offset),
	})
	if err != nil {
		fmt.Println("Error in GetEventsByGroupId:", err)
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}
	pbEvents := make([]*pb.Event, 0, len(events))
	for _, e := range events {
		var ur *wrapperspb.BoolValue
		if e.UserResponse != nil {
			ur = wrapperspb.Bool(*e.UserResponse)
		}

		pbEvents = append(pbEvents, &pb.Event{
			EventId: int64(e.EventId),
			Title:   string(e.Title),
			Body:    string(e.Body),
			User: &cm.User{
				UserId:   e.User.UserId.Int64(),
				Username: e.User.Username.String(),
				Avatar:   e.User.AvatarId.Int64(),
			},
			GroupId:       int64(e.GroupId),
			EventDate:     e.EventDate.ToProto(),
			GoingCount:    int32(e.GoingCount),
			NotGoingCount: int32(e.NotGoingCount),
			Image:         int64(e.Image),
			CreatedAt:     e.CreatedAt.ToProto(),
			UpdatedAt:     e.UpdatedAt.ToProto(),
			UserResponse:  ur,
		})
	}
	return &pb.ListEvents{Events: pbEvents}, nil
}

func (s *PostsHandler) RespondToEvent(ctx context.Context, req *pb.RespondToEventReq) (*emptypb.Empty, error) {
	fmt.Println("RespondToEvent gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	err := s.Application.RespondToEvent(ctx, models.RespondToEventReq{
		EventId:     ct.Id(req.EventId),
		ResponderId: ct.Id(req.ResponderId),
		Going:       req.Going,
	})
	if err != nil {
		fmt.Println("Error in RespondToEvent:", err)
		return nil, status.Errorf(codes.Internal, "failed to respond to event: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) RemoveEventResponse(ctx context.Context, req *pb.GenericReq) (*emptypb.Empty, error) {
	fmt.Println("RemoveEventResponse gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	err := s.Application.RemoveEventResponse(ctx, models.GenericReq{
		RequesterId: ct.Id(req.RequesterId),
		EntityId:    ct.Id(req.EntityId),
	})
	if err != nil {
		fmt.Println("Error in RemoveEventResponse:", err)
		return nil, status.Errorf(codes.Internal, "failed to remove event response: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *PostsHandler) SuggestUsersByPostActivity(ctx context.Context, req *pb.SimpleIdReq) (*cm.ListUsers, error) {
	fmt.Println("SuggestUsersByPostActivity gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	users, err := s.Application.SuggestUsersByPostActivity(ctx, models.SimpleIdReq{
		Id: ct.Id(req.Id),
	})
	if err != nil {
		fmt.Println("Error in SuggestUsersByPostActivity:", err)
		return nil, status.Errorf(codes.Internal, "failed to suggest users: %v", err)
	}
	pbUsers := make([]*cm.User, 0, len(users))
	for _, u := range users {
		pbUsers = append(pbUsers, &cm.User{
			UserId:   u.UserId.Int64(),
			Username: u.Username.String(),
			Avatar:   u.AvatarId.Int64(),
		})
	}
	return &cm.ListUsers{Users: pbUsers}, nil
}

func (s *PostsHandler) ToggleOrInsertReaction(ctx context.Context, req *pb.GenericReq) (*emptypb.Empty, error) {
	fmt.Println("ToggleOrInsertReaction gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	err := s.Application.ToggleOrInsertReaction(ctx, models.GenericReq{
		RequesterId: ct.Id(req.RequesterId),
		EntityId:    ct.Id(req.EntityId),
	})
	if err != nil {
		fmt.Println("Error in ToggleOrInsertReaction:", err)
		return nil, status.Errorf(codes.Internal, "failed to react to post: %v", err)
	}
	return &emptypb.Empty{}, nil
}

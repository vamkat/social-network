package handler

import (
	"context"
	"social-network/services/chat/internal/application"
	pb "social-network/shared/gen-go/chat"
	"social-network/shared/go/customtypes"
	"social-network/shared/go/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/lib/pq"
)

type ChatHandler struct {
	pb.UnimplementedChatServiceServer
	Application *application.ChatService
	Port        string
}

func (h *ChatHandler) CreatePrivateConversation(ctx context.Context, params *pb.CreatePrivateConvParams) (*pb.ConvId, error) {
	convId, err := h.Application.CreatePrivateConversation(ctx, models.CreatePrivateConvParams{
		UserA: customtypes.Id(params.UserA),
		UserB: customtypes.Id(params.UserB),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create private conversation %v", err)
	}
	return &pb.ConvId{ConvId: convId}, nil
}

func (h *ChatHandler) CreateGroupConversation(ctx context.Context, params *pb.CreateGroupConvParams) (*pb.ConvId, error) {
	convId, err := h.Application.CreateGroupConversation(ctx, models.CreateGroupConvParams{
		GroupId: customtypes.Id(params.GroupId),
		UserIds: customtypes.FromInt64s(params.UserIds),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create private conversation %v", err)
	}
	return &pb.ConvId{ConvId: convId}, nil
}

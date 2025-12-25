package handler

import (
	"context"

	pb "social-network/shared/gen-go/chat"
	"social-network/shared/go/ct"
	"social-network/shared/go/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/lib/pq"
)

// Creates a Conversation if and only if a conversation between the same 2 users does not exist.
// Returns NULL ConvId if a duplicate DM exists (sqlc will error if RETURNING finds no rows).
func (h *ChatHandler) CreatePrivateConversation(ctx context.Context, params *pb.CreatePrivateConvParams) (*pb.ConvId, error) {
	convId, err := h.Application.CreatePrivateConversation(ctx, models.CreatePrivateConvParams{
		UserA: ct.Id(params.UserA),
		UserB: ct.Id(params.UserB),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create private conversation %v", err)
	}
	return &pb.ConvId{ConvId: convId.Int64()}, nil
}

func (h *ChatHandler) CreateGroupConversation(ctx context.Context, params *pb.CreateGroupConvParams) (*pb.ConvId, error) {
	convId, err := h.Application.CreateGroupConversation(ctx, models.CreateGroupConvParams{
		GroupId: ct.Id(params.GroupId),
		UserIds: ct.FromInt64s(params.UserIds),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create private conversation %v", err)
	}
	return &pb.ConvId{ConvId: convId.Int64()}, nil
}

// Delete a conversation only if its members exactly match the provided list.
// Returns 0 rows if conversation doesn't exist, members don’t match exactly, conversation has extra or missing members.
func (h *ChatHandler) DeleteConversationByExactMembers(ctx context.Context, userIds *pb.UserIds) (*pb.Conversation, error) {
	resp, err := h.Application.DeleteConversationByExactMembers(ctx, ct.FromInt64s(userIds.UserIds))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete conversation %v", err)
	}
	return &pb.Conversation{
		Id:        resp.Id.Int64(),
		GroupId:   resp.GroupId.Int64(),
		CreatedAt: resp.CreatedAt.ToProto(),
		UpdatedAt: resp.UpdatedAt.ToProto(),
		DeletedAt: resp.DeletedAt.ToProto(),
	}, nil
}

// GetUserConversations

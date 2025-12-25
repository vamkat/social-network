package handler

import (
	"context"
	"social-network/shared/gen-go/chat"
	"social-network/shared/go/ct"
	"social-network/shared/go/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Find a conversation by group_id and insert the given user_ids into conversation_members.
// Existing members are ignored, new members are added.
func (h *ChatHandler) AddMembersToGroupConversation(ctx context.Context, params *chat.AddMembersToGroupConversationParams) (*chat.ConvId, error) {
	resp, err := h.Application.AddMembersToGroupConversation(ctx, models.AddMembersToGroupConversationParams{
		GroupId: ct.Id(params.GroupId),
		UserIds: ct.FromInt64s(params.UserIds),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add members %v", err)
	}
	return &chat.ConvId{ConvId: resp.Int64()}, nil
}

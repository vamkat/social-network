package application

import (
	"context"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

func (c *ChatService) AddMembersToGroupConversation(ctx context.Context,
	arg md.AddMembersToGroupConversationParams) (convId ct.Id, err error) {
	if err := ct.ValidateStruct(arg); err != nil {
		return convId, err
	}
	return c.Queries.AddMembersToGroupConversation(ctx, arg)
}

package application

import (
	"context"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

// Returns models.User type for each conversation member of 'GetConversationMembersParams.ConversationId'
// except 'GetConversationMembersParams.UserId' .
func (c *ChatService) GetConversationMembers(ctx context.Context,
	params md.GetConversationMembersParams) (members []md.User, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return members, err
	}
	ids, err := c.Queries.GetConversationMembers(ctx, params)
	if err != nil {
		return members, err
	}
	return c.Clients.UserIdsToUsers(ctx, ids)
}

func (c *ChatService) DeleteConversationMember(ctx context.Context,
	params md.DeleteConversationMemberParams,
) (dltMember md.ConversationMemberDeleted, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return dltMember, err
	}
	return c.Queries.DeleteConversationMember(ctx, params)
}

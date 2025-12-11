package application

import (
	"context"
	"database/sql"
	"fmt"
	"social-network/services/chat/internal/db/dbservice"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

func (c *ChatService) CreatePrivateConversation(ctx context.Context,
	params md.CreatePrivateConvParams) (convId ct.Id, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return 0, err
	}

	convId, err = c.Queries.CreatePrivateConv(ctx, params)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("conversation already exists")
	}
	return convId, err
}

func (c *ChatService) CreateGroupConversation(ctx context.Context,
	params md.CreateGroupConvParams) (convId ct.Id, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return 0, err
	}

	err = c.txRunner.RunTx(ctx,
		func(q dbservice.Querier) error {
			convId, err = q.CreateGroupConv(ctx, params.GroupId)
			if err != nil {
				return err
			}

			return q.AddConversationMembers(ctx,
				md.AddConversationMembersParams{
					ConversationId: ct.Id(convId),
					UserIds:        params.UserIds,
				})
		})
	return ct.Id(convId), err
}

// Delete a conversation only if its members exactly match the provided list.
// Returns 0 rows if conversation doesn't exist, members don’t match exactly, conversation has extra or missing members.
func (c *ChatService) DeleteConversationByExactMembers(ctx context.Context,
	ids ct.Ids) (conv md.ConversationDeleteResp, err error) {
	if err := ids.Validate(); err != nil {
		return conv, err
	}
	conv, err = c.Queries.DeleteConversationByExactMembers(ctx, ids)
	if err != nil {
		return conv, err
	}

	if conv == (md.ConversationDeleteResp{}) {
		err = fmt.Errorf("conversation not found")
	}
	return conv, err
}

// Find a conversation by group_id and insert the given user_ids into conversation_members.
// existing members are ignored, new members are added.
func (c *ChatService) AddMembersToGroupConversation(ctx context.Context,
	params md.AddMembersToGroupConversationParams) (convId ct.Id, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return 0, err
	}
	return c.Queries.AddMembersToGroupConversation(ctx, params)
}

func (c *ChatService) GetUserConversations(ctx context.Context,
	arg md.GetUserConversationsParams) ([]md.GetUserConversationsRow, error) {
	if err := ct.ValidateStruct(arg); err != nil {
		return nil, err
	}
	return c.Queries.GetUserConversations(ctx, arg)
}

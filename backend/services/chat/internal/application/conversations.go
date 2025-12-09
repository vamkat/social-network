package application

import (
	"context"
	"database/sql"
	"fmt"
	"social-network/services/chat/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func (c *ChatService) CreatePrivateConversation(ctx context.Context, params models.CreatePrivateConvParams) (convId int64, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return 0, err
	}

	convId, err = c.Queries.CreatePrivateConv(ctx, sqlc.CreatePrivateConvParams{UserID: params.UserA.Int64(), UserID_2: params.UserB.Int64()})
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("conversation already exists")
	}
	return convId, err
}

func (c *ChatService) CreateGroupConversation(ctx context.Context, params models.CreateGroupConvParams) (convId int64, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return 0, err
	}

	err = c.txRunner.RunTx(ctx,
		func(q sqlc.Querier) error {
			convId, err = c.Queries.CreateGroupConv(ctx, pgtype.Int8{
				Int64: params.GroupId.Int64(),
				Valid: true,
			})
			if err != nil {
				return err
			}

			return c.Queries.AddConversationMembers(ctx,
				sqlc.AddConversationMembersParams{
					ConversationID: convId,
					UserIds:        params.UserIds.Int64(),
				})
		})
	return convId, err
}

// Delete a conversation only if its members exactly match the provided list.
// Returns 0 rows if conversation doesn't exist, members don’t match exactly, conversation has extra or missing members.
func (c *ChatService) DeleteConversationByExactMembers(ctx context.Context, ids ct.Ids) (conv models.Conversation, err error) {
	if err := ids.Validate(); err != nil {
		return models.Conversation{}, err
	}
	resp, err := c.Queries.DeleteConversationByExactMembers(ctx, ids.Int64())
	if err != nil {
		return models.Conversation{}, err
	}
	return models.Conversation{
		ID:        ct.Id(resp.ID),
		GroupID:   ct.Id(resp.GroupID.Int64),
		CreatedAt: resp.CreatedAt.Time,
		UpdatedAt: resp.UpdatedAt.Time,
		DeletedAt: resp.DeletedAt.Time,
	}, nil
}

// Find a conversation by group_id and insert the given user_ids into conversation_members.
// existing members are ignored, new members are added.
func (c *ChatService) AddMembersToGroupConversation(ctx context.Context, params models.AddMembersToGroupConversationParams) (convId ct.Id, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return 0, err
	}

	arg := sqlc.AddMembersToGroupConversationParams{
		GroupID: pgtype.Int8{
			Int64: params.GroupID.Int64(),
			Valid: true,
		},
		UserIds: params.UserIds.Int64(),
	}

	resp, err := c.Queries.AddMembersToGroupConversation(ctx, arg)
	if err != nil {
		return 0, err
	}
	return ct.Id(resp), nil
}

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

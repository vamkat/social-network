package application

import (
	"context"
	"social-network/services/chat/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"

	"github.com/jackc/pgx/v5/pgtype"
)

/*
create conversation (Users Ids, grouId (nullable))
	min 1 ids if groupId == 0 minIds 2
addOrRemove conversation member (grouId, userId, bool)
*/

type CreatePrivateConvParams struct {
	UserA ct.Id
	UserB ct.Id
}

func (c *ChatService) CreatePrivateConversation(ctx context.Context, params CreatePrivateConvParams) (convId int64, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return 0, err
	}
	return c.Queries.CreatePrivateConv(ctx, sqlc.CreatePrivateConvParams{UserID: params.UserA.Int64(), UserID_2: params.UserB.Int64()})
}

type CreateGroupConvParams struct {
	GroupId ct.Id
	UserIds ct.Ids
}

func (c *ChatService) CreateGroupConversation(ctx context.Context, params CreateGroupConvParams) (convId int64, err error) {
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

package dbservice

import (
	"context"
	"errors"

	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"

	"github.com/jackc/pgx/v5"
)

// Find a conversation by group_id and insert the given user_ids into conversation_members.
// existing members are ignored, new members are added.
// func (q *Queries) AddMembersToGroupConversation(ctx context.Context,
// 	arg md.AddMembersToGroupConversationParams) (convId ct.Id, err error) {
// 	row := q.db.QueryRow(ctx,
// 		addMembersToGroupConversation,
// 		arg.GroupId,
// 		arg.UserIds,
// 	)
// 	err = row.Scan(&convId)
// 	return convId, err
// }

type AddMembersResult struct {
	ConversationID ct.Id
	RequestedCount int
	InsertedCount  int
}

func (q *Queries) AddMembersToGroupConversation(
	ctx context.Context,
	arg md.AddMembersToGroupConversationParams,
) (ct.Id, error) {
	var res AddMembersResult

	row := q.db.QueryRow(ctx, addMembersToGroupConversation,
		arg.GroupId,
		arg.UserIds,
	)

	err := row.Scan(
		&res.ConversationID,
		&res.RequestedCount,
		&res.InsertedCount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ce.Wrap(ce.ErrNotFound, errors.New("group conversation not found"))
		}
		return 0, ce.Wrap(ce.ErrInternal, errors.New("failed to add members"))
	}

	return res.ConversationID, nil
}

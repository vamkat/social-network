package dbservice

import (
	"context"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

// Find a conversation by group_id and insert the given user_ids into conversation_members.
// existing members are ignored, new members are added.
func (q *Queries) AddMembersToGroupConversation(ctx context.Context,
	arg md.AddMembersToGroupConversationParams) (convId ct.Id, err error) {
	row := q.db.QueryRow(ctx,
		addMembersToGroupConversation,
		arg.GroupId,
		arg.UserIds,
	)
	err = row.Scan(&convId)
	return convId, err
}

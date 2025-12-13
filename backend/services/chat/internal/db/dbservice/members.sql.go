package dbservice

import (
	"context"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

// Returns members of a conversation that user is a member.
func (q *Queries) GetConversationMembers(ctx context.Context,
	arg md.GetConversationMembersParams) (members ct.Ids, err error) {
	rows, err := q.db.Query(ctx,
		getConversationMembers,
		arg.ConversationId,
		arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members = ct.Ids{}
	for rows.Next() {
		var user_id int64
		if err := rows.Scan(&user_id); err != nil {
			return nil, err
		}
		members = append(members, ct.Id(user_id))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

// Deletes conversation member from conversation where user tagged as owner is a part of.
// Returnes user deleted details. If no rows returned means no deletation occured.
// Can be used for self deletation if owner and toDelete are that same id.
func (q *Queries) DeleteConversationMember(ctx context.Context,
	arg md.DeleteConversationMemberParams,
) (dltMember md.ConversationMemberDeleted, err error) {
	row := q.db.QueryRow(ctx, deleteConversationMember,
		arg.ConversationID, arg.ToDelete, arg.Owner)
	err = row.Scan(
		&dltMember.ConversationId,
		&dltMember.UserId,
		&dltMember.LastReadMessageId,
		&dltMember.JoinedAt,
		&dltMember.DeletedAt,
	)
	return dltMember, err
}

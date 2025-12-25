package dbservice

import (
	"context"
	"fmt"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

// Initiates the conversation by groupId. Group Id should be a not null value
// Use as a preparation for adding members
func (q *Queries) CreateGroupConv(ctx context.Context,
	groupId ct.Id) (convId ct.Id, err error) {
	if !groupId.IsValid() {
		return convId, fmt.Errorf("null group Id")
	}
	row := q.db.QueryRow(ctx, createGroupConv, groupId)
	err = row.Scan(&convId)
	return convId, err
}

// Add UserIDs to ConvID
func (q *Queries) AddConversationMembers(ctx context.Context,
	arg md.AddConversationMembersParams) error {
	_, err := q.db.Exec(ctx, addConversationMembers, arg.ConversationId, arg.UserIds)
	return err
}

// Creates a Conversation if and only if a conversation between the same 2 users does not exist.
// Returns NULL if a duplicate DM exists (sql will error if RETURNING finds no rows).
func (q *Queries) CreatePrivateConv(ctx context.Context,
	arg md.CreatePrivateConvParams,
) (convId ct.Id, err error) {
	row := q.db.QueryRow(ctx, createPrivateConv, arg.UserA, arg.UserB)
	err = row.Scan(&convId)
	return convId, err
}

// Delete a conversation only if its members exactly match the provided list.
// Returns 0 rows if conversation doesn't exist, members don’t match exactly, conversation has extra or missing members.
func (q *Queries) DeleteConversationByExactMembers(ctx context.Context,
	memberIds ct.Ids) (md.ConversationDeleteResp, error) {
	row := q.db.QueryRow(ctx, deleteConversationByExactMembers, memberIds)
	var i md.ConversationDeleteResp
	err := row.Scan(
		&i.Id,
		&i.GroupId,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

type GetUserConversationsRow struct {
	ConversationId    ct.Id
	CreatedAt         ct.GenDateTime
	UpdatedAt         ct.GenDateTime
	MemberIds         ct.Ids
	UnreadCount       int64
	LastReadMessageId ct.Id `validation:"nullable"`
}

// Fetches paginated conversation details, conversation members Ids
// and unread messages count for a user and a group.
// To get DMS group Id parameter must be zero.
func (q *Queries) GetUserConversations(
	ctx context.Context,
	arg md.GetUserConversationsParams,
) (conversations []GetUserConversationsRow, err error) {

	rows, err := q.db.Query(ctx,
		getUserConversations,
		arg.UserId,
		arg.Limit,
		arg.Offset,
		arg.GroupId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i GetUserConversationsRow

		err := rows.Scan(
			&i.ConversationId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.MemberIds,
			&i.UnreadCount,
			&i.LastReadMessageId,
		)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

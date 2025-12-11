package dbservice

import (
	"context"
	"fmt"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

const createGroupConv = `
INSERT INTO conversations (group_id)
VALUES ($1)
RETURNING id
`

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

const createPrivateConv = `
WITH existing AS (
    SELECT c.id
    FROM conversations c
    JOIN conversation_members cm1 ON cm1.conversation_id = c.id AND cm1.user_id = $1
    JOIN conversation_members cm2 ON cm2.conversation_id = c.id AND cm2.user_id = $2
    WHERE c.group_id IS NULL
)
INSERT INTO conversations (group_id)
SELECT NULL
WHERE NOT EXISTS (SELECT 1 FROM existing)
RETURNING id
`

// Creates a Conversation if and only if a conversation between the same 2 users does not exist.
// Returns NULL if a duplicate DM exists (sql will error if RETURNING finds no rows).
func (q *Queries) CreatePrivateConv(ctx context.Context,
	arg md.CreatePrivateConvParams,
) (convId ct.Id, err error) {
	row := q.db.QueryRow(ctx, createPrivateConv, arg.UserA, arg.UserB)
	err = row.Scan(&convId)
	return convId, err
}

const deleteConversationByExactMembers = `
WITH target_members AS (
    SELECT unnest($1::bigint[]) AS user_id
),
matched_conversation AS (
    SELECT cm.conversation_id
    FROM conversation_members cm
    JOIN target_members tm ON tm.user_id = cm.user_id
    GROUP BY cm.conversation_id
    HAVING 
        -- same count of overlapping members
        COUNT(*) = (SELECT COUNT(*) FROM target_members)
        -- and the conversation has no extra members
        AND COUNT(*) = (
            SELECT COUNT(*) 
            FROM conversation_members cm2 
            WHERE cm2.conversation_id = cm.conversation_id
              AND cm2.deleted_at IS NULL
        )
)
UPDATE conversations c
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE c.id = (SELECT conversation_id FROM matched_conversation)
RETURNING id, group_id, created_at, updated_at, deleted_at
`

// Delete a conversation only if its members exactly match the provided list.
// Returns 0 rows if conversation doesn't exist, members don’t match exactly, conversation has extra or missing members.
// OK!
func (q *Queries) DeleteConversationByExactMembers(ctx context.Context, memberIds ct.Ids) (md.ConversationDeleteResp, error) {
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

const getUserPrivateConversations = `
WITH user_conversations AS (
    SELECT c.id AS conversation_id,
		   c.created_at,
           c.updated_at,
           c.last_message_id,
           cm.last_read_message_id
    FROM conversations c
    JOIN conversation_members cm
        ON cm.conversation_id = c.id
    WHERE cm.user_id = $1
      AND cm.deleted_at IS NULL
      AND c.group_id IS NOT DISTINCT FROM $4
    ORDER BY c.last_message_id DESC
    LIMIT $2 OFFSET $3
),

member_list AS (
    SELECT uc.conversation_id,
           json_agg(cm.user_id) FILTER (WHERE cm.user_id != $1) AS member_ids
    FROM user_conversations uc
    JOIN conversation_members cm
        ON cm.conversation_id = uc.conversation_id
    GROUP BY uc.conversation_id
),

unread AS (
    SELECT uc.conversation_id,
           COUNT(m.id) AS unread_count,
           MIN(m.id) AS first_unread_message_id
    FROM user_conversations uc
    LEFT JOIN messages m
      ON m.conversation_id = uc.conversation_id
     AND m.id > COALESCE(uc.last_read_message_id, 0)
     AND m.deleted_at IS NULL
    GROUP BY uc.conversation_id
)

SELECT
	uc.conversation_id,
	uc.created_at,
	uc.updated_at,
	ml.member_ids,
	u.unread_count,
	u.first_unread_message_id
FROM user_conversations uc
JOIN member_list ml ON ml.conversation_id = uc.conversation_id
LEFT JOIN unread u ON u.conversation_id = uc.conversation_id
ORDER BY uc.last_message_id DESC;
`

// Fetches paginated conversation details, conversation members Ids and unread messages count for a user and a group
// To get DMS group Id parameter must be zero.
// OK !
func (q *Queries) GetUserConversations(
	ctx context.Context,
	arg md.GetUserConversationsParams,
) ([]md.GetUserConversationsRow, error) {

	rows, err := q.db.Query(ctx,
		getUserPrivateConversations,
		arg.UserId,
		arg.Limit,
		arg.Offset,
		arg.GroupId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []md.GetUserConversationsRow

	for rows.Next() {
		var i md.GetUserConversationsRow

		err := rows.Scan(
			&i.ConversationId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.MemberIds,
			&i.UnreadCount,
			&i.FirstUnreadMessageId,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

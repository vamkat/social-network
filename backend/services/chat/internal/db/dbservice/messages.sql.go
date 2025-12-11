package dbservice

import (
	"context"
	md "social-network/shared/go/models"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (conversation_id, sender_id, message_text)
SELECT $1, $2, $3
FROM conversation_members
WHERE conversation_id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING id, conversation_id, sender_id, message_text, created_at, updated_at, deleted_at
`

func (q *Queries) CreateMessage(ctx context.Context, arg md.CreateMessageParams) (md.MessageResp, error) {
	row := q.db.QueryRow(ctx, createMessage, arg.ConversationId, arg.SenderId, arg.MessageText)
	var i md.MessageResp
	err := row.Scan(
		&i.Id,
		&i.ConversationID,
		&i.SenderID,
		&i.MessageText,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getMessages = `-- name: GetMessages :many
SELECT m.id, m.conversation_id, m.sender_id, m.message_text, m.created_at, m.updated_at, m.deleted_at
FROM messages m
JOIN conversation_members cm 
  ON cm.conversation_id = m.conversation_id
WHERE m.conversation_id = $1
  AND cm.user_id = $2
  AND m.deleted_at IS NULL
ORDER BY m.created_at ASC
LIMIT $3 OFFSET $4
`

func (q *Queries) GetMessages(ctx context.Context, arg md.GetMessagesParams) (messages []md.MessageResp, err error) {
	rows, err := q.db.Query(ctx, getMessages,
		arg.ConversationId,
		arg.UserId,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages = []md.MessageResp{}
	for rows.Next() {
		var i md.MessageResp
		if err := rows.Scan(
			&i.Id,
			&i.ConversationID,
			&i.SenderID,
			&i.MessageText,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

const updateLastReadMessage = `-- name: UpdateLastReadMessage :one
UPDATE conversation_members cm
SET last_read_message_id = $3
WHERE cm.conversation_id = $1
  AND cm.user_id = $2
  AND cm.deleted_at IS NULL
RETURNING conversation_id, user_id, last_read_message_id, joined_at, deleted_at
`

func (q *Queries) UpdateLastReadMessage(ctx context.Context,
	arg md.UpdateLastReadMessageParams,
) (md.ConversationMember, error) {
	row := q.db.QueryRow(ctx, updateLastReadMessage, arg.ConversationId, arg.UserID, arg.LastReadMessageId)
	var i md.ConversationMember
	err := row.Scan(
		&i.ConversationID,
		&i.UserID,
		&i.LastReadMessageID,
		&i.JoinedAt,
		&i.DeletedAt,
	)
	return i, err
}

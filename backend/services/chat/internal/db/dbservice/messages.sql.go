package dbservice

import (
	"context"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

// Creates a message row with conversation id if user is a memeber.
// If user match of conversation id and user id fails no rows are returned.
func (q *Queries) CreateMessage(ctx context.Context,
	arg md.CreateMessageParams) (msg md.MessageResp, err error) {
	row := q.db.QueryRow(ctx, createMessage, arg.ConversationId, arg.SenderId, arg.MessageText)
	err = row.Scan(
		&msg.Id,
		&msg.ConversationID,
		&msg.Sender.UserId,
		&msg.MessageText,
		&msg.CreatedAt,
		&msg.UpdatedAt,
		&msg.DeletedAt,
	)
	return msg, err
}

// Returns a descending-ordered page of messages that appear chronologically
// BEFORE a given message in a conversation. This query is used for backwards
// pagination in chat history.
//
// Behavior:
//
//   - If the supplied BoundaryMessageId is NULL, the query automatically
//     substitutes the conversation's last_message_id as the boundary (inclusive).
//
//   - The caller must be a member of the conversation. Membership is enforced
//     through the conversation_members table.
//
//   - Results are ordered by m.id DESC so that the most recent messages in the
//     requested page appear last. LIMIT/OFFSET is applied after ordering.
//
// Returned fields:
//   - All message fields (id, conversation_id, sender_id, message_text, timestamps)
//   - Conversation's first_message_id
//
// Use case:
//
//	Scroll-up pagination in chat history.
func (q *Queries) GetPreviousMessages(ctx context.Context,
	args md.GetPrevMessagesParams) (resp md.GetPrevMessagesResp, err error) {
	rows, err := q.db.Query(ctx, getPrevMessages,
		args.BoundaryMessageId,
		args.ConversationId,
		args.UserId,
		args.Limit,
	)
	if err != nil {
		return resp, err
	}
	for rows.Next() {
		var msg md.MessageResp
		var firstMessageId ct.Id

		rows.Scan(
			&msg.Id,
			&msg.ConversationID,
			&msg.Sender.UserId,
			&msg.MessageText,
			&msg.CreatedAt,
			&msg.UpdatedAt,
			&msg.DeletedAt,
			&firstMessageId,
		)

		resp.FirstMessageId = firstMessageId
		resp.Messages = append(resp.Messages, msg)
	}
	return resp, nil
}

// Returns an ascending-ordered page of messages that appear chronologically
// AFTER a given message in a conversation. This query is used for forward
// pagination when loading newer messages.
//
// Behavior:
//
//   - If the supplied BoundaryMessageId ($1) is NULL, the query automatically
//     substitutes the conversation's first_message_id as the boundary.
//
//   - Only messages with id > boundary_id are returned.
//
//   - Only non-deleted messages are returned (deleted_at IS NULL).
//
//   - The caller must be a member of the conversation. Membership is enforced
//     through the conversation_members table.
//
//   - Results are ordered by m.id ASC so that the oldest messages in the
//     requested page appear first. LIMIT/OFFSET is applied after ordering.
//
// Returned fields:
//   - All message fields (id, conversation_id, sender_id, message_text, timestamps)
//   - Conversation's last_message_id (constant for all rows)
//
// Use case:
//
//	Scroll-down pagination or loading new messages after a known point.
func (q *Queries) GetNextMessages(ctx context.Context,
	args md.GetNextMessageParams,
) (resp md.GetNextMessagesResp, err error) {
	rows, err := q.db.Query(ctx, getNextMessages,
		args.BoundaryMessageId,
		args.ConversationId,
		args.UserId,
		args.Limit,
	)
	if err != nil {
		return resp, err
	}
	for rows.Next() {
		var msg md.MessageResp
		var lastMessageId ct.Id

		rows.Scan(
			&msg.Id,
			&msg.ConversationID,
			&msg.Sender.UserId,
			&msg.MessageText,
			&msg.CreatedAt,
			&msg.UpdatedAt,
			&msg.DeletedAt,
			&lastMessageId,
		)

		resp.LastMessageId = lastMessageId
		resp.Messages = append(resp.Messages, msg)
	}
	return resp, nil
}

// Updates the given users last read message in given conversation to given message id.
func (q *Queries) UpdateLastReadMessage(ctx context.Context,
	arg md.UpdateLastReadMessageParams,
) (member md.ConversationMember, err error) {
	row := q.db.QueryRow(ctx, updateLastReadMessage, arg.ConversationId, arg.UserID, arg.LastReadMessageId)
	err = row.Scan(
		&member.ConversationId,
		&member.UserId,
		&member.LastReadMessageId,
		&member.JoinedAt,
		&member.DeletedAt,
	)
	return member, err
}

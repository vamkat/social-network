package dbservice

import (
	"context"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

type Querier interface {
	// Add UserIDs to ConvID.
	AddConversationMembers(ctx context.Context, arg md.AddConversationMembersParams) error

	// Find a conversation by group_id and insert the given user_ids into conversation_members.
	// existing members are ignored, new members are added.
	AddMembersToGroupConversation(ctx context.Context, arg md.AddMembersToGroupConversationParams) (convId ct.Id, err error)

	// Initiates the conversation by groupId. Group Id should be a not null value.
	// Use as a preparation for adding members
	CreateGroupConv(ctx context.Context, groupID ct.Id) (convId ct.Id, err error)

	// Creates a message row with conversation id if user is a memeber.
	// Returns error if user match of conversation_id and user_id fails.
	CreateMessage(ctx context.Context, arg md.CreateMessageParams) (md.MessageResp, error)

	// Returns a descending-ordered page of messages that appear chronologically
	// BEFORE a given message in a conversation. This query is used for backwards
	// pagination in chat history.
	//
	// Behavior:
	//
	//   - If the supplied FirstMessageId is NULL, the query automatically
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
	GetPreviousMessages(ctx context.Context,
		args md.GetPrevMessagesParams) (resp md.GetPrevMessagesResp, err error)

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
	GetNextMessages(ctx context.Context,
		args md.GetNextMessageParams,
	) (resp md.GetNextMessagesResp, err error)

	// Creates a Conversation if and only if a conversation between the same 2 users does not exist.
	// Returns NULL if a duplicate DM exists (sqlc will error if RETURNING finds no rows).
	CreatePrivateConv(ctx context.Context, arg md.CreatePrivateConvParams) (convId ct.Id, err error)

	// Delete a conversation only if its members exactly match the provided list.
	// Returns "conversation not found" if conversation doesn't exist,
	// members don’t match exactly, conversation has extra or missing members.
	DeleteConversationByExactMembers(ctx context.Context, memberIds ct.Ids) (md.ConversationDeleteResp, error)

	// Returns memebers of a conversation that user is a member.
	GetConversationMembers(ctx context.Context, arg md.GetConversationMembersParams) (members ct.Ids, err error)

	// Fetches paginated conversation details, conversation members Ids and unread messages count for a user and a group.
	// To get DMS group Id parameter must be zero.
	GetUserConversations(ctx context.Context, arg md.GetUserConversationsParams) ([]GetUserConversationsRow, error)

	// Deletes conversation member from conversation where user tagged as owner is a part of.
	// Returns user deleted details.
	// Can be used for self deletation if owner and toDelete are that same id.
	DeleteConversationMember(ctx context.Context,
		arg md.DeleteConversationMemberParams) (md.ConversationMemberDeleted, error)

	// Updates the given users last read message in given conversation to given message id.
	UpdateLastReadMessage(ctx context.Context, arg md.UpdateLastReadMessageParams) (md.ConversationMember, error)
}

var _ Querier = (*Queries)(nil)

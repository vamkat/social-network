package dbservice

import (
	"context"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

type Querier interface {
	AddConversationMembers(ctx context.Context, arg md.AddConversationMembersParams) error

	// Find a conversation by group_id and insert the given user_ids into conversation_members.
	// existing members are ignored, new members are added.
	AddMembersToGroupConversation(ctx context.Context, arg md.AddMembersToGroupConversationParams) (convId ct.Id, err error)
	CreateGroupConv(ctx context.Context, groupID ct.Id) (convId ct.Id, err error)
	CreateMessage(ctx context.Context, arg md.CreateMessageParams) (md.MessageResp, error)

	// Creates a Conversation if and only if a conversation between the same 2 users does not exist.
	// Returns NULL if a duplicate DM exists (sqlc will error if RETURNING finds no rows).
	CreatePrivateConv(ctx context.Context, arg md.CreatePrivateConvParams) (convId ct.Id, err error)

	// Delete a conversation only if its members exactly match the provided list.
	// Returns 0 rows if conversation doesn't exist, members don’t match exactly, conversation has extra or missing members.
	DeleteConversationByExactMembers(ctx context.Context, memberIds ct.Ids) (md.ConversationDeleteResp, error)

	GetConversationMembers(ctx context.Context, arg md.GetConversationMembersParams) (members ct.Ids, err error)
	GetMessages(ctx context.Context, arg md.GetMessagesParams) (messages []md.MessageResp, err error)

	// Fetches paginated conversation details, conversation members Ids and unread messages count for a user and a group
	// To get DMS group Id parameter must be null.
	GetUserConversations(ctx context.Context, arg md.GetUserConversationsParams) ([]md.GetUserConversationsRow, error)
	SoftDeleteConversationMember(ctx context.Context,
		arg md.SoftDeleteConversationMemberParams) (md.ConversationMember, error)
	UpdateLastReadMessage(ctx context.Context, arg md.UpdateLastReadMessageParams) (md.ConversationMember, error)
}

var _ Querier = (*Queries)(nil)

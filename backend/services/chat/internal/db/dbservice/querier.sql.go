package dbservice

import (
	"context"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

type Querier interface {
	// Add UserIDs to ConvID.
	AddConversationMembers(ctx context.Context, arg md.AddConversationMembersParams) error

	// Find a conversation by group_id and insert the given user_ids into conversation_members.
	// existing members are ignored, new members are added.
	AddMembersToGroupConversation(ctx context.Context, arg md.AddMembersToGroupConversationParams) (convId ct.Id, err error)

	// Initiates the conversation by groupId. Group Id should be a not null value
	// Use as a preparation for adding members
	CreateGroupConv(ctx context.Context, groupID ct.Id) (convId ct.Id, err error)

	// Creates a message row with conversation id if user is a memeber.
	// Returns error if user match of conversation_id and user_id fails.
	CreateMessage(ctx context.Context, arg md.CreateMessageParams) (md.MessageResp, error)

	// Creates a Conversation if and only if a conversation between the same 2 users does not exist.
	// Returns NULL if a duplicate DM exists (sqlc will error if RETURNING finds no rows).
	CreatePrivateConv(ctx context.Context, arg md.CreatePrivateConvParams) (convId ct.Id, err error)

	// Delete a conversation only if its members exactly match the provided list.
	// Returns "conversation not found" if conversation doesn't exist,
	// members don’t match exactly, conversation has extra or missing members.
	DeleteConversationByExactMembers(ctx context.Context, memberIds ct.Ids) (md.ConversationDeleteResp, error)

	// Returns memebers of a conversation that user is a member.
	GetConversationMembers(ctx context.Context, arg md.GetConversationMembersParams) (members ct.Ids, err error)

	GetMessages(ctx context.Context, arg md.GetMessagesParams) (messages []md.MessageResp, err error)

	// Fetches paginated conversation details, conversation members Ids and unread messages count for a user and a group
	// To get DMS group Id parameter must be zero.
	GetUserConversations(ctx context.Context, arg md.GetUserConversationsParams) ([]md.GetUserConversationsRow, error)

	// Deletes conversation member from conversation where user tagged as owner is a part of.
	// Returns user deleted details.
	// Can be used for self deletation if owner and toDelete are that same id.
	DeleteConversationMember(ctx context.Context,
		arg md.DeleteConversationMemberParams) (md.ConversationMemberDeleted, error)

	UpdateLastReadMessage(ctx context.Context, arg md.UpdateLastReadMessageParams) (md.ConversationMember, error)
}

var _ Querier = (*Queries)(nil)

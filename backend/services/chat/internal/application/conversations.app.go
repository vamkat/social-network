package application

import (
	"context"
	"database/sql"
	"fmt"
	"social-network/services/chat/internal/db/dbservice"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

// Creates a Conversation if and only if a conversation between the same 2 users does not exist.
// Returns NULL if a duplicate DM exists (sql will error if RETURNING finds no rows).
func (c *ChatService) CreatePrivateConversation(ctx context.Context,
	params md.CreatePrivateConvParams) (convId ct.Id, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return 0, err
	}

	convId, err = c.Queries.CreatePrivateConv(ctx, params)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("conversation already exists")
	}
	return convId, err
}

func (c *ChatService) CreateGroupConversation(ctx context.Context,
	params md.CreateGroupConvParams) (convId ct.Id, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return 0, err
	}

	err = c.txRunner.RunTx(ctx,
		func(q dbservice.Querier) error {
			convId, err = q.CreateGroupConv(ctx, params.GroupId)
			if err != nil {
				return err
			}

			return q.AddConversationMembers(ctx,
				md.AddConversationMembersParams{
					ConversationId: ct.Id(convId),
					UserIds:        params.UserIds,
				})
		})
	return ct.Id(convId), err
}

// Delete a conversation only if its members exactly match the provided list.
// Returns 0 rows if conversation doesn't exist, members don’t match exactly, conversation has extra or missing members.
func (c *ChatService) DeleteConversationByExactMembers(ctx context.Context,
	ids ct.Ids) (conv md.ConversationDeleteResp, err error) {
	if err := ids.Validate(); err != nil {
		return conv, err
	}
	conv, err = c.Queries.DeleteConversationByExactMembers(ctx, ids)
	if err != nil {
		return conv, err
	}

	if conv == (md.ConversationDeleteResp{}) {
		err = fmt.Errorf("conversation not found")
	}
	return conv, err
}

// Fetches paginated conversation details, conversation members Ids and unread messages count for a user and a group
// To get DMS group Id parameter must be zero.
func (c *ChatService) GetUserConversations(ctx context.Context,
	arg md.GetUserConversationsParams,
) (conversations []md.GetUserConversationsResp, err error) {
	if err := ct.ValidateStruct(arg); err != nil {
		return nil, err
	}

	resp, err := c.Queries.GetUserConversations(ctx, arg)
	if err != nil {
		return conversations, err
	}

	allMemberIDs := make(ct.Ids, 0)
	for _, r := range resp {
		allMemberIDs = append(allMemberIDs, r.MemberIds...)
	}

	usersMap, err := c.Clients.UserIdsToMap(ctx, allMemberIDs)
	if err != nil {
		return nil, err
	}
	return ConvertConversations(ctx, usersMap, resp)
}

// Helper to convert a slice of GetUserConversationsRow containing userIds
// to a slice of GetUserConversationsResp containg User.
func ConvertConversations(
	ctx context.Context,
	usersMap map[ct.Id]md.User,
	rows []dbservice.GetUserConversationsRow,
) ([]md.GetUserConversationsResp, error) {

	// Build output
	result := make([]md.GetUserConversationsResp, len(rows))
	for i, r := range rows {

		// Build members list for this conversation
		members := make([]md.User, 0, len(r.MemberIds))
		for _, mid := range r.MemberIds {
			if u, ok := usersMap[mid]; ok {
				members = append(members, u)
			}
		}

		result[i] = md.GetUserConversationsResp{
			ConversationId:    r.ConversationId,
			CreatedAt:         r.CreatedAt,
			UpdatedAt:         r.UpdatedAt,
			Members:           members,
			UnreadCount:       r.UnreadCount,
			LastReadMessageId: r.LastReadMessageId,
		}
	}

	return result, nil
}

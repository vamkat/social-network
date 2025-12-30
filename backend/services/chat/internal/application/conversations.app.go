package application

import (
	"context"
	"database/sql"
	"fmt"
	"social-network/services/chat/internal/db/dbservice"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

// Creates a Conversation if and only if a conversation between the same 2 users does not exist.
// Returns NULL if a duplicate DM exists (sql will error if RETURNING finds no rows).
func (c *ChatService) CreatePrivateConversation(ctx context.Context,
	params md.CreatePrivateConvParams) (convId ct.Id, err error) {

	errMsg := fmt.Sprintf("create private conversation: user ids: %d, %d", params.UserA, params.UserB)

	if err := ct.ValidateStruct(params); err != nil {
		return 0, ce.Wrap(ce.ErrInvalidArgument, err, errMsg)
	}

	convId, err = c.Queries.CreatePrivateConv(ctx, params)
	if err == sql.ErrNoRows {
		return 0, ce.Wrap(ce.ErrInvalidArgument, err, errMsg).WithPublic("conversation already exists")
	} else if err != nil {
		return 0, ce.Wrap(ce.ErrInternal, err, errMsg)
	}
	return convId, nil
}

func (c *ChatService) CreateGroupConversation(ctx context.Context,
	params md.CreateGroupConvParams) (convId ct.Id, err error) {

	errMsg := fmt.Sprintf("create group conversation: group id: %d, user ids: %d", params.GroupId, params.UserIds)

	if err := ct.ValidateStruct(params); err != nil {
		return 0, ce.Wrap(ce.ErrInvalidArgument, err, errMsg)
	}

	err = c.txRunner.RunTx(ctx,
		func(q *dbservice.Queries) error {
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
	errMsg := fmt.Sprintf("delete conversation by exact members: ids: %d", ids)

	if err := ids.Validate(); err != nil {
		return conv, ce.Wrap(ce.ErrInvalidArgument, err, errMsg)
	}

	conv, err = c.Queries.DeleteConversationByExactMembers(ctx, ids)
	if err != nil {
		return conv, ce.Wrap(ce.ErrInternal, err, errMsg)
	}

	if conv == (md.ConversationDeleteResp{}) {
		return conv, ce.Wrap(ce.ErrNotFound, fmt.Errorf("db response is empty"), errMsg)
	}

	return conv, nil
}

// Fetches paginated conversation details, conversation members Ids and unread messages count for a user and a group
// To get DMS group Id parameter must be zero.
// If hydrate users is true then each user in the 'Members' field is populated with username and avatar id
// by calling the user service client.
func (c *ChatService) GetUserConversations(ctx context.Context,
	arg md.GetUserConversationsParams,
) (conversations []md.GetUserConversationsResp, err error) {

	errMsg := fmt.Sprintf("get user conversations: arg: %#v", arg)

	if err := ct.ValidateStruct(arg); err != nil {
		return nil, ce.Wrap(ce.ErrInvalidArgument, err, errMsg)
	}

	resp, err := c.Queries.GetUserConversations(ctx, arg)
	if err != nil {
		return conversations, ce.Wrap(ce.ErrInternal, err, errMsg)
	}

	if !arg.HydrateUsers {
		// Calling with nil usersMap. No hydration just conversion
		return ConvertConversations(ctx, nil, resp)
	}

	allMemberIDs := make(ct.Ids, 0)
	for _, r := range resp {
		allMemberIDs = append(allMemberIDs, r.MemberIds...)
	}
	// TODO: Check error return
	usersMap, err := c.RetriveUsers.GetUsers(ctx, allMemberIDs)
	if err != nil {
		return nil, ce.Wrap(ce.ErrUnavailable, err, errMsg)
	}

	return ConvertConversations(ctx, usersMap, resp)
}

// Helper to convert a slice of GetUserConversationsRow containing userIds
// to a slice of GetUserConversationsResp containg User.
// If nil usersMap is passed then conversion does not hydrate Members with username and avatar.
func ConvertConversations(
	ctx context.Context,
	usersMap map[ct.Id]md.User,
	rows []dbservice.GetUserConversationsRow,
) ([]md.GetUserConversationsResp, error) {
	conversations := make([]md.GetUserConversationsResp, len(rows))

	// Convert unhydrated
	if usersMap == nil {
		for i, r := range rows {
			members := make([]md.User, 0, len(r.MemberIds))
			for _, mid := range r.MemberIds {
				members = append(members, md.User{UserId: mid})
			}
			conversations[i] = md.GetUserConversationsResp{
				ConversationId:    r.ConversationId,
				CreatedAt:         r.CreatedAt,
				UpdatedAt:         r.UpdatedAt,
				Members:           members,
				UnreadCount:       r.UnreadCount,
				LastReadMessageId: r.LastReadMessageId,
			}
		}
		return conversations, nil
	}

	// Convert with hydration from map
	for i, r := range rows {

		// Build members list for this conversation
		members := make([]md.User, 0, len(r.MemberIds))
		for _, mid := range r.MemberIds {
			if u, ok := usersMap[mid]; ok {
				members = append(members, u)
			}
		}

		conversations[i] = md.GetUserConversationsResp{
			ConversationId:    r.ConversationId,
			CreatedAt:         r.CreatedAt,
			UpdatedAt:         r.UpdatedAt,
			Members:           members,
			UnreadCount:       r.UnreadCount,
			LastReadMessageId: r.LastReadMessageId,
		}
	}

	return conversations, nil
}

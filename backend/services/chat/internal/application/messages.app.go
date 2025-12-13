package application

import (
	"context"
	"fmt"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

// Creates a message row with conversation id if user is a memeber.
// If user match of conversation_id and user_id fails returns error.
// TODO: Implement call to live service
func (c *ChatService) CreateMessage(ctx context.Context,
	params md.CreateMessageParams) (msg md.MessageResp, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return msg, err
	}
	if (msg == md.MessageResp{}) {
		return msg, fmt.Errorf("user is not a member of conversation id: %v", params.ConversationId)
	}

	_, err = c.Queries.CreateMessage(ctx, params)
	if err != nil {
		return msg, err
	}
	return msg, err
}

// Returns messages with id smaller that BoundaryMessageId. If BoundaryMessageId is null
// returns previous messages from and including conversation's last_message_id.
func (c *ChatService) GetPreviousMessages(ctx context.Context,
	args md.GetPrevMessagesParams) (resp md.GetPrevMessagesResp, err error) {
	if err != ct.ValidateStruct(args) {
		return resp, err
	}

	resp, err = c.Queries.GetPreviousMessages(ctx, args)
	if err != nil {
		return resp, err
	}

	if len(resp.Messages) == 0 {
		return md.GetPrevMessagesResp{
			HaveMoreBefore: false,
			Messages:       resp.Messages,
		}, nil
	}

	if resp.FirstMessageId == resp.Messages[0].Id {
		resp.HaveMoreBefore = false
	}

	allMemberIDs := make(ct.Ids, 0)
	for _, r := range resp.Messages {
		allMemberIDs = append(allMemberIDs, r.Sender.UserId)
	}

	usersMap, err := c.Clients.UserIdsToMap(ctx, allMemberIDs)
	if err != nil {
		return resp, err // decide if should return struct with no user info or error ??
	}

	for _, s := range resp.Messages {
		s.Sender = usersMap[s.Sender.UserId]
	}
	return resp, nil
}

// Returns an ascending-ordered page of messages that appear chronologically
// AFTER a given message in a conversation. This query is used for forward
// pagination when loading newer messages.
// Can also be called with boundary as the last read message to get unread messages.
// Returns validation error if no boundary given.
func (c *ChatService) GetNextMessages(ctx context.Context,
	args md.GetNextMessageParams) (resp md.GetNextMessagesResp, err error) {
	if err != ct.ValidateStruct(args) {
		return resp, err
	}
	resp, err = c.Queries.GetNextMessages(ctx, args)
	if err != nil {
		return resp, err
	}
	if len(resp.Messages) == 0 {
		return md.GetNextMessagesResp{
			HaveMoreAfter: false,
			Messages:      resp.Messages,
		}, nil
	}

	lastIdx := len(resp.Messages) - 1
	if resp.LastMessageId == resp.Messages[lastIdx].Id {
		resp.HaveMoreAfter = false
	}

	allMemberIDs := make(ct.Ids, 0)
	for _, r := range resp.Messages {
		allMemberIDs = append(allMemberIDs, r.Sender.UserId)
	}

	usersMap, err := c.Clients.UserIdsToMap(ctx, allMemberIDs)
	if err != nil {
		return resp, err // decide if should return struct with no user info or error ??
	}

	for _, s := range resp.Messages {
		s.Sender = usersMap[s.Sender.UserId]
	}
	return resp, nil
}

// TODO: Implement call to live service
func (c *ChatService) UpdateLastReadMessage(ctx context.Context,
	params md.UpdateLastReadMessageParams,
) (member md.ConversationMember, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return member, err
	}
	if (member == md.ConversationMember{}) {
		return member, fmt.Errorf("user is not a member of conversation id: %v", params.ConversationId)
	}
	return c.Queries.UpdateLastReadMessage(ctx, params)
}

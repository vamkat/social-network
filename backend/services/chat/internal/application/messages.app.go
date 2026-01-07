package application

import (
	"context"
	"fmt"
	"social-network/services/chat/internal/db/dbservice"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

// Creates a message row with conversation id if user is a memeber.
// If user match of conversation_id and user_id fails returns error.
// TODO: Implement call to live service
func (c *ChatService) CreateMessage(ctx context.Context,
	params md.CreateMessageParams) (msg md.MessageResp, err error) {

	input := fmt.Sprintf("params: %#v", params)
	if err := ct.ValidateStruct(params); err != nil {
		return msg, ce.New(ce.ErrInvalidArgument, err, input)
	}

	msg, err = c.Queries.CreateMessageWithMembersJoin(ctx, params)

	if err != nil {
		return msg, ce.Wrap(nil, err, input)
	}
	return msg, err
}

type CreateMessageInGroupReq struct {
	GroupId     ct.Id
	SenderId    ct.Id
	MessageBody ct.MsgBody
}

func (c *ChatService) CreateMessageInGroup(ctx context.Context,
	params CreateMessageInGroupReq) (msg md.MessageResp, err error) {
	input := fmt.Sprintf("params: %#v", params)

	if err := ct.ValidateStruct(params); err != nil {
		return msg, ce.New(ce.ErrInvalidArgument, err, input)
	}

	// Call UserService to check if sender is a member of group
	isMember, err := c.Clients.IsGroupMember(ctx, params.GroupId, params.SenderId)
	if err != nil {
		return msg, ce.ParseGrpcErr(err, input)
	}
	if !isMember {
		return msg,
			ce.New(ce.ErrPermissionDenied,
				fmt.Errorf("user id: %d not a member of group id: %d", params.SenderId, params.GroupId),
				input,
			)
	}

	err = c.txRunner.RunTx(ctx,
		func(q *dbservice.Queries) error {
			// Create or get conversation
			convId, err := q.CreateGroupConv(ctx, params.GroupId)
			if err != nil {
				return ce.Wrap(nil, err, input)
			}
			// Add message
			msg, err = c.Queries.CreateMessage(ctx, md.CreateMessageParams{
				ConversationId: convId,
				SenderId:       params.SenderId,
				MessageText:    params.MessageBody,
			})

			if err != nil {
				return ce.Wrap(nil, err, input)
			}
			return nil
		})
	return msg, nil
}

// Returns messages with id smaller that BoundaryMessageId. If BoundaryMessageId is null
// returns previous messages from and including conversation's last_message_id.
// If hydrate users is true then the messages.Sender field is populated with username and avatar id
// by calling the user service client.
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

	if resp.FirstMessageId != resp.Messages[0].Id {
		resp.HaveMoreBefore = true
	}

	if !args.HydrateUsers {
		return resp, nil
	}

	allMemberIDs := make(ct.Ids, 0)
	for _, r := range resp.Messages {
		allMemberIDs = append(allMemberIDs, r.Sender.UserId)
	}

	usersMap, err := c.RetriveUsers.GetUsers(ctx, allMemberIDs)
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
// If hydrate users is true then the messages.Sender field is populated with username and avatar id
// by calling the user service client.
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
	if resp.LastMessageId != resp.Messages[lastIdx].Id {
		resp.HaveMoreAfter = true
	}

	if !args.RetrieveUsers {
		return resp, nil
	}

	allMemberIDs := make(ct.Ids, 0)
	for _, r := range resp.Messages {
		allMemberIDs = append(allMemberIDs, r.Sender.UserId)
	}

	usersMap, err := c.RetriveUsers.GetUsers(ctx, allMemberIDs)
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

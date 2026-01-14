package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"
)

var (
	ErrNotConnected = errors.New("users are not connected")
)

// Returns a sorted paginated list of private conversations
// older that the given BeforeDate where user with UserId is a member.
// Respose per PC includes last message and unread count from users side.
func (c *ChatService) GetPrivateConversations(ctx context.Context,
	arg md.GetPrivateConvsReq,
) ([]md.PrivateConvsPreview, *ce.Error) {

	input := fmt.Sprintf("arg: %#v", arg)

	err := ct.ValidateStruct(arg)
	if err != nil {
		return nil, ce.Wrap(ce.ErrInvalidArgument, err, input)
	}

	conversations, err := c.Queries.GetPrivateConvs(ctx, arg)
	if err != nil {
		return conversations, ce.Wrap(nil, err, input)
	}

	allMemberIDs := make(ct.Ids, 0)
	for _, r := range conversations {
		allMemberIDs = append(allMemberIDs, r.Interlocutor.UserId)
	}

	usersMap, err := c.RetriveUsers.GetUsers(ctx, allMemberIDs)
	if err != nil {
		return nil, ce.Wrap(nil, err, input)
	}

	for i := range conversations {
		retrieved := usersMap[conversations[i].Interlocutor.UserId]
		conversations[i].Interlocutor.Username = retrieved.Username
		conversations[i].Interlocutor.AvatarId = retrieved.AvatarId
		conversations[i].Interlocutor.AvatarURL = retrieved.AvatarURL
	}

	return conversations, nil
}

// Creates a private message and returns an id
func (c *ChatService) CreatePrivateMessage(ctx context.Context,
	params md.CreatePrivateMsgReq) (msg md.PrivateMsg, Err *ce.Error) {

	input := fmt.Sprintf("params: %#v", params)
	err := ct.ValidateStruct(params)
	if err != nil {
		return msg, ce.New(ce.ErrInvalidArgument, err, input)
	}

	msg, err = c.Queries.CreateNewPrivateMessage(ctx, params)

	if err != nil {
		return msg, ce.Wrap(nil, err, input)
	}

	messageBytes, err := json.Marshal(msg)
	if err != nil {
		err = ce.New(ce.ErrInternal, err, input)
		tele.Error(ctx, "failed to publish private message to nats: @1", "error", err.Error())
	}

	err = c.NatsConn.Publish(ct.PrivateMessageKey(params.InterlocutorId), messageBytes)
	if err != nil {
		err = ce.New(ce.ErrInternal, err, input)
		tele.Error(ctx, "failed to publish private message to nats: @1", "error", err.Error())
	}

	return msg, nil
}

func (c *ChatService) GetPreviousPMs(ctx context.Context,
	arg md.GetPrivateMsgsReq) (res md.GetPrivateMsgsResp, Err *ce.Error) {
	input := fmt.Sprintf("arg: %#v", arg)

	if err := ct.ValidateStruct(arg); err != nil {
		return res, ce.New(ce.ErrInvalidArgument, err, input)
	}

	res, err := c.Queries.GetPrevPrivateMsgs(ctx, arg)
	if err != nil {
		return res, ce.Wrap(nil, err, input)
	}

	if arg.RetrieveUsers {
		if err := c.retrievePrivatMessageSenders(ctx, res.Messages, input); err != nil {
			tele.Error(ctx, "failed to retrieve users for messages", "input", input, "error", err)
		}
	}

	return res, nil
}

func (c *ChatService) GetNextPMs(ctx context.Context,
	arg md.GetPrivateMsgsReq) (res md.GetPrivateMsgsResp, Err *ce.Error) {
	input := fmt.Sprintf("arg: %#v", arg)

	if err := ct.ValidateStruct(arg); err != nil {
		return res, ce.New(ce.ErrInvalidArgument, err, input)
	}

	res, err := c.Queries.GetNextPrivateMsgs(ctx, arg)
	if err != nil {
		return res, ce.Wrap(nil, err, input)
	}

	if arg.RetrieveUsers {
		if err := c.retrievePrivatMessageSenders(ctx, res.Messages, input); err != nil {
			tele.Error(ctx, "failed to retrieve users for messages", "input", input, "error", err)
		}
	}
	return res, nil
}

func (c *ChatService) UpdateLastReadPrivateMsg(ctx context.Context, arg md.UpdateLastReadMsgParams) *ce.Error {
	input := fmt.Sprintf("arg: %#v", arg)

	if err := ct.ValidateStruct(arg); err != nil {
		return ce.New(ce.ErrInvalidArgument, err, input)
	}

	err := c.Queries.UpdateLastReadPrivateMsg(ctx, arg)
	if err != nil {
		tele.Error(ctx, "failed to publish private message to nats: @1", "error", err.Error())
		return ce.Wrap(nil, err, input)
	}
	return nil
}

func (c *ChatService) retrievePrivatMessageSenders(ctx context.Context, msgs []md.PrivateMsg, input string) error {
	allMemberIDs := make(ct.Ids, 0)
	for _, r := range msgs {
		allMemberIDs = append(allMemberIDs, r.Sender.UserId)
	}

	usersMap, err := c.RetriveUsers.GetUsers(ctx, allMemberIDs)
	if err != nil {
		return ce.Wrap(nil, err, input)
	}

	for i := range msgs {
		retrieved := usersMap[msgs[i].Sender.UserId]
		msgs[i].Sender.Username = retrieved.Username
		msgs[i].Sender.AvatarId = retrieved.AvatarId
		msgs[i].Sender.AvatarURL = retrieved.AvatarURL
	}
	return nil
}

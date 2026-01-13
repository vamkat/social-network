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

// Creates new conversation between two users or fetches an existing.
// Returns convesation id, last read message id (if existing) and other user basic info if opted via RetrieveOther.
func (c *ChatService) GetOrCreatePrivateConv(ctx context.Context,
	params md.GetOrCreatePrivateConvReq) (res md.GetOrCreatePrivateConvResp, Err *ce.Error) {

	input := fmt.Sprintf("user ids: %d, %d", params.UserId, params.OtherUserId)

	if err := ct.ValidateStruct(params); err != nil {
		return res, ce.Wrap(ce.ErrInvalidArgument, err, input)
	}

	connected, Err := c.Clients.AreConnected(ctx, params.UserId, params.OtherUserId)
	if Err != nil {
		return res, ce.Wrap(nil, Err, input)
	}

	if !connected {
		return res, ce.New(ce.ErrPermissionDenied, ErrNotConnected, input)
	}

	newPC, err := c.Queries.GetOrCreatePrivateConv(ctx, params)
	if err != nil {
		return res, ce.Wrap(ce.ErrInternal, Err, input)
	}

	var isNew bool
	if newPC.LastReadMessageIdA == 0 && newPC.LastReadMessageIdB == 0 {
		isNew = true
	}

	var lastRead ct.Id
	if newPC.UserA == params.UserId {
		lastRead = newPC.LastReadMessageIdA
	} else {
		lastRead = newPC.LastReadMessageIdB
	}

	var otherUser md.User = md.User{UserId: params.OtherUserId}
	if params.RetrieveOtherUser {
		receiver, err := c.RetriveUsers.GetUser(ctx, params.OtherUserId)
		if err != nil {
			return res, ce.Wrap(nil, err, input)
		}
		otherUser = receiver
	}

	res = md.GetOrCreatePrivateConvResp{
		ConversationId:  newPC.Id,
		OtherUser:       otherUser,
		LastReadMessage: lastRead,
		IsNew:           isNew,
	}
	return res, nil
}

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
		allMemberIDs = append(allMemberIDs, r.OtherUser.UserId)
	}

	usersMap, err := c.RetriveUsers.GetUsers(ctx, allMemberIDs)
	if err != nil {
		return nil, ce.Wrap(nil, err, input)
	}

	for i := range conversations {
		retrieved := usersMap[conversations[i].OtherUser.UserId]
		conversations[i].OtherUser.Username = retrieved.Username
		conversations[i].OtherUser.AvatarId = retrieved.AvatarId
		conversations[i].OtherUser.AvatarURL = retrieved.AvatarURL
	}

	return conversations, nil
}

// Creates a message row with conversation id if user is a memeber.
// If user match of conversation_id and user_id fails returns error.
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

	err = c.NatsConn.Publish(ct.PrivateMessageKey(msg.ReceiverId), messageBytes)
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
	return res, nil
}

func (c *ChatService) UpdateLastReadPrivateMsg(ctx context.Context, arg md.UpdateLastReadMsgParams) *ce.Error {
	input := fmt.Sprintf("arg: %#v", arg)

	if err := ct.ValidateStruct(arg); err != nil {
		return ce.New(ce.ErrInvalidArgument, err, input)
	}

	err := c.Queries.UpdateLastReadPrivateMsg(ctx, arg)
	if err != nil {
		return ce.Wrap(nil, err, input)
	}
	return nil
}

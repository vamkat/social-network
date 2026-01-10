package application

import (
	"context"
	"fmt"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

// Returns a conversation id of a newly created or an existing conversation.
func (c *ChatService) CreateGroupConversation(ctx context.Context,
	params md.CreateGroupConvReq) (err error) {

	input := fmt.Sprintf("group id: %d, user ids: %d", params.GroupId, params.UserIds)

	if err := ct.ValidateStruct(params); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, input)
	}

	err = c.Queries.CreateGroupConv(ctx, params.GroupId)
	if err != nil {
		return ce.Wrap(ce.ErrInternal, err, input)
	}

	return nil
}

type CreateMessageInGroupReq struct {
	GroupId     ct.Id
	SenderId    ct.Id
	MessageBody ct.MsgBody
}

func (c *ChatService) CreateMessageInGroup(ctx context.Context,
	req CreateMessageInGroupReq) (res md.GroupMsg, err error) {
	input := fmt.Sprintf("params: %#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return res, ce.New(ce.ErrInvalidArgument, err, input)
	}

	if err := c.isMember(ctx, req.GroupId, req.SenderId, input); err != nil {
		return res, ce.Wrap(nil, err)
	}

	// Add message
	res, err = c.Queries.CreateNewGroupMessage(ctx, md.CreateGroupMsgReq{
		GroupId:     req.GroupId,
		SenderId:    req.SenderId,
		MessageText: req.MessageBody,
	})

	if err != nil {
		return res, ce.Wrap(nil, err, input)
	}

	return res, nil
}

func (c *ChatService) GetPrevGroupMessages(ctx context.Context,
	req md.GetGroupMsgsReq) (res md.GetGetGroupMsgsResp, err error) {

	input := fmt.Sprintf("req: %#v", req)
	if err := ct.ValidateStruct(req); err != nil {
		return res, ce.Wrap(ce.ErrInvalidArgument, err, input)
	}

	if err := c.isMember(ctx, req.GroupId, req.MemberId, input); err != nil {
		return res, ce.Wrap(nil, err)
	}

	res, err = c.Queries.GetPrevGroupMessages(ctx, req)
	if err != nil {
		return res, ce.Wrap(nil, err, input)
	}
	return res, nil
}

func (c *ChatService) GetNextGroupMessages(ctx context.Context,
	req md.GetGroupMsgsReq) (res md.GetGetGroupMsgsResp, err error) {

	input := fmt.Sprintf("req: %#v", req)
	if err := ct.ValidateStruct(req); err != nil {
		return res, ce.Wrap(ce.ErrInvalidArgument, err, input)
	}

	if err := c.isMember(ctx, req.GroupId, req.MemberId, input); err != nil {
		return res, ce.Wrap(nil, err)
	}

	res, err = c.Queries.GetNextGroupMessages(ctx, req)
	if err != nil {
		return res, ce.Wrap(nil, err, input)
	}

	return res, nil
}

// Returns a commonerrors Error type with public message if user is not a group member.
// Input is the caller function contextual details.
func (c *ChatService) isMember(ctx context.Context, groupId, memberId ct.Id, input string) error {
	isMember, err := c.Clients.IsGroupMember(ctx, groupId, memberId)
	if err != nil {
		return ce.ParseGrpcErr(err, input)
	}

	if !isMember {
		return ce.New(ce.ErrPermissionDenied,
			fmt.Errorf("user id: %d not a member of group id: %d",
				memberId, groupId),
			input,
		).WithPublic("user not a group member")
	}
	return nil
}

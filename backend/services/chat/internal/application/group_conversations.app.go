package application

import (
	"context"
	"fmt"
	"social-network/services/chat/internal/db/dbservice"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

// Returns a conversation id of a newly created or an existing conversation.
func (c *ChatService) CreateGroupConversation(ctx context.Context,
	params md.CreateGroupConvReq) (convId ct.Id, err error) {

	input := fmt.Sprintf("group id: %d, user ids: %d", params.GroupId, params.UserIds)

	if err := ct.ValidateStruct(params); err != nil {
		return 0, ce.Wrap(ce.ErrInvalidArgument, err, input)
	}

	convId, err = c.Queries.CreateGroupConv(ctx, params.GroupId)
	if err != nil {
		return 0, ce.Wrap(ce.ErrInternal, err, input)
	}

	return ct.Id(convId), err
}

type CreateMessageInGroupReq struct {
	GroupId     ct.Id
	SenderId    ct.Id
	MessageBody ct.MsgBody
}

func (c *ChatService) CreateMessageInGroup(ctx context.Context,
	params CreateMessageInGroupReq) (msg md.PrivateMsg, err error) {
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
				fmt.Errorf("user id: %d not a member of group id: %d",
					params.SenderId, params.GroupId),
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
			msg, err = c.Queries.CreateNewGroupMessage(ctx, md.CreateGroupMsgReq{
				GroupId:     convId,
				SenderId:    params.SenderId,
				MessageText: params.MessageBody,
			})

			if err != nil {
				return ce.Wrap(nil, err, input)
			}
			return nil
		})
	return msg, nil
}

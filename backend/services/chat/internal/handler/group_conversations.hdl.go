package handler

import (
	"context"
	"social-network/services/chat/internal/application"
	pb "social-network/shared/gen-go/chat"
	ce "social-network/shared/go/commonerrors"
	"social-network/shared/go/ct"
	mp "social-network/shared/go/mapping"
	md "social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *ChatHandler) CreateGroupConversation(
	ctx context.Context,
	params *pb.CreateGroupConversationRequest,
) (*emptypb.Empty, error) {

	tele.Info(ctx, "create group conversation called",
		"request", params,
	)

	err := h.Application.CreateGroupConversation(
		ctx,
		md.CreateGroupConvReq{
			GroupId: ct.Id(params.GroupId),
		},
	)
	if err != nil {
		tele.Error(ctx, "create group conversation error",
			"request", params,
			"error", err.Error(),
		)
		return nil, ce.GRPCStatus(err)
	}
	resp := &emptypb.Empty{}

	tele.Info(ctx, "create group conversation success",
		"groupId", params.GroupId,
	)

	return resp, nil
}

func (h *ChatHandler) CreateGroupMessage(
	ctx context.Context,
	params *pb.CreateGroupMessageRequest,
) (*pb.GroupMessage, error) {

	tele.Info(ctx, "create group message called",
		"request", params,
	)

	res, err := h.Application.CreateMessageInGroup(ctx,
		application.CreateMessageInGroupReq{
			GroupId:     ct.Id(params.GroupId),
			SenderId:    ct.Id(params.SenderId),
			MessageBody: ct.MsgBody(params.MessageText),
		},
	)
	if err != nil {
		tele.Error(ctx, "create group message error",
			"request", params,
			"error", err.Error(),
		)
		return nil, ce.GRPCStatus(err)
	}

	resp := &pb.GroupMessage{
		Id:          res.Id.Int64(),
		GroupId:     res.GroupId.Int64(),
		Sender:      mp.MapUserToProto(res.Sender),
		MessageText: string(res.MessageText),
		CreatedAt:   res.CreatedAt.ToProto(),
		UpdatedAt:   res.UpdatedAt.ToProto(),
		DeletedAt:   res.DeletedAt.ToProto(),
	}

	tele.Info(ctx, "create group message success",
		"response", resp,
	)

	return resp, nil
}

func (h *ChatHandler) GetPreviousGroupMessages(
	ctx context.Context,
	params *pb.GetGroupMessagesRequest,
) (*pb.GetGroupMessagesResponse, error) {

	tele.Info(ctx, "get previous group messages called",
		"request", params,
	)

	res, err := h.Application.GetPrevGroupMessages(ctx,
		md.GetGroupMsgsReq{
			GroupId:           ct.Id(params.GroupId),
			MemberId:          ct.Id(params.MemberId),
			BoundaryMessageId: ct.Id(params.BoundaryMessageId),
			Limit:             ct.Limit(params.Limit),
		},
	)
	if err != nil {
		tele.Error(ctx, "get previous group messages error",
			"request", params,
			"error", err.Error(),
		)
		return nil, ce.GRPCStatus(err)
	}

	resp := &pb.GetGroupMessagesResponse{
		HaveMore: res.HaveMore,
		Messages: mp.MapGroupMessagesToProto(res.Messages),
	}

	return resp, nil
}

func (h *ChatHandler) GetNextGroupMessages(
	ctx context.Context,
	params *pb.GetGroupMessagesRequest,
) (*pb.GetGroupMessagesResponse, error) {

	tele.Info(ctx, "get next group messages called",
		"request", params,
	)

	res, err := h.Application.GetNextGroupMessages(ctx,
		md.GetGroupMsgsReq{
			GroupId:           ct.Id(params.GroupId),
			MemberId:          ct.Id(params.MemberId),
			BoundaryMessageId: ct.Id(params.BoundaryMessageId),
			Limit:             ct.Limit(params.Limit),
		},
	)
	if err != nil {
		tele.Error(ctx, "get next group messages error",
			"request", params,
			"error", err.Error(),
		)
		return nil, ce.GRPCStatus(err)
	}

	resp := &pb.GetGroupMessagesResponse{
		HaveMore: res.HaveMore,
		Messages: mp.MapGroupMessagesToProto(res.Messages),
	}

	return resp, nil
}

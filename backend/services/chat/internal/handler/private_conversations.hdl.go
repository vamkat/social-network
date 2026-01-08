package handler

import (
	"context"

	pb "social-network/shared/gen-go/chat"
	ce "social-network/shared/go/commonerrors"
	"social-network/shared/go/ct"
	mp "social-network/shared/go/mapping"
	md "social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"

	"google.golang.org/protobuf/types/known/emptypb"

	_ "github.com/lib/pq"
)

// GetOrCreatePrivateConv creates a new private conversation between two users
// or returns an existing one if it already exists.
func (h *ChatHandler) GetOrCreatePrivateConv(
	ctx context.Context,
	params *pb.GetOrCreatePrivateConvRequest,
) (*pb.GetOrCreatePrivateConvResponse, error) {

	tele.Info(ctx, "get or create private conversation called @1", "request", params)

	// Call application layer
	res, err := h.Application.GetOrCreatePrivateConv(ctx, md.GetOrCreatePrivateConvReq{
		UserId:            ct.Id(params.User),
		OtherUserId:       ct.Id(params.OtherUser),
		RetrieveOtherUser: params.RetrieveOtherUser,
	})
	if err != nil {
		tele.Error(ctx, "get or create private conversation error",
			"request", params,
			"error", err.Error(),
		)
		return nil, ce.GRPCStatus(err)
	}

	resp := &pb.GetOrCreatePrivateConvResponse{
		ConversationId:  res.ConversationId.Int64(),
		OtherUser:       mp.MapUserToProto(res.OtherUser),
		LastReadMessage: res.LastReadMessage.Int64(),
		IsNew:           res.IsNew,
	}

	tele.Info(ctx, "get or create private conversation success. @1 @2",
		"request", params,
		"response", resp,
	)

	return resp, nil
}

// CreatePrivateMessage creates a new private message in a conversation.
func (h *ChatHandler) CreatePrivateMessage(
	ctx context.Context,
	params *pb.CreatePrivateMessageRequest,
) (*pb.PrivateMessage, error) {

	tele.Info(ctx, "create private message called @1", "request", params)

	// Call application layer
	msg, err := h.Application.CreatePrivateMessage(ctx, md.CreatePrivatMsgReq{
		ConversationId: ct.Id(params.ConversationId),
		SenderId:       ct.Id(params.SenderId),
		MessageText:    ct.MsgBody(params.MessageText),
	})
	if err != nil {
		tele.Error(ctx, "create private message error",
			"request", params,
			"error", err.Error(),
		)
		return nil, ce.GRPCStatus(err)
	}

	resp := mp.MapPMToProto(msg)

	tele.Info(ctx, "create private message success. @1 @2",
		"request", params,
		"response", resp,
	)

	return resp, nil
}

// GetPreviousPrivateMessages returns messages older than the boundary message ID.
func (h *ChatHandler) GetPreviousPrivateMessages(
	ctx context.Context,
	params *pb.GetPrivateMessagesRequest,
) (*pb.GetPrivateMessagesResponse, error) {

	tele.Info(ctx, "get previous private messages called @1", "request", params)

	// Call application layer
	res, err := h.Application.GetPreviousPMs(ctx, md.GetPrivatMsgsReq{
		ConversationId:    ct.Id(params.ConversationId),
		UserId:            ct.Id(params.UserId),
		BoundaryMessageId: ct.Id(params.BoundaryMessageId),
		Limit:             ct.Limit(params.Limit),
		RetrieveUsers:     params.RetrieveUsers,
	})
	if err != nil {
		tele.Error(ctx, "get previous private messages error",
			"request", params,
			"error", err.Error(),
		)
		return nil, ce.GRPCStatus(err)
	}

	resp := mp.MapGetPMsResp(res)

	tele.Info(ctx, "get previous private messages success. @1 @2",
		"request", params,
		"response", resp,
	)

	return resp, nil
}

// GetNextPrivateMessages returns messages newer than the boundary message ID.
func (h *ChatHandler) GetNextPrivateMessages(
	ctx context.Context,
	params *pb.GetPrivateMessagesRequest,
) (*pb.GetPrivateMessagesResponse, error) {

	tele.Info(ctx, "get next private messages called @1", "request", params)

	// Call application layer
	res, err := h.Application.GetNextPMs(ctx, md.GetPrivatMsgsReq{
		ConversationId:    ct.Id(params.ConversationId),
		UserId:            ct.Id(params.UserId),
		BoundaryMessageId: ct.Id(params.BoundaryMessageId),
		Limit:             ct.Limit(params.Limit),
		RetrieveUsers:     params.RetrieveUsers,
	})
	if err != nil {
		tele.Error(ctx, "get next private messages error",
			"request", params,
			"error", err.Error(),
		)
		return nil, ce.GRPCStatus(err)
	}

	resp := mp.MapGetPMsResp(res)

	tele.Info(ctx, "get next private messages success. @1 @2",
		"request", params,
		"response", resp,
	)

	return resp, nil
}

// UpdateLastReadPrivateMessage updates the last read message pointer
// for a user in a private conversation.
func (h *ChatHandler) UpdateLastReadPrivateMessage(
	ctx context.Context,
	params *pb.UpdateLastReadPrivateMessageRequest,
) (*emptypb.Empty, error) {

	tele.Info(ctx, "update last read private message called @1", "request", params)

	// Call application layer
	err := h.Application.UpdateLastReadPrivateMsg(ctx, md.UpdateLastReadMsgParams{
		ConversationId:    ct.Id(params.ConversationId),
		UserId:            ct.Id(params.UserId),
		LastReadMessageId: ct.Id(params.LastReadMessageId),
	})
	if err != nil {
		tele.Error(ctx, "update last read private message error",
			"request", params,
			"error", err.Error(),
		)
		return nil, ce.GRPCStatus(err)
	}

	resp := &emptypb.Empty{}

	tele.Info(ctx, "update last read private message success. @1 @2",
		"request", params,
		"response", resp,
	)

	return resp, nil
}

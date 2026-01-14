package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"social-network/shared/gen-go/chat"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/gorpc"
	utils "social-network/shared/go/http-utils"
	"social-network/shared/go/jwt"
	"social-network/shared/go/mapping"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"
	"time"
)

// Returns an existing or creates a new conversation for a user and a target user.
// If 'retrieveOther' is true the target user's basic info is also fetched.
// DEPRECATED
func (h *Handlers) GetOrCreatePrivateConversation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			tele.Error(ctx, "problem fetching claims")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "can't find claims")
			return
		}
		userId := claims.UserId

		type req struct {
			OtherUserId   ct.Id `json:"other_user_id"`
			RetrieveOther bool  `json:"retrieve_other"`
		}
		httpReq := models.GetOrCreatePrivateConvReq{}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		if err := ct.ValidateStruct(httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "bad url params: "+err.Error())
			return
		}

		res, err := h.ChatService.GetOrCreatePrivateConv(ctx, &chat.GetOrCreatePrivateConvRequest{
			User:                 userId,
			Interlocutor:         httpReq.InterlocutorId.Int64(),
			RetrieveInterlocutor: httpReq.RetrieveInterlocutor,
		})

		httpCode, _ := gorpc.Classify(err)
		if err != nil {
			err = ce.ParseGrpcErr(err)
			utils.ErrorJSON(ctx, w, httpCode, err.Error())
			return
		}

		err = utils.WriteJSON(ctx, w,
			httpCode,
			&models.GetOrCreatePrivateConvResp{
				ConversationId:  ct.Id(res.ConversationId),
				Interlocutor:    mapping.MapUserFromProto(res.Interlocutor),
				LastReadMessage: ct.Id(res.LastReadMessage),
				IsNew:           res.IsNew,
			})
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
		}
	}
}

// Creates a new message in a conversation. If the conversation does not exist it creates a new one.
func (h *Handlers) CreatePrivateMsg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			tele.Error(ctx, "problem fetching claims")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "can't find claims")
			return
		}

		userId := claims.UserId
		type req struct {
			InterlocutorId ct.Id      `json:"interlocutor_id"`
			Message        ct.MsgBody `json:"message_body"`
		}
		httpReq := req{}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		if err := ct.ValidateStruct(httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "bad url params: "+err.Error())
			return
		}

		grpcResponse, err := h.ChatService.CreatePrivateMessage(ctx,
			&chat.CreatePrivateMessageRequest{
				SenderId:       userId,
				InterlocutorId: httpReq.InterlocutorId.Int64(),
				MessageText:    httpReq.Message.String(),
			})

		httpCode, _ := gorpc.Classify(err)
		if err != nil {
			err = ce.ParseGrpcErr(err)
			utils.ErrorJSON(ctx, w, httpCode, err.Error())
			return
		}

		err = utils.WriteJSON(ctx, w, httpCode,
			mapping.MapPMFromProto(grpcResponse))
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
		}
	}
}

func (h *Handlers) GetPrivateConversations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			tele.Error(ctx, "problem fetching claims")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "can't find claims")
			return
		}

		v := r.URL.Query()
		userId := claims.UserId

		end := time.Now().AddDate(100, 0, 0)
		beforeDate, err1 := utils.ParamGet(v, "before-date", end, false)

		limit, err2 := utils.ParamGet(v, "limit", 100, true)
		beforeDateCt := ct.GenDateTime(beforeDate)

		if err := errors.Join(err1, err2); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "bad url params: "+err.Error())
			return
		}

		grpcResponse, err := h.ChatService.GetPrivateConversations(ctx, &chat.GetPrivateConversationsRequest{
			UserId:     userId,
			BeforeDate: beforeDateCt.ToProto(),
			Limit:      int32(limit),
		})

		httpCode, _ := gorpc.Classify(err)
		if err != nil {
			err = ce.ParseGrpcErr(err)
			utils.ErrorJSON(ctx, w, httpCode, err.Error())
			return
		}

		err = utils.WriteJSON(ctx, w,
			httpCode,
			mapping.MapConversationsFromProto(grpcResponse.Conversations))
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
		}
	}
}

func (h *Handlers) GetPrivateMessagesPag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tele.Info(ctx, "get private messages paginated called")

		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			tele.Error(ctx, "problem fetching claims")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "can't find claims")
			return
		}

		tele.Debug(ctx, "1. get private messages reached here ")
		v := r.URL.Query()
		tele.Debug(ctx, "1.1 get private messages reached here ")
		userId := claims.UserId
		interlocutorId, err1 := utils.ParamGet(v, "interlocutor-id", ct.Id(0), true)
		boundary, err2 := utils.ParamGet(v, "boundary", int64(0), false)
		limit, err3 := utils.ParamGet(v, "limit", 100, true)
		tele.Debug(ctx, "1.5 get private messages reached here ")
		retrieveusers, err4 := utils.ParamGet(v, "retrieve-users", false, false)
		tele.Debug(ctx, "1.6 get private messages reached here ")
		getPrevious, err5 := utils.ParamGet(v, "get-previous", true, false)

		tele.Debug(ctx, "2. get private messages reached here ")
		if err := errors.Join(err1, err2, err3, err4, err5); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "bad url params: "+err.Error())
			return
		}

		tele.Debug(ctx, "3. get private messages reached here ")
		if err := ct.ValidateBatch(interlocutorId, ct.Limit(limit)); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "bad url params: "+err.Error())
			return
		}

		tele.Debug(ctx, "4. get private messages reached here ")
		getFunc := h.ChatService.GetPreviousPrivateMessages
		if !getPrevious {
			getFunc = h.ChatService.GetNextPrivateMessages
		}

		grpcResponse, err := getFunc(ctx, &chat.GetPrivateMessagesRequest{
			UserId:            userId,
			InterlocutorId:    interlocutorId.Int64(),
			BoundaryMessageId: boundary,
			Limit:             int32(limit),
			RetrieveUsers:     retrieveusers,
		})

		tele.Debug(ctx, "get private messages @1 @2", "grpcRes", grpcResponse, "error", err)
		httpCode, _ := gorpc.Classify(err)
		if err != nil {
			err = ce.ParseGrpcErr(err)
			utils.ErrorJSON(ctx, w, httpCode, err.Error())
			return
		}

		err = utils.WriteJSON(ctx, w, httpCode, mapping.MapGetPMsRespFromProto(grpcResponse))
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
		}
	}
}

func (h *Handlers) CreateGroupMsg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			tele.Error(ctx, "problem fetching claims")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "can't find claims")
			return
		}

		userId := claims.UserId
		type req struct {
			GroupId ct.Id      `json:"group_id"`
			Message ct.MsgBody `json:"message_body"`
		}
		httpReq := req{}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		if err := ct.ValidateStruct(httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "bad url params: "+err.Error())
			return
		}

		grpcResponse, err := h.ChatService.CreateGroupMessage(ctx, &chat.CreateGroupMessageRequest{
			SenderId:    userId,
			GroupId:     httpReq.GroupId.Int64(),
			MessageText: httpReq.Message.String(),
		})

		httpCode, _ := gorpc.Classify(err)
		if err != nil {
			err = ce.ParseGrpcErr(err)
			utils.ErrorJSON(ctx, w, httpCode, err.Error())
			return
		}

		err = utils.WriteJSON(ctx, w, httpCode,
			mapping.MapGroupMessageFromProto(grpcResponse))
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
		}
	}
}

func (h *Handlers) GetGroupMessagesPag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			tele.Error(ctx, "problem fetching claims")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "can't find claims")
			return
		}

		v := r.URL.Query()
		userId := claims.UserId
		groupId, err1 := utils.ParamGet(v, "group-id", ct.Id(0), true)
		boundary, err2 := utils.ParamGet(v, "boundary", int64(0), true)
		limit, err3 := utils.ParamGet(v, "limit", int32(100), true)
		getPrevious, err4 := utils.ParamGet(v, "get-previous", true, false)

		if err := errors.Join(err1, err2, err3, err4); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "bad url params: "+err.Error())
			return
		}

		if err := ct.ValidateBatch(groupId, ct.Limit(limit)); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "bad url params: "+err.Error())
		}

		getFunc := h.ChatService.GetPreviousGroupMessages
		if !getPrevious {
			getFunc = h.ChatService.GetNextGroupMessages
		}

		grpcResponse, err := getFunc(ctx, &chat.GetGroupMessagesRequest{
			GroupId:           groupId.Int64(),
			MemberId:          userId,
			BoundaryMessageId: boundary,
			Limit:             limit,
		})

		httpCode, _ := gorpc.Classify(err)
		if err != nil {
			err = ce.ParseGrpcErr(err)
			utils.ErrorJSON(ctx, w, httpCode, err.Error())
			return
		}

		err = utils.WriteJSON(ctx, w, httpCode, &models.GetGroupMsgsResp{
			HaveMore: grpcResponse.HaveMore,
			Messages: mapping.MapGroupMessagesFromProto(grpcResponse.Messages),
		})
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
		}
	}
}

func (h *Handlers) UpdateLastRead() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			tele.Error(ctx, "problem fetching claims")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "can't find claims")
			return
		}

		userId := claims.UserId
		type req struct {
			ConversationId    ct.Id `json:"conversation_id"`
			LastReadMessageId ct.Id `json:"last_read_message_id"`
		}
		httpReq := req{}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		if err := ct.ValidateStruct(httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "bad url params: "+err.Error())
			return
		}

		grpcResponse, err := h.ChatService.UpdateLastReadPrivateMessage(ctx,
			&chat.UpdateLastReadPrivateMessageRequest{
				UserId:            userId,
				ConversationId:    httpReq.ConversationId.Int64(),
				LastReadMessageId: httpReq.LastReadMessageId.Int64(),
			})

		httpCode, _ := gorpc.Classify(err)
		if err != nil {
			err = ce.ParseGrpcErr(err)
			utils.ErrorJSON(ctx, w, httpCode, err.Error())
			return
		}

		err = utils.WriteJSON(ctx, w, httpCode, grpcResponse)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
		}
	}
}

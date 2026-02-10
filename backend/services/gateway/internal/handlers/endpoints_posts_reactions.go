package handlers

import (
	"errors"
	"net/http"
	"social-network/shared/gen-go/posts"
	ct "social-network/shared/go/ct"
	utils "social-network/shared/go/http-utils"
	"social-network/shared/go/jwt"
	"social-network/shared/go/models"
)

func (s *Handlers) toggleOrInsertReaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GenericReq{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := posts.GenericReq{
			RequesterId: claims.UserId,
			EntityId:    int64(body.EntityId),
		}

		_, err = s.PostsService.ToggleOrInsertReaction(ctx, &req)
		if err != nil {
			utils.ReturnHttpError(ctx, w, err)
			//utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("Could not toggle or insert reaction to entity with id %v: %v ", body.EntityId, err.Error()))
			return
		}

		utils.WriteJSON(ctx, w, http.StatusOK, nil)
	}
}

func (s *Handlers) getWhoLikedEntityId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		// if !ok {
		// 	panic(1)
		// }

		v := r.URL.Query()

		limit, err1 := utils.ParamGet(v, "limit", int32(1), false)

		offset, err2 := utils.ParamGet(v, "offset", int32(0), false)

		entityId, err3 := utils.PathValueGet(r, "entity_id", ct.Id(0), true)

		if err := errors.Join(err1, err2, err3); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "bad url params: "+err.Error())
			return
		}

		req := posts.GenericEntityPaginatedReq{
			EntityId: int64(entityId),
			Limit:    limit,
			Offset:   offset,
		}

		grpcResp, err := s.PostsService.GetWhoLikedEntityId(ctx, &req)
		if err != nil {
			utils.ReturnHttpError(ctx, w, err)
			//utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("Could not toggle or insert reaction to entity with id %v: %v ", body.EntityId, err.Error()))
			return
		}

		resp := make([]models.User, 0, len(grpcResp.Users))
		for _, user := range grpcResp.Users {
			newUser := models.User{
				UserId:    ct.Id(user.UserId),
				Username:  ct.Username(user.Username),
				AvatarId:  ct.Id(user.Avatar),
				AvatarURL: user.AvatarUrl,
			}
			resp = append(resp, newUser)
		}

		utils.WriteJSON(ctx, w, http.StatusOK, resp)
	}
}

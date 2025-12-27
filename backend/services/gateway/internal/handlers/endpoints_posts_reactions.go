package handlers

import (
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/posts"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
)

func (s *Handlers) toggleOrInsertReaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
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
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("Could not toggle or insert reaction to entity with id %v: %v ", body.EntityId, err.Error()))
			return
		}

		utils.WriteJSON(ctx, w, http.StatusOK, nil)
	}
}

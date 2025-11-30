package handlers

import (
	"net/http"
	"social-network/gateway/internal/security"
	"social-network/gateway/internal/utils"
	"social-network/shared/gen-go/users"
	"strconv"
)

func (h *Handlers) getUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userIdStr := r.URL.Query().Get("user_id")
		if userIdStr == "" {
			utils.ErrorJSON(w, http.StatusBadRequest, "missing user_id query param")
			return
		}

		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "invalid user_id query param")
			return
		}

		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		requesterId := int64(claims.UserId)

		grpcReq := users.GetUserProfileRequest{
			UserId:      userId,
			RequesterId: requesterId,
		}

		basicUserResp, err := h.Services.Users.GetUserProfile(r.Context(), &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to get user info")
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, basicUserResp)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send user info")
			return
		}

	}
}

package handlers

import (
	"fmt"
	"net/http"
	"social-network/gateway/internal/utils"
	"social-network/shared/gen-go/users"
	"strconv"
	"strings"
)

func (h *Handlers) getUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getUserProfile handler called")

		pathParts := strings.Split(r.URL.Path, "/")
		if pathParts[len(pathParts)-1] == "" {
			utils.ErrorJSON(w, http.StatusBadRequest, "missing user_id in URL path")
		}

		userId, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "invalid user_id query param")
			return
		}

		// claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		// if !ok {
		// 	panic(1)
		// }
		// requesterId := int64(claims.UserId)

		var requesterId int64 = 0

		grpcReq := users.GetUserProfileRequest{
			UserId:      userId,
			RequesterId: requesterId,
		}

		basicUserResp, err := h.Services.Users.GetUserProfile(r.Context(), &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to get user info: "+err.Error())
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, basicUserResp)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send user info")
			return
		}

	}
}

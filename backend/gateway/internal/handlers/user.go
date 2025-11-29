package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/gateway/internal/utils"
	"social-network/shared/gen-go/users"
)

func (h *Handlers) getBasicUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type basicUserRequest struct {
			UserId int64 `json:"user_id"`
		}

		httpReq := basicUserRequest{}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&httpReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		grpcReq := users.UserBasicInfoRequest{
			Id: httpReq.UserId,
		}

		basicUserResp, err := h.Services.Users.GetBasicUserInfo(r.Context(), &grpcReq)
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

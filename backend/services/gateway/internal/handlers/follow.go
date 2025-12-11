package handlers

import (
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Handlers) GetFollowSuggestions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		// if !ok {
		// 	panic(1)
		// }

		type reqBody struct {
			Value int64 `json:"value"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := wrapperspb.Int64Value{Value: body.Value}

		out, err := s.App.Users.GetFollowSuggestions(ctx, &req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch suggestions: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, out)
	}
}

func (s *Handlers) GetFollowersPaginated() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		// if !ok {
		// 	panic(1)
		// }

		type reqBody struct {
			UserId int64 `json:"user_id"`
			Limit  int32 `json:"limit"`
			Offset int32 `json:"offset"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := users.Pagination{
			UserId: body.UserId,
			Limit:  body.Limit,
			Offset: body.Offset,
		}

		out, err := s.App.Users.GetFollowersPaginated(ctx, &req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch followers: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, out)
	}
}

func (s *Handlers) FollowUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			TargetUserId int64 `json:"target_user_id"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := users.FollowUserRequest{
			FollowerId:   claims.UserId,
			TargetUserId: body.TargetUserId,
		}

		resp, err := s.App.Users.FollowUser(ctx, &req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not follow user: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp) //TODO check if returned values need to be removed
	}
}

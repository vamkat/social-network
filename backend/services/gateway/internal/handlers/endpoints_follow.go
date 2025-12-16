package handlers

import (
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/common"
	"social-network/shared/gen-go/posts"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Handlers) GetFollowSuggestions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}
		requesterId := int64(claims.UserId)

		req := wrapperspb.Int64Value{Value: requesterId}

		part1, err := s.UsersService.GetFollowSuggestions(ctx, &req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch suggestions from users: "+err.Error())
			return
		}

		part2, err := s.PostsService.SuggestUsersByPostActivity(ctx, &posts.SimpleIdReq{Id: requesterId})
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch suggestions from posts: "+err.Error())
			return
		}

		myMap := make(map[int64]*common.User)
		for _, user := range part1.Users {
			myMap[user.UserId] = user
		}
		for _, user := range part2.Users {
			myMap[user.UserId] = user
		}
		dedupedUsers := make([]models.User, 0, len(part1.Users)+len(part2.Users))
		for _, user := range myMap {
			newUser := models.User{
				UserId:   ct.Id(user.UserId),
				Username: ct.Username(user.Username),
				AvatarId: ct.Id(user.Avatar),
			}
			dedupedUsers = append(dedupedUsers, newUser)
		}
		resp := models.Users{
			Users: dedupedUsers,
		}
		utils.WriteJSON(w, http.StatusOK, resp)
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

		out, err := s.UsersService.GetFollowersPaginated(ctx, &req)
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

		body, err := utils.JSON2Struct(&models.FollowUserReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := users.FollowUserRequest{
			FollowerId:   claims.UserId,
			TargetUserId: body.TargetUserId.Int64(),
		}

		resp, err := s.UsersService.FollowUser(ctx, &req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not follow user: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp) //TODO check if returned values need to be removed
	}
}

package handlers

import (
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	cm "social-network/shared/gen-go/common"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (h *Handlers) getUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getUserProfile handler called")

		pathParts := strings.Split(r.URL.Path, "/")
		if pathParts[len(pathParts)-1] == "" {
			utils.ErrorJSON(w, http.StatusBadRequest, "missing user_id in URL path")
			return
		}

		userId, err := ct.DecryptId(pathParts[len(pathParts)-1])
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "invalid user_id query param")
			return
		}

		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}
		requesterId := int64(claims.UserId)

		grpcReq := users.GetUserProfileRequest{
			UserId:      userId.Int64(),
			RequesterId: requesterId,
		}

		grpcResp, err := h.App.Users.GetUserProfile(r.Context(), &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to get user info: "+err.Error())
			return
		}

		fmt.Println("retrieved user profile: ", grpcResp)

		type userProfile struct {
			UserId            ct.EncryptedId `json:"user_id"`
			Username          ct.Username    `json:"username"`
			FirstName         ct.Name        `json:"first_name"`
			LastName          ct.Name        `json:"last_name"`
			DateOfBirth       ct.DateOfBirth `json:"date_of_birth"`
			Avatar            ct.EncryptedId `json:"avatar,omitempty"`
			About             ct.About       `json:"about,omitempty"`
			Public            bool           `json:"public"`
			CreatedAt         time.Time      `json:"created_at"`
			FollowersCount    int64          `json:"followers_count"`
			FollowingCount    int64          `json:"following_count"`
			GroupsCount       int64          `json:"groups_count"`
			OwnedGroupsCount  int64          `json:"owned_groups_count"`
			ViewerIsFollowing bool           `json:"viewer_is_following"`
			OwnProfile        bool           `json:"own_profile"`
			IsPending         bool           `json:"is_pending"`
		}

		userProfileResponse := userProfile{
			UserId:            ct.EncryptedId(grpcResp.UserId),
			Username:          ct.Username(grpcResp.Username),
			FirstName:         ct.Name(grpcResp.FirstName),
			LastName:          ct.Name(grpcResp.LastName),
			DateOfBirth:       ct.DateOfBirth(grpcResp.DateOfBirth.AsTime()),
			Avatar:            ct.EncryptedId(grpcResp.Avatar),
			About:             ct.About(grpcResp.About),
			Public:            grpcResp.Public,
			CreatedAt:         grpcResp.CreatedAt.AsTime(),
			FollowersCount:    grpcResp.FollowersCount,
			FollowingCount:    grpcResp.FollowingCount,
			GroupsCount:       grpcResp.GroupsCount,
			OwnedGroupsCount:  grpcResp.OwnedGroupsCount,
			ViewerIsFollowing: grpcResp.ViewerIsFollowing,
			OwnProfile:        grpcResp.OwnProfile,
			IsPending:         grpcResp.IsPending,
		}

		fmt.Println("transformed profile struct: ", userProfileResponse)

		err = utils.WriteJSON(w, http.StatusOK, userProfileResponse)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send user info")
			return
		}

	}
}

func (s *Handlers) GetBasicUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		// if !ok {
		// 	panic(1)
		// }

		type reqBody struct {
			UserId int64 `json:"user_id"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := wrapperspb.Int64Value{Value: body.UserId}

		out, err := s.App.Users.GetBasicUserInfo(ctx, &req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch user: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, out)
	}
}

func (s *Handlers) GetBatchBasicUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		// if !ok {
		// 	panic(1)
		// }

		type reqBody struct {
			Values []int64 `json:"values"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := cm.Int64Arr{Values: body.Values}

		out, err := s.App.Users.GetBatchBasicUserInfo(ctx, &req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch users: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, out)
	}
}

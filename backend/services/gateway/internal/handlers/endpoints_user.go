package handlers

import (
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
	"strings"
	"time"
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

		grpcResp, err := h.UsersService.GetUserProfile(r.Context(), &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to get user info: "+err.Error())
			return
		}

		fmt.Println("retrieved user profile: ", grpcResp)

		type userProfile struct {
			UserId            ct.Id          `json:"user_id"`
			Username          ct.Username    `json:"username"`
			FirstName         ct.Name        `json:"first_name"`
			LastName          ct.Name        `json:"last_name"`
			DateOfBirth       ct.DateOfBirth `json:"date_of_birth"`
			Avatar            ct.Id          `json:"avatar,omitempty"`
			About             ct.About       `json:"about,omitempty"`
			Public            bool           `json:"public"`
			CreatedAt         time.Time      `json:"created_at"`
			Email             fmt.Stringer   `json:"email"`
			FollowersCount    int64          `json:"followers_count"`
			FollowingCount    int64          `json:"following_count"`
			GroupsCount       int64          `json:"groups_count"`
			OwnedGroupsCount  int64          `json:"owned_groups_count"`
			ViewerIsFollowing bool           `json:"viewer_is_following"`
			OwnProfile        bool           `json:"own_profile"`
			IsPending         bool           `json:"is_pending"`
		}

		userProfileResponse := userProfile{
			UserId:            ct.Id(grpcResp.UserId),
			Username:          ct.Username(grpcResp.Username),
			FirstName:         ct.Name(grpcResp.FirstName),
			LastName:          ct.Name(grpcResp.LastName),
			DateOfBirth:       ct.DateOfBirth(grpcResp.DateOfBirth.AsTime()),
			Avatar:            ct.Id(grpcResp.Avatar),
			About:             ct.About(grpcResp.About),
			Public:            grpcResp.Public,
			CreatedAt:         grpcResp.CreatedAt.AsTime(),
			Email:             ct.Email(grpcResp.Email),
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

// OK?
func (s *Handlers) SearchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type reqBody struct {
			Query string `json:"query"`
			Limit int32  `json:"limit"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.UserSearchRequest{
			SearchTerm: body.Query,
			Limit:      body.Limit,
		}

		grpcResp, err := s.UsersService.SearchUsers(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not search users: "+err.Error())
			return
		}

		resp := models.Users{
			Users: make([]models.User, 0, len(grpcResp.Users)),
		}

		for _, user := range grpcResp.Users {
			newUser := models.User{
				UserId:   ct.Id(user.UserId),
				Username: ct.Username(user.Username),
				AvatarId: ct.Id(user.Avatar),
			}

			resp.Users = append(resp.Users, newUser)
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

// OK
func (s *Handlers) UpdateProfilePrivacy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			Public bool `json:"public"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.UpdateProfilePrivacyRequest{
			UserId: claims.UserId,
			Public: body.Public,
		}

		_, err = s.UsersService.UpdateProfilePrivacy(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not update privacy: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, nil)
	}
}

// OK
func (s *Handlers) UpdateUserEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			Email string `json:"email"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.UpdateEmailRequest{
			UserId: claims.UserId,
			Email:  body.Email,
		}

		_, err = s.UsersService.UpdateUserEmail(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not update email: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, nil)
	}
}

// TODO should probably be done using a specific link / needs extra validation
func (s *Handlers) UpdateUserPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}
		_ = body

		fmt.Println("old password:", body.OldPassword, " new password:", body.NewPassword)

		oldPassword, err := ct.Password(body.OldPassword).Hash()
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "could not hash password")
			return
		}

		newPassword, err := ct.Password(body.NewPassword).Hash()
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "could not hash password")
			return
		}

		fmt.Println("hashed old password:", oldPassword.String(), " hashed new password:", newPassword.String())

		req := &users.UpdatePasswordRequest{
			UserId:      claims.UserId,
			OldPassword: oldPassword.String(),
			NewPassword: newPassword.String(),
		}

		_, err = s.UsersService.UpdateUserPassword(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not update password: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, nil)
	}
}

func (s *Handlers) UpdateUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type UpdateProfileJSONRequest struct {
			Username    ct.Username    `json:"username"`
			FirstName   ct.Name        `json:"first_name"`
			LastName    ct.Name        `json:"last_name"`
			DateOfBirth ct.DateOfBirth `json:"date_of_birth"`
			AvatarId    ct.Id          `json:"avatar_id" validate:"nullable"`
			About       ct.About       `json:"about" validate:"nullable"`
		}

		body, err := utils.JSON2Struct(&UpdateProfileJSONRequest{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcRequest := &users.UpdateProfileRequest{
			UserId:      claims.UserId,
			Username:    body.Username.String(),
			FirstName:   body.FirstName.String(),
			LastName:    body.LastName.String(),
			DateOfBirth: body.DateOfBirth.ToProto(),
			Avatar:      body.AvatarId.Int64(),
			About:       body.About.String(),
		}

		grpcResp, err := s.UsersService.UpdateUserProfile(ctx, grpcRequest)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not update profile: "+err.Error())
			return
		}

		resp := models.UserProfileResponse{
			UserId:      ct.Id(grpcResp.UserId),
			Username:    ct.Username(grpcResp.Username),
			FirstName:   ct.Name(grpcResp.FirstName),
			LastName:    ct.Name(grpcResp.LastName),
			DateOfBirth: ct.DateOfBirth(grpcResp.DateOfBirth.AsTime()), //TODO @vag lmao, is this correct?
			AvatarId:    ct.Id(grpcResp.Avatar),
			About:       ct.About(grpcResp.About),
			Public:      grpcResp.Public,
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

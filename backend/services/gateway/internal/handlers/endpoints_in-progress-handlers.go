package handlers

import (
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
)

// unfollowUser
// updateProfileInfo
// updateEmail
// updatePassword
// updatePrivacy
// searchUsers
// getUserGroupsPaginated -> if i could get the total num of my groups as a first response that would be nice
// searchGroups

// OK...?
func (s *Handlers) GetFollowingPaginated() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			Limit  int32 `json:"limit"`
			Offset int32 `json:"offset"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.Pagination{
			UserId: claims.UserId,
			Limit:  body.Limit,
			Offset: body.Offset,
		}

		grpcResp, err := s.UsersService.GetFollowingPaginated(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch following users: "+err.Error())
			return
		}

		resp := &models.Users{}

		for _, grpcUser := range grpcResp.Users {
			user := models.User{
				UserId:   ct.Id(grpcUser.UserId),
				Username: ct.Username(grpcUser.Username),
				AvatarId: ct.Id(grpcUser.Avatar),
			}
			resp.Users = append(resp.Users, user)
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

// OK?
func (s *Handlers) GetGroupInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			GroupId int64 `json:"group_id"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.GeneralGroupRequest{
			GroupId: body.GroupId,
			UserId:  claims.UserId,
		}

		grpcResp, err := s.UsersService.GetGroupInfo(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch group info: "+err.Error())
			return
		}

		resp := models.Group{
			GroupId:          ct.Id(grpcResp.GroupId),
			GroupOwnerId:     ct.Id(grpcResp.GroupOwnerId),
			GroupTitle:       ct.Title(grpcResp.GroupTitle),
			GroupDescription: ct.About(grpcResp.GroupDescription),
			GroupImage:       grpcResp.GroupImage,
			MembersCount:     grpcResp.MembersCount,
			IsMember:         grpcResp.IsMember,
			IsOwner:          grpcResp.IsOwner,
			IsPending:        grpcResp.IsPending,
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

// OK?
func (s *Handlers) GetGroupMembers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			GroupId int64 `json:"group_id"`
			Limit   int32 `json:"limit"`
			Offset  int32 `json:"offset"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.GroupMembersRequest{
			UserId:  claims.UserId,
			GroupId: body.GroupId,
			Limit:   body.Limit,
			Offset:  body.Offset,
		}

		grpcResp, err := s.UsersService.GetGroupMembers(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch group members: "+err.Error())
			return
		}

		resp := models.GroupUsers{}

		for _, group := range grpcResp.GroupUserArr {
			newGroup := models.GroupUser{
				UserId:    ct.Id(group.UserId),
				Username:  ct.Username(group.Username),
				AvatarId:  ct.Id(group.Avatar),
				GroupRole: group.GroupRole,
			}
			resp.GroupUsers = append(resp.GroupUsers, newGroup)
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

// OK?
func (s *Handlers) GetUserGroupsPaginated() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			Limit  int32 `json:"limit"`
			Offset int32 `json:"offset"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.Pagination{
			UserId: claims.UserId,
			Limit:  body.Limit,
			Offset: body.Offset,
		}

		grpcResp, err := s.UsersService.GetUserGroupsPaginated(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch user groups: "+err.Error())
			return
		}

		resp := GroupsT{}
		for _, group := range grpcResp.GroupArr {
			newGroup := GroupT{
				GroupId:          ct.Id(group.GroupId),
				GroupOwnerId:     ct.Id(group.GroupOwnerId),
				GroupTitle:       group.GroupTitle,
				GroupDescription: group.GroupDescription,
				GroupImage:       group.GroupImage,
				MembersCount:     group.MembersCount,
				IsMember:         group.IsMember,
				IsOwner:          group.IsOwner,
				IsPending:        group.IsPending,
			}
			resp.Groups = append(resp.Groups, newGroup)
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

// OK?
func (s *Handlers) HandleFollowRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			RequesterId int64 `json:"requester_id"`
			Accept      bool  `json:"accept"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.HandleFollowRequestRequest{
			UserId:      claims.UserId,
			RequesterId: body.RequesterId,
			Accept:      body.Accept,
		}

		_, err = s.UsersService.HandleFollowRequest(ctx, req)
		if err != nil { //soft TODO better error?
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not handle follow request: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, nil)
	}
}

// OK?
func (s *Handlers) HandleGroupJoinRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.HandleJoinRequest{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.HandleJoinRequest{
			OwnerId:     claims.UserId,
			GroupId:     body.GroupId.Int64(),
			RequesterId: body.RequesterId.Int64(),
			Accepted:    body.Accepted,
		}

		_, err = s.UsersService.HandleGroupJoinRequest(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not handle group join request: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, nil)
	}
}

// OK?
func (s *Handlers) InviteToGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.InviteToGroupReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.InviteToGroupRequest{
			InviterId: claims.UserId,
			InvitedId: body.InvitedId.Int64(),
			GroupId:   body.GroupId.Int64(),
		}

		_, err = s.UsersService.InviteToGroup(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not invite user to group: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, nil)
	}
}

// OK?
func (s *Handlers) LeaveGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GeneralGroupReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.GeneralGroupRequest{
			UserId:  claims.UserId,
			GroupId: body.GroupId.Int64(),
		}

		_, err = s.UsersService.LeaveGroup(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not leave group: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, nil)
	}
}

// OK?
func (s *Handlers) RequestJoinGroupOrCancel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GroupJoinRequest{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.GroupJoinRequest{
			RequesterId: claims.UserId,
			GroupId:     body.GroupId.Int64(),
		}

		_, err = s.UsersService.RequestJoinGroup(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not process join request: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, nil)
	}
}

// OK?
func (s *Handlers) RespondToGroupInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			GroupId int64 `json:"group_id"`
			Accept  bool  `json:"accept"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.HandleGroupInviteRequest{
			InvitedId: claims.UserId,
			GroupId:   body.GroupId,
			Accepted:  body.Accept,
		}

		_, err = s.UsersService.RespondToGroupInvite(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not respond to invite: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, nil)
	}
}

// OK?
func (s *Handlers) SearchGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			Query  string `json:"query"`
			Limit  int32  `json:"limit"`
			Offset int32  `json:"offset"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.GroupSearchRequest{
			SearchTerm: body.Query,
			Limit:      body.Limit,
			Offset:     body.Offset,
			UserId:     claims.UserId,
		}

		grpcResp, err := s.UsersService.SearchGroups(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not search groups: "+err.Error())
			return
		}

		resp := models.Groups{
			Groups: make([]models.Group, 0, len(grpcResp.GroupArr)),
		}
		for _, group := range grpcResp.GroupArr {
			newGroup := models.Group{
				GroupId:          ct.Id(group.GroupId),
				GroupOwnerId:     ct.Id(group.GroupOwnerId),
				GroupTitle:       ct.Title(group.GroupTitle),
				GroupDescription: ct.About(group.GroupDescription),
				GroupImage:       group.GroupImage,
				MembersCount:     group.MembersCount,
				IsMember:         group.IsMember,
				IsOwner:          group.IsOwner,
				IsPending:        group.IsPending,
			}

			resp.Groups = append(resp.Groups, newGroup)
		}

		utils.WriteJSON(w, http.StatusOK, resp)
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
func (s *Handlers) UnFollowUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			UserId int64 `json:"user_id"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.FollowUserRequest{
			FollowerId:   claims.UserId,
			TargetUserId: body.UserId,
		}

		resp, err := s.UsersService.UnFollowUser(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not unfollow user: "+err.Error())
			return
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

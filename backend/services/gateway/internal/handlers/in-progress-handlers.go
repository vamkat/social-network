package handlers

import (
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// func (s *Handlers) GetAllGroupsPaginated() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}

// 		type reqBody struct {
// 			Limit  int32 `json:"limit"`
// 			Offset int32 `json:"offset"`
// 		}

// 		body, err := utils.JSON2Struct(&reqBody{}, r)
// 		if err != nil {
// 			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
// 			return
// 		}

// 		req := users.Pagination{
// 			UserId: claims.UserId,
// 			Limit:  body.Limit,
// 			Offset: body.Offset,
// 		}

// 		out, err := s.App.Users.GetAllGroupsPaginated(ctx, &req)
// 		if err != nil {
// 			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch groups: "+err.Error())
// 			return
// 		}

// 		utils.WriteJSON(w, http.StatusOK, out)
// 	}
// }

// func (s *Handlers) GetBatchBasicUserInfo() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		// claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		// if !ok {
// 		// 	panic(1)
// 		// }

// 		type reqBody struct {
// 			Values []int64 `json:"values"`
// 		}

// 		body, err := utils.JSON2Struct(&reqBody{}, r)
// 		if err != nil {
// 			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
// 			return
// 		}

// 		req := cm.UserIds{Values: body.Values}

// 		out, err := s.App.Users.GetBatchBasicUserInfo(ctx, &req)
// 		if err != nil {
// 			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch users: "+err.Error())
// 			return
// 		}

// 		utils.WriteJSON(w, http.StatusOK, out)
// 	}
// }

// func (s *Handlers) FollowUser() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}

// 		type reqBody struct {
// 			TargetUserId int64 `json:"target_user_id"`
// 		}

// 		body, err := utils.JSON2Struct(&reqBody{}, r)
// 		if err != nil {
// 			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
// 			return
// 		}

// 		req := users.FollowUserRequest{
// 			FollowerId:   claims.UserId,
// 			TargetUserId: body.TargetUserId,
// 		}

// 		resp, err := s.App.Users.FollowUser(ctx, &req)
// 		if err != nil {
// 			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not follow user: "+err.Error())
// 			return
// 		}

// 		utils.WriteJSON(w, http.StatusOK, resp) //TODO check if returned values need to be removed
// 	}
// }

func (s *Handlers) GetFollowingIds() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		req := &wrapperspb.Int64Value{Value: claims.UserId}
		resp, err := s.App.Users.GetFollowingIds(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch following ids: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

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

		resp, err := s.App.Users.GetFollowingPaginated(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch following users: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

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

		resp, err := s.App.Users.GetGroupInfo(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch group info: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

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

		resp, err := s.App.Users.GetGroupMembers(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch group members: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

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

		resp, err := s.App.Users.GetUserGroupsPaginated(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch user groups: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

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

		resp, err := s.App.Users.HandleFollowRequest(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not handle follow request: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

//@kv unsure what to do here in regard to what data the endpoint wants/needs

// func (s *Handlers) HandleGroupJoinRequest() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}

// 		body, err := utils.JSON2Struct(&models.GroupJoinRequest{}, r)
// 		if err != nil {
// 			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
// 			return
// 		}

// 		models.GroupJoinRequest

// 		req := &users.HandleJoinRequest{
// 			OwnerId:     claims.UserId,
// 			GroupId:     body.GroupId.Int64(),
// 			RequesterId: claims.UserId,
// 			Accepted:    body.Accept,
// 		}

// 		_, err = s.App.Users.HandleGroupJoinRequest(ctx, req)
// 		if err != nil {
// 			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not handle group join request: "+err.Error())
// 			return
// 		}

// 		utils.WriteJSON(w, http.StatusOK, struct{}{})
// 	}
// }

func (s *Handlers) InviteToGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			GroupId int64 `json:"group_id"`
			UserId  int64 `json:"user_id"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.InviteToGroupRequest{
			InviterId: claims.UserId,
			// InvitedId: claims.UserId,
			GroupId: body.GroupId,
			// UserId:    body.UserId,
		}

		_, err = s.App.Users.InviteToGroup(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not invite user to group: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, struct{}{})
	}
}

func (s *Handlers) IsFollowing() http.HandlerFunc {
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

		// req := &users.IsFollowingRequest{
		// 	FollowerId:   claims.UserId,
		// 	TargetUserId: body.TargetUserId,
		// }

		req := &users.FollowUserRequest{
			FollowerId:   claims.UserId,
			TargetUserId: body.TargetUserId,
		}

		resp, err := s.App.Users.IsFollowing(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not check following status: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (s *Handlers) IsGroupMember() http.HandlerFunc {
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
			UserId:  claims.UserId,
			GroupId: body.GroupId,
		}

		resp, err := s.App.Users.IsGroupMember(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not check group membership: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (s *Handlers) LeaveGroup() http.HandlerFunc {
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
			UserId:  claims.UserId,
			GroupId: body.GroupId,
		}

		_, err = s.App.Users.LeaveGroup(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not leave group: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, struct{}{})
	}
}

func (s *Handlers) RequestJoinGroupOrCancel() http.HandlerFunc {
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

		req := &users.GroupJoinRequest{
			RequesterId: claims.UserId,
			GroupId:     body.GroupId,
		}

		_, err = s.App.Users.RequestJoinGroupOrCancel(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not process join request: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, struct{}{})
	}
}

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
			// Accept:  body.Accept,
		}

		_, err = s.App.Users.RespondToGroupInvite(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not respond to invite: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, struct{}{})
	}
}

func (s *Handlers) SearchGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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
		}

		resp, err := s.App.Users.SearchGroups(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not search groups: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (s *Handlers) SearchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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

		req := &users.UserSearchRequest{
			SearchTerm: body.Query,
			Limit:      body.Limit,
			// Offset: body.Offset,
		}

		resp, err := s.App.Users.SearchUsers(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not search users: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (s *Handlers) UnFollowUser() http.HandlerFunc {
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

		req := &users.FollowUserRequest{
			FollowerId:   claims.UserId,
			TargetUserId: body.TargetUserId,
		}

		resp, err := s.App.Users.UnFollowUser(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not unfollow user: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (s *Handlers) UpdateProfilePrivacy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			Private bool `json:"private"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.UpdateProfilePrivacyRequest{
			UserId: claims.UserId,
			Public: !body.Private,
		}

		_, err = s.App.Users.UpdateProfilePrivacy(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not update privacy: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, struct{}{})
	}
}

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

		_, err = s.App.Users.UpdateUserEmail(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not update email: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, struct{}{})
	}
}

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

		req := &users.UpdatePasswordRequest{
			UserId:      claims.UserId,
			OldPassword: body.OldPassword,
			NewPassword: body.NewPassword,
		}

		_, err = s.App.Users.UpdateUserPassword(ctx, req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not update password: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, struct{}{})
	}
}

func (s *Handlers) UpdateUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&users.UpdateProfileRequest{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		body.UserId = claims.UserId

		resp, err := s.App.Users.UpdateUserProfile(ctx, body)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not update profile: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

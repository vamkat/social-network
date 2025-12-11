package handlers

// func (s *Handlers) GetFollowingIds() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.GetFollowingIds()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) GetFollowingPaginated() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.GetFollowingPaginated()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) GetGroupInfo() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.GetGroupInfo()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) GetGroupMembers() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.GetGroupMembers()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) GetUserGroupsPaginated() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.GetUserGroupsPaginated()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) HandleFollowRequest() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.HandleFollowRequest()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) HandleGroupJoinRequest() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.HandleGroupJoinRequest()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) InviteToGroup() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.InviteToGroup()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) IsFollowing(context.Context, *users.IsFollowingRequest) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.IsFollowing()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) IsGroupMember(context.Context, *users.GeneralGroupRequest) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.IsGroupMember()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) LeaveGroup() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.LeaveGroup()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) RequestJoinGroupOrCancel() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.RequestJoinGroupOrCancel()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) RespondToGroupInvite() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.RespondToGroupInvite()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) SearchGroups() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.SearchGroups()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) SearchUsers() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.SearchUsers()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) UnFollowUser() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.UnFollowUser()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) UpdateProfilePrivacy() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.UpdateProfilePrivacy()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) UpdateUserEmail() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.UpdateUserEmail()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) UpdateUserPassword() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.UpdateUserPassword()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

// func (s *Handlers) UpdateUserProfile() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			panic(1)
// 		}
// 		// s.App.Users.UpdateUserProfile()
// 		utils.WriteJSON(w, http.StatusOK, resp)
// 	}
// }

package application

// ctx := context.Background()
// app := NewApplication(db, chatSvc, notifSvc)

// // ======================
// // USER REGISTRATION & LOGIN
// // ======================
// u1, _ := app.RegisterUser(ctx, models.RegisterUserRequest{Username: "alice"})            // create new user
// u2, _ := app.RegisterUser(ctx, models.RegisterUserRequest{Username: "bob"})              // create another user
// _, _ = app.LoginUser(ctx, models.LoginRequest{Username: "alice"})                         // login user
// _, _ = app.LoginUser(ctx, models.LoginRequest{Username: "bob"})                           // login second user

// // ======================
// // BASIC INFO & PROFILES
// // ======================
// _, _ = app.GetBasicUserInfo(ctx, u1.UserId)                                              // fetch basic info
// _, _ = app.GetBatchBasicUserInfo(ctx, customtypes.Ids{u1.UserId, u2.UserId})             // fetch multiple users
// _, _ = app.GetUserProfile(ctx, models.UserProfileRequest{UserId: u1.UserId, RequesterId: u2.UserId}) // fetch profile
// _ = app.UpdateUserProfile(ctx, models.UpdateProfileRequest{UserId: u1.UserId, FirstName: "Alice"})   // update profile
// _ = app.UpdateUserEmail(ctx, models.UpdateEmailRequest{UserId: u1.UserId, Email: "alice@example.com"}) // update email
// _ = app.UpdateUserPassword(ctx, models.UpdatePasswordRequest{UserId: u1.UserId, OldPassword: "old", NewPassword: "new"}) // update password
// _ = app.UpdateProfilePrivacy(ctx, models.UpdateProfilePrivacyRequest{UserId: u1.UserId, Public: false}) // toggle privacy

// // ======================
// // FOLLOWERS & SOCIAL
// // ======================
// _, _ = app.IsFollowing(ctx, models.FollowUserReq{UserID: u1.UserId, FollowerID: u2.UserId}) // check if following
// _, _ = app.AreFollowingEachOther(ctx, models.FollowUserReq{UserID: u1.UserId, FollowerID: u2.UserId}) // mutual follow
// _, _ = app.isFollowRequestPending(ctx, models.FollowUserReq{UserID: u1.UserId, FollowerID: u2.UserId}) // pending follow

// _, _ = app.FollowUser(ctx, models.FollowUserReq{UserID: u1.UserId, FollowerID: u2.UserId}) // follow action
// _, _ = app.UnFollowUser(ctx, models.FollowUserReq{UserID: u1.UserId, FollowerID: u2.UserId}) // unfollow action
// _ = app.HandleFollowRequest(ctx, models.HandleFollowRequestReq{UserID: u1.UserId, RequesterID: u2.UserId}) // accept/reject request
// _, _ = app.GetFollowersPaginated(ctx, models.Pagination{Limit: 10}) // fetch followers
// _, _ = app.GetFollowingPaginated(ctx, models.Pagination{Limit: 10}) // fetch following
// _, _ = app.GetFollowingIds(ctx, u1.UserId)                           // get following ids
// _, _ = app.GetFollowSuggestions(ctx, u1.UserId)                      // suggest followers

// // ======================
// // GROUP CREATION & MANAGEMENT
// // ======================
// gID, _ := app.CreateGroup(ctx, &models.CreateGroupRequest{Name: "TestGroup"}) // create group
// _, _ = app.GetGroupInfo(ctx, models.GeneralGroupReq{GroupId: gID})            // fetch group info
// _, _ = app.GetGroupMembers(ctx, models.GroupMembersReq{GroupId: gID})         // list members
// _, _ = app.GetAllGroupsPaginated(ctx, models.Pagination{Limit: 10})           // all groups paginated
// _, _ = app.GetUserGroupsPaginated(ctx, models.Pagination{Limit: 10})          // user groups paginated
// _, _ = app.IsGroupMember(ctx, models.GeneralGroupReq{GroupId: gID, UserId: u1.UserId}) // check membership
// _, _ = app.isGroupOwner(ctx, models.GeneralGroupReq{GroupId: gID, UserId: u1.UserId})  // check ownership
// _, _ = app.isGroupMembershipPending(ctx, models.GeneralGroupReq{GroupId: gID, UserId: u2.UserId}) // pending membership
// _, _ = app.userInRelationToGroup(ctx, models.GeneralGroupReq{GroupId: gID, UserId: u2.UserId}) // relation type

// // ======================
// // GROUP INVITES & REQUESTS
// // ======================
// _ = app.InviteToGroup(ctx, models.InviteToGroupReq{GroupId: gID, UserId: u2.UserId})              // send invite
// _ = app.RespondToGroupInvite(ctx, models.HandleGroupInviteRequest{GroupId: gID, UserId: u2.UserId, Accept: true}) // respond to invite
// _ = app.RequestJoinGroup(ctx, models.GroupJoinRequest{GroupId: gID, UserId: u2.UserId})           // request join
// _ = app.HandleGroupJoinRequest(ctx, models.HandleJoinRequest{GroupId: gID, UserId: u2.UserId, Accept: true}) // handle join request
// _ = app.CancelInviteToGroup(ctx, models.InviteToGroupReq{GroupId: gID, UserId: u2.UserId})       // cancel invite
// _ = app.CancelJoinGroupRequest(ctx, models.GroupJoinRequest{GroupId: gID, UserId: u2.UserId})    // cancel join
// _ = app.RemoveFromGroup(ctx, models.RemoveFromGroupRequest{GroupId: gID, UserId: u2.UserId})     // remove member
// _ = app.LeaveGroup(ctx, models.GeneralGroupReq{GroupId: gID, UserId: u1.UserId})                // leave group

// // ======================
// // SEARCH
// // ======================
// _, _ = app.SearchUsers(ctx, models.UserSearchReq{SearchTerm: "alice", Limit: 10}) // search users
// _, _ = app.SearchGroups(ctx, models.GroupSearchReq{SearchTerm: "Test", Limit: 10}) // search groups

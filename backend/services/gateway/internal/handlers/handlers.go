package handlers

import (
	"net/http"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/media"
	"social-network/shared/gen-go/notifications"
	"social-network/shared/gen-go/posts"
	"social-network/shared/gen-go/users"
	middleware "social-network/shared/go/http-middleware"
	"social-network/shared/go/ratelimit"
)

type Handlers struct {
	CacheService CacheService
	UsersService users.UserServiceClient
	PostsService posts.PostsServiceClient
	ChatService  chat.ChatServiceClient
	MediaService media.MediaServiceClient
	NotifService notifications.NotificationServiceClient
}

func NewHandlers(
	serviceName string,
	CacheService CacheService,
	UsersService users.UserServiceClient,
	PostsService posts.PostsServiceClient,
	ChatService chat.ChatServiceClient,
	MediaService media.MediaServiceClient,
	NotifService notifications.NotificationServiceClient,
) *http.ServeMux {
	handlers := Handlers{
		CacheService: CacheService,
		UsersService: UsersService,
		PostsService: PostsService,
		ChatService:  ChatService,
		MediaService: MediaService,
		NotifService: NotifService,
	}
	return handlers.BuildMux(serviceName)
}

//TODO remove endpoint from chain, and find another way

// BuildMux builds and returns the HTTP request multiplexer with all routes and middleware applied
func (h *Handlers) BuildMux(serviceName string) *http.ServeMux {
	mux := http.NewServeMux()
	ratelimiter := ratelimit.NewRateLimiter(serviceName+":", h.CacheService)
	middlewareObj := middleware.NewMiddleware(ratelimiter, serviceName)
	Chain := middlewareObj.Chain

	IP := middleware.IPLimit
	USERID := middleware.UserLimit

	mux.HandleFunc("/test/{yo}/hello",
		Chain("/test/{yo}/hello").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			EnrichContext().
			Finalize(h.testHandler()))

	mux.HandleFunc("/login",
		Chain("/login").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			EnrichContext().
			Finalize(h.loginHandler()))

	mux.HandleFunc("/register",
		Chain("/register").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			EnrichContext().
			Finalize(h.registerHandler()))

	//fileid url
	mux.HandleFunc("/files/{file_id}/validate",
		Chain("/files/{file_id}/validate").
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.validateFileUpload()))

	//image id variant from url
	mux.HandleFunc("/files/images/{image_id}/{variant}",
		Chain("/files/images/{image_id}/{variant}").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.getImageUrl()))

	mux.HandleFunc("/logout",
		Chain("/logout").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.logoutHandler()))

	// mux.HandleFunc("/auth-status",
	// 	Chain("/auth-status").
	// 		AllowedMethod("GET").
	// 		RateLimit(IP, 20, 5).
	// 		Auth().
	// 		EnrichContext().
	// 		RateLimit(USERID, 20, 5).
	// 		Finalize(h.authStatus()))

	//NEW ENDPOINTS BELOW:

	//START GROUPS ======================================
	//START GROUPS ======================================
	//START GROUPS ======================================
	mux.HandleFunc("/groups",
		Chain("/groups").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.createGroup()))

	//TODO make params
	mux.HandleFunc("/groups",

		Chain("/groups").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getAllGroupsPaginated()))

	//TODO cut url, and json
	mux.HandleFunc("/groups/{group_id}",
		Chain("/groups/{group_id}").
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.updateGroup()))

	//TODO slice url
	mux.HandleFunc("/groups/{group_id}",
		Chain("/groups/{group_id}").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getGroupInfo()))

	//TODO group id from url
	mux.HandleFunc("/groups/{group_id}/popular-post",
		Chain("/groups/{group_id}/pospular-post").
			AllowedMethod("GET").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.getMostPopularPostInGroup()))

	//TODO get params, and groupid from url
	mux.HandleFunc("/groups/{group_id}/members",
		Chain("/groups/{group_id}/members").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getGroupMembers()))

	//TODO group id from url
	mux.HandleFunc("/groups/{group_id}/join-response",
		Chain("/groups/{group_id}/join-response").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.handleGroupJoinRequest()))

	//TODO group id from url
	mux.HandleFunc("/groups/{group_id}/invite",
		Chain("/groups/{group_id}/invite").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.inviteToGroup()))

	//TODO get group id from url
	mux.HandleFunc("/groups/{group_id}/leave",
		Chain("/groups/{group_id}/leave").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.leaveGroup()))

	//TODO get group id from url
	mux.HandleFunc("/groups/{group_id}/remove-member",
		Chain("/groups/{group_id}/remove-member").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.removeFromGroup()))

	//TODO group id url
	mux.HandleFunc("/groups/{group_id}/join-request",
		Chain("/groups/{group_id}/join-request").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.requestJoinGroup()))

	//TODO group id url
	mux.HandleFunc("/groups/{group_id}/cancel-join-request",
		Chain("/groups/{group_id}/cancel-join-request").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.cancelGroupJoinRequest()))

	//groupid from url
	mux.HandleFunc("/groups/{group_id}/invite-response",
		Chain("/groups/{group_id}/invite-response").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.respondToGroupInvite()))

	//todo params and groupid from url
	mux.HandleFunc("/groups/{group_id}/pending-requests",
		Chain("/groups/{group_id}/pending-requests").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPendingGroupJoinRequests()))

	//TODO group id, params
	mux.HandleFunc("/groups/{group_id}/pending-count",
		Chain("/groups/{group_id}/pending-count").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPendingGroupJoinRequestsCount()))

	//group id from url, params
	mux.HandleFunc("groups/{group_id}/invitable-followers",
		Chain("groups/{group_id}/invitable-followers").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetFollowersNotInvitedToGroup()))

	//params, groupid url
	mux.HandleFunc("/groups/{group_id}/search",
		Chain("/groups/{group_id}/search").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.searchGroups()))

	//params, groupid url
	mux.HandleFunc("/groups/{group_id}/posts",
		Chain("/groups/{group_id}/posts").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getGroupPostsPaginated()))

	//TODO get params
	mux.HandleFunc("/my/groups",
		Chain("/my/groups").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getUserGroupsPaginated()))

	//params, groupid url
	mux.HandleFunc("/groups/{group_id}/events",
		Chain("/groups/{group_id}/events").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getEventsByGroupId()))

	//END GROUPS ======================================
	//END GROUPS ======================================
	//END GROUPS ======================================

	//END USERS ======================================
	//END USERS ======================================
	//END USERS ======================================

	// users_id url
	mux.HandleFunc("/users/{users_id}/follow",
		Chain("/users/{users_id}/follow").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.followUser()))

	//params, userid url
	mux.HandleFunc("/users/{user_id}/followers",
		Chain("/users/{user_id}/followers").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getFollowersPaginated()))

	//TODO currently unused in front
	mux.HandleFunc("/my/follow-suggestions",
		Chain("/my/follow-suggestions").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getFollowSuggestions()))

	//TODO not used by front yet
	//params, userid url
	mux.HandleFunc("/users/{user_id}/following",
		Chain("/users/{user_id}/following").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getFollowingPaginated()))

	//get user id url
	mux.HandleFunc("/users/{user_id}/follow-response",
		Chain("/users/{user_id}/follow-response").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.handleFollowRequest()))

	//params
	mux.HandleFunc("/users/search",
		Chain("/users/search").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.searchUsers()))

	//userid url
	mux.HandleFunc("/users/{users_id}/unfollow",
		Chain("/user/{users_id}/unfollow").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.unFollowUser()))

	//user id
	mux.HandleFunc("users/{user_id}/profile",
		Chain("users/{user_id}/profile").
			AllowedMethod("GET").
			RateLimit(IP, 40, 20).
			Auth().
			EnrichContext().
			RateLimit(USERID, 40, 20).
			Finalize(h.getUserProfile()))

	//TODO params, user_id url
	mux.HandleFunc("/users/{user_id}/posts",
		Chain("/users/{user_id}/posts").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getUserPostsPaginated()))

	//done
	mux.HandleFunc("my/profile/privacy",
		Chain("my/profile/privacy").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.updateProfilePrivacy()))

	mux.HandleFunc("my/profile/email",
		Chain("my/profile/email").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.updateUserEmail()))

	mux.HandleFunc("/my/profile/password",
		Chain("/my/profile/password").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.updateUserPassword()))

	mux.HandleFunc("my/profile",
		Chain("my/profile").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.updateUserProfile()))

	//TODO params
	mux.HandleFunc("posts/public",
		Chain("posts/public").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPublicFeed()))

	//TODO params
	mux.HandleFunc("posts/friends",
		Chain("posts/friends").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPersonalizedFeed()))

	//params, postsid url
	mux.HandleFunc("/posts/{post_id}",
		Chain("/post/{post_id}").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPostById()))

	//done
	mux.HandleFunc("/posts",
		Chain("/posts").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.createPost()))

	//post id url
	mux.HandleFunc("/posts/{post_id}",
		Chain("/posts/{post_id}").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.editPost()))

	//post id url
	mux.HandleFunc("/posts/{post_id}",
		Chain("/posts/{post_id}").
			AllowedMethod("DELETE").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.deletePost()))

	mux.HandleFunc("/comments/create",
		Chain("/comments/create").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.createComment()))

	mux.HandleFunc("/comments/edit",
		Chain("/comments/edit").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.editComment()))

	mux.HandleFunc("/comments/delete",
		Chain("/comments/delete").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.deleteComment()))

	mux.HandleFunc("/comments/",
		Chain("/comments/").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getCommentsByParentId()))

	mux.HandleFunc("/events/create",
		Chain("/events/create").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.createEvent()))

	mux.HandleFunc("/events/edit",
		Chain("/events/edit").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.editEvent()))

	mux.HandleFunc("/events/delete",
		Chain("/events/delete").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.deleteEvent()))

	mux.HandleFunc("/events/respond",
		Chain("/events/respond").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.respondToEvent()))

	mux.HandleFunc("/events/remove-response",
		Chain("/events/remove-response").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.RemoveEventResponse()))

	mux.HandleFunc("/reactions/",
		Chain("/reactions/").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.toggleOrInsertReaction()))

	mux.HandleFunc("/notifications/all",
		Chain("/notifications/all").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetUserNotifications()))

	mux.HandleFunc("/notifications/unread",
		Chain("/notifications/unread").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetUnreadNotificationsCount()))

	mux.HandleFunc("/notifications/read",
		Chain("/notifications/read").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.MarkNotificationAsRead()))

	mux.HandleFunc("/notifications/allread",
		Chain("/notifications/allread").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.MarkAllAsRead()))

	mux.HandleFunc("/notifications/delete",
		Chain("/notifications/delete").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.DeleteNotification()))

	mux.HandleFunc("/notifications/preferences",
		Chain("/notifications/preferences").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetNotificationPreferences()))

	// CHAT
	mux.HandleFunc("/chat/create-pm",
		Chain("/chat/create-pm").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.CreatePrivateMsg()))

	mux.HandleFunc("/chat/get-private-conversation-by-id",
		Chain("/chat/get-private-conversation-by-id").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetPrivateConversationById()))

	mux.HandleFunc("/chat/get-private-conversations",
		Chain("/chat/get-private-conversations").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetPrivateConversations()))

	mux.HandleFunc("/chat/get-pms-paginated",
		Chain("/chat/get-pms-paginated").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetPrivateMessagesPag()))

	mux.HandleFunc("/chat/create-group-message",
		Chain("/chat/create-group-message").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.CreateGroupMsg()))

	mux.HandleFunc("/chat/get-group-messages-paginated",
		Chain("/chat/get-group-messages-paginated").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetGroupMessagesPag()))

	mux.HandleFunc("/chat/update-last-read-pm",
		Chain("/chat/update-last-read-pm").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.UpdateLastRead()))

	return mux
}

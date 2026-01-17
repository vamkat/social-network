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

	// mux.HandleFunc("/test/{yo}/hello",
	// 	Chain("/test/{yo}/hello").
	// 		AllowedMethod("GET").
	// 		RateLimit(IP, 20, 5).
	// 		EnrichContext().
	// 		Finalize(h.testHandler()))

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

	//fileid url, --DONE
	mux.HandleFunc("/files/{file_id}/validate",
		Chain("/files/{file_id}/validate").
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.validateFileUpload()))

	//image id variant from url --DONE
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

	//TODO make params --DONE
	mux.HandleFunc("/groups",

		Chain("/groups").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getAllGroupsPaginated()))

	//TODO cut url, and json --DONE
	mux.HandleFunc("/groups/{group_id}",
		Chain("/groups/{group_id}").
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.updateGroup()))

	//TODO slice url --DONE
	mux.HandleFunc("/groups/{group_id}",
		Chain("/groups/{group_id}").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getGroupInfo()))

	//TODO group id from url --DONE
	mux.HandleFunc("/groups/{group_id}/popular-post",
		Chain("/groups/{group_id}/popular-post").
			AllowedMethod("GET").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.getMostPopularPostInGroup()))

	//TODO get params, and groupid from url --DONE
	mux.HandleFunc("/groups/{group_id}/members",
		Chain("/groups/{group_id}/members").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getGroupMembers()))

	//TODO group id from url --DONE
	mux.HandleFunc("/groups/{group_id}/join-response",
		Chain("/groups/{group_id}/join-response").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.handleGroupJoinRequest()))

	//TODO group id from url --DONE
	mux.HandleFunc("/groups/{group_id}/invite",
		Chain("/groups/{group_id}/invite").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.inviteToGroup()))

	//TODO get group id from url --DONE
	mux.HandleFunc("/groups/{group_id}/leave",
		Chain("/groups/{group_id}/leave").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.leaveGroup()))

	//TODO get group id from url --DONE
	mux.HandleFunc("/groups/{group_id}/remove-member",
		Chain("/groups/{group_id}/remove-member").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.removeFromGroup()))

	//TODO group id url --DONE
	mux.HandleFunc("/groups/{group_id}/join-request",
		Chain("/groups/{group_id}/join-request").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.requestJoinGroup()))

	//TODO group id url --DONE
	mux.HandleFunc("/groups/{group_id}/cancel-join-request",
		Chain("/groups/{group_id}/cancel-join-request").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.cancelGroupJoinRequest()))

	//groupid from url --DONE
	mux.HandleFunc("/groups/{group_id}/invite-response",
		Chain("/groups/{group_id}/invite-response").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.respondToGroupInvite()))

	//todo params and groupid from url --DONE
	mux.HandleFunc("/groups/{group_id}/pending-requests",
		Chain("/groups/{group_id}/pending-requests").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPendingGroupJoinRequests()))

	//TODO group id, params --DONE
	mux.HandleFunc("/groups/{group_id}/pending-count",
		Chain("/groups/{group_id}/pending-count").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPendingGroupJoinRequestsCount()))

	//group id from url, params --DONE
	mux.HandleFunc("groups/{group_id}/invitable-followers",
		Chain("groups/{group_id}/invitable-followers").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetFollowersNotInvitedToGroup()))

	//params, groupid url --DONE
	mux.HandleFunc("/groups/search",
		Chain("/groups/search").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.searchGroups()))

	//params, groupid url --DONE
	mux.HandleFunc("/groups/{group_id}/posts",
		Chain("/groups/{group_id}/posts").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getGroupPostsPaginated()))

	//TODO get params --DONE
	mux.HandleFunc("/my/groups",
		Chain("/my/groups").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getUserGroupsPaginated()))

	//params, groupid url --DONE
	mux.HandleFunc("/groups/{group_id}/events",
		Chain("/groups/{group_id}/events").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getEventsByGroupId()))

	//groupid param --DONE
	mux.HandleFunc("groups/{group_id}/events",
		Chain("groups/{group_id}/events").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.createEvent()))

	// USERS ======================================
	// USERS ======================================
	// USERS ======================================

	// users_id url --DONE
	mux.HandleFunc("/users/{user_id}/follow",
		Chain("/users/{user_id}/follow").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.followUser()))

	//params, userid url --DONE
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
	//params, userid url --DONE
	mux.HandleFunc("/users/{user_id}/following",
		Chain("/users/{user_id}/following").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getFollowingPaginated()))

	//get user id url --DONE
	mux.HandleFunc("/users/{user_id}/follow-response",
		Chain("/users/{user_id}/follow-response").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.handleFollowRequest()))

	//params --DONE
	mux.HandleFunc("/users/search",
		Chain("/users/search").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.searchUsers()))

	//userid url --DONE
	mux.HandleFunc("/users/{users_id}/unfollow",
		Chain("/user/{users_id}/unfollow").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.unFollowUser()))

	//user id --DONE
	mux.HandleFunc("users/{user_id}/profile",
		Chain("users/{user_id}/profile").
			AllowedMethod("GET").
			RateLimit(IP, 40, 20).
			Auth().
			EnrichContext().
			RateLimit(USERID, 40, 20).
			Finalize(h.getUserProfile()))

	//TODO params, user_id url --DONE
	mux.HandleFunc("/users/{user_id}/posts",
		Chain("/users/{user_id}/posts").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getUserPostsPaginated()))

	// MY =====================
	// MY =====================
	// MY =====================
	// MY =====================

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

	// POST =====================
	// POST =====================
	// POST =====================
	// POST =====================

	//TODO params --DONE
	mux.HandleFunc("posts/public",
		Chain("posts/public").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPublicFeed()))

	//TODO params --DONE
	mux.HandleFunc("posts/friends",
		Chain("posts/friends").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPersonalizedFeed()))

	//params, postsid url --DONE
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

	//post id url --DONE
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

	// COMMENTS ===================
	// COMMENTS ===================
	// COMMENTS ===================

	mux.HandleFunc("/comments",
		Chain("/comments").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.createComment()))

	//todo commnetid url --DONE
	mux.HandleFunc("/comments/{comment_id}",
		Chain("/comments/{comment_id}").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.editComment()))

	//commentid ulr --DONE
	mux.HandleFunc("/comments/{comment_id}",
		Chain("/comments/{comment_id}").
			AllowedMethod("DELETE").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.deleteComment()))

	//params --DONE
	mux.HandleFunc("/comments",
		Chain("/comments").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getCommentsByParentId()))

	//EVENTS ===========================
	//EVENTS ===========================
	//EVENTS ===========================
	//EVENTS ===========================

	//  eventid url --DONE
	mux.HandleFunc("events/{event_id}",
		Chain("events/{event_id}").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.editEvent()))

	// eventid url --DONE
	mux.HandleFunc("events/{event_id}",
		Chain("events/{event_id}").
			AllowedMethod("DELETE").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.deleteEvent()))

	// --DONE
	mux.HandleFunc("/events/{event_id}/response",
		Chain("/events/{event_id}/response").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.respondToEvent()))

	// --DONE
	mux.HandleFunc("/events/{event_id}/response",
		Chain("/events/{event_id}/response").
			AllowedMethod("DELETE").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.RemoveEventResponse()))

	// --DONE
	mux.HandleFunc("/reactions",
		Chain("/reactions").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.toggleOrInsertReaction()))

	// NOTIFICATIONS =====================
	// NOTIFICATIONS =====================
	// NOTIFICATIONS =====================
	// NOTIFICATIONS =====================

	//TODO remove notification type

	//params //add unread parameter// and only read// maybe read type? --DONE
	mux.HandleFunc("/notifications",
		Chain("/notifications").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetUserNotifications()))

	//--DONE
	mux.HandleFunc("/notifications/mark-all",
		Chain("/notifications/mark-all").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.MarkAllAsRead()))

	//params
	mux.HandleFunc("/notifications/{notification_id}",
		Chain("/notifications/{notification_id}").
			AllowedMethod("DELETE").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.DeleteNotification()))

	// CHAT ============================================
	// CHAT ============================================
	// CHAT ============================================
	// CHAT ============================================
	//--DONE
	mux.HandleFunc("my/chat/{interlocutor_id}",
		Chain("my/chat/{interlocutor_id}").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.CreatePrivateMsg()))

	//--DONE
	mux.HandleFunc("my/chat/{interlocutor_id}",
		Chain("my/chat/{interlocutor_id}").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetPrivateMessagesPag()))

	//conv id url, params --DONE
	mux.HandleFunc("my/chat/{conversation_id}/preview",
		Chain("my/chat/{conversation_id}/preview").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetPrivateConversationById()))

	//params
	mux.HandleFunc("/my/chat/previews",
		Chain("/my/chat/previews").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetPrivateConversations()))

	//group id url
	mux.HandleFunc("groups/{group_id}/chat",
		Chain("groups/{group_id}/chat").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.CreateGroupMsg()))

	mux.HandleFunc("groups/{group_id}/chat",
		Chain("groups/{group_id}/chat").
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetGroupMessagesPag()))

	mux.HandleFunc("my/chat/read",
		Chain("my/chat/read").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.UpdateLastRead()))

	return mux
}

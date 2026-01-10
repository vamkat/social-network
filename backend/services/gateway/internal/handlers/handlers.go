package handlers

import (
	"net/http"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/media"
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
}

func NewHandlers(
	serviceName string,
	CacheService CacheService,
	UsersService users.UserServiceClient,
	PostsService posts.PostsServiceClient,
	ChatService chat.ChatServiceClient,
	MediaService media.MediaServiceClient,
) *http.ServeMux {
	handlers := Handlers{
		CacheService: CacheService,
		UsersService: UsersService,
		PostsService: PostsService,
		ChatService:  ChatService,
		MediaService: MediaService,
	}
	return handlers.BuildMux(serviceName)
}

// BuildMux builds and returns the HTTP request multiplexer with all routes and middleware applied
func (h *Handlers) BuildMux(serviceName string) *http.ServeMux {
	mux := http.NewServeMux()
	ratelimiter := ratelimit.NewRateLimiter(serviceName+":", h.CacheService)
	middlewareObj := middleware.NewMiddleware(ratelimiter, serviceName)
	Chain := middlewareObj.Chain

	IP := middleware.IPLimit
	USERID := middleware.UserLimit

	mux.HandleFunc("/test",
		Chain("/test").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			EnrichContext().
			Finalize(h.testHandler()))

	mux.HandleFunc("/profile/",
		Chain("/profile/").
			AllowedMethod("POST").
			RateLimit(IP, 40, 20).
			Auth().
			EnrichContext().
			RateLimit(USERID, 40, 20).
			Finalize(h.getUserProfile()))

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

	mux.HandleFunc("/validate-file-upload",
		Chain("/validate-file-upload").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			EnrichContext().
			Finalize(h.validateFileUpload()))

	// Test handler for media package
	mux.HandleFunc("/get-image",
		Chain("/get-image").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			EnrichContext().
			Finalize(h.getImageUrl()))

	mux.HandleFunc("/logout",
		Chain("/logout").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.logoutHandler()))

	mux.HandleFunc("/auth-status",
		Chain("/auth-status").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.authStatus()))

	//NEW ENDPOINTS BELOW:

	// Groups
	mux.HandleFunc("/groups/create",
		Chain("/groups/create").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.createGroup()))

	mux.HandleFunc("/groups/update",
		Chain("/groups/update").
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.updateGroup()))

	mux.HandleFunc("/groups/popular",
		Chain("/groups/popular").
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.getMostPopularPostInGroup()))

	mux.HandleFunc("/groups/paginated",
		Chain("/groups/paginated").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getAllGroupsPaginated()))

	// Follow actions
	mux.HandleFunc("/user/follow",
		Chain("/user/follow").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.followUser()))

	mux.HandleFunc("/users/followers/paginated",
		Chain("/users/followers/paginated").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getFollowersPaginated()))

	mux.HandleFunc("/users/follow-suggestions",
		Chain("/users/follow-suggestions").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getFollowSuggestions()))

	mux.HandleFunc("/following",
		Chain("/following").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getFollowingPaginated()))

	mux.HandleFunc("/group/",
		Chain("/group/").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getGroupInfo()))

	mux.HandleFunc("/group/members",
		Chain("/group/members").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getGroupMembers()))

	mux.HandleFunc("/groups/user/",
		Chain("/groups/user/").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getUserGroupsPaginated()))

	mux.HandleFunc("/follow/response",
		Chain("/follow/response").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.handleFollowRequest()))

	mux.HandleFunc("/group/handle-request",
		Chain("/group/handle-request").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.handleGroupJoinRequest()))

	mux.HandleFunc("/group/invite/user",
		Chain("/group/invite/user").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.inviteToGroup()))

	mux.HandleFunc("/group/leave",
		Chain("/group/leave").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.leaveGroup()))

	mux.HandleFunc("/group/remove",
		Chain("/group/remove").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.removeFromGroup()))

	mux.HandleFunc("/group/join",
		Chain("/group/join").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.requestJoinGroup()))

	mux.HandleFunc("/group/cancel-request",
		Chain("/group/cancek-request").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.cancelGroupJoinRequest()))

	mux.HandleFunc("/group/invite/response",
		Chain("/group/invite/response").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.respondToGroupInvite()))

	mux.HandleFunc("/group/pending",
		Chain("/group/pending").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPendingGroupJoinRequests()))

	mux.HandleFunc("/group/pending-count",
		Chain("/group/pending-count").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPendingGroupJoinRequestsCount()))

	mux.HandleFunc("/group/notinvited",
		Chain("/group/notinvited").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.GetFollowersNotInvitedToGroup()))

	mux.HandleFunc("/search/group",
		Chain("/search/group").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.searchGroups()))

	mux.HandleFunc("/users/search",
		Chain("/users/search").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.searchUsers()))

	mux.HandleFunc("/user/unfollow",
		Chain("/user/unfollow").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.unFollowUser()))

	mux.HandleFunc("/account/update/public",
		Chain("/account/update/public").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.updateProfilePrivacy()))

	mux.HandleFunc("/account/update/email",
		Chain("/account/update/email").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.updateUserEmail()))

	mux.HandleFunc("/account/update/password",
		Chain("/account/update/password").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.updateUserPassword()))

	mux.HandleFunc("/profile/update",
		Chain("/profile/update").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.updateUserProfile()))

	// POSTS
	mux.HandleFunc("/public-feed",
		Chain("/public-feed").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPublicFeed()))

	mux.HandleFunc("/friends-feed",
		Chain("/friends-feed").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPersonalizedFeed()))

	mux.HandleFunc("/user/posts",
		Chain("/user/posts").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getUserPostsPaginated()))

	mux.HandleFunc("/group/posts",
		Chain("/group/posts").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getGroupPostsPaginated()))

	mux.HandleFunc("/post/",
		Chain("/post/").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getPostById()))

	mux.HandleFunc("/posts/create",
		Chain("/posts/create").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.createPost()))

	mux.HandleFunc("/posts/edit",
		Chain("/posts/edit").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.editPost()))

	mux.HandleFunc("/posts/delete/",
		Chain("/posts/delete/").
			AllowedMethod("POST").
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

	mux.HandleFunc("/events/",
		Chain("/events/").
			AllowedMethod("POST").
			RateLimit(IP, 20, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 20, 5).
			Finalize(h.getEventsByGroupId()))

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

	return mux
}

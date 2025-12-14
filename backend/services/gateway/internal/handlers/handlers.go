package handlers

import (
	"net/http"
	"social-network/services/gateway/internal/application"
	"social-network/services/gateway/internal/middleware"
	"social-network/shared/go/ratelimit"
)

type Handlers struct {
	App application.GatewayApp
}

func NewHandlers(app application.GatewayApp) (*Handlers, error) {
	handlers := &Handlers{
		App: app,
	}
	return handlers, nil
}

// BuildMux builds and returns the HTTP request multiplexer with all routes and middleware applied
func (h *Handlers) BuildMux(serviceName string) *http.ServeMux {
	mux := http.NewServeMux()
	ratelimiter := ratelimit.NewRateLimiter(serviceName+":", h.App.Redis)
	middlewareObj := middleware.NewMiddleware(ratelimiter, "gateway")
	Chain := middlewareObj.Chain

	IP := middleware.IPLimit
	USERID := middleware.UserLimit

	mux.HandleFunc("/test",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			EnrichContext().
			Finalize(h.testHandler()))

	mux.HandleFunc("/profile/",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 40, 20).
			Auth().
			EnrichContext().
			RateLimit(USERID, 40, 20).
			Finalize(h.getUserProfile()))

	mux.HandleFunc("/login",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			EnrichContext().
			Finalize(h.loginHandler()))

	mux.HandleFunc("/register",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			EnrichContext().
			Finalize(h.registerHandler()))

	mux.HandleFunc("/logout",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.logoutHandler()))

	mux.HandleFunc("/auth-status",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.authStatus()))

	mux.HandleFunc("/public-feed",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.getPublicFeed()))

	//NEW ENDPOINTS BELOW:

	// Groups
	mux.HandleFunc("/groups/create",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.CreateGroup()))

	mux.HandleFunc("/groups/paginated",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetAllGroupsPaginated()))

	// Follow actions
	mux.HandleFunc("/user/follow",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.FollowUser()))

	mux.HandleFunc("/users/followers/paginated",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetFollowersPaginated()))

	mux.HandleFunc("/users/follow-suggestions",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetFollowSuggestions()))

	mux.HandleFunc("/following",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetFollowingPaginated()))

	mux.HandleFunc("/group/",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetGroupInfo()))

	mux.HandleFunc("/group/members",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetGroupMembers()))

	mux.HandleFunc("/groups/user/",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetUserGroupsPaginated()))

	mux.HandleFunc("/follow/response",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.HandleFollowRequest()))

	mux.HandleFunc("/group/handle-request",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.HandleGroupJoinRequest()))

	mux.HandleFunc("/group/invite/user",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.InviteToGroup()))

	mux.HandleFunc("/group/leave",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.LeaveGroup()))

	mux.HandleFunc("/group/join",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.RequestJoinGroupOrCancel()))

	mux.HandleFunc("/group/invite/response",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.RespondToGroupInvite()))

	mux.HandleFunc("/search/group",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.SearchGroups()))

	mux.HandleFunc("/users/search",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.SearchUsers()))

	mux.HandleFunc("/user/unfollow",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.UnFollowUser()))

	mux.HandleFunc("/account/update/public",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.UpdateProfilePrivacy()))

	mux.HandleFunc("/account/update/email",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.UpdateUserEmail()))

	mux.HandleFunc("/account/update/password",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.UpdateUserPassword()))

	mux.HandleFunc("/profile/update",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.UpdateUserProfile()))

	return mux
}

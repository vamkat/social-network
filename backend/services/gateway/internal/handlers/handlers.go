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
	middleware := middleware.NewMiddleware(ratelimiter)
	Chain := middleware.Chain

	mux.HandleFunc("/test",
		Chain().
			AllowedMethod("GET").
			RateLimitIP(5, 5).
			EnrichContext().
			Finalize(h.testHandler()))

	mux.HandleFunc("/profile/",
		Chain().
			AllowedMethod("GET").
			RateLimitIP(40, 20).
			Auth().
			RateLimitUser(40, 20).
			EnrichContext().
			Finalize(h.getUserProfile()))

	mux.HandleFunc("/login",
		Chain().
			AllowedMethod("POST").
			RateLimitIP(5, 5).
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.loginHandler()))

	mux.HandleFunc("/register",
		Chain().
			AllowedMethod("POST").
			RateLimitIP(5, 5).
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.registerHandler()))

	mux.HandleFunc("/logout",
		Chain().
			AllowedMethod("POST").
			RateLimitIP(5, 5).
			Auth().
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.logoutHandler()))

	mux.HandleFunc("/auth-status",
		Chain().
			AllowedMethod("POST").
			RateLimitIP(5, 5).
			Auth().
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.authStatus()))

	mux.HandleFunc("/public-feed",
		Chain().
			AllowedMethod("GET").
			RateLimitIP(5, 5).
			Auth().
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.getPublicFeed()))

	//NEW ENDPOINTS BELOW:

	// Groups
	mux.HandleFunc("/groups/create",
		Chain().
			AllowedMethod("POST").
			RateLimitIP(5, 5).
			Auth().
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.CreateGroup()))

	mux.HandleFunc("/groups/paginated",
		Chain().
			AllowedMethod("GET").
			RateLimitIP(5, 5).
			Auth().
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.GetAllGroupsPaginated()))

	// Follow actions
	mux.HandleFunc("/users/follow",
		Chain().
			AllowedMethod("POST").
			RateLimitIP(5, 5).
			Auth().
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.FollowUser()))

	mux.HandleFunc("/users/followers/paginated",
		Chain().
			AllowedMethod("GET").
			RateLimitIP(5, 5).
			Auth().
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.GetFollowersPaginated()))

	mux.HandleFunc("/users/follow-suggestions",
		Chain().
			AllowedMethod("GET").
			RateLimitIP(5, 5).
			Auth().
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.GetFollowSuggestions()))

	// Basic user info
	mux.HandleFunc("/users/basic-info",
		Chain().
			AllowedMethod("GET").
			RateLimitIP(5, 5).
			Auth().
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.GetBasicUserInfo()))

	mux.HandleFunc("/users/basic-info/batch",
		Chain().
			AllowedMethod("GET").
			RateLimitIP(5, 5).
			Auth().
			RateLimitUser(5, 5).
			EnrichContext().
			Finalize(h.GetBatchBasicUserInfo()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.GetFollowingIds()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.GetFollowingPaginated()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.GetGroupInfo()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.GetGroupMembers()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.GetUserGroupsPaginated()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.HandleFollowRequest()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.HandleGroupJoinRequest()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.InviteToGroup()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.LeaveGroup()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.RequestJoinGroupOrCancel()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.RespondToGroupInvite()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.SearchGroups()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.SearchUsers()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.UnFollowUser()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.UpdateProfilePrivacy()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.UpdateUserEmail()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.UpdateUserPassword()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimitIP(5, 5).
	// 		Auth().
	// 		RateLimitUser(5, 5).
	// 		EnrichContext().
	// 		Finalize(h.UpdateUserProfile()))

	return mux
}

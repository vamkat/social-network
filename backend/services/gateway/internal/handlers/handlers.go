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

	IP := middleware.GlobalLimit
	USERID := middleware.UserLimit

	mux.HandleFunc("/test",
		Chain().
			AllowedMethod("GET").
			RateLimit(IP, 5, 5).
			EnrichContext().
			Finalize(h.testHandler()))

	mux.HandleFunc("/profile/",
		Chain().
			AllowedMethod("GET").
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
			RateLimit(USERID, 5, 5).
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
			AllowedMethod("GET").
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
			AllowedMethod("GET").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetAllGroupsPaginated()))

	// Follow actions
	mux.HandleFunc("/users/follow",
		Chain().
			AllowedMethod("POST").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.FollowUser()))

	mux.HandleFunc("/users/followers/paginated",
		Chain().
			AllowedMethod("GET").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetFollowersPaginated()))

	mux.HandleFunc("/users/follow-suggestions",
		Chain().
			AllowedMethod("GET").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetFollowSuggestions()))

	// Basic user info
	mux.HandleFunc("/users/basic-info",
		Chain().
			AllowedMethod("GET").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetBasicUserInfo()))

	mux.HandleFunc("/users/basic-info/batch",
		Chain().
			AllowedMethod("GET").
			RateLimit(IP, 5, 5).
			Auth().
			EnrichContext().
			RateLimit(USERID, 5, 5).
			Finalize(h.GetBatchBasicUserInfo()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.GetFollowingIds()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.GetFollowingPaginated()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.GetGroupInfo()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.GetGroupMembers()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.GetUserGroupsPaginated()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.HandleFollowRequest()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.HandleGroupJoinRequest()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.InviteToGroup()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.LeaveGroup()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.RequestJoinGroupOrCancel()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.RespondToGroupInvite()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.SearchGroups()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.SearchUsers()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.UnFollowUser()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.UpdateProfilePrivacy()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.UpdateUserEmail()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.UpdateUserPassword()))

	// mux.HandleFunc("/public-feed",
	// 	Chain().
	// 		AllowedMethod("GET").
	// 		RateLimit(IP,5, 5).
	// 		Auth()
	// 		EnrichContext().
	// 		RateLimit(USERID, 5, 5).
	// 		Finalize(h.UpdateUserProfile()))

	return mux
}

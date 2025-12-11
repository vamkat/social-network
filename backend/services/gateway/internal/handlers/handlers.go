package handlers

import (
	"errors"
	"net/http"
	"social-network/services/gateway/internal/application"
	"social-network/services/gateway/internal/middleware"
	"social-network/shared/go/ratelimit"
)

var ErrMinIO = errors.New("minIO failed to be created")

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
			RateLimitIP(1, 4).
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
			RateLimitIP(1, 4).
			RateLimitUser(1, 4).
			EnrichContext().
			Finalize(h.loginHandler()))

	mux.HandleFunc("/register",
		Chain().
			AllowedMethod("POST").
			RateLimitIP(1, 4).
			RateLimitUser(1, 4).
			EnrichContext().
			Finalize(h.registerHandler()))

	mux.HandleFunc("/logout",
		Chain().
			AllowedMethod("POST").
			RateLimitIP(1, 4).
			Auth().
			RateLimitUser(1, 4).
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

	return mux
}

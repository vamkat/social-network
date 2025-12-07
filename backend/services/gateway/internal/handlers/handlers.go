package handlers

import (
	"net/http"
	"social-network/services/gateway/internal/middleware"
	remoteservices "social-network/services/gateway/internal/remote_services"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/ratelimit"
	redis_connector "social-network/shared/go/redis"

	"github.com/minio/minio-go/v7"
)

type Handlers struct {
	Services    *remoteservices.GRpcServices
	MinIOClient *minio.Client
	redisClient *redis_connector.RedisClient
}

func NewHandlers(redisClient *redis_connector.RedisClient) *Handlers {
	handlers := &Handlers{
		redisClient: redisClient,
		Services:    remoteservices.NewServices([]ct.CtxKey{ct.UserId, ct.ReqID, ct.TraceId}),
		MinIOClient: remoteservices.NewMinIOConn(),
	}

	return handlers
}

// BuildMux builds and returns the HTTP request multiplexer with all routes and middleware applied
func (h *Handlers) BuildMux() *http.ServeMux {
	mux := http.NewServeMux()
	ratelimiter := ratelimit.NewRateLimiter("api:", h.redisClient)
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

	return mux
}

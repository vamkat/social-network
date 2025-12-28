package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"slices"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	ct "social-network/shared/go/ct"
	tele "social-network/shared/go/telemetry"

	"strings"
)

type ratelimiter interface {
	Allow(ctx context.Context, key string, limit int, durationSeconds int64) (bool, error)
}

type middleware struct {
	ratelimiter ratelimiter
	serviceName string
}

func NewMiddleware(ratelimiter ratelimiter, serviceName string) *middleware {
	return &middleware{
		ratelimiter: ratelimiter,
		serviceName: serviceName,
	}
}

// MiddleSystem holds the middleware chain
type MiddleSystem struct {
	middlewareChain []func(http.ResponseWriter, *http.Request) (bool, *http.Request)
	ratelimiter     ratelimiter
	serviceName     string
	endpoint        string
}

// Chain initializes a new middleware chain
func (m *middleware) Chain(endpoint string) *MiddleSystem {
	return &MiddleSystem{
		ratelimiter: m.ratelimiter,
		serviceName: m.serviceName,
		endpoint:    endpoint,
	}
}

// add appends a middleware function to the chain
func (m *MiddleSystem) add(f func(http.ResponseWriter, *http.Request) (bool, *http.Request)) {
	m.middlewareChain = append(m.middlewareChain, f)
}

// AllowedMethod sets allowed HTTP methods and handles CORS preflight requests
func (m *MiddleSystem) AllowedMethod(methods ...string) *MiddleSystem {
	m.add(func(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
		ctx := r.Context()
		tele.Debug(ctx, fmt.Sprint("endpoint called:", r.URL.Path, " with method: ", r.Method))
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", ")+", OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-Id, X-Timestamp, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// TODO fix this, return cors to be
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", ")+", OPTIONS")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-Id, X-Timestamp")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			tele.Debug(ctx, "Method in options")
			w.WriteHeader(http.StatusNoContent) // 204
			return false, nil
		}

		if slices.Contains(methods, r.Method) {
			return true, r
		}

		// method not allowed
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		tele.Debug(ctx, "method not allowed")
		return false, nil
	})
	return m
}

// EnrichContext adds request ID and trace ID to the request context
func (m *MiddleSystem) EnrichContext() *MiddleSystem {
	m.add(func(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
		r = utils.RequestWithValue(r, ct.ReqID, utils.GenUUID())
		r = utils.RequestWithValue(r, ct.TraceId, utils.GenUUID())
		return true, r
	})
	return m
}

// Auth middleware to validate JWT and enrich context with claims
func (m *MiddleSystem) Auth() *MiddleSystem {
	m.add(func(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
		ctx := r.Context()

		tele.Debug(ctx, "in auth")
		// tele.Info(ctx, "Cookies received:", r.Cookies())
		cookie, err := r.Cookie("jwt")
		if err != nil {
			tele.Warn(ctx, "no cookie")
			utils.ErrorJSON(ctx, w, http.StatusUnauthorized, "missing auth cookie")
			return false, nil
		}
		// tele.Info(ctx, "JWT cookie value:", cookie.Value)
		claims, err := security.ParseAndValidate(cookie.Value)
		if err != nil {
			tele.Warn(ctx, "unauthorized")
			utils.ErrorJSON(ctx, w, http.StatusUnauthorized, err.Error())
			return false, nil
		}
		// enrich request with claims
		tele.Debug(ctx, "auth ok")
		r = utils.RequestWithValue(r, ct.ClaimsKey, claims)
		r = utils.RequestWithValue(r, ct.UserId, claims.UserId)
		tele.Debug(ctx, fmt.Sprint("adding these to context -> claims:", claims, ", userid: ", claims.UserId))
		return true, r
	})
	return m
}

// // BindReqMeta binds request metadata to context
// func (m *MiddleSystem) BindReqMeta() *MiddleSystem {
// 	m.add(func(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
// 		r = utils.RequestWithValue(r, ct.ReqActionDetails, r.Header.Get("X-Action-Details"))
// 		r = utils.RequestWithValue(r, ct.ReqTimestamp, r.Header.Get("X-Timestamp"))
// 		return true, r
// 	})
// 	return m
// }

type rateLimitType string

const (
	UserLimit rateLimitType = "user"
	IPLimit   rateLimitType = "ip"
)

func (m *MiddleSystem) RateLimit(rateLimitType rateLimitType, limit int, durationSeconds int64) *MiddleSystem {
	m.add(func(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
		ctx := r.Context()
		tele.Debug(ctx, "in ratelimit, type: "+fmt.Sprint(rateLimitType))
		rateLimitKey := ""
		switch rateLimitType {
		case IPLimit:
			remoteIp, err := getRemoteIpKey(r)
			if err != nil {
				tele.Debug(ctx, "rate limited remoteIp:"+remoteIp)
				utils.ErrorJSON(ctx, w, http.StatusNotAcceptable, "your IP is absolutely WACK")
				return false, nil
			}
			rateLimitKey = fmt.Sprintf("%s:%s:ip:%s", m.serviceName, m.endpoint, remoteIp)
		case UserLimit:
			userId, ok := ctx.Value(ct.UserId).(int64)
			if !ok {
				tele.Warn(ctx, "err or no userId:", userId, " ok:", ok)
				utils.ErrorJSON(ctx, w, http.StatusNotAcceptable, "how the hell did you end up here without a user id?")
				return false, nil
			}
			tele.Debug(ctx, "rate limited userId:"+fmt.Sprint(userId))
			rateLimitKey = fmt.Sprintf("%s:%s:id:%d", m.serviceName, m.endpoint, userId)
		default:
			panic("bad rate limit type argument!")
		}

		ok, err := m.ratelimiter.Allow(ctx, rateLimitKey, limit, durationSeconds)
		if err != nil {
			tele.Debug(ctx, "rate limited userId:"+fmt.Sprint(rateLimitKey))
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "you broke the rate limiter")
			return false, nil
		}
		if !ok {
			tele.Debug(ctx, fmt.Sprintf("rate limited for key: %s, reached max: %d per %d sec \n", rateLimitKey, limit, durationSeconds), "rateLimitKey", rateLimitKey, "limit", limit, "durationSeconds", durationSeconds)
			utils.ErrorJSON(ctx, w, http.StatusTooManyRequests, "429: stop it, get some help")
			return false, nil
		}
		return true, r
	})
	return m
}

func getRemoteIpKey(r *http.Request) (string, error) {
	remoteIp, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIp = r.RemoteAddr
	}
	if remoteIp == "" {
		//ip is broken somehow
		return "", nil
	}
	return remoteIp, nil
}

// Finalize constructs the final http.HandlerFunc with all middleware applied
func (m *MiddleSystem) Finalize(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, mw := range m.middlewareChain {
			proceed, newReq := mw(w, r)
			r = newReq
			if !proceed {
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

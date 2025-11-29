package middleware

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"social-network/gateway/internal/security"
	"social-network/gateway/internal/utils"
	"strings"
)

type MiddleSystem struct {
	middlewareChain []func(http.ResponseWriter, *http.Request) (bool, *http.Request)
}

func Chain() *MiddleSystem {
	return &MiddleSystem{}
}

func (m *MiddleSystem) add(f func(http.ResponseWriter, *http.Request) (bool, *http.Request)) {
	m.middlewareChain = append(m.middlewareChain, f)
}

// CORS + method gating
func (m *MiddleSystem) AllowedMethod(methods ...string) *MiddleSystem {
	m.add(func(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
		fmt.Println("endpoint: ", r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "https://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", ")+", OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-Id, X-Timestamp")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			fmt.Println("Method in options")
			w.WriteHeader(http.StatusNoContent) // 204
			return false, r
		}

		if slices.Contains(methods, r.Method) {
			return true, r
		}

		// method not allowed
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		fmt.Println("method not allowed")
		return false, r
	})
	return m
}

// Enrich context with request ID
func (m *MiddleSystem) EnrichContext() *MiddleSystem {
	m.add(func(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
		ctx := context.WithValue(r.Context(), "requestId", utils.GenUUID())
		return true, r.WithContext(ctx)
	})
	return m
}

// Bearer auth
func (m *MiddleSystem) Auth() *MiddleSystem {
	m.add(func(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
		fmt.Println("in auth")
		// fmt.Println("Cookies received:", r.Cookies())
		cookie, err := r.Cookie("jwt")
		if err != nil {
			fmt.Println("no cookie")
			utils.ErrorJSON(w, http.StatusUnauthorized, "missing auth cookie")
			return false, r
		}
		// fmt.Println("JWT cookie value:", cookie.Value)
		claims, err := security.ParseAndValidate(cookie.Value)
		if err != nil {
			fmt.Println("unauthorized")
			utils.ErrorJSON(w, http.StatusUnauthorized, err.Error())
			return false, r
		}
		// enrich request with claims
		fmt.Println("auth ok")
		r = utils.RequestWithValue(r, utils.ClaimsKey, claims)
		return true, r
	})
	return m
}

// Bind request meta into context
func (m *MiddleSystem) BindReqMeta() *MiddleSystem {
	m.add(func(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
		rid := r.Header.Get("X-Request-Id")
		act := r.Header.Get("X-Action-Details")
		ts := r.Header.Get("X-Timestamp")

		ctx := context.WithValue(r.Context(), utils.ReqId, rid)
		ctx = context.WithValue(ctx, utils.ReqActionDetails, act)
		ctx = context.WithValue(ctx, utils.ReqTimestamp, ts)

		return true, r.WithContext(ctx)
	})
	return m
}

// Build the final handler
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

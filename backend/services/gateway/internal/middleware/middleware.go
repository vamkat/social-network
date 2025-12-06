package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	ct "social-network/shared/go/customtypes"

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
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", ")+", OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-Id, X-Timestamp")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

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
		r = utils.RequestWithValue(r, ct.ReqID, utils.GenUUID())
		r = utils.RequestWithValue(r, ct.TraceId, utils.GenUUID())
		return true, r
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
		r = utils.RequestWithValue(r, ct.ClaimsKey, claims)
		r = utils.RequestWithValue(r, ct.UserId, claims.UserId)
		return true, r
	})
	return m
}

// Bind request meta into context
func (m *MiddleSystem) BindReqMeta() *MiddleSystem {
	m.add(func(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
		r = utils.RequestWithValue(r, ct.ReqID, r.Header.Get(ct.ReqID))
		r = utils.RequestWithValue(r, ct.ReqActionDetails, r.Header.Get("X-Action-Details"))
		r = utils.RequestWithValue(r, ct.ReqTimestamp, r.Header.Get("X-Timestamp"))
		return true, r
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

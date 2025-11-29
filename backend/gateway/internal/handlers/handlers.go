package handlers

import (
	"net/http"
	"social-network/gateway/internal/middleware"
	remoteservices "social-network/gateway/internal/remote_services"
)

type Handlers struct {
	Services remoteservices.GRpcServices
}

func NewHandlers() Handlers {
	handlers := Handlers{}
	handlers.Services = remoteservices.NewServices()
	return handlers
}

// BuildMux builds and returns the HTTP request multiplexer with all routes and middleware applied
func (h *Handlers) BuildMux() *http.ServeMux {
	mux := http.NewServeMux()
	Chain := middleware.Chain

	mux.HandleFunc("/test", Chain().AllowedMethod("GET").EnrichContext().Finalize(h.testHandler()))
	mux.HandleFunc("/user", Chain().AllowedMethod("GET").EnrichContext().Finalize(h.getBasicUserInfo()))
	mux.HandleFunc("/login", Chain().AllowedMethod("POST").EnrichContext().Finalize(h.loginHandler()))
	return mux
}

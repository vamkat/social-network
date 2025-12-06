package handlers

import (
	"net/http"
	"social-network/gateway/internal/middleware"
	remoteservices "social-network/gateway/internal/remote_services"
	ct "social-network/shared/go/customtypes"

	"github.com/minio/minio-go/v7"
)

type Handlers struct {
	Services    remoteservices.GRpcServices
	MinIOClient *minio.Client
}

func NewHandlers() Handlers {
	handlers := Handlers{}
	handlers.Services = remoteservices.NewServices([]ct.CtxKey{ct.UserId, ct.ReqID, ct.TraceId})
	handlers.MinIOClient = remoteservices.NewMinIOConn()
	return handlers
}

// BuildMux builds and returns the HTTP request multiplexer with all routes and middleware applied
func (h *Handlers) BuildMux() *http.ServeMux {
	mux := http.NewServeMux()
	Chain := middleware.Chain

	mux.HandleFunc("/test", Chain().AllowedMethod("GET").EnrichContext().Finalize(h.testHandler()))
	mux.HandleFunc("/profile/", Chain().AllowedMethod("GET").EnrichContext().Auth().Finalize(h.getUserProfile()))
	mux.HandleFunc("/login", Chain().AllowedMethod("POST").EnrichContext().Finalize(h.loginHandler()))
	mux.HandleFunc("/register", Chain().AllowedMethod("POST").EnrichContext().Finalize(h.registerHandler()))
	mux.HandleFunc("/logout", Chain().AllowedMethod("POST").EnrichContext().Auth().Finalize(h.logoutHandler()))
	mux.HandleFunc("/auth-status", Chain().AllowedMethod("POST").EnrichContext().Auth().Finalize(h.authStatus()))

	return mux
}

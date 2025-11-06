package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.Health)
	r.Get("/test", h.Test)

	return r
}

func (h *AuthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Auth service is healthy"))
}

func (h *AuthHandler) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Auth test endpoint"))
}

package handlers

import (
	"net/http"
	"social-network/gateway/middleware"
)

type Handlers struct {
}

func (h *Handlers) SetHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	Chain := middleware.Chain

	// conn, _ := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	// client := pb.NewYourServiceClient(conn)

	mux.HandleFunc("/test", Chain().AllowedMethod("GET").Finalize(h.testHandler()))
	mux.HandleFunc("/user", Chain().AllowedMethod("GET").Finalize(h.getBasicUserInfo()))
	return mux
}

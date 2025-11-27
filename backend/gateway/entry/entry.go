package entry

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"social-network/gateway/handlers"
	"syscall"
	"time"
)

// server starting sequence
func Start() {

	// set handlers
	handlers := handlers.Handlers{}

	// set server
	var server http.Server
	server.Handler = handlers.SetHandlers()
	server.Addr = "localhost:8081" //todo get from a config file or something

	go func() {
		log.Printf("Server running on https://%s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServeTLS failed: %v", err)
		}
	}()

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful server Shutdown Failed: %v", err)
	}
	log.Println("Server stopped")
}

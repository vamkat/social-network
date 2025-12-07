package entry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"social-network/services/gateway/internal/handlers"
	redis_connector "social-network/shared/go/redis"

	"syscall"
	"time"
)

// server starting sequence
func Start() {

	redisClient := redis_connector.NewRedisClient("redis:6379", "", 0)
	err := redisClient.TestRedisConnection()
	if err != nil {
		log.Fatalf("connection test failed: %v", err)
	}
	fmt.Println("redis connection started correctly")

	// set handlers
	handlers := handlers.NewHandlers(redisClient)

	// start gRPC connections
	deferMe, err := handlers.Services.StartConnections()
	if err != nil {
		log.Fatalf("failed to start gRPC services connections: %v", err)
	}
	defer deferMe()

	fmt.Println("gRPC services connections started")

	// set server
	var server http.Server
	server.Handler = handlers.BuildMux()
	server.Addr = "0.0.0.0:8081" //todo get from a config file or something

	go func() {
		log.Printf("Starting server on http://%s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve: %v", err)
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

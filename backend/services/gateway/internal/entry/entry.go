package entry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"social-network/services/gateway/internal/handlers"
	remoteservices "social-network/services/gateway/internal/remote_services"
	ct "social-network/shared/go/customtypes"
	redis_connector "social-network/shared/go/redis"
	tele "social-network/shared/go/telemetry"

	"syscall"
	"time"
)

// server starting sequence
func Start() {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	/*

		==============================
		         OPEN TELEMETRY
		==============================
	*/
	otelShutdown, err := tele.SetupOTelSDK(ctx)
	if err != nil {
		log.Fatal("open telemetry sdk failed, ERROR:", err.Error())
	}
	fmt.Println("open telemetry ready")

	defer func() {
		err := otelShutdown(context.Background())
		if err != nil {
			log.Println("otel shutdown ungracefully! ERROR: " + err.Error())
		} else {
			log.Println("otel shutdown gracefully")
		}
	}()

	/*

		==============================
		          REDIS SETUP
		==============================
	*/
	redisClient := redis_connector.NewRedisClient("redis:6379", "", 0)
	err = redisClient.TestRedisConnection()
	if err != nil {
		log.Fatalf("connection test failed, ERROR: %v", err)
	}
	fmt.Println("redis connection started correctly")

	/*

		==============================
		      REMOTE gRPC SERVICES
		==============================
	*/
	gRpcServices := remoteservices.NewServices([]ct.CtxKey{ct.UserId, ct.ReqID, ct.TraceId})
	deferMe, err := gRpcServices.StartConnections()
	if err != nil {
		log.Fatalf("failed to start gRPC services connections: %v", err)
	}
	defer deferMe()
	fmt.Println("gRPC services connections started")

	/*

		==============================
		        HANDLER + ROUTER
		==============================
	*/
	apiHandlers, err := handlers.NewHandlers(redisClient, gRpcServices)
	if err != nil {
		log.Fatal("Can't create handlers, ERROR:", err)
	}
	apiMux := apiHandlers.BuildMux()

	/*

		==============================
		         HTTP SERVER
		==============================
	*/
	var server = &http.Server{
		Handler:     apiMux,
		Addr:        "0.0.0.0:8081",
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	srvErr := make(chan error, 1)
	go func() {
		log.Printf("Starting server on http://%s\n", server.Addr)
		srvErr <- server.ListenAndServe()
	}()

	/*

		==============================
		         SHUTDOWN LOGIC
		==============================
	*/
	select {
	case err = <-srvErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	case <-ctx.Done():
		stopSignal() //may be redundant
	}

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Graceful server Shutdown Failed: %v", err)
	}

	log.Println("Server stopped")
}

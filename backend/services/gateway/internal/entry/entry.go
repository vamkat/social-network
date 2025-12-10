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
	"social-network/services/gateway/internal/application"
	"social-network/services/gateway/internal/handlers"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/gorpc"

	redis_connector "social-network/shared/go/redis"
	tele "social-network/shared/go/telemetry"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

var logger = otelslog.NewLogger("api-gateway")

// server starting sequence
func Start() {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	gatewayApplication := application.GatewayApp{}

	defer setupOpenTelemetry(ctx)

	/* ==============================
	          REDIS SETUP
	==============================*/
	gatewayApplication.Redis = redis_connector.NewRedisClient("redis:6379", "", 0)
	err := gatewayApplication.Redis.TestRedisConnection()
	if err != nil {
		log.Fatalf("connection test failed, ERROR: %v", err)
	}
	fmt.Println("redis connection started correctly")

	/*==============================
	      REMOTE gRPC SERVICES
	==============================*/
	gatewayApplication.Users, err = gorpc.GetGRpcClient(users.NewUserServiceClient, "users:50051", []ct.CtxKey{ct.UserId, ct.ReqID, ct.TraceId})
	if err != nil {
		log.Fatalf("failed to connect to users service: %v", err)
	}
	gatewayApplication.Chat, err = gorpc.GetGRpcClient(chat.NewChatServiceClient, "chat:50051", []ct.CtxKey{ct.UserId, ct.ReqID, ct.TraceId})
	if err != nil {
		log.Fatalf("failed to connect to chat service: %v", err)
	}

	/*

		==============================
		         PROMETHEUS EXPORTER SERVER
		==============================
	*/
	// TODO finish setting up promethius exporting endpoint
	// metricsServer := &http.Server{
	// 	Addr:    ":2222",
	// 	Handler: nil, // Will be set below
	// }

	/*

		==============================
		        HANDLER + ROUTER
		==============================
	*/
	fmt.Println(gatewayApplication)
	apiHandlers, err := handlers.NewHandlers(gatewayApplication)
	if err != nil {
		log.Fatal("Can't create handlers, ERROR:", err)
	}
	apiMux := apiHandlers.BuildMux("gateway")

	/*
		==============================
		         API GATEWAY HTTP SERVER
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

func setupOpenTelemetry(ctx context.Context) func() {
	otelShutdown, err := tele.SetupOTelSDK(ctx)
	if err != nil {
		log.Fatal("open telemetry sdk failed, ERROR:", err.Error())
	}
	fmt.Println("open telemetry ready")

	return func() {
		err := otelShutdown(context.Background())
		if err != nil {
			log.Println("otel shutdown ungracefully! ERROR: " + err.Error())
		} else {
			log.Println("otel shutdown gracefully")
		}
	}
}

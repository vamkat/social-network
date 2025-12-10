package entry

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"social-network/services/users/internal/application"
	"social-network/services/users/internal/db/sqlc"
	"social-network/services/users/internal/handler"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/gorpc"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run() error {
	pool, err := connectToDb(context.Background())
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to users-db database")

	clients := InitClients()
	app := application.NewApplication(sqlc.New(pool), pool, clients)
	service := *handler.NewUsersHanlder(app)

	log.Println("Running gRpc service...")
	//UserServiceServer
	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	grpcServer, err := gorpc.RunGRPCServer(service, []ct.CtxKey{ct.UserId, ct.ReqID, ct.TraceId})
	if err != nil {
		log.Fatalf("couldn't start gRpc Server: %s", err.Error())
	}

	users.RegisterUserServiceServer(grpcServer, &service)

	log.Printf("gRPC server listening on %s", port)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
	return nil

}

func connectToDb(ctx context.Context) (pool *pgxpool.Pool, err error) {
	connStr := os.Getenv("DATABASE_URL")
	for i := range 10 {
		pool, err = pgxpool.New(ctx, connStr)
		if err == nil {
			break
		}
		log.Printf("DB not ready yet (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return pool, err
}

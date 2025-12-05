package entry

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"social-network/services/posts/internal/application"
	"social-network/services/posts/internal/db/sqlc"
	"social-network/services/posts/internal/handler"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func Run() error {
	pool, err := connectToDb(context.Background())
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to posts-db database")

	clients := InitClients()
	app := application.NewApplication(sqlc.New(pool), pool, clients)

	service := handler.NewPostsHandler(app)

	log.Println("Running gRpc service...")
	grpc := RunGRPCServer(service)

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	grpc.GracefulStop()
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

// RunGRPCServer starts the gRPC server and blocks
func RunGRPCServer(s *handler.PostsHandler) *grpc.Server {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	grpcServer := grpc.NewServer()

	// pb.RegisterPostsServiceServer(grpcServer, s)

	log.Printf("gRPC server listening on %s", s.Port)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	return grpcServer
}

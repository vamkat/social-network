package entry

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"social-network/services/media/internal/application"
	"social-network/services/media/internal/client"
	"social-network/services/media/internal/configs"
	"social-network/services/media/internal/db/dbservice"
	"social-network/services/media/internal/handler"
	pb "social-network/shared/gen-go/media"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/gorpc"

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

	log.Println("Connected to chat database")

	cfgs := configs.Config{
		Port: ":50051",
		FileService: configs.FileService{
			Buckets: configs.Buckets{
				Originals: "uploads-originals",
				Variants:  "uploads-variants",
			},
			Endpoint:  os.Getenv("MINIO_ENDPOINT"),
			AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
			Secret:    os.Getenv("MINIO_SECRET_KEY"),
		},
	}

	fileServiceClient, err := NewMinIOConn(cfgs.FileService)
	if err != nil {
		return err
	}

	app := application.NewMediaService(
		pool,
		&client.Clients{MinIOClient: fileServiceClient},
		dbservice.New(pool),
		cfgs,
	)

	service := &handler.MediaHandler{
		Application: app,
		Port:        cfgs.Port,
	}

	log.Println("Running gRpc service...")

	grpc, err := RunGRPCServer(service)
	if err != nil {
		return err
	}

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
func RunGRPCServer(s *handler.MediaHandler) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	customUnaryInterceptor, err := gorpc.UnaryServerInterceptorWithContextKeys([]gorpc.StringableKey{ct.UserId, ct.ReqID, ct.TraceId}...)
	if err != nil {
		return nil, err
	}
	customStreamInterceptor, err := gorpc.StreamServerInterceptorWithContextKeys([]gorpc.StringableKey{ct.UserId, ct.ReqID, ct.TraceId}...)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(customUnaryInterceptor),
		grpc.StreamInterceptor(customStreamInterceptor),
	)

	pb.RegisterMediaServiceServer(grpcServer, s)

	log.Printf("gRPC server listening on %s", s.Port)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	return grpcServer, nil
}

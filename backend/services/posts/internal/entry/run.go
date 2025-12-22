package entry

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"social-network/services/posts/internal/application"
	"social-network/services/posts/internal/db/sqlc"
	"social-network/services/posts/internal/handler"
	"social-network/shared/gen-go/posts"
	contextkeys "social-network/shared/go/context-keys"
	"social-network/shared/go/gorpc"
	postgresql "social-network/shared/go/postgre"
	"syscall"
)

func Run() error {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

	dbUrl := os.Getenv("DATABASE_URL")
	pool, err := postgresql.NewPool(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to posts-db database")

	clients := InitClients()
	app, err := application.NewApplication(sqlc.New(pool), pool, clients)
	if err != nil {
		return fmt.Errorf("failed to create posts application: %v", err)
	}

	service := handler.NewPostsHandler(app)

	log.Println("Running gRpc service...")
	startServerFunc, endServerFunc, err := gorpc.CreateGRpcServer[posts.PostsServiceServer](posts.RegisterPostsServiceServer, service, ":50051", contextkeys.CommonKeys())
	if err != nil {
		return err
	}
	defer endServerFunc()

	go func() {
		err := startServerFunc()
		if err != nil {
			log.Fatal("server failed to start")
		}
		fmt.Println("server finished")
	}()

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	endServerFunc()
	log.Println("Server stopped")
	return nil

}

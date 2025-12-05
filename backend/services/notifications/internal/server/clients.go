package server

import (
	"fmt"
	"log"
	"social-network/shared/ports"
	"time"

	usersPb "social-network/shared/gen-go/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

// Initialize connections to other services. Each one is called with dial options
func (s *Server) InitClients() {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
        "loadBalancingConfig": [{"round_robin":{}}]
    	}`),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 2 * time.Second,
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.2,
				Jitter:     0.5,
				MaxDelay:   5 * time.Second,
			},
		}),
	}

	// List of initializer functions for connecting to other services
	initializers := []func(opts []grpc.DialOption) error{
		// Add initialization functions for services this one needs to connect to
		s.InitUsersClient,
		// s.InitPostsClient,  // Comment out until posts service exists
	}

	for _, initFn := range initializers {
		if err := initFn(dialOpts); err != nil {
			log.Printf("Failed to initialize client: %v", err)
		}
	}
}

// Connects to users service
func (s *Server) InitUsersClient(opts []grpc.DialOption) (err error) {
	conn, err := grpc.NewClient(ports.Users, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial user service: %v", err)
	}

	s.Clients.UsersClient = usersPb.NewUserServiceClient(conn)
	return nil
}

// Connects to posts service (comment out until posts service exists)
// func (s *Server) InitPostsClient(opts []grpc.DialOption) (err error) {
// 	conn, err := grpc.NewClient(ports.Posts, opts...)
// 	if err != nil {
// 		return fmt.Errorf("failed to dial posts service: %v", err)
// 	}
//
//	s.Clients.PostsClient = postsPb.NewPostsServiceClient(conn)
// 	return nil
// }

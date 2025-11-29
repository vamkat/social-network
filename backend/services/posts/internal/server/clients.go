/*
Establish connections to other services
*/

package server

import (
	"fmt"
	// userpb "social-network/shared/gen/users"
	userpb "social-network/shared/gen-go/users"
	"social-network/shared/ports"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

// Initialize connections to clients
func (s *Server) InitClients() {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
        "loadBalancingConfig": [{"round_robin":{}}]
    	}`),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 2 * time.Second,
			Backoff: backoff.Config{
				BaseDelay: 1 * time.Second,
				MaxDelay:  5 * time.Second,
			},
		}),
	}

	// List of initializer functions
	initializers := []func(opts []grpc.DialOption) error{
		s.InitUserClient,
		// Add more here as you add more clients
	}

	for _, initFn := range initializers {
		if err := initFn(dialOpts); err != nil {
			fmt.Println(err)
		}
	}
}

// Connects to client and adds connection to s.Clients
func (s *Server) InitUserClient(opts []grpc.DialOption) (err error) {
	conn, err := grpc.NewClient(ports.Users, opts...)
	if err != nil {
		err = fmt.Errorf("failed to dial user service: %v", err)
	}
	s.Clients.UserClient = userpb.NewUserServiceClient(conn)
	return err
}

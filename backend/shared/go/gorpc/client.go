package gorpc

import (
	"fmt"
	ct "social-network/shared/go/customtypes"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

// GetGRpcClient creates a gRPC client of type T. Pass it the grpc generated constructor that creates the client service, the full address of the service, and the context keys to propagate.
func GetGRpcClient[T any](constructor func(grpc.ClientConnInterface) T, fullAddress string, contextKeys []ct.CtxKey) (T, error) {
	customUnaryInterceptor, err := UnaryClientInterceptorWithContextKeys(contextKeys...)
	if err != nil {
		return *new(T), err
	}
	customStreamInterceptor, err := StreamClientInterceptorWithContextKeys(contextKeys...)
	if err != nil {
		return *new(T), err
	}
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]	}`),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 2 * time.Second,
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.2,
				Jitter:     0.5,
				MaxDelay:   5 * time.Second,
			},
		}),
		grpc.WithUnaryInterceptor(customUnaryInterceptor),
		grpc.WithStreamInterceptor(customStreamInterceptor),
	}

	conn, err := grpc.NewClient(fullAddress, dialOpts...)
	if err != nil {
		return *new(T), fmt.Errorf("failed to dial user service: %v", err)
	}

	return constructor(conn), nil
}

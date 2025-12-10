package gorpc

import (
	ct "social-network/shared/go/customtypes"

	"google.golang.org/grpc"
)

// RunGRPCServer starts the gRPC server and blocks
func RunGRPCServer[T any](handler T, contextKeys []ct.CtxKey) (*grpc.Server, error) {

	customUnaryInterceptor, err := UnaryServerInterceptorWithContextKeys(contextKeys...)
	if err != nil {
		return nil, err
	}
	customStreamInterceptor, err := StreamServerInterceptorWithContextKeys(contextKeys...)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(customUnaryInterceptor),
		grpc.StreamInterceptor(customStreamInterceptor),
	)

	return grpcServer, nil
}

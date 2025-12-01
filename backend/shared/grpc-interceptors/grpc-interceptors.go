package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// WARNING THIS PACKAGE IS A WORK IN PROGRESS
// WARNING THIS PACKAGE IS A WORK IN PROGRESS
// WARNING THIS PACKAGE IS A WORK IN PROGRESS
// WARNING THIS PACKAGE IS A WORK IN PROGRESS
// WARNING THIS PACKAGE IS A WORK IN PROGRESS
// WARNING THIS PACKAGE IS A WORK IN PROGRESS
// WARNING THIS PACKAGE IS A WORK IN PROGRESS
// WARNING THIS PACKAGE IS A WORK IN PROGRESS
// WARNING THIS PACKAGE IS A WORK IN PROGRESS
// WARNING THIS PACKAGE IS A WORK IN PROGRESS

// this package holds gRPC interceptors for both client and server sides
//  they can add specified values from context to metadata and vice versa

/*
=====  SERVER ==================
*/

type wrappedStream struct {
	grpc.ServerStream
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func (w *wrappedStream) RecvMsg(m any) error {
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m any) error {
	return w.ServerStream.SendMsg(m)
}

// UnaryServerInterceptorWithContextKeys returns a server interceptor that extracts specified keys from metadata and adds them to the context.
func UnaryServerInterceptorWithContextKeys(keys []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			fmt.Println(md)
		}

		m, err := handler(ctx, req)
		return m, err
	}
}

// StreamServerInterceptorWithContextKeys returns a server interceptor that extracts specified keys from metadata and adds them to the context.
func StreamServerInterceptorWithContextKeys(keys []string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		handler(srv, ss)
		return nil
	}
}

/*
=====  CLIENT ==================
*/

// UnaryClientInterceptorWithContextKeys returns a client interceptor that adds specified context values to outgoing metadata.
func UnaryClientInterceptorWithContextKeys(keys []string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		metadata.AppendToOutgoingContext(ctx)
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

// StreamClientInterceptorWithContextKeys returns a client interceptor that adds specified context values to outgoing metadata.
func StreamClientInterceptorWithContextKeys(keys []string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		return clientStream, err
	}
}

package interceptor

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// this package holds gRPC interceptors for both client and server sides
// they can add specified values from context to metadata and vice versa, so that we can seamlessly propagate context values from one service to another
// the propagation only happens from client to server, cause context can only go one direction.

/*
=====  SERVER ==================
	ATM SERVER INTERCEPTOR ISNT NEEDED
*/

// IMPORTANT: Only "a-z", "0-9", and "-_." characters allowed for keys
func UnaryServerInterceptorWithContextKeys(keys ...string) grpc.UnaryServerInterceptor {
	if !validateContextKeys(keys...) {
		panic("bad context keys passed to interceptor creator, keys don't follow the validation requirements")
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		fmt.Println("metadata:", md)

		for _, key := range keys {
			vals := md.Get(key)
			for _, val := range vals {
				ctx = context.WithValue(ctx, key, val)
			}
		}

		m, err := handler(ctx, req)
		return m, err
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
}

func newWrappedServerStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedServerStream{s}
}

func (w *wrappedServerStream) RecvMsg(m any) error {
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedServerStream) SendMsg(m any) error {
	return w.ServerStream.SendMsg(m)
}

func StreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, newWrappedServerStream(ss))
	return err
}

/*
=====  CLIENT ==================
*/

// UnaryClientInterceptorWithContextKeys returns a client interceptor that adds specified context values to outgoing metadata.
//
// IMPORTANT: Only "a-z", "0-9", and "-_." characters allowed for keys
func UnaryClientInterceptorWithContextKeys(keys ...string) grpc.UnaryClientInterceptor {
	if !validateContextKeys(keys...) {
		panic("bad context keys passed to interceptor creator, keys don't follow the validation requirements")
	}
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.NewOutgoingContext(ctx, context2Metadata(ctx, keys...))
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

type wrappedClientStream struct {
	grpc.ClientStream
}

func newWrappedClientStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedClientStream{s}
}

func (w *wrappedClientStream) RecvMsg(m any) error {
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedClientStream) SendMsg(m any) error {
	return w.ClientStream.SendMsg(m)
}

// StreamClientInterceptorWithContextKeys returns a client interceptor that adds specified context values to outgoing metadata.
//
// IMPORTANT: Only "a-z", "0-9", and "-_." characters allowed for keys
func StreamClientInterceptorWithContextKeys(keys ...string) grpc.StreamClientInterceptor {
	if !validateContextKeys(keys...) {
		panic("bad context keys passed to interceptor creator, keys don't follow the validation requirements")
	}
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.NewOutgoingContext(ctx, context2Metadata(ctx, keys...))
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		return newWrappedClientStream(clientStream), err
	}
}

/*
===== OTHER ==================
*/

// context2Metadata is needed cause when making a grpc call, the context values don't automatically propagate, you need to manually create metadata and add them to the call
func context2Metadata(ctx context.Context, keys ...string) metadata.MD {
	md := metadata.New(make(map[string]string))
	for _, key := range keys {
		val, ok := ctx.Value(key).(string)
		if !ok {
			continue
		}
		md.Set(key, val)
	}

	fmt.Printf("[DEBUG] metadata from ctx: %v\n", md)
	return md
}

/* FROM GRPC DOCUMENTATION, relating to metadata keys

Only the following ASCII characters are allowed in keys:
- digits: 0-9
- uppercase letters: A-Z (normalized to lower)
- lowercase letters: a-z
- special characters: -_.

Uppercase letters are automatically converted to lowercase.

Keys beginning with "grpc-" are reserved for grpc-internal use only and may
result in errors if set in metadata.

*/

// validateContextKeys validates the keys so that nothing bad happens during context value propagation due to the above limitations
func validateContextKeys(keys ...string) bool {
	for _, key := range keys {
		for _, r := range []rune(key) {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' && r != '_' && r != '.' {
				return false
			}
		}

		if strings.HasPrefix(key, "grpc-") {
			return false
		}
	}

	return true
}

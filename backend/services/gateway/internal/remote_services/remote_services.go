package remoteservices

import (
	"fmt"
	"log"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"
	interceptor "social-network/shared/go/grpc-interceptors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRpcServices struct {
	Users       users.UserServiceClient
	Chat        chat.ChatServiceClient
	contextKeys []ct.CtxKey
}

// NewService creates a services object, it contains grpc clients and you use it to talk to grpc servers
//
// CtxKeys are the context keys that will be propagated to other services through the context object
func NewServices(contextKeys []ct.CtxKey) *GRpcServices {
	return &GRpcServices{
		contextKeys: contextKeys,
	}
}

func (g *GRpcServices) StartConnections() (func(), error) {

	customUnaryInterceptor, err := interceptor.UnaryClientInterceptorWithContextKeys(g.contextKeys...)
	if err != nil {
		return nil, err
	}
	customStreamInterceptor, err := interceptor.StreamClientInterceptorWithContextKeys(g.contextKeys...)
	if err != nil {
		return nil, err
	}

	usersConn, err := grpc.NewClient(
		"users:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		grpc.WithUnaryInterceptor(customUnaryInterceptor),
		grpc.WithStreamInterceptor(customStreamInterceptor),
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("[DEBUG]", usersConn.CanonicalTarget())
	fmt.Println("[DEBUG]", usersConn.GetState())
	fmt.Println("[DEBUG]", usersConn.Target())

	g.Users = users.NewUserServiceClient(usersConn)

	chatConn, err := grpc.NewClient("chat:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}
	g.Chat = chat.NewChatServiceClient(chatConn)

	deferMe := func() {
		usersConn.Close()
		chatConn.Close()
	}
	return deferMe, nil
}

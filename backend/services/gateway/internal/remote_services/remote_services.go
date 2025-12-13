package remoteservices

import (
	"log"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/users"
	contextkeys "social-network/shared/go/context-keys"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/gorpc"

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

	customUnaryInterceptor, err := gorpc.UnaryClientInterceptorWithContextKeys(contextkeys.CommonKeys()...)
	if err != nil {
		return nil, err
	}
	customStreamInterceptor, err := gorpc.StreamClientInterceptorWithContextKeys(contextkeys.CommonKeys()...)
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
		return nil, err
	}
	g.Users = users.NewUserServiceClient(usersConn)

	chatConn, err := grpc.NewClient("chat:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	g.Chat = chat.NewChatServiceClient(chatConn)

	deferMe := func() {
		fails := 0
		err := usersConn.Close()
		if err != nil {
			fails++
			log.Println("usersConn failed to close, ERROR:", err.Error())
		}
		err = chatConn.Close()
		if err != nil {
			fails++
			log.Println("chatConn failed to close, ERROR:", err.Error())
		}
		if fails == 0 {
			log.Println("Grpc connections closed correctly")
		}
	}
	return deferMe, nil
}

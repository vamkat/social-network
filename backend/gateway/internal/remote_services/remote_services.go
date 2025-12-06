package remoteservices

import (
	"fmt"
	"log"
	"social-network/gateway/internal/utils"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/users"
	interceptor "social-network/shared/go/grpc-interceptors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRpcServices struct {
	Users users.UserServiceClient
	Chat  chat.ChatServiceClient
}

func NewServices() GRpcServices {
	return GRpcServices{}
}

func (g *GRpcServices) StartConnections() (func(), error) {
	usersConn, err := grpc.NewClient(
		"users:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// Add context keys to be propagated via interceptors
		grpc.WithStreamInterceptor(interceptor.StreamClientInterceptorWithContextKeys(utils.UserId, utils.ReqUUID, utils.TraceId)),
		grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptorWithContextKeys(utils.UserId, utils.ReqUUID, utils.TraceId)),
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

package remoteservices

import (
	"fmt"
	"log"
	"social-network/shared/gen-go/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRpcServices struct {
	Users users.UserServiceClient
}

func NewServices() GRpcServices {
	return GRpcServices{}
}

func (g *GRpcServices) StartConnections() (func(), error) {
	usersConn, err := grpc.NewClient("users:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(usersConn.CanonicalTarget())
	fmt.Println(usersConn.GetState())
	fmt.Println(usersConn.Target())
	g.Users = users.NewUserServiceClient(usersConn)

	deferMe := func() {
		usersConn.Close()
	}
	return deferMe, nil
}

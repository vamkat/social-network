package gateway_test

import (
	"context"
	"fmt"
	"math/rand/v2"

	"social-network/services/testing/internal/configs"
	"social-network/services/testing/internal/utils"
	"social-network/shared/gen-go/users"
	contextkeys "social-network/shared/go/context-keys"
	"social-network/shared/go/gorpc"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var UsersService users.UserServiceClient

func StartTest(ctx context.Context, cfgs configs.Configs) {
	var err error
	UsersService, err = gorpc.GetGRpcClient(
		users.NewUserServiceClient,
		cfgs.UsersGRPCAddr,
		contextkeys.CommonKeys(),
	)
	if err != nil {
		panic("failed to connect to users service: %v" + err.Error())
	}

	randomRegister(ctx)

}

var fail = "FAIL TEST -> register, err:"

func randomRegister(ctx context.Context) {
	fmt.Println("starting register test")
	for range 50 {
		req := users.RegisterUserRequest{
			Username:    utils.RandomString(10),
			FirstName:   utils.RandomString(10),
			LastName:    utils.RandomString(10),
			DateOfBirth: timestamppb.New(time.Unix(rand.Int64N(1000000), 0)),
			Avatar:      0,
			About:       utils.RandomString(10),
			Public:      false,
			Email:       utils.RandomString(10),
			Password:    utils.RandomString(10),
		}
		resp, err := UsersService.RegisterUser(ctx, &req)
		if err != nil {
			panic(fail + err.Error())
		}

		if resp.UserId < 1 {
			panic(fail)
		}
	}

	fmt.Println("register test passed")
}

package users_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"

	"social-network/services/testing/internal/configs"
	"social-network/services/testing/internal/utils"
	"social-network/shared/gen-go/users"
	"social-network/shared/go/ct"
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
		ct.CommonKeys(),
	)
	if err != nil {
		panic("failed to connect to users service: %v" + err.Error())
	}

	var wg sync.WaitGroup
	wg.Go(func() { randomRegister(ctx) })
	wg.Go(func() { randomLogin(ctx) })
	wg.Go(func() { registerLogin(ctx) })
	wg.Wait()

}

var fail = "FAIL TEST: err ->"

func randomRegister(ctx context.Context) {
	fmt.Println("starting register test")
	for range 100 {
		req := newRegisterReq()
		resp, err := UsersService.RegisterUser(ctx, req)
		if err != nil {
			panic(fail + err.Error())
		}

		if resp.UserId < 1 {
			panic(fail)
		}

	}

	fmt.Println("random register test passed")
}

func randomLogin(ctx context.Context) {
	fmt.Println("starting Login test")
	for range 100 {
		req := newLoginReq()
		_, err := UsersService.LoginUser(ctx, req)
		if err != nil && !strings.Contains(err.Error(), "invalid identifier or password") {
			panic(fail + "wrong type of error!")
		}
		if err == nil {
			panic(fail + "expected error! cause these random logins should all be failing!")
		}
	}

	fmt.Println("random login test passed")
}

func registerLogin(ctx context.Context) {
	reg := newRegisterReq()

	_, err := UsersService.RegisterUser(ctx, reg)
	if err != nil {
		panic(fail + err.Error())
	}

	log := newLoginReq()
	log.Identifier = reg.Email
	log.Password = reg.Password

	resp, err := UsersService.LoginUser(ctx, log)
	if err != nil {
		panic(fail + "should have worked! err should be nil:" + err.Error())
	}
	if resp.Username != reg.Username {
		panic(fail + "incorrect login, these two should be the same: `" + resp.Username + "` <-> `" + reg.Username + "`")
	}

	if resp.UserId == 0 || resp.Username == "" {
		panic("found empty values")
	}
	fmt.Println("passed simple reg login test")
}

func newRegisterReq() *users.RegisterUserRequest {
	req := users.RegisterUserRequest{
		Username:    strings.Title(utils.RandomString(10, false)),
		FirstName:   strings.Title(utils.RandomString(10, false)),
		LastName:    strings.Title(utils.RandomString(10, false)),
		DateOfBirth: timestamppb.New(time.Unix(rand.Int64N(1000000), 0)),
		Avatar:      0,
		About:       utils.RandomString(300, true),
		Public:      false,
		Email:       utils.RandomString(10, false) + "@hotmail.com",
		Password:    utils.RandomString(10, true),
	}
	return &req
}

func newLoginReq() *users.LoginRequest {
	req := users.LoginRequest{
		Identifier: strings.Title(utils.RandomString(10, false)),
		Password:   strings.Title(utils.RandomString(10, false)),
	}
	return &req
}

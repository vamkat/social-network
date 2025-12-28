package entry

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"social-network/services/testing/internal/configs"
	gateway_test "social-network/services/testing/internal/gateway_testing"
	users_test "social-network/services/testing/internal/users_testing"
	tele "social-network/shared/go/telemetry"
	"sync"
	"syscall"
)

var cfgs configs.Configs

func Run() {
	fmt.Println("start run")
	cfgs = configs.GetConfigs()

	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	wg := sync.WaitGroup{}

	wg.Go(func() {
		if err := users_test.StartTest(ctx, cfgs); err != nil {
			tele.Fatal("ERROR WTF" + err.Error())
		}
	})

	wg.Go(func() {
		if err := gateway_test.StartTest(ctx, cfgs); err != nil {
			tele.Fatal("ERROR WTF" + err.Error())
		}
	})

	wg.Wait()

	stopSignal()
	fmt.Println("end run")
}

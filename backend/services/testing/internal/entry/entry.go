package entry

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"social-network/services/testing/internal/configs"
	gateway_test "social-network/services/testing/internal/gateway_testing"
	"syscall"
)

var cfgs configs.Configs

func Run() {
	fmt.Println("start run")
	cfgs = configs.GetConfigs()
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	gateway_test.StartTest(ctx, cfgs)

	stopSignal()
	fmt.Println("end run")
}

package entry

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"social-network/services/testing/internal/configs"
	users_test "social-network/services/testing/internal/users_testing"
	"syscall"
)

var cfgs configs.Configs

func Run() {
	fmt.Println("start run")
	cfgs = configs.GetConfigs()
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	users_test.StartTest(ctx, cfgs)

	stopSignal()
	fmt.Println("end run")
}

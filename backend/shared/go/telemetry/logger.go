package tele

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log/global"
)

var ErrUnevenArgs = errors.New("passed arguments aren't even")

type LogLevel struct {
	tag   string
	level int
}

// TODO
//get stack info
//extra context info

type logging struct {
	serviceName string
	enableDebug bool     //if debug prints will be shown or not
	contextKeys []string //the context keys that will be added to logs as metadata
	slog        *slog.Logger
	simplePrint bool //if it should print logs in a simple way, or a super verbose way with all details
}

// newLogger returns a logger that actually logs, uses a handler that taken from a global provider created by the otel sdk
func newLogger(serviceName string, contextKeys contextKeys, enableDebug bool, simplePrint bool) *logging {
	handler := otelslog.NewHandler(
		serviceName,
		otelslog.WithLoggerProvider(global.GetLoggerProvider()),
		otelslog.WithSource(true),
	)

	fmt.Println("handle record")
	for range 4 {
		err := handler.Handle(context.Background(), slog.NewRecord(time.Unix(1, 1), slog.LevelInfo, "HANLDE MESSAGE", uintptr(1)))
		fmt.Println(err)
		time.Sleep(time.Second * 1)
	}

	logger := slog.New(handler)
	fmt.Println("logger.info")
	for range 4 {
		logger.Info("logger.Info(msg) A TEEEEEEEEEEST")
		time.Sleep(time.Second * 1)
	}

	return &logging{
		serviceName: serviceName,
		contextKeys: contextKeys.GetKeys(),
		slog:        logger,
		enableDebug: enableDebug,
		simplePrint: simplePrint,
	}
}

func newBasicLogger() *logging {
	return &logging{
		serviceName: "not-initalized",
		slog:        slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func (l *logging) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if level == slog.LevelDebug && l.enableDebug == false {
		return
	}

	ctxArgs := l.context2Args(ctx)
	for _, ctxArg := range ctxArgs {
		args = append(args, ctxArg)
	}

	if !l.simplePrint {
		args = []any{}
	}

	//maybe not use context
	l.slog.Log(ctx, level, msg, args...)
}

func (l *logging) context2Args(ctx context.Context) []string {
	args := []string{}
	for _, key := range l.contextKeys {
		val, ok := ctx.Value(key).(string)
		if !ok {
			continue
		}
		args = append(args, key)
		args = append(args, val)
	}
	return args
}

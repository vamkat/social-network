package tele

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
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

func NewLogger(serviceName string, contextKeys contextKeys, enableDebug bool, simplePrint bool) logging {

	logger := otelslog.NewLogger(serviceName, otelslog.WithSource(true))
	slog.SetDefault(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	return logging{
		serviceName: serviceName,
		contextKeys: contextKeys.GetKeys(),
		slog:        logger,
		enableDebug: enableDebug,
		simplePrint: simplePrint,
	}
}

func (l *logging) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if level == slog.LevelDebug && l.enableDebug == false {
		return
	}

	if l.simplePrint {
		fmt.Printf("%s: %s\n", level.String(), msg)
		return
	}

	ctxArgs := l.context2Args(ctx)
	for _, ctxArg := range ctxArgs {
		args = append(args, ctxArg)
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

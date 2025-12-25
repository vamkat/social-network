package tele

import (
	"context"
	"errors"
	"log/slog"
	"os"
)

var ErrUnevenArgs = errors.New("passed arguments aren't even")

type LogLevel struct {
	tag   string
	level int
}

var (
	DEBUG  = LogLevel{"DEBUG", 0}
	INFO   = LogLevel{"INFO", 1}
	WARN   = LogLevel{"WARN", 2}
	SEVERE = LogLevel{"SEVERE", 3}
	FATAL  = LogLevel{"FATAL", 4}
)

//get stack info
//extra context info

type logging struct {
	serviceName string
	enableDebug bool
	contextKeys []string
	slog        *slog.Logger
}

func NewLogger(serviceName string, contextKeys contextKeys, enableDebug bool) logging {
	slogHandler := &slog.HandlerOptions{
		AddSource: true,
	}
	slog := slog.New(slog.NewJSONHandler(os.Stderr, slogHandler))
	return logging{
		serviceName: serviceName,
		contextKeys: contextKeys.GetKeys(),
		slog:        slog,
		enableDebug: enableDebug,
	}
}

func (l *logging) log(ctx context.Context, level slog.Level, msg string, args ...any) error {
	if len(args)%2 != 0 {
		return ErrUnevenArgs
	}
	ctxArgs := l.context2Args(ctx)
	for _, ctxArg := range ctxArgs {
		args = append(args, ctxArg)
	}
	//maybe not use context
	l.slog.Log(ctx, level, msg, args...)
	return nil
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

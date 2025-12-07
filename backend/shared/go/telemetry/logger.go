package tele

import (
	"context"
	"errors"
	"log/slog"
	"os"
)

var ErrUnevenArgs = errors.New("passed arguments aren't even")

var (
	DEBUG  = "DEBUG"
	INFO   = "INFO"
	WARN   = "WARN"
	SEVERE = "SEVERE"
	FATAL  = "FATAL"
)

//get stack info
//extra context info

type logger struct {
	serviceName string
	outputLevel int
	contextKeys []string
	slog        *slog.Logger
}

func NewLogger(serviceName string, contextKeys []string) logger {
	slogHandler := &slog.HandlerOptions{
		AddSource: true,
	}
	slog := slog.New(slog.NewJSONHandler(os.Stderr, slogHandler))
	return logger{
		serviceName: serviceName,
		contextKeys: contextKeys,
		slog:        slog,
	}
}

func (l *logger) Info(ctx context.Context, args ...string) {
	l.log(slog.LevelInfo, ctx, args...)
}

func (l *logger) log(level slog.Level, ctx context.Context, msg string, args ...string) error {
	if len(args)%2 != 0 {
		return ErrUnevenArgs
	}
	ctxArgs := l.context2Args(ctx)
	args = append(args, ctxArgs...)
	//maybe not use context
	l.slog.Log(nil, msg, level, args)
}

func (l *logger) context2Args(ctx context.Context) []string {
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

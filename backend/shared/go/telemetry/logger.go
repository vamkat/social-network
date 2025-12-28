package tele

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
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
//log the 3 functions that called the log

type logging struct {
	serviceName string
	enableDebug bool     //if debug prints will be shown or not
	contextKeys []string //the context keys that will be added to logs as metadata
	slog        *slog.Logger
	simplePrint bool   //if it should print logs in a simple way, or a super verbose way with all details
	prefix      string //this will be added at the start of logs that appear in the local terminal only, suggestion: keep it 3 letters CAPITAL, ex: API, MED, SOC, NOT, POS, RED, USE, CHA
	hasPrefix   bool
}

// newLogger returns a logger that actually logs, uses a handler that taken from a global provider created by the otel sdk
func newLogger(serviceName string, contextKeys contextKeys, enableDebug bool, simplePrint bool, prefix string) *logging {
	handler := otelslog.NewHandler(
		serviceName,
		otelslog.WithLoggerProvider(global.GetLoggerProvider()),
		otelslog.WithSource(true),
	)

	logger := slog.New(handler)

	return &logging{
		serviceName: serviceName,
		contextKeys: contextKeys.GetKeys(),
		slog:        logger,
		enableDebug: enableDebug,
		simplePrint: simplePrint,
		prefix:      prefix,
		hasPrefix:   "" != prefix,
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

	ctxArgs := []any{}
	if ctx == nil {
		ctx = context.Background()
	} else {
		for _, ctxArg := range l.context2Args(ctx) {
			ctxArgs = append(ctxArgs, ctxArg)
		}
	}

	l.slog.Log(ctx, level, msg, append(args, ctxArgs)...)

	if !l.simplePrint {
		args = append(args, ctxArgs)
	}

	time := fmt.Sprint(time.Now().Format("15:04:05.000"))
	if len(args) == 0 {
		if l.hasPrefix {
			fmt.Printf("%s [%s]: %s - %s\n", time, l.prefix, level.String(), msg)
		} else {
			fmt.Printf("%s [%s]: %s - %s - args: %s\n", time, l.serviceName, level.String(), msg, fmt.Sprint(args...))
		}
		return
	}

	var prefix string
	if l.hasPrefix {
		prefix = l.prefix
	} else {
		prefix = l.serviceName
	}

	var argsPart string
	if len(args) > 0 {
		argsPart = fmt.Sprintf(" - args: %s", formatArgs(args...))
	}

	fmt.Printf("%s [%s]: %s - %s%s\n", time, prefix, level.String(), msg, argsPart)
}

func formatArgs(args ...any) any {
	parts := make([]any, 0, len(args))

	for _, arg := range args {
		v := reflect.ValueOf(arg)

		// Handle pointers
		if v.Kind() == reflect.Pointer && !v.IsNil() {
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct {
			parts = append(parts, fmt.Sprintf("%#v", arg))
		} else {
			parts = append(parts, fmt.Sprint(arg))
		}
	}

	return fmt.Sprint(parts...)
}

// TODO think what to do here, not only context keys
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

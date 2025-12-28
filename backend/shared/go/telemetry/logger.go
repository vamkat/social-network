package tele

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
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

	callerInfo := functionCallers()

	ctxArgs := []slog.Attr{}
	if ctx == nil {
		ctx = context.Background()
	} else {
		ctxArgs = l.context2Attributes(ctx)
	}

	var prefix string
	if l.hasPrefix {
		prefix = l.prefix
	} else {
		prefix = l.serviceName
	}

	l.slog.Log(
		ctx,
		level,
		msg,
		slog.GroupAttrs("customArgs", kvPairsToAttrs(args)...),
		slog.GroupAttrs("context", ctxArgs...),
		slog.String("callers", callerInfo),
		slog.String("prefix", prefix),
	)

	if !l.simplePrint {
		args = append(args, ctxArgs)
	}

	var argsPart string
	if len(args) > 0 {
		argsPart = fmt.Sprintf(" - args: %s", formatArgs(args...))
	}

	time := fmt.Sprint(time.Now().Format("15:04:05.000"))
	fmt.Printf("%s [%s]: %s - %s%s\n", time, prefix, level.String(), msg, argsPart)
}

func kvPairsToAttrs(pairs []any) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			key = "invalid_key"
		}
		attrs = append(attrs, slog.Any(key, pairs[i+1]))
	}
	return attrs
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

func (l *logging) context2Attributes(ctx context.Context) []slog.Attr {
	args := []slog.Attr{}
	for _, key := range l.contextKeys {
		val, ok := ctx.Value(key).(string)
		if !ok {
			continue
		}
		args = append(args, slog.Any(key, val))
	}
	return args
}

func functionCallers() string {
	var callers = []string{}
	pc := make([]uintptr, 3)
	n := runtime.Callers(4, pc)
	if n == 0 {
		return "(no caller data)"
	}
	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.Function, "runtime_test") {
			break
		}
		start := strings.LastIndex(frame.Func.Name(), "/")
		callers = append(callers, fmt.Sprintf("by %s at %d ", frame.Func.Name()[start+1:], frame.Line))
		if !more {
			break
		}
	}

	return strings.Join(callers, "\n")
}

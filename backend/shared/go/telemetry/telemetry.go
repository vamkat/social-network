package tele

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
)

var (
	telemeter *telemetry
	tracer    *tracing
)

func init() {
	// we create an empty no op telemetry so that we don't panic if this packages logging is called without initalizing it properly
	telemeter = &telemetry{
		logger: newBasicLogger(),
	}

	telemeter.logger.slog.Log(context.Background(), slog.LevelDebug, "this is a test message")
}

// this interface exists to enforce usage of typed context keys instead of just strings
type contextKeys interface {
	GetKeys() []string
}

type telemetry struct {
	logger      *logging
	tracer      *tracing
	meterer     *metering
	serviceName string
	enableDebug bool
}

// Will only show up in dev environment
func Debug(ctx context.Context, message string, args ...any) {
	telemeter.logger.log(ctx, slog.LevelDebug, message, args...)
}

// General info
func Info(ctx context.Context, message string, args ...any) {
	telemeter.logger.log(ctx, slog.LevelInfo, message, args...)
}

// Something that isn't really breaking something, but if it happens a lot that could mean something bad is going on and should be looked into
func Warn(ctx context.Context, message string, args ...any) {
	telemeter.logger.log(ctx, slog.LevelWarn, message, args...)
}

// Something severe has happened that shouldn't have, it needs to be looked at immediately and addressed!
func Error(ctx context.Context, message string, args ...any) {
	telemeter.logger.log(ctx, slog.LevelError, message, args...)
}

func Fatal(message string) {
	telemeter.logger.log(context.Background(), slog.LevelError, message)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	telemeter.logger.log(context.Background(), slog.LevelError, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func TraceStart(ctx context.Context, message string, args any) (context.Context, trace.Span) {
	return telemeter.tracer.tracer.Start(ctx, message,
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(),
	)
}

//TODO handle cancillation from ctx?

// actually activates the functionality of open telemetry
func InitTelemetry(ctx context.Context, serviceName string, servicePrefix string, collectorAddress string, contextKeys contextKeys, enableDebug bool, simplePrint bool) (func(), error) {

	otelShutdown, err := SetupOTelSDK(ctx, collectorAddress, serviceName)
	if err != nil {
		Fatalf("open telemetry sdk failed, ERROR: %s", err.Error())
	}
	Info(ctx, "open telemetry sdk initialized with args: @1 @2 @3 @4 @5", "name", serviceName, "prefix", servicePrefix, "address", collectorAddress, "keys", contextKeys, "debug", enableDebug, "simple", simplePrint)

	logger := newLogger(serviceName, contextKeys, enableDebug, simplePrint, servicePrefix)
	slog.SetDefault(logger.slog)
	// rollCnt metric.Int64Counter

	tracer := NewTracer(serviceName)

	telemeter = &telemetry{
		logger:      logger,
		tracer:      tracer,
		meterer:     nil,
		serviceName: serviceName,
		enableDebug: enableDebug,
	}

	close := func() {
		err := otelShutdown(context.Background())
		if err != nil {
			Info(ctx, "otel shutdown ungracefully! ERROR: "+err.Error())
		} else {
			Info(ctx, "otel shutdown gracefully")
		}
	}

	return close, nil
}

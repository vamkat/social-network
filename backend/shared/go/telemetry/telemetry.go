package tele

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/exporters/prometheus"
)

var (
	telemeter *telemetry
)

func init() {
	telemeter = &telemetry{}
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

// actually activates the functionality of
func InitTelemetry(ctx context.Context, serviceName string, contextKeys contextKeys, enableDebug bool, simplePrint bool) func() {

	logger := NewLogger(serviceName, contextKeys, enableDebug, simplePrint)

	// rollCnt metric.Int64Counter

	telemeter = &telemetry{
		logger:      &logger,
		tracer:      nil,
		meterer:     nil,
		serviceName: serviceName,
		enableDebug: enableDebug,
	}

	close := initOpenTelemetrySDK(ctx)
	return close
}

func initOpenTelemetrySDK(ctx context.Context) func() {
	otelShutdown, err := SetupOTelSDK(ctx)
	if err != nil {
		Fatalf("open telemetry sdk failed, ERROR: %s", err.Error())
	}
	Info(ctx, "open telemetry ready")

	return func() {
		err := otelShutdown(context.Background())
		if err != nil {
			Info(ctx, "otel shutdown ungracefully! ERROR: "+err.Error())
		} else {
			Info(ctx, "otel shutdown gracefully")
		}
	}
}

func newPrometheus() *prometheus.Exporter {
	x, _ := prometheus.New(nil)
	return x
}

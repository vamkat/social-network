package tele

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"go.opentelemetry.io/otel/exporters/prometheus"
)

var (
	telemeter *telemetery
)

type telemetery struct {
	logger  *logging
	tracer  *tracing
	meterer *metering
}

// will only show up in dev environment
func Debug(ctx context.Context, message string, args ...any) error {
	return telemeter.logger.log(ctx, slog.LevelDebug, message, args...)
}

// general info
func Info(ctx context.Context, message string, args ...any) error {
	return telemeter.logger.log(ctx, slog.LevelInfo, message, args...)
}

// something that isn't really breaking something, but if it happens a lot that could mean something bad is going on and should be looked into
func Warn(ctx context.Context, message string, args ...any) error {
	return telemeter.logger.log(ctx, slog.LevelWarn, message, args...)
}

// inits telemeter to be no-op
func init() {

}

// this interface exists to enforce usage of typed context keys instead of just strings
type contextKeys interface {
	GetKeys() []string
}

// actually activates the functionality of
func InitTelemetry(ctx context.Context, serviceName string, contextKeys contextKeys, enableDebug bool) (telemetery, func()) {
	close := initOpenTelemetrySDK(ctx)
	logger := NewLogger(serviceName, contextKeys, enableDebug)

	telemeter := telemetery{
		logger: &logger,
	}

	return telemeter, close
}

func initOpenTelemetrySDK(ctx context.Context) func() {
	otelShutdown, err := SetupOTelSDK(ctx)
	if err != nil {
		log.Fatal("open telemetry sdk failed, ERROR:", err.Error())
	}
	fmt.Println("open telemetry ready")

	return func() {
		err := otelShutdown(context.Background())
		if err != nil {
			log.Println("otel shutdown ungracefully! ERROR: " + err.Error())
		} else {
			log.Println("otel shutdown gracefully")
		}
	}
}

func newPrometheus() *prometheus.Exporter {
	x, _ := prometheus.New(nil)
	return x
}

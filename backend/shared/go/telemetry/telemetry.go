package tele

import "go.opentelemetry.io/otel/exporters/prometheus"

func NewPrometheus() *prometheus.Exporter {
	x, _ := prometheus.New(nil)
	return x
}

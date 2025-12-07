package tele

import (
	"log/slog"
	"os"
)

const (
	DEBUG  = 0
	INFO   = 1
	WARN   = 2
	SEVERE = 3
	FATAL  = 4
)

//get stack info
//extra context info

type Logger struct {
	serviceName string
	outputLevel int
	contextKeys []string
	slog        any
}

func NewLogger(serviceName string, contextKeys []string) Logger {
	return Logger{
		serviceName: serviceName,
		contextKeys: contextKeys,
		slog:        slog.New(slog.NewTextHandler(os.Stderr)),
	}
}

func process() {

}

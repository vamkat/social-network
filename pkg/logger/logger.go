package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	once   sync.Once
	Logger *zap.Logger
)

func Init() {
	once.Do(func() {
		l, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		Logger = l
	})
}

func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

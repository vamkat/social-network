package application

import (
	"context"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/users"
	"time"
)

type RedisClient interface {
	IncrEx(ctx context.Context, key string, expSeconds int64) (int, error)
	SetStr(ctx context.Context, key string, value string, exp time.Duration) error
	GetStr(ctx context.Context, key string) (any, error)
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
	GetObj(ctx context.Context, key string, dest any) error
	Del(ctx context.Context, key string) error
	TestRedisConnection() error
}

type GatewayApp struct {
	//telemetry
	Redis RedisClient
	Users users.UserServiceClient
	Chat  chat.ChatServiceClient
}

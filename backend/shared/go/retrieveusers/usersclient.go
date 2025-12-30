package retrieveusers

import (
	"context"
	cm "social-network/shared/gen-go/common"
	"time"
)

type GetBatchBasicUserInfo func(ctx context.Context, req *cm.UserIds) (*cm.ListUsers, error)

// RedisCache defines the minimal Redis operations used by the hydrator.
type RedisCache interface {
	GetObj(ctx context.Context, key string, dest any) error
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
}

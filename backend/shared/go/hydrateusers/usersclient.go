package userhydrate

import (
	"context"
	cm "social-network/shared/gen-go/common"
	"time"
)

// UsersBatchClient is the subset the hydrator needs.
type UsersBatchClient interface {
	GetBatchBasicUserInfo(ctx context.Context, userIds []int64) (*cm.ListUsers, error)
}

// RedisCache defines the minimal Redis operations used by the hydrator.
type RedisCache interface {
	GetObj(ctx context.Context, key string, dest any) error
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
}

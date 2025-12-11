package userhydrate

import (
	"context"
	userpb "social-network/shared/gen-go/users"
	"time"
)

// UsersBatchClient is the subset the hydrator needs.
type UsersBatchClient interface {
	GetBatchBasicUserInfo(ctx context.Context, userIds []int64) (*userpb.ListUsers, error)
}

// RedisCache defines the minimal Redis operations used by the hydrator.
type RedisCache interface {
	GetObj(ctx context.Context, key string, dest any) error
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
}

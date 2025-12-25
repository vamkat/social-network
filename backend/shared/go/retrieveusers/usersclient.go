package retrieveusers

import (
	"context"
	cm "social-network/shared/gen-go/common"
	"social-network/shared/gen-go/media"
	ct "social-network/shared/go/ct"
	"time"
)

// UsersBatchClient is the subset the hydrator needs.
type UsersBatchClient interface {
	GetBatchBasicUserInfo(ctx context.Context, userIds ct.Ids) (*cm.ListUsers, error)
	GetImages(ctx context.Context, imageIds ct.Ids, variant media.FileVariant) (map[int64]string, []int64, error)
}

// RedisCache defines the minimal Redis operations used by the hydrator.
type RedisCache interface {
	GetObj(ctx context.Context, key string, dest any) error
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
}

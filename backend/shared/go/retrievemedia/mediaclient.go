package retrievemedia

import (
	"context"
	"social-network/shared/gen-go/media"
	"time"
)

// UsersBatchClient is the subset the hydrator needs.
type Client interface {
	GetImages(ctx context.Context, req *media.GetImagesRequest, variant *media.FileVariant) (*media.GetImagesResponse, error)
}

// RedisCache defines the minimal Redis operations used by the hydrator.
type RedisCache interface {
	GetObj(ctx context.Context, key string, dest any) error
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
}

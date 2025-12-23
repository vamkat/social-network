package application

import (
	"context"
	"net/url"
	md "social-network/services/media/internal/models"
	ct "social-network/shared/go/customtypes"
	"time"
)

type Clients interface {
	GenerateDownloadURL(
		ctx context.Context,
		bucket string,
		objectKey string,
		expiry time.Duration,
	) (*url.URL, error)

	GenerateUploadURL(
		ctx context.Context,
		bucket string,
		objectKey string,
		expiry time.Duration,
	) (*url.URL, error)

	ValidateUpload(
		ctx context.Context,
		upload md.FileMeta,
	) error

	GenerateVariant(
		ctx context.Context,
		srcBucket string,
		srcObjectKey string,
		trgBucket string,
		trgObjectKey string,
		variant ct.FileVariant,
	) (size int64, err error)

	DeleteFile(ctx context.Context,
		bucket string,
		objectKey string,
	) error
}

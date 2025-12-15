package application

import (
	"context"
	"net/url"
	md "social-network/shared/go/models"
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
		fm md.FileMeta,
	) error
}

package client

import (
	"context"
	"io"
	"social-network/services/media/internal/configs"

	"github.com/minio/minio-go/v7"
)

type Clients struct {
	Configs           configs.FileService
	MinIOClient       *minio.Client
	PublicMinIOClient *minio.Client
	Validator         Validator
}

type Validator interface {
	ValidateImage(ctx context.Context, r io.Reader) error
}

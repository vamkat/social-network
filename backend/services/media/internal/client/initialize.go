package client

import (
	"social-network/services/media/internal/configs"

	"github.com/minio/minio-go/v7"
)

type Clients struct {
	Configs     configs.FileService
	MinIOClient *minio.Client
}

package client

import (
	"github.com/minio/minio-go/v7"
)

type Clients struct {
	MinIOClient *minio.Client
}

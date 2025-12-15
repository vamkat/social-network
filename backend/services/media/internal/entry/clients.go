package entry

import (
	"context"
	"log"
	"reflect"
	"social-network/services/media/internal/configs"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOConn(cfgs configs.FileService) (*minio.Client, error) {
	var minioClient *minio.Client
	var err error

	endpoint := cfgs.Endpoint
	accessKey := cfgs.AccessKey
	secret := cfgs.Secret

	for range 10 {
		minioClient, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secret, ""),
			Secure: false,
		})
		if err == nil {
			break
		}
		log.Println("MinIO not ready, retrying in 2s...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	// Ensure bucket exists
	ctx := context.Background()
	if err := EnsureBuckets(ctx,
		minioClient, cfgs.Buckets); err != nil {
		return nil, err
	}

	return minioClient, nil
}

func EnsureBuckets(ctx context.Context, client *minio.Client, buckets configs.Buckets) error {
	v := reflect.ValueOf(buckets)

	for i := 0; i < v.NumField(); i++ {
		bucketName := v.Field(i).String()

		if bucketName == "" {
			continue
		}

		exists, err := client.BucketExists(ctx, bucketName)
		if err != nil {
			return err
		}

		if !exists {
			if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
				return err
			}
		}
	}

	return nil
}

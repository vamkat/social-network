package entry

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"social-network/services/media/internal/configs"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

func NewMinIOConn(cfgs configs.FileService, endpoint string, skipBucketCreation bool) (*minio.Client, error) {
	var minioClient *minio.Client
	var err error

	accessKey := cfgs.AccessKey
	secret := cfgs.Secret

	for range 10 {
		minioClient, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secret, ""),
			Secure: false,
			Region: "us-east-1",
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

	log.Println("Connected to minio client")

	if skipBucketCreation {
		return minioClient, nil
	}

	// Ensure bucket exists
	ctx := context.Background()
	if err := EnsureBuckets(ctx,
		minioClient, cfgs.Buckets); err != nil {
		return nil, err
	}

	log.Println("Setting up lifecycle rules")

	lcfg := lifecycle.NewConfiguration()

	rule := lifecycle.Rule{
		ID:     "delete-unvalidated",
		Status: "Enabled",
		RuleFilter: lifecycle.Filter{
			Tag: lifecycle.Tag{
				Key:   "validated",
				Value: "false",
			},
		},
		Expiration: lifecycle.Expiration{
			Days: lifecycle.ExpirationDays(1),
		},
	}

	lcfg.Rules = append(lcfg.Rules, rule)

	err = minioClient.SetBucketLifecycle(ctx, cfgs.Buckets.Originals, lcfg)
	if err != nil {
		log.Println("Error setting lifecycle:", err)
		// We might still continue
	}

	return minioClient, nil
}

func EnsureBuckets(ctx context.Context, client *minio.Client, buckets configs.Buckets) error {
	v := reflect.ValueOf(buckets)
	log.Println("Creating buckets")
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
			fmt.Printf("Creating bucket: %v\n", bucketName)
			if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
				return err
			}
		}
	}
	log.Println("Buckets created!")

	return nil
}

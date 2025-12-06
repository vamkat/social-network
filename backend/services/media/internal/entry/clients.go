package entry

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOConn() *minio.Client {
	var minioClient *minio.Client
	var err error

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secret := os.Getenv("MINIO_SECRET_KEY")

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
		log.Fatalln(err)
	}

	// Ensure bucket exists
	bucket := "images"
	ctx := context.Background()
	exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
	if errBucketExists != nil {
		log.Fatalln(errBucketExists)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Created bucket:", bucket)
	}
	return minioClient
}

package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
)

func (c *Clients) UploadToMinIO(
	ctx context.Context,
	fileContent []byte,
	filename string,
	bucket string,
	contentType string,
) (minio.UploadInfo, error) {

	tmp, err := os.CreateTemp("", "upload-*")
	if err != nil {
		return minio.UploadInfo{}, err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	reader := bytes.NewReader(fileContent)
	if _, err = io.Copy(tmp, reader); err != nil {
		return minio.UploadInfo{}, err
	}

	info, err := c.MinIOClient.FPutObject(
		ctx,
		bucket,
		filename,
		tmp.Name(),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return minio.UploadInfo{}, err
	}

	return info, nil
}

// GetFromMiniIo returns an object from MinIO
func (c *Clients) GetFromMiniIo(ctx context.Context, info minio.UploadInfo, destPath string) (*minio.Object, error) {
	bucket := info.Bucket
	objectName := info.Key

	// Get the object from MinIO
	obj, err := c.MinIOClient.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return obj, fmt.Errorf("failed to get object from MinIO: %w", err)
	}
	// defer obj.Close()

	fmt.Printf("Successfully downloaded %s from bucket %s to %s\n", objectName, bucket, destPath)
	return obj, nil
}

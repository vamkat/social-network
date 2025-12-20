package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
	"strings"
	"time"

	// _ "image/gif"
	// _ "image/jpeg"
	// _ "image/png"

	// _ "golang.org/x/image/webp"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/tags"
)

var ErrMinIO = errors.New("minio error")

func (c *Clients) GenerateDownloadURL(
	ctx context.Context,
	bucket string,
	objectKey string,
	expiry time.Duration,
) (*url.URL, error) {

	client := c.MinIOClient

	// Only for development
	if c.PublicMinIOClient != nil {
		client = c.PublicMinIOClient
	}

	return client.PresignedGetObject(
		ctx,
		bucket,
		objectKey,
		expiry,
		nil,
	)
}

func (c *Clients) GenerateUploadURL(
	ctx context.Context,
	bucket string,
	objectKey string,
	expiry time.Duration,
) (*url.URL, error) {

	client := c.MinIOClient

	// Only for development
	if c.PublicMinIOClient != nil {
		client = c.PublicMinIOClient
	}

	return client.PresignedPutObject(
		ctx,
		bucket,
		objectKey,
		expiry,
	)
}

func (c *Clients) ValidateUpload(
	ctx context.Context,
	fm md.FileMeta,
) error {
	fileCnstr := c.Configs.FileConstraints

	info, err := c.MinIOClient.StatObject(
		ctx,
		fm.Bucket,
		fm.ObjectKey,
		minio.StatObjectOptions{},
	)
	if err != nil {
		return errors.Join(ErrMinIO, err) // upload never completed
	}

	if info.Size != fm.SizeBytes {
		return fmt.Errorf(
			"upload size mismatch: expected=%d actual=%d",
			fm.SizeBytes,
			info.Size,
		)
	}

	ext := strings.ToLower(filepath.Ext(fm.Filename))
	if ok := fileCnstr.AllowedExt[ext]; !ok {
		return fmt.Errorf("invalid file ext %v", ext)
	}

	switch {
	case fileCnstr.AllowedMIMEs[fm.MimeType]:
		if info.Size > fileCnstr.MaxImageUpload {
			return fmt.Errorf("image size %v exceedes allowed size %v",
				info.Size,
				fileCnstr.MaxImageUpload,
			)
		}
		obj, err := c.MinIOClient.GetObject(ctx, fm.Bucket, fm.ObjectKey, minio.GetObjectOptions{})
		if err != nil {
			return errors.Join(ErrMinIO, err)
		}
		defer obj.Close()
		c.Validator.ValidateImage(ctx, obj)
	default:
		return fmt.Errorf("unsuported mime type %v", fm.MimeType)

	}

	tagSet, err := tags.NewTags(map[string]string{
		"validated": "true",
	},
		true,
	)
	if err != nil {
		return errors.Join(ErrMinIO, err)
	}

	err = c.MinIOClient.PutObjectTagging(
		ctx,
		fm.Bucket,
		fm.ObjectKey,
		tagSet,
		minio.PutObjectTaggingOptions{},
	)
	if err != nil {
		return errors.Join(ErrMinIO, err)
	}
	return nil
}

func (c *Clients) DeleteFile(ctx context.Context,
	bucket string,
	objectKey string,
) error {
	return c.MinIOClient.RemoveObject(
		ctx,
		bucket,
		objectKey,
		minio.RemoveObjectOptions{},
	)
}

func (c *Clients) GenerateVariant(
	ctx context.Context,
	bucket string,
	objectKey string,
	variant ct.FileVariant,
) (size int64, err error) {

	obj, err := c.MinIOClient.GetObject(ctx,
		bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return 0, errors.Join(ErrMinIO, err)
	}
	defer obj.Close()

	outBuf, err := c.ImageConvertor.ConvertImageToVariant(obj, variant)

	info, err := c.MinIOClient.PutObject(
		ctx,
		c.Configs.Buckets.Variants,
		objectKey+"/"+variant.String(),
		&outBuf,
		int64(outBuf.Len()),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)
	size = info.Size
	if err != nil {
		return 0, errors.Join(ErrMinIO, err)
	}
	return size, nil
}

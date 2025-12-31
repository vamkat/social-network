package client

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	md "social-network/services/media/internal/models"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/tags"
)

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

	url, err := client.PresignedGetObject(
		ctx,
		bucket,
		objectKey,
		expiry,
		nil,
	)

	if err != nil {
		return nil, ce.Wrap(ce.ErrInternal, err)
	}

	return url, nil
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

	url, err := client.PresignedPutObject(
		ctx,
		bucket,
		objectKey,
		expiry,
	)

	if err != nil {
		return nil, ce.Wrap(ce.ErrInternal, err)
	}

	return url, nil
}

func (c *Clients) ValidateUpload(
	ctx context.Context,
	fm md.FileMeta,
) error {
	errMsg := fmt.Sprintf("S3 client: validate upload: file object key: %v", fm.ObjectKey)

	validated, _ := c.CheckValidationStatus(ctx, fm)
	if validated {
		return nil
	}

	fileCnstr := c.Configs.FileConstraints

	info, err := c.MinIOClient.StatObject(
		ctx,
		fm.Bucket,
		fm.ObjectKey,
		minio.StatObjectOptions{},
	)
	if err != nil {
		return ce.Wrap(ce.ErrNotFound, err, errMsg+": get object stats") // upload never completed
	}

	if info.Size != fm.SizeBytes {
		return ce.Wrap(
			ce.ErrPermissionDenied,
			fmt.Errorf("upload size mismatch: expected=%d actual=%d", fm.SizeBytes, info.Size),
			errMsg+": promised and actual size check",
		).WithPublic("size mismatch")
	}

	ext := strings.ToLower(filepath.Ext(fm.Filename))
	if ok := fileCnstr.AllowedExt[ext]; !ok {
		return ce.Wrap(
			ce.ErrPermissionDenied,
			fmt.Errorf("invalid file ext %v", ext),
			errMsg+": allowed ext check",
		).WithPublic("invalid file extension")
	}

	switch {
	case fileCnstr.AllowedMIMEs[fm.MimeType]:
		if info.Size > fileCnstr.MaxImageUpload {
			return ce.Wrap(ce.ErrPermissionDenied,
				fmt.Errorf("image size %v exceedes allowed size %v",
					info.Size,
					fileCnstr.MaxImageUpload,
				),
				errMsg+": size check",
			).WithPublic("file too big")
		}
		obj, err := c.MinIOClient.GetObject(ctx, fm.Bucket, fm.ObjectKey, minio.GetObjectOptions{})
		if err != nil {
			return ce.Wrap(ce.ErrInternal, err)
		}
		defer obj.Close()
		if err := c.Validator.ValidateImage(ctx, obj); err != nil {
			return err // Validate returns customerrors type with public message
		}
	default:
		return ce.Wrap(ce.ErrPermissionDenied,
			fmt.Errorf("unsuported mime type %v", fm.MimeType),
			errMsg,
		).WithPublic(fmt.Sprintf("unsuported mime type %v", fm.MimeType))
	}

	tagSet, err := tags.NewTags(map[string]string{
		"validated": "true",
	},
		true,
	)
	if err != nil {
		return ce.Wrap(ce.ErrInternal, err, errMsg+": tagSet")
	}

	err = c.MinIOClient.PutObjectTagging(
		ctx,
		fm.Bucket,
		fm.ObjectKey,
		tagSet,
		minio.PutObjectTaggingOptions{},
	)
	if err != nil {
		return ce.Wrap(ce.ErrInternal, err, errMsg+": putObjectTagging")
	}
	return nil
}

func (c *Clients) CheckValidationStatus(ctx context.Context,
	fm md.FileMeta) (bool, error) {
	errMsg := "s3 client: check validation status"
	tagging, err := c.MinIOClient.GetObjectTagging(
		ctx,
		fm.Bucket,
		fm.ObjectKey,
		minio.GetObjectTaggingOptions{},
	)
	if err != nil {
		if minio.ToErrorResponse(err).Code != "NoSuchTagSet" {
			return false, ce.Wrap(ce.ErrInternal, err, errMsg)
		}
	}

	existingTags := tagging.ToMap()
	if existingTags["validated"] == "true" {
		return true, nil
	}
	return false, nil
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

// TODO: Generate many variants for the same file
func (c *Clients) GenerateVariant(
	ctx context.Context,
	srcBucket string,
	srcObjectKey string,
	trgBucket string,
	trgObjectKey string,
	variant ct.FileVariant,
) (size int64, err error) {

	obj, err := c.MinIOClient.GetObject(ctx,
		srcBucket, srcObjectKey, minio.GetObjectOptions{})
	if err != nil {
		return 0, ce.Wrap(ce.ErrNotFound, err)
	}
	defer obj.Close()

	outBuf, err := c.ImageConvertor.ConvertImageToVariant(obj, variant)

	info, err := c.MinIOClient.PutObject(
		ctx,
		trgBucket,
		trgObjectKey,
		&outBuf,
		int64(outBuf.Len()),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)
	size = info.Size
	if err != nil {
		return 0, ce.Wrap(ce.ErrInternal, err)
	}
	return size, nil
}

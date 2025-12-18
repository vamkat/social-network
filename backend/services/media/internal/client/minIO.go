package client

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"math"
	"net/url"
	"path/filepath"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/chai2010/webp"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/tags"
	"golang.org/x/image/draw"
)

func (c *Clients) GenerateDownloadURL(
	ctx context.Context,
	bucket string,
	objectKey string,
	expiry time.Duration,
) (*url.URL, error) {

	client := c.MinIOClient
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
		return err // upload never completed
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
			return err
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
		return err
	}
	return c.MinIOClient.PutObjectTagging(
		ctx,
		fm.Bucket,
		fm.ObjectKey,
		tagSet,
		minio.PutObjectTaggingOptions{},
	)
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
		return 0, err
	}
	defer obj.Close()

	img, _, err := image.Decode(obj)
	if err != nil {
		return 0, err
	}

	resized := resizeForVariant(img, variant)

	var buf bytes.Buffer
	if err := webp.Encode(&buf, resized, &webp.Options{Quality: 80}); err != nil {
		return 0, err
	}

	info, err := c.MinIOClient.PutObject(
		ctx,
		c.Configs.Buckets.Variants,
		objectKey,
		&buf,
		int64(buf.Len()),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)
	size = info.Size
	return size, err
}

func resizeForVariant(src image.Image, variant ct.FileVariant) image.Image {
	maxWidth, maxHeight := variantToSize(variant)
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	ratioW := float64(maxWidth) / float64(w)
	ratioH := float64(maxHeight) / float64(h)
	ratio := math.Min(ratioW, ratioH)

	newW := int(float64(w) * ratio)
	newH := int(float64(h) * ratio)

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	draw.CatmullRom.Scale(
		dst,
		dst.Bounds(),
		src,
		bounds,
		draw.Over,
		nil,
	)

	return dst
}

func variantToSize(variant ct.FileVariant) (maxWidth, maxHeight int) {
	switch variant {
	case ct.ImgLarge:
		return 1600, 1600

	case ct.ImgMedium:
		return 800, 800

	case ct.ImgSmall:
		return 400, 400

	case ct.ImgThumbnail:
		return 150, 150

	default:
		// fallback (treat as medium)
		return 800, 800
	}
}

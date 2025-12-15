package client

import (
	"bytes"
	"context"
	"errors"
	"image"
	"math"
	"net/url"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
	"time"

	"github.com/chai2010/webp"

	"github.com/minio/minio-go/v7"
	"golang.org/x/image/draw"
)

func (c *Clients) GenerateDownloadURL(
	ctx context.Context,
	bucket string,
	objectKey string,
	expiry time.Duration,
) (*url.URL, error) {

	return c.MinIOClient.PresignedGetObject(
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

	return c.MinIOClient.PresignedPutObject(
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
		return errors.New("size mismatch")
	}

	// TODO: Compress and copy to file variant

	// TODO: deep validation (images)
	// obj, err := c.MinIOClient.GetObject(ctx, fm.Bucket, fm.ObjectKey, minio.GetObjectOptions{})
	// if err != nil {
	// 	return err
	// }
	// defer obj.Close()

	// if _, _, err := image.DecodeConfig(obj); err != nil {
	// 	return errors.New("invalid image")
	// }
	if fm.Variant != ct.Large {
		// Generate Variant
		// go c.generateVariant(context.Background(),)
	}
	return nil
}

func (c *Clients) GenerateVariant(
	ctx context.Context,
	fm md.FileMeta,
) error {

	obj, err := c.MinIOClient.GetObject(ctx, fm.Bucket, fm.ObjectKey, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer obj.Close()

	img, _, err := image.Decode(obj)
	if err != nil {
		return err
	}

	resized := resizeForVariant(img, fm.Variant)

	var buf bytes.Buffer
	if err := webp.Encode(&buf, resized, &webp.Options{Quality: 80}); err != nil {
		return err
	}

	_, err = c.MinIOClient.PutObject(
		ctx,
		c.Configs.Buckets.Variants,
		fm.ObjectKey,
		&buf,
		int64(buf.Len()),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)
	return err
}

func resizeForVariant(src image.Image, variant ct.ImgVariant) image.Image {
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

func variantToSize(variant ct.ImgVariant) (maxWidth, maxHeight int) {
	switch variant {
	case ct.Large:
		return 1600, 1600

	case ct.Medium:
		return 800, 800

	case ct.Small:
		return 400, 400

	case ct.Thumbnail:
		return 150, 150

	default:
		// fallback (treat as medium)
		return 800, 800
	}
}

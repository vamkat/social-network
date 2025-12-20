package convertor

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"io"
	"math"
	"social-network/services/media/internal/configs"
	ct "social-network/shared/go/customtypes"

	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
)

type ImageConvertor struct {
	Configs configs.FileConstraints
}

func NewImageconvertor(c configs.FileConstraints) *ImageConvertor {
	return &ImageConvertor{
		Configs: c,
	}
}

// ConvertImageToVariant reads an image from r, resizes it according to the specified variant,
// and encodes it as a WebP image. Variants control the target width and height (e.g., large, medium, small, thumbnail).
// Returns a bytes.Buffer containing the converted image or an error if reading, decoding, resizing,
// or encoding fails. Ensures the input does not exceed the maximum allowed upload size.
func (i *ImageConvertor) ConvertImageToVariant(
	r io.Reader, variant ct.FileVariant,
) (out bytes.Buffer, err error) {

	buf, err := io.ReadAll(r)
	if err != nil {
		return out, fmt.Errorf("failed to read object: %w", err)
	}

	if int64(len(buf)) > i.Configs.MaxImageUpload {
		return out, fmt.Errorf("image size exceeds limit")
	}

	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return out, fmt.Errorf("failed to decode image: %w", err)
	}

	resized := resizeForVariant(img, variant)

	if err := webp.Encode(&out, resized, &webp.Options{Quality: 80}); err != nil {
		return out, err
	}
	return out, nil
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

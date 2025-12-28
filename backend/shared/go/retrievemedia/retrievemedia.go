package retrievemedia

import (
	"context"
	"fmt"
	"maps"
	"time"

	"social-network/shared/gen-go/media"
	"social-network/shared/go/ct"

	"google.golang.org/grpc"
)

// MediaInfoRetriever defines the interface for fetching media info (usually gRPC client).
type MediaInfoRetriever interface {
	GetImages(ctx context.Context, in *media.GetImagesRequest, opts ...grpc.CallOption) (*media.GetImagesResponse, error)
}

type MediaRetriever struct {
	client MediaInfoRetriever
	cache  RedisCache
	ttl    time.Duration
}

func NewMediaRetriever(client MediaInfoRetriever, cache RedisCache, ttl time.Duration) *MediaRetriever {
	return &MediaRetriever{client: client, cache: cache, ttl: ttl}
}

// GetImages returns a map[imageId]imageUrl, using cache + batch RPC.
func (h *MediaRetriever) GetImages(ctx context.Context, imageIds ct.Ids, variant media.FileVariant) (map[int64]string, []int64, error) {

	uniqueImageIds := imageIds.Unique()
	images := make(map[int64]string, len(uniqueImageIds))
	var missingImages ct.Ids

	// Redis lookup for images
	for _, imageId := range uniqueImageIds {
		var imageURL string
		// Cache key: img_<variant_string>:<id>
		// We use the string representation of the variant for clarity in Redis
		// or just the enum value. The existing retrieveusers used "img_thumbnail:<id>".
		// To keep it clean and robust, let's use the variant name.
		key := fmt.Sprintf("img_%s:%d", variant.String(), imageId)

		if err := h.cache.GetObj(ctx, key, &imageURL); err == nil {
			images[imageId.Int64()] = imageURL
		} else {
			missingImages = append(missingImages, imageId)
		}
	}

	var imagesToDelete []int64
	// Batch RPC for missing images
	if len(missingImages) > 0 {
		// fmt.Println("calling media for these images", missingImages)
		req := &media.GetImagesRequest{
			ImgIds:  &media.ImageIds{ImgIds: missingImages.Int64()},
			Variant: variant,
		}

		resp, err := h.client.GetImages(ctx, req)
		if err != nil {
			return nil, nil, err
		}

		for _, failedImage := range resp.FailedIds {
			if failedImage.GetStatus() == 4 || failedImage.GetStatus() == 0 {
				imagesToDelete = append(imagesToDelete, failedImage.FileId)
			}
		}

		// merge with redis map
		maps.Copy(images, resp.DownloadUrls)

		// Cache the new results
		for id, url := range resp.DownloadUrls {
			key := fmt.Sprintf("img_%s:%d", variant.String(), id)
			_ = h.cache.SetObj(ctx, key, url, h.ttl)
		}
	}

	return images, imagesToDelete, nil
}

// GetImage returns a single image url, using cache + batch RPC.
func (h *MediaRetriever) GetImage(ctx context.Context, imageId int64, variant media.FileVariant) (string, error) {
	images, _, err := h.GetImages(ctx, ct.Ids{ct.Id(imageId)}, variant)
	if err != nil {
		return "", err
	}
	return images[imageId], nil
}

package retrieveusers

import (
	"context"
	"fmt"
	"maps"

	"social-network/shared/gen-go/media"
	"social-network/shared/go/ct"
	redis_connector "social-network/shared/go/redis"
	"time"
)

type MediaRetriever struct {
	clients Client
	cache   RedisCache
	ttl     time.Duration
}

func NewMediaRetriever(clients Client, cache *redis_connector.RedisClient, ttl time.Duration) *MediaRetriever {
	return &MediaRetriever{clients: clients, cache: cache, ttl: ttl}
}

// GetImages returns a map[imageId]imageUrl, using cache + batch RPC.
func (h *MediaRetriever) GetImages(ctx context.Context, imageIds ct.Ids, variant media.FileVariant) (map[int64]string, []int64, error) {

	uniqueImageIds := imageIds.Unique()

	//fmt.Println("Unique image ids", uniqueImageIds)

	images := make(map[int64]string, len(uniqueImageIds))
	var missingImages ct.Ids

	// Redis lookup for images
	for _, imageId := range uniqueImageIds {
		var imageURL string
		if err := h.cache.GetObj(ctx, fmt.Sprintf("img_%s:%d", variant, imageId), &imageURL); err == nil {
			images[imageId.Int64()] = imageURL
		} else {
			missingImages = append(missingImages, imageId)
		}
	}

	//fmt.Println("found on redis", images)

	var imagesToDelete []int64
	// // Batch RPC for missing images
	if len(missingImages) > 0 {
		fmt.Println("calling media for these images", missingImages)
		req := &media.GetImagesRequest{
			ImgIds:  &media.ImageIds{ImgIds: imageIds.Int64()},
			Variant: variant,
		}
		resp, err := h.clients.GetImages(ctx, req, &variant)
		if err != nil {
			return nil, nil, err
		}
		for _, failedImage := range resp.FailedIds {
			if failedImage.GetStatus() == 4 || failedImage.GetStatus() == 0 {
				imagesToDelete = append(imagesToDelete, failedImage.FileId)
			}
		}
		//merge with redis map
		maps.Copy(images, resp.DownloadUrls)
		fmt.Println("failed", imagesToDelete)
		//fmt.Println("map", images)
		for id, image := range images {

			_ = h.cache.SetObj(ctx, fmt.Sprintf("img_%s:%d", variant, id), &image, h.ttl)
		}
	}

	//TODO batch call to delete missing image ids

	return images, imagesToDelete, nil
}

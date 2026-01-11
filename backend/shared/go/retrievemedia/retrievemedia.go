package retrievemedia

import (
	"context"
	"fmt"
	"maps"
	"time"

	"social-network/shared/gen-go/media"
	ce "social-network/shared/go/commonerrors"
	"social-network/shared/go/ct"
	"social-network/shared/go/mapping"

	"google.golang.org/grpc"
)

// MediaInfoRetriever defines the interface for fetching media info (usually gRPC client).
type MediaInfoRetriever interface {
	GetImages(ctx context.Context, in *media.GetImagesRequest, opts ...grpc.CallOption) (*media.GetImagesResponse, error)
	GetImage(ctx context.Context, in *media.GetImageRequest, opts ...grpc.CallOption) (*media.GetImageResponse, error)
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
	input := fmt.Sprintf("ids %v, variant: %v", imageIds, variant)

	uniqueImageIds := imageIds.Unique()
	images := make(map[int64]string, len(uniqueImageIds))
	var missingImages ct.Ids

	ctVariant := mapping.PbToCtFileVariant(variant)
	if err := ctVariant.Validate(); err != nil {
		// If variant is invalid, we probably can't do anything meaningful.
		// Returning error or empty map? Returning error seems safest.
		return nil, nil, ce.New(ce.ErrInvalidArgument, err, input)
	}

	// Redis lookup for images
	for _, imageId := range uniqueImageIds {
		key, err := ct.ImageKey{Id: imageId, Variant: ctVariant}.String()
		if err != nil {
			fmt.Printf("RETRIEVE MEDIA - failed to construct redis key for image %v: %v\n", imageId, err)
			missingImages = append(missingImages, imageId)
			continue
		}

		imageURL, err := h.cache.GetStr(ctx, key)
		if err == nil {
			fmt.Println("Got Image from redis")
			images[imageId.Int64()] = imageURL.(string)
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
			return nil, nil, ce.ParseGrpcErr(err, input)
		}

		for _, failedImage := range resp.FailedIds {
			if failedImage.GetStatus() == media.UploadStatus_UPLOAD_STATUS_FAILED {
				imagesToDelete = append(imagesToDelete, failedImage.FileId)
			}
		}

		// merge with redis map
		maps.Copy(images, resp.DownloadUrls)

		// Cache the new results
		for id, url := range resp.DownloadUrls {
			key, err := ct.ImageKey{Id: ct.Id(id), Variant: ctVariant}.String()
			if err == nil {
				_ = h.cache.SetStr(ctx, key, url, h.ttl)
			} else {
				fmt.Printf("RETRIEVE MEDIA - failed to construct redis key for caching image %v: %v\n", id, err)
			}
		}

		//==============pinpoint not found images============
		for _, imageId := range uniqueImageIds {
			id := imageId.Int64()

			// skip if download succeeded
			if _, exists := resp.DownloadUrls[id]; exists {
				continue
			}

			// skip if download failed
			foundFailed := false
			for _, failedImage := range resp.FailedIds {
				if failedImage.FileId == id {
					foundFailed = true
					break
				}
			}
			if foundFailed {
				continue
			}

			imagesToDelete = append(imagesToDelete, id) //now imagesToDelete includes failed and not found
		}
	}

	return images, imagesToDelete, nil
}

// func toCtVariant(v media.FileVariant) (ct.FileVariant, error) {
// 	switch v {
// 	case media.FileVariant_THUMBNAIL:
// 		return ct.ImgThumbnail, nil
// 	case media.FileVariant_SMALL:
// 		return ct.ImgSmall, nil
// 	case media.FileVariant_MEDIUM:
// 		return ct.ImgMedium, nil
// 	case media.FileVariant_LARGE:
// 		return ct.ImgLarge, nil
// 	case media.FileVariant_ORIGINAL:
// 		return ct.Original, nil
// 	default:
// 		return "", fmt.Errorf("unknown media variant: %v", v)
// 	}
// }

// GetImage returns a single image url, using cache + batch RPC.
func (h *MediaRetriever) GetImage(ctx context.Context, imageId int64, variant media.FileVariant) (string, error) {
	input := fmt.Sprintf("id %v, variant: %v", imageId, variant)

	ctVariant := mapping.PbToCtFileVariant(variant)
	if err := ctVariant.Validate(); err != nil {
		return "", ce.New(ce.ErrInvalidArgument, err, input)
	}

	// Redis lookup for images
	key, err := ct.ImageKey{Id: ct.Id(imageId), Variant: ctVariant}.String()
	if err != nil {
		fmt.Printf("RETRIEVE MEDIA - failed to construct redis key for image %v: %v\n", imageId, err)
	}

	imageURL, err := h.cache.GetStr(ctx, key)
	if err == nil {
		fmt.Println("Got Image from redis")
		return imageURL.(string), nil
	}

	imageURL = &media.GetImageRequest{
		ImageId: imageId,
		Variant: variant,
	}

	resp, err := h.client.GetImage(ctx, &media.GetImageRequest{
		ImageId: imageId,
		Variant: variant,
	})
	if err != nil {
		return "", ce.ParseGrpcErr(err, input)
	}

	//if err, check if need to delete image

	// Cache the new result
	key, err = ct.ImageKey{Id: ct.Id(imageId), Variant: ctVariant}.String()
	if err == nil {
		_ = h.cache.SetStr(ctx, key, resp.DownloadUrl, h.ttl)
	} else {
		fmt.Printf("RETRIEVE MEDIA - failed to construct redis key for caching image %v: %v\n", imageId, err)
	}

	return resp.DownloadUrl, nil
}

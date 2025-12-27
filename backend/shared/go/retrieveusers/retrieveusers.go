package retrieveusers

import (
	"context"
	"fmt"
	"maps"

	"social-network/shared/gen-go/media"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	redis_connector "social-network/shared/go/redis"
	"time"
)

type UserRetriever struct {
	clients UsersBatchClient
	cache   RedisCache
	ttl     time.Duration
}

func NewUserRetriever(clients UsersBatchClient, cache *redis_connector.RedisClient, ttl time.Duration) *UserRetriever {
	return &UserRetriever{clients: clients, cache: cache, ttl: ttl}
}

// GetUsers returns a map[userID]User, using cache + batch RPC.
func (h *UserRetriever) GetUsers(ctx context.Context, userIDs ct.Ids) (map[ct.Id]models.User, error) {

	//========================== STEP 1 : get user info from users ===============================================

	ids := userIDs.Unique()

	users := make(map[ct.Id]models.User, len(ids))
	var missing ct.Ids

	// Redis lookup
	for _, id := range ids {
		var u models.User
		if err := h.cache.GetObj(ctx, fmt.Sprintf("basic_user_info:%d", id), &u); err == nil {
			users[id] = u
			fmt.Println("RETRIEVE USERS - found user on redis:", u)
		} else {
			missing = append(missing, id)
		}
	}

	// Batch RPC for missing users
	if len(missing) > 0 {
		resp, err := h.clients.GetBatchBasicUserInfo(ctx, missing)
		if err != nil {
			return nil, err
		}

		for _, u := range resp.Users {
			user := models.User{
				UserId:   ct.Id(u.UserId),
				Username: ct.Username(u.Username),
				AvatarId: ct.Id(u.Avatar),
			}
			users[user.UserId] = user
			_ = h.cache.SetObj(ctx,
				fmt.Sprintf("basic_user_info:%d", u.UserId),
				user,
				h.ttl,
			)
		}
	}
	//========================== STEP 2 : get avatars from media ===============================================
	// Get image urls for users
	var imageIds ct.Ids
	for _, user := range users {
		if user.AvatarId > 0 { //exclude 0 imageIds
			imageIds = append(imageIds, user.AvatarId)
		}
	}

	uniqueImageIds := imageIds.Unique()

	//fmt.Println("Unique image ids", uniqueImageIds)

	images := make(map[int64]string, len(uniqueImageIds))
	var missingImages ct.Ids

	// Redis lookup for images
	for _, imageId := range uniqueImageIds {
		var imageURL string
		if err := h.cache.GetObj(ctx, fmt.Sprintf("img_thumbnail:%d", imageId), &imageURL); err == nil {
			images[imageId.Int64()] = imageURL
			fmt.Println("RETRIEVE USERS - found image on redis:", imageURL)
		} else {
			missingImages = append(missingImages, imageId)
		}
	}

	//fmt.Println("found on redis", images)
	//fmt.Println("users before image urls", users)
	imageMap := make(map[int64]string)
	failedImages := []int64{}
	var err error
	// // Batch RPC for missing images
	if len(missingImages) > 0 {
		fmt.Println("asking media for these avatars", missingImages)
		imageMap, failedImages, err = h.clients.GetImages(ctx, missingImages, media.FileVariant_THUMBNAIL)
		if err != nil {
			return nil, err
		}
	}
	//merge with redis map
	maps.Copy(images, imageMap)

	for id, u := range users {

		if url, ok := images[u.AvatarId.Int64()]; ok {
			u.AvatarURL = url
			users[id] = u

			_ = h.cache.SetObj(ctx,
				fmt.Sprintf("img_thumbnail:%d", u.AvatarId.Int64()),
				url,
				h.ttl,
			)
		}
	}
	//fmt.Println("users after image urls", users)

	//TODO batch call to delete missing image ids
	fmt.Println("RETRIEVE USERS - failed avatars:", failedImages)
	fmt.Println("RETRIEVE USERS - found avatars:", images)

	return users, nil
}

func (h *UserRetriever) GetImages(ctx context.Context, imageIds ct.Ids, variant media.FileVariant) (map[int64]string, []int64, error) {
	return h.clients.GetImages(ctx, imageIds, variant)
}

package retrieveusers

import (
	"context"
	"fmt"

	ct "social-network/shared/go/customtypes"
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
func (h *UserRetriever) GetUsers(ctx context.Context, userIDs []int64) (map[int64]models.User, error) {

	//========================== STEP 1 : get user info from users ===============================================

	idSet := make(map[int64]struct{}, len(userIDs))
	for _, id := range userIDs {
		idSet[id] = struct{}{}
	}

	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	users := make(map[int64]models.User, len(ids))
	var missing []int64

	// Redis lookup
	for _, id := range ids {
		var u models.User
		if err := h.cache.GetObj(ctx, fmt.Sprintf("basic_user_info:%d", id), &u); err == nil {
			users[id] = u
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
			users[u.UserId] = user
			_ = h.cache.SetObj(ctx,
				fmt.Sprintf("basic_user_info:%d", u.UserId),
				user,
				h.ttl,
			)
		}
	}
	//========================== STEP 2 : get avatars from media ===============================================
	// Get image urls for users
	var imageIds []int64
	for _, user := range users {
		imageIds = append(imageIds, user.AvatarId.Int64())
	}

	//there shouldn't be duplicates but making sure
	imageIdSet := make(map[int64]struct{}, len(imageIds))
	for _, imageId := range imageIds {
		imageIdSet[imageId] = struct{}{}
	}

	uniqueImageIds := make([]int64, 0, len(imageIdSet))
	for imageId := range imageIdSet {
		uniqueImageIds = append(uniqueImageIds, imageId)
	}

	images := make(map[int64]string, len(uniqueImageIds))
	var missingImages []int64

	// Redis lookup for images
	for _, imageId := range imageIds {
		var imageURL string
		if err := h.cache.GetObj(ctx, fmt.Sprintf("img_thumbnail:%d", imageId), &imageURL); err == nil {
			images[imageId] = imageURL
		} else {
			missingImages = append(missingImages, imageId)
		}
	}

	//TODO:
	// Init client for media in users
	// send it to retrieve users
	// change the basic user model to include url, or add a new model
	// decide on redis keys
	// test

	// // Batch RPC for missing images
	if len(missingImages) > 0 {
		imageMap, failedImages, err := h.clients.GetImages(ctx, missingImages)
		if err != nil {
			return nil, err
		}

		for id, u := range users {

			if url, ok := imageMap[u.AvatarId.Int64()]; ok {
				u.AvatarURL = url
				users[id] = u

				_ = h.cache.SetObj(ctx,
					fmt.Sprintf("img_thumbnail:%d", u.AvatarId.Int64()),
					url,
					h.ttl,
				)
			}
		}

		//batch call to delete missing image ids
		fmt.Println(failedImages)
	}

	return users, nil
}

func (h *UserRetriever) GetImages(ctx context.Context, imageIds []int64) (map[int64]string, []int64, error) {
	return h.clients.GetImages(ctx, imageIds)
}

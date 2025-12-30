package retrieveusers

import (
	"context"
	"fmt"

	cm "social-network/shared/gen-go/common"
	"social-network/shared/gen-go/media"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	redis_connector "social-network/shared/go/redis"
	"social-network/shared/go/retrievemedia"
	"time"
)

type UserRetriever struct {
	GetBatchBasicUserInfo GetBatchBasicUserInfo
	cache                 RedisCache
	mediaRetriever        *retrievemedia.MediaRetriever
	ttl                   time.Duration
}

// UserRetriever provides a function that abstracts the process of populating a map[ct.Id]models.User
// from a slice of user ids. It depends on:
//
//  1. GetBatchBasicUserInfo function provided by an initiator that holds a connection to social-network user service.
//  2. A cache interface that implements GetObj() and SetObj() methods.
//  3. The retrievemedia package.
func NewUserRetriever(userClient GetBatchBasicUserInfo, cache *redis_connector.RedisClient, mediaRetriever *retrievemedia.MediaRetriever, ttl time.Duration) *UserRetriever {
	return &UserRetriever{GetBatchBasicUserInfo: userClient, cache: cache, mediaRetriever: mediaRetriever, ttl: ttl}
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

		key, err := ct.BasicUserInfoKey{Id: id}.String()
		if err != nil {
			fmt.Printf("RETRIEVE USERS - failed to construct redis key for id %v: %v\n", id, err)
			missing = append(missing, id)
			continue
		}

		if err := h.cache.GetObj(ctx, key, &u); err == nil {
			users[id] = u
			fmt.Println("RETRIEVE USERS - found user on redis:", u)
		} else {
			missing = append(missing, id)
		}
	}

	// Batch RPC for missing users
	if len(missing) > 0 {
		resp, err := h.GetBatchBasicUserInfo(ctx, &cm.UserIds{Values: missing.Int64()})
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

			key, err := ct.BasicUserInfoKey{Id: user.UserId}.String()
			if err == nil {
				_ = h.cache.SetObj(ctx,
					key,
					user,
					h.ttl,
				)
			} else {
				fmt.Printf("RETRIEVE USERS - failed to construct redis key for user %v: %v\n", user.UserId, err)
			}
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
	imageIds = imageIds.Unique()
	if len(imageIds) > 0 {
		// Use shared MediaRetriever for images (handles caching and fetching)
		imageMap, _, err := h.mediaRetriever.GetImages(ctx, imageIds, media.FileVariant_THUMBNAIL)
		if err != nil {
			return nil, err
		}

		for id, u := range users {
			if url, ok := imageMap[u.AvatarId.Int64()]; ok {
				u.AvatarURL = url
				users[id] = u
			}
		}
	}

	return users, nil
}

// func (h *UserRetriever) GetImages(ctx context.Context, imageIds ct.Ids, variant media.FileVariant) (map[int64]string, []int64, error) {
// 	return h.mediaRetriever.GetImages(ctx, imageIds, variant)
// }

package userhydrate

import (
	"context"
	"fmt"

	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
	redis_connector "social-network/shared/go/redis"
	"time"
)

type UserHydrator struct {
	clients UsersBatchClient
	cache   RedisCache
	ttl     time.Duration
}

func NewUserHydrator(clients UsersBatchClient, cache *redis_connector.RedisClient, ttl time.Duration) *UserHydrator {
	return &UserHydrator{clients: clients, cache: cache, ttl: ttl}
}

// GetUsers returns a map[userID]User, using cache + batch RPC.
func (h *UserHydrator) GetUsers(ctx context.Context, userIDs []int64) (map[int64]models.User, error) {
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

	return users, nil
}

// // hydrateUsers fills in user data for all HasUser items from redis and/or user service
// func (h *UserHydrator) HydrateUsers(ctx context.Context, items []models.HasUser) error {
// 	//collect unique user IDs
// 	idSet := make(map[int64]struct{}, len(items))

// 	for _, item := range items {
// 		idSet[item.GetUserId()] = struct{}{}
// 	}

// 	ids := make([]int64, 0, len(idSet))
// 	for id := range idSet {
// 		ids = append(ids, id)
// 	}

// 	//check Redis first
// 	cachedUsers := make(map[int64]models.User)
// 	var missingIDs []int64

// 	// for _, id := range ids {
// 	// 	var u models.User
// 	// 	err := h.cache.GetObj(ctx, fmt.Sprintf("basic_user_info:%d", id), &u)
// 	// 	if err != nil {
// 	// 		if err == redis_connector.ErrNotFound {
// 	// 			missingIDs = append(missingIDs, id)
// 	// 		} else {
// 	// 			return err
// 	// 		}
// 	// 	} else {
// 	// 		cachedUsers[id] = u
// 	// 	}
// 	// }

// 	// TODO: When Redis is disabled for testing, fetch all users from service
// 	missingIDs = ids

// 	// fetch missing users from Users service
// 	if len(missingIDs) > 0 {
// 		resp, err := h.clients.GetBatchBasicUserInfo(ctx, missingIDs)
// 		if err != nil {
// 			return err
// 		}

// 		for _, u := range resp.Users {
// 			user := models.User{
// 				UserId:   ct.Id(u.UserId),
// 				Username: ct.Username(u.Username),
// 				AvatarId: ct.Id(u.Avatar),
// 			}
// 			cachedUsers[u.UserId] = user

// 			// cache it
// 			// if err := h.cache.SetObj(ctx, fmt.Sprintf("basic_user_info:%d", u.UserId), user, h.ttl); err != nil {
// 			// 	fmt.Printf("failed to cache user %d: %v\n", u.UserId, err)
// 			// }
// 		}
// 	}

// 	// fill original items
// 	for _, item := range items {
// 		id := item.GetUserId()
// 		if u, ok := cachedUsers[id]; ok {
// 			item.SetUser(u)
// 		}
// 	}

// 	return nil
// }

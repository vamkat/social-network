package userhydrate

import (
	"context"

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

// hydrateUsers fills in user data for all HasUser items from redis and/or user service
func (h *UserHydrator) HydrateUsers(ctx context.Context, items []models.HasUser) error {
	//collect unique user IDs
	idSet := make(map[int64]struct{}, len(items))

	for _, item := range items {
		idSet[item.GetUserId()] = struct{}{}
	}

	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	//check Redis first
	cachedUsers := make(map[int64]models.User)
	var missingIDs []int64

	// for _, id := range ids {
	// 	var u models.User
	// 	err := h.cache.GetObj(ctx, fmt.Sprintf("basic_user_info:%d", id), &u)
	// 	if err != nil {
	// 		if err == redis_connector.ErrNotFound {
	// 			missingIDs = append(missingIDs, id)
	// 		} else {
	// 			return err
	// 		}
	// 	} else {
	// 		cachedUsers[id] = u
	// 	}
	// }

	// TODO: When Redis is disabled for testing, fetch all users from service
	missingIDs = ids

	// fetch missing users from Users service
	if len(missingIDs) > 0 {
		resp, err := h.clients.GetBatchBasicUserInfo(ctx, missingIDs)
		if err != nil {
			return err
		}

		for _, u := range resp.Users {
			user := models.User{
				UserId:   ct.Id(u.UserId),
				Username: ct.Username(u.Username),
				AvatarId: ct.Id(u.Avatar),
			}
			cachedUsers[u.UserId] = user

			// cache it
			// if err := h.cache.SetObj(ctx, fmt.Sprintf("basic_user_info:%d", u.UserId), user, h.ttl); err != nil {
			// 	fmt.Printf("failed to cache user %d: %v\n", u.UserId, err)
			// }
		}
	}

	// fill original items
	for _, item := range items {
		id := item.GetUserId()
		if u, ok := cachedUsers[id]; ok {
			item.SetUser(u)
		}
	}

	return nil
}

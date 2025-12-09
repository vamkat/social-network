package application

import (
	"context"
	"fmt"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
	redis_connector "social-network/shared/go/redis"
)

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

	for _, id := range ids {
		var u models.User
		err := h.cache.GetObj(ctx, fmt.Sprintf("basic_user_info:%d", id), &u)
		if err != nil {
			if err == redis_connector.ErrNotFound {
				missingIDs = append(missingIDs, id)
			} else {
				return err
			}
		} else {
			cachedUsers[id] = u
		}
	}

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
			if err := h.cache.SetObj(ctx, fmt.Sprintf("basic_user_info:%d", u.UserId), user, h.ttl); err != nil {
				fmt.Printf("failed to cache user %d: %v\n", u.UserId, err)
			}
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

func (s *Application) hydratePost(ctx context.Context, post *models.Post) error {
	items := []models.HasUser{post}

	return s.hydrator.HydrateUsers(ctx, items)
}

func (s *Application) hydratePosts(ctx context.Context, posts []models.Post) error {
	items := make([]models.HasUser, 0, len(posts)*2)

	for i := range posts {
		// Always include post creator
		items = append(items, &posts[i])

	}

	return s.hydrator.HydrateUsers(ctx, items)
}

func (s *Application) hydrateComments(ctx context.Context, comments []models.Comment) error {
	items := make([]models.HasUser, 0, len(comments)*2)

	for i := range comments {
		// Always include post creator
		items = append(items, &comments[i])

	}

	return s.hydrator.HydrateUsers(ctx, items)
}

func (s *Application) hydrateEvents(ctx context.Context, events []models.Event) error {
	items := make([]models.HasUser, 0, len(events)*2)

	for i := range events {
		// Always include post creator
		items = append(items, &events[i])

	}

	return s.hydrator.HydrateUsers(ctx, items)
}

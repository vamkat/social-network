package application

import (
	"context"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
)

func (h *UserHydrator) HydrateUsers(ctx context.Context, items []models.HasUser) error {
	idSet := make(map[int64]struct{}, len(items))

	for _, item := range items {
		idSet[item.GetUserId()] = struct{}{}
	}

	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	resp, err := h.clients.GetBatchBasicUserInfo(ctx, ids)
	if err != nil {
		return err
	}

	userMap := make(map[int64]models.User, len(resp.Users))
	for _, u := range resp.Users {
		userMap[u.UserId] = models.User{
			UserId:   ct.Id(u.UserId),
			Username: ct.Username(u.Username),
			AvatarId: ct.Id(u.Avatar),
		}
	}

	for _, item := range items {
		id := item.GetUserId()
		if u, ok := userMap[id]; ok {
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

package application

// import (
// 	"context"
// 	"social-network/shared/go/models"
// )

// func (s *Application) hydrateComments(ctx context.Context, comments []models.Comment) error {
// 	items := make([]models.HasUser, 0, len(comments)*2)

// 	for i := range comments {
// 		// Always include post creator
// 		items = append(items, &comments[i])

// 	}

// 	return s.hydrator.HydrateUsers(ctx, items)
// }

// func (s *Application) hydrateEvents(ctx context.Context, events []models.Event) error {
// 	items := make([]models.HasUser, 0, len(events)*2)

// 	for i := range events {
// 		// Always include post creator
// 		items = append(items, &events[i])

// 	}

// 	return s.hydrator.HydrateUsers(ctx, items)
// }

// func (s *Application) hydratePost(ctx context.Context, post *models.Post) error {
// 	items := []models.HasUser{post}

// 	return s.hydrator.HydrateUsers(ctx, items)
// }

// func (s *Application) hydratePosts(ctx context.Context, posts []models.Post) error {
// 	items := make([]models.HasUser, 0, len(posts)*2)

// 	for i := range posts {
// 		// Always include post creator
// 		items = append(items, &posts[i])

// 	}

// 	return s.hydrator.HydrateUsers(ctx, items)
// }

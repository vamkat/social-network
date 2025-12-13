package application

import (
	"context"
	"database/sql"
	"errors"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
)

// Returns five random ids that fit one of the following criteria:
// Users who liked one or more of *your public posts*
// Users who commented on your public posts
// Users who liked the same posts as you
// Users who commented on the same posts as you
// Actual Basic User Info will be retrieved by HANDLER from users
func (s *Application) SuggestUsersByPostActivity(ctx context.Context, req models.SimpleIdReq) ([]models.User, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}
	ids, err := s.db.SuggestUsersByPostActivity(ctx, req.Id.Int64())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	userMap, err := s.hydrator.GetUsers(ctx, ids)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0, len(ids))
	for _, id := range ids {
		if u, ok := userMap[id]; ok {
			users = append(users, u)
		}
	}

	// err = s.hydrator.HydrateUserSlice(ctx, users)
	// if err != nil {
	// 	return nil, err
	// }

	return users, nil
}

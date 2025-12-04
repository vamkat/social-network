package application

import (
	"context"
	"database/sql"
	"errors"
	ct "social-network/shared/go/customtypes"
)

// Returns five random ids that fit one of the following criteria:
// Users who liked one or more of *your public posts*
// Users who commented on your public posts
// Users who liked the same posts as you
// Users who commented on the same posts as you
// Actual Basic User Info will be retrieved by HANDLER from users
func (s *Application) SuggestUsersByPostActivity(ctx context.Context, userId ct.Id) (ct.Ids, error) {
	if err := ct.ValidateStruct(userId); err != nil {
		return nil, err
	}
	ids, err := s.db.SuggestUsersByPostActivity(ctx, userId.Int64())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return ct.FromInt64s(ids), nil
}

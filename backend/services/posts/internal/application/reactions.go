package application

import (
	"context"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
)

func (s *Application) ToggleOrInsertReaction(ctx context.Context, req GenericReq) error {
	// check requester can actually view parent entity?
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	rowsAffected, err := s.db.ToggleOrInsertReaction(ctx, sqlc.ToggleOrInsertReactionParams{
		ContentID: req.EntityId.Int64(),
		UserID:    req.RequesterId.Int64(),
	})
	if err != nil || rowsAffected != 1 {
		return err
	}

	return nil
}

// SKIP FOR NOW
func (s *Application) GetWhoLikedEntityId(ctx context.Context, req GenericReq) ([]int64, error) {
	return nil, nil
}

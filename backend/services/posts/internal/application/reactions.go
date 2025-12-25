package application

import (
	"context"
	ds "social-network/services/posts/internal/db/dbservice"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
)

func (s *Application) ToggleOrInsertReaction(ctx context.Context, req models.GenericReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrNotAllowed
	}

	rowsAffected, err := s.db.ToggleOrInsertReaction(ctx, ds.ToggleOrInsertReactionParams{
		ContentID: req.EntityId.Int64(),
		UserID:    req.RequesterId.Int64(),
	})
	if err != nil || rowsAffected != 1 {
		return err
	}

	return nil
}

// SKIP FOR NOW
func (s *Application) GetWhoLikedEntityId(ctx context.Context, req models.GenericReq) ([]int64, error) {
	return nil, nil
}

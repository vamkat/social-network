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

	action, err := s.db.ToggleOrInsertReaction(ctx, ds.ToggleOrInsertReactionParams{
		ContentID: req.EntityId.Int64(),
		UserID:    req.RequesterId.Int64(),
	})
	if err != nil {
		return err
	}
	if action == "added" {
		//create notification
		userMap, err := s.userRetriever.GetUsers(ctx, ct.Ids{req.RequesterId})
		if err != nil {
			//log error
		}
		var likerUsername string
		if u, ok := userMap[req.RequesterId]; ok {
			likerUsername = u.Username.String()
		}
		row, err := s.db.GetEntityCreatorAndGroup(ctx, req.EntityId.Int64())
		if err != nil {
			//log and don't proceed to notif
		}
		err = s.clients.CreatePostLike(ctx, row.CreatorID, req.RequesterId.Int64(), req.EntityId.Int64(), likerUsername)
	} else {
		//remove notification or not? how?
	}
	return nil
}

// SKIP FOR NOW
func (s *Application) GetWhoLikedEntityId(ctx context.Context, req models.GenericReq) ([]int64, error) {
	return nil, nil
}

package application

import (
	"context"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
)

// group and post audience=group: only members can see
// post audience=everyone: everyone can see (can we check this before all the fetches from users?)
// post audience=followers: requester can see if they follow creator
// post audience=selected: requester can see if they are in post audience table
func (s *Application) hasRightToView(ctx context.Context, req AccessContext) (bool, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}
	canSee, err := s.db.CanUserSeeEntity(ctx, sqlc.CanUserSeeEntityParams{
		UserID:       req.RequesterId.Int64(),
		FollowingIds: req.RequesterFollowsIds.Int64(),
		GroupIds:     req.RequesterGroups.Int64(),
		EntityID:     req.ParentEntityId.Int64(),
	})
	if err != nil {
		return false, err
	}
	return canSee, nil
}

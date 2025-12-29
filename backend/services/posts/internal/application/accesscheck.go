package application

import (
	"context"
	ds "social-network/services/posts/internal/db/dbservice"
)

// group and post audience=group: only members can see
// post audience=everyone: everyone can see (can we check this before all the fetches from users?)
// post audience=followers: requester can see if they follow creator
// post audience=selected: requester can see if they are in post audience table
func (s *Application) hasRightToView(ctx context.Context, req accessContext) (bool, error) {

	row, err := s.db.GetEntityCreatorAndGroup(ctx, req.entityId)
	if err != nil {
		return false, err
	}


	isFollowing, err := s.clients.IsFollowing(ctx, req.requesterId, row.CreatorID)
	if err != nil {
		return false, err
	}

	var isMember bool
	if row.GroupID > 0 {
		isMember, err = s.clients.IsGroupMember(ctx, req.requesterId, row.GroupID)
		if err != nil {
			return false, err
		}
	}

	canSee, err := s.db.CanUserSeeEntity(ctx, ds.CanUserSeeEntityParams{
		UserID:      req.requesterId,
		EntityID:    req.entityId,
		IsFollowing: isFollowing,
		IsMember:    isMember,
	})
	if err != nil {
		return false, err
	}
	return canSee, nil
}

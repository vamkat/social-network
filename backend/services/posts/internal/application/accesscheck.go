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

	var targetUserId int64
	if row.ParentCreatorID > 0 { //in case of comment, we need the parent creator for isFollowing
		targetUserId = row.ParentCreatorID
	} else {
		targetUserId = row.CreatorID
	}
	isFollowing, err := s.clients.IsFollowing(ctx, req.requesterId, targetUserId)
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

	entityID := req.entityId //this is the event or post id - in case of a comment we take the parent post id
	if row.ParentID > 0 {
		entityID = row.ParentID
	}

	canSee, err := s.db.CanUserSeeEntity(ctx, ds.CanUserSeeEntityParams{
		UserID:      req.requesterId,
		EntityID:    entityID,
		IsFollowing: isFollowing,
		IsMember:    isMember,
	})
	if err != nil {
		return false, err
	}
	return canSee, nil
}

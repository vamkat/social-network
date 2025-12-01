package application

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
)

func (s *Application) GetFollowersPaginated(ctx context.Context, req Pagination) ([]User, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return []User{}, err
	}
	//paginated, sorted by newest first
	rows, err := s.db.GetFollowers(ctx, sqlc.GetFollowersParams{
		FollowingID: req.UserId.Int64(),
		Limit:       req.Limit.Int32(),
		Offset:      req.Offset.Int32(),
	})
	if err != nil {
		return []User{}, err
	}
	users := make([]User, 0, len(rows))
	for _, r := range rows {
		users = append(users, User{
			UserId:   ct.Id(r.ID),
			Username: ct.Username(r.Username),
			Avatar:   r.Avatar,
		})
	}

	return users, nil

}

func (s *Application) GetFollowingPaginated(ctx context.Context, req Pagination) ([]User, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return []User{}, err
	}

	//paginated, sorted by newest first
	rows, err := s.db.GetFollowing(ctx, sqlc.GetFollowingParams{
		FollowerID: req.UserId.Int64(),
		Limit:      req.Limit.Int32(),
		Offset:     req.Offset.Int32(),
	})
	if err != nil {
		return []User{}, err
	}
	users := make([]User, 0, len(rows))
	for _, r := range rows {
		users = append(users, User{
			UserId:   ct.Id(r.ID),
			Username: ct.Username(r.Username),
			Avatar:   r.Avatar,
		})
	}

	return users, nil

}

// CHAT SERVICE EVENT should trigger event that creates conversation between two users for chat service (unless the target user was already following, so conversation exists)
func (s *Application) FollowUser(ctx context.Context, req FollowUserReq) (resp FollowUserResp, err error) {
	if err := ct.ValidateStruct(req); err != nil {
		return FollowUserResp{}, err
	}
	status, err := s.db.FollowUser(ctx, sqlc.FollowUserParams{
		PFollower: req.FollowerId.Int64(),
		PTarget:   req.TargetUserId.Int64(),
	})
	if err != nil {
		return FollowUserResp{}, err
	}
	if status == "requested" { //I don't love hardcoding this, we'll see
		resp.IsPending = true
		resp.ViewerIsFollowing = false
	} else {
		resp.IsPending = false
		resp.ViewerIsFollowing = true
	}
	return resp, nil
}

// CHAT SERVICE EVENT should it trigger event to delete conversation if none of the two follow each other any more? Or just make it inactive?
func (s *Application) UnFollowUser(ctx context.Context, req FollowUserReq) (viewerIsFollowing bool, err error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}
	err = s.db.UnfollowUser(ctx, sqlc.UnfollowUserParams{
		FollowerID:  req.FollowerId.Int64(),
		FollowingID: req.TargetUserId.Int64(),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// CHAT SERVICE EVENT accepting a follow request should also trigger event to create conversation
func (s *Application) HandleFollowRequest(ctx context.Context, req HandleFollowRequestReq) error {
	var err error
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	if req.Accept {
		err = s.db.AcceptFollowRequest(ctx, sqlc.AcceptFollowRequestParams{
			RequesterID: req.RequesterId.Int64(),
			TargetID:    req.UserId.Int64(),
		})

	} else {
		err = s.db.RejectFollowRequest(ctx, sqlc.RejectFollowRequestParams{
			RequesterID: req.RequesterId.Int64(),
			TargetID:    req.UserId.Int64(),
		})

	}
	if err != nil {
		return err
	}
	return nil
}

// returns ids of people a user follows for posts service so that the feed can be fetched
func (s *Application) GetFollowingIds(ctx context.Context, userId ct.Id) ([]int64, error) {
	if err := userId.Validate(); err != nil {
		return []int64{}, err
	}
	ids, err := s.db.GetFollowingIds(ctx, userId.Int64())
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// returns ten random users that people you follow follow, or are in your groups
// TODO for extra suggestions (have liked your public posts,
// or have commented on your public posts, or have liked posts you liked) need to ask posts service
func (s *Application) GetFollowSuggestions(ctx context.Context, userId ct.Id) ([]User, error) {
	if err := userId.Validate(); err != nil {
		return []User{}, err
	}
	rows, err := s.db.GetFollowSuggestions(ctx, userId.Int64())
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(rows))
	for _, r := range rows {
		users = append(users, User{
			UserId:   ct.Id(r.ID),
			Username: ct.Username(r.Username),
			Avatar:   r.Avatar,
		})
	}
	return users, nil
}

// NOT GRPC
func (s *Application) isFollowRequestPending(ctx context.Context, req FollowUserReq) (bool, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}
	isPending, err := s.db.IsFollowRequestPending(ctx, sqlc.IsFollowRequestPendingParams{
		RequesterID: req.FollowerId.Int64(),
		TargetID:    req.TargetUserId.Int64(),
	})
	if err != nil {
		return false, err
	}
	return isPending, nil
}

// SKIP GRPC FOR NOW
func (s *Application) IsFollowing(ctx context.Context, req FollowUserReq) (bool, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}
	isfollowing, err := s.db.IsFollowing(ctx, sqlc.IsFollowingParams{
		FollowerID:  req.FollowerId.Int64(),
		FollowingID: req.TargetUserId.Int64(),
	})
	if err != nil {
		return false, err
	}
	return isfollowing, nil
}

// SKIP GRPC FOR NOW
func (s *Application) IsFollowingEither(ctx context.Context, req FollowUserReq) (bool, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}

	atLeastOneIsFollowing, err := s.db.IsFollowingEither(ctx, sqlc.IsFollowingEitherParams{
		FollowerID:  req.FollowerId.Int64(),
		FollowingID: req.TargetUserId.Int64(),
	})
	if err != nil {
		return false, err
	}
	return atLeastOneIsFollowing, nil
}

// ---------------------------------------------------------------------
// low priority
// ---------------------------------------------------------------------
func GetMutualFollowers() {}

//get pending follow requests for user

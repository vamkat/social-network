package application

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
)

func (s *UserService) GetFollowersPaginated(ctx context.Context, req Pagination) ([]User, error) {

	//paginated, sorted by newest first
	rows, err := s.db.GetFollowers(ctx, sqlc.GetFollowersParams{
		FollowingID: req.UserId,
		Limit:       req.Limit,
		Offset:      req.Offset,
	})
	if err != nil {
		return []User{}, err
	}
	users := make([]User, 0, len(rows))
	for _, r := range rows {
		users = append(users, User{
			UserId:   r.ID,
			Username: r.Username,
			Avatar:   r.Avatar,
		})
	}

	return users, nil

}

func (s *UserService) GetFollowingPaginated(ctx context.Context, req Pagination) ([]User, error) {

	//paginated, sorted by newest first
	rows, err := s.db.GetFollowing(ctx, sqlc.GetFollowingParams{
		FollowerID: req.UserId,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})
	if err != nil {
		return []User{}, err
	}
	users := make([]User, 0, len(rows))
	for _, r := range rows {
		users = append(users, User{
			UserId:   r.ID,
			Username: r.Username,
			Avatar:   r.Avatar,
		})
	}

	return users, nil

}

// CHAT SERVICE EVENT should trigger event that creates conversation between two users for chat service (unless the target user was already following, so conversation exists)
func (s *UserService) FollowUser(ctx context.Context, req FollowUserReq) (resp FollowUserResp, err error) {
	status, err := s.db.FollowUser(ctx, sqlc.FollowUserParams{
		PFollower: req.FollowerId,
		PTarget:   req.TargetUserId,
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
func (s *UserService) UnFollowUser(ctx context.Context, req FollowUserReq) (viewerIsFolling bool, err error) {
	err = s.db.UnfollowUser(ctx, sqlc.UnfollowUserParams{
		FollowerID:  req.FollowerId,
		FollowingID: req.TargetUserId,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// CHAT SERVICE EVENT accepting a follow request should also trigger event to create conversation
func (s *UserService) HandleFollowRequest(ctx context.Context, req HandleFollowRequestReq) error {
	var err error
	if req.Accept {
		err = s.db.AcceptFollowRequest(ctx, sqlc.AcceptFollowRequestParams{
			RequesterID: req.RequesterId,
			TargetID:    req.UserId,
		})

	} else {
		err = s.db.RejectFollowRequest(ctx, sqlc.RejectFollowRequestParams{
			RequesterID: req.RequesterId,
			TargetID:    req.UserId,
		})

	}
	if err != nil {
		return err
	}
	return nil
}

// TODO returns ids of people a user follows so that the feed can be fetched
func (s *UserService) GetFollowingIds(ctx context.Context, userId int64) ([]int64, error) {
	return nil, nil
}

// NOT GRPC
func (s *UserService) isFollowRequestPending(ctx context.Context, req FollowUserReq) (bool, error) {
	isPending, err := s.db.IsFollowRequestPending(ctx, sqlc.IsFollowRequestPendingParams{
		RequesterID: req.FollowerId,
		TargetID:    req.TargetUserId,
	})
	if err != nil {
		return false, err
	}
	return isPending, nil
}

// SKIP GRPC FOR NOW
func (s *UserService) IsFollowing(ctx context.Context, req FollowUserReq) (bool, error) {
	isfollowing, err := s.db.IsFollowing(ctx, sqlc.IsFollowingParams{
		FollowerID:  req.FollowerId,
		FollowingID: req.TargetUserId,
	})
	if err != nil {
		return false, err
	}
	return isfollowing, nil
}

// SKIP GRPC FOR NOW
func (s *UserService) IsFollowingEither(ctx context.Context, req FollowUserReq) (bool, error) {

	atLeastOneIsFollowing, err := s.db.IsFollowingEither(ctx, sqlc.IsFollowingEitherParams{
		FollowerID:  req.FollowerId,
		FollowingID: req.TargetUserId,
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

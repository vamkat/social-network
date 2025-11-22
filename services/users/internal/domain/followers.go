package userservice

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
)

func (s *UserService) GetFollowersPaginated(ctx context.Context, req GetFollowersReq) ([]User, error) {
	userId, err := stringToUUID(req.FollowingID)
	if err != nil {
		return nil, err
	}

	//paginated, sorted by newest first
	rows, err := s.db.GetFollowers(ctx, sqlc.GetFollowersParams{
		Pub:    userId,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return []User{}, err
	}
	users := make([]User, 0, len(rows))
	for _, r := range rows {
		users = append(users, User{
			UserId:   r.PublicID.String(),
			Username: r.Username,
			Avatar:   r.Avatar,
		})
	}

	return users, nil

}

func (s *UserService) GetFollowingPaginated(ctx context.Context, req GetFollowingReq) ([]User, error) {
	userId, err := stringToUUID(req.FollowerID)
	if err != nil {
		return nil, err
	}

	//paginated, sorted by newest first
	rows, err := s.db.GetFollowing(ctx, sqlc.GetFollowingParams{
		Pub:    userId,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return []User{}, err
	}
	users := make([]User, 0, len(rows))
	for _, r := range rows {
		users = append(users, User{
			UserId:   r.PublicID.String(),
			Username: r.Username,
			Avatar:   r.Avatar,
		})
	}

	return users, nil

}

func (s *UserService) FollowUser(ctx context.Context, req FollowUserReq) (pending bool, err error) {
	followerId, err := stringToUUID(req.FollowerId)
	if err != nil {
		return false, err
	}
	followingId, err := stringToUUID(req.TargetUserId)
	if err != nil {
		return false, err
	}

	status, err := s.db.FollowUser(ctx, sqlc.FollowUserParams{
		Pub:   followerId,
		Pub_2: followingId,
	})
	if err != nil {
		return false, err
	}
	if status == "requested" { //I don't love hardcoding this, we'll see
		pending = true
	} else {
		pending = false
	}
	return pending, nil
}

func (s *UserService) UnFollowUser(ctx context.Context, req FollowUserReq) error {
	followerId, err := stringToUUID(req.FollowerId)
	if err != nil {
		return err
	}
	followingId, err := stringToUUID(req.TargetUserId)
	if err != nil {
		return err
	}

	err = s.db.UnfollowUser(ctx, sqlc.UnfollowUserParams{
		Pub:   followerId,
		Pub_2: followingId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) HandleFollowRequest(ctx context.Context, req HandleFollowRequestReq) error {
	requesterId, err := stringToUUID(req.RequesterId)
	if err != nil {
		return err
	}
	targetId, err := stringToUUID(req.UserId)
	if err != nil {
		return err
	}

	if req.Accept {
		err = s.db.AcceptFollowRequest(ctx, sqlc.AcceptFollowRequestParams{
			Pub:   requesterId,
			Pub_2: targetId,
		})

	} else {
		err = s.db.RejectFollowRequest(ctx, sqlc.RejectFollowRequestParams{
			Pub:   requesterId,
			Pub_2: targetId,
		})

	}
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) IsFollowing(ctx context.Context, req FollowUserReq) (bool, error) {
	followerId, err := stringToUUID(req.FollowerId)
	if err != nil {
		return false, err
	}
	followingId, err := stringToUUID(req.TargetUserId)
	if err != nil {
		return false, err
	}

	isfollowing, err := s.db.IsFollowing(ctx, sqlc.IsFollowingParams{
		Pub:   followerId,
		Pub_2: followingId,
	})
	if err != nil {
		return false, err
	}
	return isfollowing, nil
}

func (s *UserService) IsFollowingEither(ctx context.Context, req FollowUserReq) (bool, error) {
	followerId, err := stringToUUID(req.FollowerId)
	if err != nil {
		return false, err
	}
	followingId, err := stringToUUID(req.TargetUserId)
	if err != nil {
		return false, err
	}
	atLeastOneIsFollowing, err := s.db.IsFollowingEither(ctx, sqlc.IsFollowingEitherParams{
		Pub:   followerId,
		Pub_2: followingId,
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

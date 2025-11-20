package userservice

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
)

func (s *UserService) GetFollowersPaginated(ctx context.Context, req GetFollowersReq) ([]User, error) {
	//paginated, sorted by newest first
	rows, err := s.db.GetFollowers(ctx, sqlc.GetFollowersParams{
		FollowingID: req.FollowingID,
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
			Public:   r.ProfilePublic,
		})
	}

	return users, nil

}

func (s *UserService) GetFollowingPaginated(ctx context.Context, req GetFollowingReq) ([]User, error) {
	//paginated, sorted by newest first

	rows, err := s.db.GetFollowing(ctx, sqlc.GetFollowingParams{
		FollowerID: req.FollowerID,
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
			Public:   r.ProfilePublic,
		})
	}

	return users, nil

}

func FollowRequest() {

}

func HandleFollowRequest() {

}

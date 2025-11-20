package userservice

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
	"time"
)

func (s *UserService) GetBasicUserInfo(ctx context.Context, userId int64) (resp User, err error) {
	row, err := s.db.GetUserBasic(ctx, userId)
	if err != nil {
		return User{}, err
	}
	u := User{
		UserId:   row.ID,
		Username: row.Username,
		Avatar:   row.Avatar,
		Public:   row.ProfilePublic,
	}
	return u, nil

}

func (s *UserService) GetUserProfile(ctx context.Context, req UserProfileRequest) (UserProfileResponse, error) {
	var profile UserProfileResponse
	err := s.runTx(ctx, func(q *sqlc.Queries) error { //TODO consider not using a transaction here?
		// TODO helper: check if user has permission to see (public profile or isFollower)
		row, err := q.GetUserProfile(ctx, req.UserId)
		if err != nil {
			return err
		}

		dob := time.Time{}
		if row.DateOfBirth.Valid {
			dob = row.DateOfBirth.Time
		}

		profile = UserProfileResponse{
			UserId:      row.ID,
			Username:    row.Username,
			FirstName:   row.FirstName,
			LastName:    row.LastName,
			DateOfBirth: dob,
			Avatar:      row.Avatar,
			About:       row.AboutMe,
			Public:      row.ProfilePublic,
		}

		return nil
	})

	if err != nil {
		return UserProfileResponse{}, err
	}

	profile.FollowersCount, err = s.db.GetFollowerCount(ctx, profile.UserId)
	if err != nil {
		return UserProfileResponse{}, err
	}
	profile.FollowingCount, err = s.db.GetFollowingCount(ctx, profile.UserId)
	if err != nil {
		return UserProfileResponse{}, err
	}

	profile.Groups, err = s.GetUserGroupsPaginated(ctx, profile.UserId) //if pagination then it should be separate call and I should probably have a groups count
	if err != nil {
		return UserProfileResponse{}, err
	}

	return profile, nil

	// possibly usergroups also a different call
	// from forum service get all posts paginated (and number of posts)
	// and within all posts check each one if viewer has permission
}

func (s *UserService) SearchUsers(ctx context.Context, req UserSearchReq) ([]User, error) {

	rows, err := s.db.SearchUsers(ctx, sqlc.SearchUsersParams{
		Username: req.SearchTerm,
		Limit:    req.Limit,
	})

	if err != nil {
		return nil, err
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

func UpdateUserProfile() {
	//called with user_id and any of: username (TODO), first_name, last_name, date_of_birth, avatar, about_me
	//returns full profile
	//request needs to come from same user
	//---------------------------------------------------------------------

	//UpdateUserProflie
	//TODO check how to not update all fields but only changes (nil pointers?)
}

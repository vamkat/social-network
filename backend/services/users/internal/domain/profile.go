package userservice

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *UserService) GetBasicUserInfo(ctx context.Context, userId int64) (resp User, err error) {

	row, err := s.db.GetUserBasic(ctx, userId)
	if err != nil {
		return User{}, err
	}
	u := User{
		UserId:   userId,
		Username: row.Username,
		Avatar:   row.Avatar,
	}
	return u, nil

}

func (s *UserService) GetUserProfile(ctx context.Context, req UserProfileRequest) (UserProfileResponse, error) {

	var profile UserProfileResponse

	row, err := s.db.GetUserProfile(ctx, req.UserId)
	if err != nil {
		return UserProfileResponse{}, err
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
		CreatedAt:   row.CreatedAt.Time,
	}

	profile.ViewerIsFollowing, err = s.IsFollowing(ctx, FollowUserReq{
		FollowerId:   req.RequesterId,
		TargetUserId: req.UserId,
	})
	if err != nil {
		return UserProfileResponse{}, err
	}

	if req.RequesterId == req.UserId {
		profile.OwnProfile = true
	}

	if !profile.Public && !profile.ViewerIsFollowing && !profile.OwnProfile {
		return UserProfileResponse{}, ErrProfilePrivate
	}

	profile.FollowersCount, err = s.db.GetFollowerCount(ctx, req.UserId)
	if err != nil {
		return UserProfileResponse{}, err
	}
	profile.FollowingCount, err = s.db.GetFollowingCount(ctx, req.UserId)
	if err != nil {
		return UserProfileResponse{}, err
	}

	groupsRow, err := s.db.UserGroupCountsPerRole(ctx, req.UserId)
	if err != nil {
		return UserProfileResponse{}, err
	}
	profile.GroupsCount = groupsRow.TotalMemberships //owner and member, can change to member only
	profile.OwnedGroupsCount = groupsRow.OwnerCount

	return profile, nil

	// usergroups a different call
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
		})
	}

	return users, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, req UpdateProfileRequest) (UserProfileResponse, error) {
	//NOTE front needs to send everything, not just changed fields

	dob := pgtype.Date{
		Time:  req.DateOfBirth,
		Valid: true,
	}

	row, err := s.db.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
		ID:          req.UserId,
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DateOfBirth: dob,
		Avatar:      req.Avatar,
		AboutMe:     req.About,
	})
	if err != nil {
		return UserProfileResponse{}, err
	}

	newDob := time.Time{}
	if row.DateOfBirth.Valid {
		newDob = row.DateOfBirth.Time
	}

	profile := UserProfileResponse{
		UserId:      req.UserId,
		Username:    row.Username,
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		DateOfBirth: newDob,
		Avatar:      row.Avatar,
		About:       row.AboutMe,
		Public:      row.ProfilePublic,
	}

	return profile, nil

}

func (s *UserService) UpdateProfilePrivacy(ctx context.Context, req UpdateProfilePrivacyRequest) error {

	err := s.db.UpdateProfilePrivacy(ctx, sqlc.UpdateProfilePrivacyParams{
		ID:            req.UserId,
		ProfilePublic: req.Public,
	})
	if err != nil {
		return err
	}

	return nil
}

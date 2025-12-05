package application

import (
	"context"
	"fmt"
	"social-network/services/users/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// TODO add checks for post events (registration, updates, group titles and descriptions) - date:(valid format, over 13, not over 110),text fields:(length, special characters)

func (s *Application) GetBasicUserInfo(ctx context.Context, userId ct.Id) (resp User, err error) {
	if err := userId.Validate(); err != nil {
		return User{}, err
	}

	row, err := s.db.GetUserBasic(ctx, userId.Int64())
	if err != nil {
		return User{}, err
	}
	u := User{
		UserId:   ct.Id(userId),
		Username: ct.Username(row.Username),
		AvatarId: ct.Id(row.AvatarID),
	}
	return u, nil

}

func (s *Application) GetBatchBasicUserInfo(ctx context.Context, userIds ct.Ids) ([]User, error) {
	if err := userIds.Validate(); err != nil {
		return nil, err
	}

	rows, err := s.db.GetBatchUsersBasic(ctx, userIds.Int64())
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(rows))
	for _, r := range rows {
		users = append(users, User{
			UserId:   ct.Id(r.ID),
			Username: ct.Username(r.Username),
			AvatarId: ct.Id(r.AvatarID),
		})
	}
	return users, nil
}

func (s *Application) GetUserProfile(ctx context.Context, req UserProfileRequest) (UserProfileResponse, error) {
	var profile UserProfileResponse
	if err := ct.ValidateStruct(req); err != nil {
		return profile, err
	}

	row, err := s.db.GetUserProfile(ctx, req.UserId.Int64())
	if err != nil {
		fmt.Println(err)
		return UserProfileResponse{}, err
	}
	dob := time.Time{}
	if row.DateOfBirth.Valid {
		dob = row.DateOfBirth.Time
	}

	profile = UserProfileResponse{
		UserId:      ct.Id(row.ID),
		Username:    ct.Username(row.Username),
		FirstName:   ct.Name(row.FirstName),
		LastName:    ct.Name(row.LastName),
		DateOfBirth: ct.DateOfBirth(dob),
		AvatarId:    ct.Id(row.AvatarID),
		About:       ct.About(row.AboutMe),
		Public:      row.ProfilePublic,
		CreatedAt:   row.CreatedAt.Time,
	}

	followingParams := FollowUserReq{
		FollowerId:   req.RequesterId,
		TargetUserId: req.UserId,
	}

	profile.ViewerIsFollowing, err = s.IsFollowing(ctx, followingParams)
	if err != nil {
		return UserProfileResponse{}, err
	}

	profile.IsPending, err = s.isFollowRequestPending(ctx, followingParams)
	if err != nil {
		return UserProfileResponse{}, err
	}

	if req.RequesterId == req.UserId {
		profile.OwnProfile = true
	}

	if !profile.Public && !profile.ViewerIsFollowing && !profile.OwnProfile {
		return UserProfileResponse{}, ErrProfilePrivate
	}

	profile.FollowersCount, err = s.db.GetFollowerCount(ctx, req.UserId.Int64())
	if err != nil {
		return UserProfileResponse{}, err
	}
	profile.FollowingCount, err = s.db.GetFollowingCount(ctx, req.UserId.Int64())
	if err != nil {
		return UserProfileResponse{}, err
	}

	groupsRow, err := s.db.UserGroupCountsPerRole(ctx, req.UserId.Int64())
	if err != nil {
		return UserProfileResponse{}, err
	}

	profile.GroupsCount = groupsRow.TotalMemberships //owner and member, can change to member only
	profile.OwnedGroupsCount = groupsRow.OwnerCount

	return profile, nil

	// usergroups a different call
	// from posts service get all posts paginated (and number of posts)
	// and within all posts check each one if viewer has permission
}

func (s *Application) SearchUsers(ctx context.Context, req UserSearchReq) ([]User, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return []User{}, err
	}

	rows, err := s.db.SearchUsers(ctx, sqlc.SearchUsersParams{
		Username: req.SearchTerm.String(),
		Limit:    req.Limit.Int32(),
	})

	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(rows))
	for _, r := range rows {
		users = append(users, User{
			UserId:   ct.Id(r.ID),
			Username: ct.Username(r.Username),
			AvatarId: ct.Id(r.AvatarID),
		})
	}

	return users, nil
}

func (s *Application) UpdateUserProfile(ctx context.Context, req UpdateProfileRequest) (UserProfileResponse, error) {
	//NOTE front needs to send everything, not just changed fields

	if err := ct.ValidateStruct(req); err != nil {
		return UserProfileResponse{}, err
	}

	dob := pgtype.Date{
		Time:  req.DateOfBirth.Time(),
		Valid: true,
	}

	row, err := s.db.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
		ID:          req.UserId.Int64(),
		Username:    req.Username.String(),
		FirstName:   req.FirstName.String(),
		LastName:    req.LastName.String(),
		DateOfBirth: dob,
		AvatarID:    req.AvatarId.Int64(),
		AboutMe:     req.About.String(),
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
		Username:    ct.Username(row.Username),
		FirstName:   ct.Name(row.FirstName),
		LastName:    ct.Name(row.LastName),
		DateOfBirth: ct.DateOfBirth(newDob),
		AvatarId:    ct.Id(row.AvatarID),
		About:       ct.About(row.AboutMe),
		Public:      row.ProfilePublic,
	}

	return profile, nil

}

func (s *Application) UpdateProfilePrivacy(ctx context.Context, req UpdateProfilePrivacyRequest) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	err := s.db.UpdateProfilePrivacy(ctx, sqlc.UpdateProfilePrivacyParams{
		ID:            req.UserId.Int64(),
		ProfilePublic: req.Public,
	})
	if err != nil {
		return err
	}

	return nil
}

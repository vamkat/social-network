package application

import (
	"context"
	"fmt"
	ds "social-network/services/users/internal/db/dbservice"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// TODO add checks for post events (registration, updates, group titles and descriptions) - date:(valid format, over 13, not over 110),text fields:(length, special characters)

func (s *Application) GetBasicUserInfo(ctx context.Context, userId ct.Id) (resp models.User, err error) {
	if err := userId.Validate(); err != nil {
		return models.User{}, err
	}

	row, err := s.db.GetUserBasic(ctx, userId.Int64())
	if err != nil {
		return models.User{}, err
	}
	u := models.User{
		UserId:   ct.Id(userId),
		Username: ct.Username(row.Username),
		AvatarId: ct.Id(row.AvatarID),
	}
	return u, nil

}

// TO CONSIDER: who calls this and the above? Should fetching the url happen here or in retrieve users?
func (s *Application) GetBatchBasicUserInfo(ctx context.Context, userIds ct.Ids) ([]models.User, error) {
	if err := userIds.Validate(); err != nil {
		return nil, err
	}

	rows, err := s.db.GetBatchUsersBasic(ctx, userIds.Int64())
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0, len(rows))
	for _, r := range rows {
		users = append(users, models.User{
			UserId:   ct.Id(r.ID),
			Username: ct.Username(r.Username),
			AvatarId: ct.Id(r.AvatarID),
		})
	}
	return users, nil
}

func (s *Application) GetUserProfile(ctx context.Context, req models.UserProfileRequest) (models.UserProfileResponse, error) {
	var profile models.UserProfileResponse
	if err := ct.ValidateStruct(req); err != nil {
		return profile, err
	}

	row, err := s.db.GetUserProfile(ctx, req.UserId.Int64())
	if err != nil {
		fmt.Println(err)
		return models.UserProfileResponse{}, err
	}
	dob := time.Time{}
	if row.DateOfBirth.Valid {
		dob = row.DateOfBirth.Time
	}

	profile = models.UserProfileResponse{
		UserId:      ct.Id(row.ID),
		Username:    ct.Username(row.Username),
		FirstName:   ct.Name(row.FirstName),
		LastName:    ct.Name(row.LastName),
		DateOfBirth: ct.DateOfBirth(dob),
		AvatarId:    ct.Id(row.AvatarID),
		About:       ct.About(row.AboutMe),
		Public:      row.ProfilePublic,
		CreatedAt:   ct.GenDateTime(row.CreatedAt.Time),
		Email:       ct.Email(row.Email),
	}

	followingParams := models.FollowUserReq{
		FollowerId:   req.RequesterId,
		TargetUserId: req.UserId,
	}

	profile.ViewerIsFollowing, err = s.IsFollowing(ctx, followingParams)
	if err != nil {
		return models.UserProfileResponse{}, err
	}

	profile.IsPending, err = s.isFollowRequestPending(ctx, followingParams)
	if err != nil {
		return models.UserProfileResponse{}, err
	}

	if req.RequesterId == req.UserId {
		profile.OwnProfile = true
	}

	// if !profile.Public && !profile.ViewerIsFollowing && !profile.OwnProfile {
	// 	return models.UserProfileResponse{}, ErrProfilePrivate
	// }

	profile.FollowersCount, err = s.db.GetFollowerCount(ctx, req.UserId.Int64())
	if err != nil {
		return models.UserProfileResponse{}, err
	}
	profile.FollowingCount, err = s.db.GetFollowingCount(ctx, req.UserId.Int64())
	if err != nil {
		return models.UserProfileResponse{}, err
	}

	groupsRow, err := s.db.UserGroupCountsPerRole(ctx, req.UserId.Int64())
	if err != nil {
		return models.UserProfileResponse{}, err
	}

	profile.GroupsCount = groupsRow.TotalMemberships //owner and member, can change to member only
	profile.OwnedGroupsCount = groupsRow.OwnerCount

	if profile.AvatarId > 0 {
		imageUrl, err := s.clients.GetImage(ctx, profile.AvatarId.Int64())
		if err != nil {
			return models.UserProfileResponse{}, err
		}

		profile.AvatarURL = imageUrl
	}

	return profile, nil

	// usergroups a different call
	// from posts service get all posts paginated (and number of posts)
	// and within all posts check each one if viewer has permission
}

func (s *Application) SearchUsers(ctx context.Context, req models.UserSearchReq) ([]models.User, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return []models.User{}, err
	}

	rows, err := s.db.SearchUsers(ctx, ds.SearchUsersParams{
		Query: req.SearchTerm.String(),
		Limit: req.Limit.Int32(),
	})

	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0, len(rows))
	imageIds := make([]int64, 0, len(rows))
	for _, r := range rows {
		users = append(users, models.User{
			UserId:   ct.Id(r.ID),
			Username: ct.Username(r.Username),
			AvatarId: ct.Id(r.AvatarID),
		})
		if r.AvatarID > 0 {
			imageIds = append(imageIds, r.AvatarID)
		}
	}

	//get avatar urls
	if len(imageIds) > 0 {
		avatarMap, _, err := s.clients.GetImages(ctx, imageIds) //TODO delete failed
		if err != nil {
			return []models.User{}, err
		}
		for i := range users {
			users[i].AvatarURL = avatarMap[users[i].AvatarId.Int64()]
		}
	}

	return users, nil
}

func (s *Application) UpdateUserProfile(ctx context.Context, req models.UpdateProfileRequest) (models.UserProfileResponse, error) {
	//NOTE front needs to send everything, not just changed fields

	if err := ct.ValidateStruct(req); err != nil {
		return models.UserProfileResponse{}, err
	}

	dob := pgtype.Date{
		Time:  req.DateOfBirth.Time(),
		Valid: true,
	}

	row, err := s.db.UpdateUserProfile(ctx, ds.UpdateUserProfileParams{
		ID:          req.UserId.Int64(),
		Username:    req.Username.String(),
		FirstName:   req.FirstName.String(),
		LastName:    req.LastName.String(),
		DateOfBirth: dob,
		AvatarID:    req.AvatarId.Int64(),
		AboutMe:     req.About.String(),
	})
	if err != nil {
		return models.UserProfileResponse{}, err
	}

	//update redis basic user info
	basicUserInfo := models.User{
		UserId:   req.UserId,
		Username: req.Username,
		AvatarId: req.AvatarId,
	}
	_ = s.clients.SetObj(ctx,
		fmt.Sprintf("basic_user_info:%d", req.UserId),
		basicUserInfo,
		3*time.Minute,
	)

	newDob := time.Time{}
	if row.DateOfBirth.Valid {
		newDob = row.DateOfBirth.Time
	}

	profile := models.UserProfileResponse{
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

func (s *Application) UpdateProfilePrivacy(ctx context.Context, req models.UpdateProfilePrivacyRequest) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	err := s.db.UpdateProfilePrivacy(ctx, ds.UpdateProfilePrivacyParams{
		ID:            req.UserId.Int64(),
		ProfilePublic: req.Public,
	})
	if err != nil {
		return err
	}

	return nil
}

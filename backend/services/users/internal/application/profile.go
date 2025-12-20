package application

import (
	"context"
	"fmt"
	"social-network/services/users/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// TODO add checks for post events (registration, updates, group titles and descriptions) - date:(valid format, over 13, not over 110),text fields:(length, special characters)

func (app *Application) GetBasicUserInfo(ctx context.Context, userId ct.Id) (resp models.User, err error) {
	if err := userId.Validate(); err != nil {
		return models.User{}, err
	}

	row, err := app.db.Queries().GetUserBasic(ctx, userId.Int64())
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

func (app *Application) GetBatchBasicUserInfo(ctx context.Context, userIds ct.Ids) ([]models.User, error) {
	if err := userIds.Validate(); err != nil {
		return nil, err
	}

	rows, err := app.db.Queries().GetBatchUsersBasic(ctx, userIds.Int64())
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

func (app *Application) GetUserProfile(ctx context.Context, req models.UserProfileRequest) (models.UserProfileResponse, error) {
	var profile models.UserProfileResponse
	if err := ct.ValidateStruct(req); err != nil {
		return profile, err
	}

	row, err := app.db.Queries().GetUserProfile(ctx, req.UserId.Int64())
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

	profile.ViewerIsFollowing, err = app.IsFollowing(ctx, followingParams)
	if err != nil {
		return models.UserProfileResponse{}, err
	}

	profile.IsPending, err = app.isFollowRequestPending(ctx, followingParams)
	if err != nil {
		return models.UserProfileResponse{}, err
	}

	if req.RequesterId == req.UserId {
		profile.OwnProfile = true
	}

	// if !profile.Public && !profile.ViewerIsFollowing && !profile.OwnProfile {
	// 	return models.UserProfileResponse{}, ErrProfilePrivate
	// }

	profile.FollowersCount, err = app.db.Queries().GetFollowerCount(ctx, req.UserId.Int64())
	if err != nil {
		return models.UserProfileResponse{}, err
	}
	profile.FollowingCount, err = app.db.Queries().GetFollowingCount(ctx, req.UserId.Int64())
	if err != nil {
		return models.UserProfileResponse{}, err
	}

	groupsRow, err := app.db.Queries().UserGroupCountsPerRole(ctx, req.UserId.Int64())
	if err != nil {
		return models.UserProfileResponse{}, err
	}

	profile.GroupsCount = groupsRow.TotalMemberships //owner and member, can change to member only
	profile.OwnedGroupsCount = groupsRow.OwnerCount

	return profile, nil

	// usergroups a different call
	// from posts service get all posts paginated (and number of posts)
	// and within all posts check each one if viewer has permission
}

func (app *Application) SearchUsers(ctx context.Context, req models.UserSearchReq) ([]models.User, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return []models.User{}, err
	}

	rows, err := app.db.Queries().SearchUsers(ctx, sqlc.SearchUsersParams{
		Username: req.SearchTerm.String(),
		Limit:    req.Limit.Int32(),
	})

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

func (app *Application) UpdateUserProfile(ctx context.Context, req models.UpdateProfileRequest) (models.UserProfileResponse, error) {
	//NOTE front needs to send everything, not just changed fields

	if err := ct.ValidateStruct(req); err != nil {
		return models.UserProfileResponse{}, err
	}

	dob := pgtype.Date{
		Time:  req.DateOfBirth.Time(),
		Valid: true,
	}

	row, err := app.db.Queries().UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
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

func (app *Application) UpdateProfilePrivacy(ctx context.Context, req models.UpdateProfilePrivacyRequest) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	err := app.db.Queries().UpdateProfilePrivacy(ctx, sqlc.UpdateProfilePrivacyParams{
		ID:            req.UserId.Int64(),
		ProfilePublic: req.Public,
	})
	if err != nil {
		return err
	}

	return nil
}

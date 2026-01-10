package application

import (
	"context"
	"errors"
	"fmt"
	ds "social-network/services/users/internal/db/dbservice"
	"social-network/shared/gen-go/media"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Application) GetFollowersPaginated(ctx context.Context, req models.Pagination) ([]models.User, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return []models.User{}, ce.Wrap(ce.ErrInvalidArgument, err, input).WithPublic("invalid data received")
	}
	//paginated, sorted by newest first
	rows, err := s.db.GetFollowers(ctx, ds.GetFollowersParams{
		FollowingID: req.UserId.Int64(),
		Limit:       req.Limit.Int32(),
		Offset:      req.Offset.Int32(),
	})
	if err != nil {
		return []models.User{}, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	users := make([]models.User, 0, len(rows))
	var imageIds ct.Ids
	for _, r := range rows {
		users = append(users, models.User{
			UserId:   ct.Id(r.ID),
			Username: ct.Username(r.Username),
			AvatarId: ct.Id(r.AvatarID),
		})
		if r.AvatarID > 0 {
			imageIds = append(imageIds, ct.Id(r.AvatarID))
		}
	}
	//get avatar urls
	if len(imageIds) > 0 {
		avatarMap, _, err := s.mediaRetriever.GetImages(ctx, imageIds, media.FileVariant(1)) //TODO delete failed
		if err != nil {
			return []models.User{}, ce.Wrap(nil, err, input).WithPublic("error retrieving user avatars")
		}
		for i := range users {
			users[i].AvatarURL = avatarMap[users[i].AvatarId.Int64()]
		}
	}

	//tele.Debug(ctx, "get followers paginated: @1", "users", users)

	return users, nil

}

func (s *Application) GetFollowingPaginated(ctx context.Context, req models.Pagination) ([]models.User, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return []models.User{}, ce.Wrap(ce.ErrInvalidArgument, err, input).WithPublic("invalid data received")
	}

	//paginated, sorted by newest first
	rows, err := s.db.GetFollowing(ctx, ds.GetFollowingParams{
		FollowerID: req.UserId.Int64(),
		Limit:      req.Limit.Int32(),
		Offset:     req.Offset.Int32(),
	})
	if err != nil {
		return []models.User{}, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	users := make([]models.User, 0, len(rows))
	var imageIds ct.Ids
	for _, r := range rows {
		users = append(users, models.User{
			UserId:   ct.Id(r.ID),
			Username: ct.Username(r.Username),
			AvatarId: ct.Id(r.AvatarID),
		})
		if r.AvatarID > 0 {
			imageIds = append(imageIds, ct.Id(r.AvatarID))
		}
	}

	//get avatar urls
	if len(imageIds) > 0 {
		avatarMap, _, err := s.mediaRetriever.GetImages(ctx, imageIds, media.FileVariant(1)) //TODO delete failed
		if err != nil {
			return []models.User{}, ce.Wrap(nil, err, input).WithPublic("error retrieving user avatars")
		}
		for i := range users {
			users[i].AvatarURL = avatarMap[users[i].AvatarId.Int64()]
		}
	}

	return users, nil

}

func (s *Application) FollowUser(ctx context.Context, req models.FollowUserReq) (resp models.FollowUserResp, err error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return models.FollowUserResp{}, ce.Wrap(ce.ErrInvalidArgument, err, input).WithPublic("invalid data received")
	}
	status, err := s.db.FollowUser(ctx, ds.FollowUserParams{
		PFollower: req.FollowerId.Int64(),
		PTarget:   req.TargetUserId.Int64(),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "22023": // invalid_parameter_value
				return models.FollowUserResp{}, ce.New(ce.ErrInvalidArgument, err, input).WithPublic("user cannot follow self")
			case "P0002": // custom "not found"
				return models.FollowUserResp{}, ce.New(ce.ErrNotFound, err, input).WithPublic("user not found")
			}
		}
		return models.FollowUserResp{}, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	if status == "requested" { //Profile was private, request sent
		resp.IsPending = true
		resp.ViewerIsFollowing = false

		// create notification
		follower, err := s.GetBasicUserInfo(ctx, req.FollowerId)
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		err = s.clients.CreateFollowRequestNotification(ctx, req.TargetUserId.Int64(), req.FollowerId.Int64(), follower.Username.String())
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
	} else {
		resp.IsPending = false
		resp.ViewerIsFollowing = true

		// create notification
		follower, err := s.GetBasicUserInfo(ctx, req.FollowerId)
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		err = s.clients.CreateNewFollower(ctx, req.TargetUserId.Int64(), req.FollowerId.Int64(), follower.Username.String())
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		//try and create conversation in chat service if none exists
		//condition that exactly one user follows the other is checked before call is made
		// err = s.createPrivateConversation(ctx, req)
		// if err != nil {
		// 	tele.Info("conversation couldn't be created", err)
		// }
	}

	return resp, nil
}

func (s *Application) UnFollowUser(ctx context.Context, req models.FollowUserReq) (err error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, input).WithPublic("invalid data received")
	}
	//if already following, unfollows
	// if request pending, cancels request

	rowsAffected, err := s.db.UnfollowUser(ctx, ds.UnfollowUserParams{
		FollowerID:  req.FollowerId.Int64(),
		FollowingID: req.TargetUserId.Int64(),
	})
	if err != nil {
		return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	switch rowsAffected {
	case 0:
		// either idempotent success OR error
		tele.Warn(ctx, "no follow or follow request found with user @1 and target user @2", "userId", req.FollowerId, "targetUserId", req.TargetUserId)

	case 1:
		// Expected success

	default:
		// Inconsistent DB state
		// Log aggressively, but don't fail the user
		tele.Warn(ctx, "unexpected rows affected: expected @1, got @2", "expectedRows", 1, "affectedRows", rowsAffected)

	}
	// err = s.deletePrivateConversation(ctx, req)
	// if err != nil {
	// 	tele.Info("conversation couldn't be deleted", err)
	// }

	return nil
}

func (s *Application) HandleFollowRequest(ctx context.Context, req models.HandleFollowRequestReq) error {
	input := fmt.Sprintf("%#v", req)

	var err error
	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	if req.Accept {
		err = s.db.AcceptFollowRequest(ctx, ds.AcceptFollowRequestParams{
			RequesterID: req.RequesterId.Int64(),
			TargetID:    req.UserId.Int64(),
		})
		if err != nil {
			return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
		}

		//create notification
		targetUser, err := s.GetBasicUserInfo(ctx, req.UserId)
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		err = s.clients.CreateFollowRequestAccepted(ctx, req.RequesterId.Int64(), req.UserId.Int64(), targetUser.Username.String())
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		//try and create conversation in chat service if none exists
		//condition that exactly one user follows the other is checked before call is made
		// err = s.createPrivateConversation(ctx, models.FollowUserReq{
		// 	FollowerId:   req.RequesterId,
		// 	TargetUserId: req.UserId,
		// })
		// if err != nil {
		// 	tele.Info("conversation couldn't be created", err)
		// }

	} else {
		err = s.db.RejectFollowRequest(ctx, ds.RejectFollowRequestParams{
			RequesterID: req.RequesterId.Int64(),
			TargetID:    req.UserId.Int64(),
		})
		if err != nil {
			return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
		}
		//create notification
		targetUser, err := s.GetBasicUserInfo(ctx, req.UserId)
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		err = s.clients.CreateFollowRequestRejected(ctx, req.RequesterId.Int64(), req.UserId.Int64(), targetUser.Username.String())
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
	}
	return nil
}

// returns ids of people a user follows for posts service so that the feed can be fetched
func (s *Application) GetFollowingIds(ctx context.Context, userId ct.Id) ([]int64, error) {
	input := fmt.Sprintf("%#v", userId)

	if err := userId.Validate(); err != nil {
		return []int64{}, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	ids, err := s.db.GetFollowingIds(ctx, userId.Int64())
	if err != nil {
		return nil, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	return ids, nil
}

// returns ten random users that people you follow follow, or are in your groups
func (s *Application) GetFollowSuggestions(ctx context.Context, userId ct.Id) ([]models.User, error) {
	input := fmt.Sprintf("%#v", userId)

	if err := userId.Validate(); err != nil {
		return []models.User{}, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	rows, err := s.db.GetFollowSuggestions(ctx, userId.Int64())
	if err != nil {
		return nil, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}

	users := make([]models.User, 0, len(rows))
	var imageIds ct.Ids
	for _, r := range rows {
		users = append(users, models.User{
			UserId:   ct.Id(r.ID),
			Username: ct.Username(r.Username),
			AvatarId: ct.Id(r.AvatarID),
		})
		if r.AvatarID > 0 {
			imageIds = append(imageIds, ct.Id(r.AvatarID))
		}
	}

	//get avatar urls
	if len(imageIds) > 0 {
		avatarMap, _, err := s.mediaRetriever.GetImages(ctx, imageIds, media.FileVariant(1)) //TODO delete failed
		if err != nil {
			return []models.User{}, ce.Wrap(nil, err, input).WithPublic("error retrieving user avatars")
		}
		for i := range users {
			users[i].AvatarURL = avatarMap[users[i].AvatarId.Int64()]
		}
	}

	return users, nil
}

// NOT GRPC
func (s *Application) isFollowRequestPending(ctx context.Context, req models.FollowUserReq) (bool, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return false, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	isPending, err := s.db.IsFollowRequestPending(ctx, ds.IsFollowRequestPendingParams{
		RequesterID: req.FollowerId.Int64(),
		TargetID:    req.TargetUserId.Int64(),
	})
	if err != nil {
		return false, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	return isPending, nil
}

func (s *Application) IsFollowing(ctx context.Context, req models.FollowUserReq) (bool, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return false, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	isfollowing, err := s.db.IsFollowing(ctx, ds.IsFollowingParams{
		FollowerID:  req.FollowerId.Int64(),
		FollowingID: req.TargetUserId.Int64(),
	})
	if err != nil {
		return false, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	return isfollowing, nil
}

func (s *Application) AreFollowingEachOther(ctx context.Context, req models.FollowUserReq) (models.FollowRelationship, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return models.FollowRelationship{}, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	row, err := s.db.AreFollowingEachOther(ctx, ds.AreFollowingEachOtherParams{
		FollowerID:  req.FollowerId.Int64(),
		FollowingID: req.TargetUserId.Int64(),
	})
	if err != nil {
		return models.FollowRelationship{}, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}

	followRelationship := models.FollowRelationship{
		FollowerFollowsTarget: row.User1FollowsUser2,
		TargetFollowsFollower: row.User2FollowsUser1,
	}

	return followRelationship, nil //neither follows the other
}

// func (s *Application) createPrivateConversation(ctx context.Context, req models.FollowUserReq) error {
// 	atLeastOneIsFollowing, err := s.AreFollowingEachOther(ctx, req)
// 	if err != nil {
// 		return err
// 	}
// 	if atLeastOneIsFollowing != nil && !*atLeastOneIsFollowing { //I need exactly one follower, so false
// 		err := s.clients.CreatePrivateConversation(ctx, req.FollowerId.Int64(), req.TargetUserId.Int64())
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (s *Application) deletePrivateConversation(ctx context.Context, req models.FollowUserReq) error {
// 	atLeastOneIsFollowing, err := s.AreFollowingEachOther(ctx, req)
// 	if err != nil {
// 		return err
// 	}
// 	if atLeastOneIsFollowing == nil { //neither follows the other
// 		err := s.clients.DeleteConversationByExactMembers(ctx, []int64{req.FollowerId.Int64(), req.TargetUserId.Int64()})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// ---------------------------------------------------------------------
// low priority
// ---------------------------------------------------------------------
func GetMutualFollowers() {}

//get pending follow requests for user

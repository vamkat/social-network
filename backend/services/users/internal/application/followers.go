package application

import (
	"context"
	ds "social-network/services/users/internal/db/dbservice"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
)

func (s *Application) GetFollowersPaginated(ctx context.Context, req models.Pagination) ([]models.User, error) {
	//don't I need to check viewer has right to see? Or is this open to anyone?
	if err := ct.ValidateStruct(req); err != nil {
		return []models.User{}, err
	}
	//paginated, sorted by newest first
	rows, err := s.db.GetFollowers(ctx, ds.GetFollowersParams{
		FollowingID: req.UserId.Int64(),
		Limit:       req.Limit.Int32(),
		Offset:      req.Offset.Int32(),
	})
	if err != nil {
		return []models.User{}, err
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

	//fmt.Println(users)

	return users, nil

}

func (s *Application) GetFollowingPaginated(ctx context.Context, req models.Pagination) ([]models.User, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return []models.User{}, err
	}

	//paginated, sorted by newest first
	rows, err := s.db.GetFollowing(ctx, ds.GetFollowingParams{
		FollowerID: req.UserId.Int64(),
		Limit:      req.Limit.Int32(),
		Offset:     req.Offset.Int32(),
	})
	if err != nil {
		return []models.User{}, err
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

func (s *Application) FollowUser(ctx context.Context, req models.FollowUserReq) (resp models.FollowUserResp, err error) {
	if err := ct.ValidateStruct(req); err != nil {
		return models.FollowUserResp{}, err
	}
	status, err := s.db.FollowUser(ctx, ds.FollowUserParams{
		PFollower: req.FollowerId.Int64(),
		PTarget:   req.TargetUserId.Int64(),
	})
	if err != nil {
		return models.FollowUserResp{}, err
	}
	if status == "requested" { //Profile was private, request sent
		resp.IsPending = true
		resp.ViewerIsFollowing = false
		//TODO CREATE NOTIFICATION EVENT
	} else {
		resp.IsPending = false
		resp.ViewerIsFollowing = true

		//try and create conversation in chat service if none exists
		//condition that exactly one user follows the other is checked before call is made
		// err = s.createPrivateConversation(ctx, req)
		// if err != nil {
		// 	fmt.Println("conversation couldn't be created", err)
		// }
	}

	return resp, nil
}

func (s *Application) UnFollowUser(ctx context.Context, req models.FollowUserReq) (err error) {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	//if already following, unfollows
	// if request pending, cancels request

	rowsAffected, err := s.db.UnfollowUser(ctx, ds.UnfollowUserParams{
		FollowerID:  req.FollowerId.Int64(),
		FollowingID: req.TargetUserId.Int64(),
	})
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrNotFound
	}

	// err = s.deletePrivateConversation(ctx, req)
	// if err != nil {
	// 	fmt.Println("conversation couldn't be deleted", err)
	// }

	return nil
}

func (s *Application) HandleFollowRequest(ctx context.Context, req models.HandleFollowRequestReq) error {
	var err error
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	if req.Accept {
		err = s.db.AcceptFollowRequest(ctx, ds.AcceptFollowRequestParams{
			RequesterID: req.RequesterId.Int64(),
			TargetID:    req.UserId.Int64(),
		})
		if err != nil {
			return err
		}

		//try and create conversation in chat service if none exists
		//condition that exactly one user follows the other is checked before call is made
		// err = s.createPrivateConversation(ctx, models.FollowUserReq{
		// 	FollowerId:   req.RequesterId,
		// 	TargetUserId: req.UserId,
		// })
		// if err != nil {
		// 	fmt.Println("conversation couldn't be created", err)
		// }

	} else {
		err = s.db.RejectFollowRequest(ctx, ds.RejectFollowRequestParams{
			RequesterID: req.RequesterId.Int64(),
			TargetID:    req.UserId.Int64(),
		})
		if err != nil {
			return err
		}

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
func (s *Application) GetFollowSuggestions(ctx context.Context, userId ct.Id) ([]models.User, error) {
	if err := userId.Validate(); err != nil {
		return []models.User{}, err
	}
	rows, err := s.db.GetFollowSuggestions(ctx, userId.Int64())
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

// NOT GRPC
func (s *Application) isFollowRequestPending(ctx context.Context, req models.FollowUserReq) (bool, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}
	isPending, err := s.db.IsFollowRequestPending(ctx, ds.IsFollowRequestPendingParams{
		RequesterID: req.FollowerId.Int64(),
		TargetID:    req.TargetUserId.Int64(),
	})
	if err != nil {
		return false, err
	}
	return isPending, nil
}

func (s *Application) IsFollowing(ctx context.Context, req models.FollowUserReq) (bool, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}
	isfollowing, err := s.db.IsFollowing(ctx, ds.IsFollowingParams{
		FollowerID:  req.FollowerId.Int64(),
		FollowingID: req.TargetUserId.Int64(),
	})
	if err != nil {
		return false, err
	}
	return isfollowing, nil
}

// returns a bool pointer. Nil: neither is following the other, false: one is following the other, true: both are following each other
func (s *Application) AreFollowingEachOther(ctx context.Context, req models.FollowUserReq) (*bool, error) {
	var mutualFollow *bool // default: nil
	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}

	row, err := s.db.AreFollowingEachOther(ctx, ds.AreFollowingEachOtherParams{
		FollowerID:  req.FollowerId.Int64(),
		FollowingID: req.TargetUserId.Int64(),
	})
	if err != nil {
		return nil, err
	}
	if row.User1FollowsUser2 && row.User2FollowsUser1 { //both follow each other

		v := true
		mutualFollow = &v
		return mutualFollow, nil
	}

	if row.User1FollowsUser2 || row.User2FollowsUser1 { //one follows the other

		v := false
		mutualFollow = &v
		return mutualFollow, nil
	}
	return nil, nil //neither follows the other
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

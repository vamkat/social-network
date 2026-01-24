package retrieveusers

import (
	"context"
	"errors"
	"fmt"
	"time"

	cm "social-network/shared/gen-go/common"
	"social-network/shared/gen-go/media"
	userpb "social-network/shared/gen-go/users"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// GetUsers returns a map[userID]User, using cache + batch RPC.
func (h *UserRetriever) GetUsers(ctx context.Context, userIDs ct.Ids) (map[ct.Id]models.User, error) {
	input := fmt.Sprintf("user retriever: get users: uses ids: %v", userIDs)
	//========================== STEP 1 : get user info from users ===============================================
	if len(userIDs) == 0 {
		tele.Warn(ctx, "get users called with empty ids slice")
		return nil, nil
	}

	if err := userIDs.Validate(); err != nil {
		return nil, ce.New(ce.ErrInvalidArgument, err, input)
	}

	tele.Debug(ctx, "get users called with ids @1", "ids", userIDs)

	ids := userIDs.Unique()

	users := make(map[ct.Id]models.User, len(ids))

	if h.LocalCache != nil {
		for _, id := range ids {
			lcu, ok := h.LocalCache.Get(id)
			if ok && lcu != nil {
				users[id] = *lcu
				continue
			}
		}
	}

	var missing ct.Ids
	// Redis lookup
	for _, id := range ids {
		var u models.User

		// Check redis
		key, err := ct.BasicUserInfoKey{Id: id}.String()
		if err != nil {
			tele.Warn(ctx, "failed to construct redis key for id @1: @2", "userId", id, "error", err.Error())
			missing = append(missing, id)
			continue
		}

		if err := h.cache.GetObj(ctx, key, &u); err == nil {
			users[id] = u
			tele.Info(ctx, "found user on redis: @1", "user", u)
		} else {
			missing = append(missing, id)
		}
	}

	// Batch RPC for missing users
	if len(missing) > 0 {
		resp, err := h.client.GetBatchBasicUserInfo(ctx, &cm.UserIds{Values: missing.Int64()})
		if err != nil {
			return nil, ce.DecodeProto(err, input)
		}

		for _, u := range resp.Users {
			user := models.User{
				UserId:   ct.Id(u.UserId),
				Username: ct.Username(u.Username),
				AvatarId: ct.Id(u.Avatar),
			}
			users[user.UserId] = user

			key, err := ct.BasicUserInfoKey{Id: user.UserId}.String()
			if err == nil {
				_ = h.cache.SetObj(ctx,
					key,
					user,
					h.ttl,
				)
				tele.Debug(ctx, "user set on redis: @1 with key @2", "user", user, "key", key)
			} else {
				tele.Warn(ctx, "failed to construct redis key for user @1: @2", "userId", user.UserId, "error", err.Error())
			}
		}
	}

	//========================== STEP 2 : get avatars from media ===============================================
	// Get image urls for users
	var imageIds ct.Ids
	for _, user := range users {
		if user.AvatarId > 0 { //exclude 0 imageIds
			imageIds = append(imageIds, user.AvatarId)
		}
	}
	imageIds = imageIds.Unique()
	if len(imageIds) > 0 {
		// Use shared MediaRetriever for images (handles caching and fetching)
		imageMap, imagesToDelete, err := h.mediaRetriever.GetImages(ctx, imageIds, media.FileVariant_THUMBNAIL)
		if err != nil {
			tele.Error(ctx, "media retriever failed for @1", "request", imageIds, "error", err.Error()) //log error instead of returning
			//return nil, ce.Wrap(nil, err, input) // keep the code from retrieve media by wrapping the error and add errMsg for context
		} else {

			for id, u := range users {
				if url, ok := imageMap[u.AvatarId.Int64()]; ok {
					u.AvatarURL = url
					users[id] = u
				}
			}
		}

		if len(imagesToDelete) > 0 {
			msg := &userpb.FailedImageIds{
				ImgIds: imagesToDelete,
			}
			go func(m *userpb.FailedImageIds) {
				tele.Info(ctx, "removing avatar ids @1 from users", "failedImageIds", imagesToDelete)
				_, err := h.client.RemoveImages(context.WithoutCancel(ctx), m)
				if err != nil {
					tele.Warn(ctx, "failed  to delete failed images @1 from users: @2", "failedImageIds", imagesToDelete, "error", err.Error())
				}
			}(msg)
		}
	}

	if h.LocalCache != nil {
		for _, user := range users {
			h.LocalCache.SetWithTTL(user.UserId, &user, 1, time.Duration(10*time.Second))
		}
	}

	return users, nil
}

func (h *UserRetriever) GetUser(ctx context.Context, userId ct.Id) (models.User, error) {
	input := fmt.Sprintf("user retriever: get user: id: %v", userId)

	//========================== STEP 1 : get user info from users ===============================================

	if err := userId.Validate(); err != nil {
		return models.User{}, ce.New(ce.ErrInvalidArgument, err, input)
	}

	tele.Debug(ctx, "retrieve user called with user id @1", "userId", userId)

	// Local cache lookup
	if h.LocalCache != nil {
		u, ok := h.LocalCache.Get(userId)
		if ok && u != nil {
			tele.Debug(ctx, "found user on local cache", "user", *u)
			return *u, nil
		}
	}

	// Redis lookup
	key, err := ct.BasicUserInfoKey{Id: userId}.String()
	if err != nil {
		tele.Warn(ctx, "failed to construct redis key for id @1: @2", "userId", userId, "error", err.Error())
	}
	tele.Debug(ctx, "redis key constructed: @1", "redisKey", key)

	var user models.User
	if err := h.cache.GetObj(ctx, key, &user); err == nil {
		tele.Info(ctx, "found user on redis: @1 using key @2", "user", user, "key", key)
		return user, nil
	}
	resp, err := h.client.GetBasicUserInfo(ctx, wrapperspb.Int64(userId.Int64()))
	if err != nil {
		return models.User{}, ce.DecodeProto(err, input)
	}

	user = models.User{
		UserId:   ct.Id(resp.UserId),
		Username: ct.Username(resp.Username),
		AvatarId: ct.Id(resp.Avatar),
	}

	key, err = ct.BasicUserInfoKey{Id: user.UserId}.String()
	if err == nil {
		_ = h.cache.SetObj(ctx,
			key,
			user,
			h.ttl,
		)
		tele.Debug(ctx, "user set on redis: @1 with key @2", "user", user, "key", key)
	} else {
		tele.Warn(ctx, "failed to construct redis key for user @1: @1", "userId", user.UserId, "error", err.Error())
	}

	//========================== STEP 2 : get avatar from media ===============================================
	// Get image url for users

	if user.AvatarId > 0 { //exclude 0 imageIds

		// Use shared MediaRetriever for images (handles caching and fetching)
		imageUrl, err := h.mediaRetriever.GetImage(ctx, user.AvatarId.Int64(), media.FileVariant_THUMBNAIL)
		if err != nil {
			var commonError *ce.Error
			if errors.As(err, &commonError) {
				if err.(*ce.Error).IsClass(ce.ErrNotFound) {
					msg := &userpb.FailedImageIds{
						ImgIds: []int64{user.AvatarId.Int64()},
					}
					go func(m *userpb.FailedImageIds) {
						tele.Info(ctx, "removing avatar id @1 for user @2", "failedImageId", user.AvatarId, "userId", user.UserId)
						_, err := h.client.RemoveImages(context.WithoutCancel(ctx), m)
						if err != nil {
							tele.Warn(ctx, "failed to delete failed image @1 from users: @2", "failedImageId", user.AvatarId, "error", err.Error())
						}
					}(msg)
					return models.User{}, ce.Wrap(nil, err, input) // keep the code from retrieve media by wrapping the error and add errMsg for context
				}
			}
		}

		user.AvatarURL = imageUrl
	}

	if h.LocalCache != nil {
		h.LocalCache.SetWithTTL(userId, &user, 1, time.Duration(10*time.Second))
	}

	return user, nil
}

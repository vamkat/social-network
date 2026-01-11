package retrieveusers

import (
	"context"
	"fmt"

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

	ids := userIDs.Unique()

	users := make(map[ct.Id]models.User, len(ids))
	var missing ct.Ids

	// Redis lookup
	for _, id := range ids {
		var u models.User

		key, err := ct.BasicUserInfoKey{Id: id}.String()
		if err != nil {
			fmt.Printf("RETRIEVE USERS - failed to construct redis key for id %v: %v\n", id, err)
			missing = append(missing, id)
			continue
		}

		if err := h.cache.GetObj(ctx, key, &u); err == nil {
			users[id] = u
			fmt.Println("RETRIEVE USERS - found user on redis:", u)
		} else {
			missing = append(missing, id)
		}
	}

	// Batch RPC for missing users
	if len(missing) > 0 {
		resp, err := h.client.GetBatchBasicUserInfo(ctx, &cm.UserIds{Values: missing.Int64()})
		if err != nil {
			return nil, ce.ParseGrpcErr(err, input)
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
			} else {
				fmt.Printf("RETRIEVE USERS - failed to construct redis key for user %v: %v\n", user.UserId, err)
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
		imageMap, failedImageIds, err := h.mediaRetriever.GetImages(ctx, imageIds, media.FileVariant_THUMBNAIL)
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
		//==============pinpoint not found images============
		var imagesToDelete []int64
		//add if in failed
		imagesToDelete = append(imagesToDelete, failedImageIds...)

		for _, imageId := range imageIds {
			id := imageId.Int64()

			// skip if download succeeded
			if _, exists := imageMap[id]; exists {
				continue
			} else {

				imagesToDelete = append(imagesToDelete, id) //now imagesToDelete includes failed and not found
			}
		}

		msg := &userpb.FailedImageIds{
			ImgIds: imagesToDelete,
		}
		go func(m *userpb.FailedImageIds) {
			_, err := h.client.RemoveImages(ctx, m)
			if err != nil {
				tele.Warn(context.WithoutCancel(ctx), "failed  to delete failed images @1 from users", "failedImageIds", failedImageIds)
			}
		}(msg)
	}

	return users, nil
}

func (h *UserRetriever) GetUser(ctx context.Context, userID ct.Id) (models.User, error) {
	input := fmt.Sprintf("user retriever: get user: id: %v", userID)

	//========================== STEP 1 : get user info from users ===============================================

	// Redis lookup
	var u models.User

	key, err := ct.BasicUserInfoKey{Id: userID}.String()
	if err != nil {
		fmt.Printf("RETRIEVE USERS - failed to construct redis key for id %v: %v\n", userID, err)
	}

	var user models.User
	if err := h.cache.GetObj(ctx, key, &user); err == nil {
		fmt.Println("RETRIEVE USERS - found user on redis:", u)
		return user, nil
	}
	resp, err := h.client.GetBasicUserInfo(ctx, wrapperspb.Int64(userID.Int64()))
	if err != nil {
		return models.User{}, ce.ParseGrpcErr(err, input)
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
	} else {
		fmt.Printf("RETRIEVE USERS - failed to construct redis key for user %v: %v\n", user.UserId, err)
	}

	//========================== STEP 2 : get avatars from media ===============================================
	// Get image urls for users

	if user.AvatarId > 0 { //exclude 0 imageIds

		// Use shared MediaRetriever for images (handles caching and fetching)
		imageUrl, err := h.mediaRetriever.GetImage(ctx, user.AvatarId.Int64(), media.FileVariant_THUMBNAIL)
		if err != nil {
			return models.User{}, ce.Wrap(nil, err, input) // keep the code from retrieve media by wrapping the error and add errMsg for context
		}

		u.AvatarURL = imageUrl
	}

	return user, nil
}

package application

import (
	"context"
	ds "social-network/services/posts/internal/db/dbservice"
	"social-network/shared/gen-go/media"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Application) GetPersonalizedFeed(ctx context.Context, req models.GetPersonalizedFeedReq) ([]models.Post, error) {

	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}

	idsRequesterFollows, err := s.clients.GetFollowingIds(ctx, req.RequesterId.Int64())
	if err != nil {
		return nil, err
	}

	rows, err := s.db.GetPersonalizedFeed(ctx, ds.GetPersonalizedFeedParams{
		UserID:  req.RequesterId.Int64(),
		Column2: idsRequesterFollows,
		Offset:  req.Offset.Int32(),
		Limit:   req.Limit.Int32(),
	})
	if err != nil {
		return nil, err
	}
	posts := make([]models.Post, 0, len(rows))
	userIDs := make(ct.Ids, 0, len(rows))
	PostImageIds := make(ct.Ids, 0, len(rows))

	for _, r := range rows {
		uid := r.CreatorID
		userIDs = append(userIDs, ct.Id(uid))

		posts = append(posts, models.Post{

			PostId: ct.Id(r.ID),
			Body:   ct.PostBody(r.PostBody),
			User: models.User{
				UserId: ct.Id(uid),
			},
			CommentsCount:   int(r.CommentsCount),
			ReactionsCount:  int(r.ReactionsCount),
			LastCommentedAt: ct.GenDateTime(r.LastCommentedAt.Time),
			CreatedAt:       ct.GenDateTime(r.CreatedAt.Time),
			UpdatedAt:       ct.GenDateTime(r.UpdatedAt.Time),
			LikedByUser:     r.LikedByUser,
			ImageId:         ct.Id(r.Image),
		})

		if r.Image > 0 {
			PostImageIds = append(PostImageIds, ct.Id(r.Image))
		}

	}

	if len(posts) == 0 {
		return posts, nil
	}

	userMap, err := s.userRetriever.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	var imageMap map[int64]string
	if len(PostImageIds) > 0 {
		imageMap, _, err = s.clients.GetImages(ctx, PostImageIds, media.FileVariant_MEDIUM)
	}

	for i := range posts {
		uid := posts[i].User.UserId
		if u, ok := userMap[uid]; ok {
			posts[i].User = u
		}
		posts[i].ImageUrl = imageMap[posts[i].ImageId.Int64()]
	}

	return posts, nil
}

func (s *Application) GetPublicFeed(ctx context.Context, req models.GenericPaginatedReq) ([]models.Post, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}
	rows, err := s.db.GetPublicFeed(ctx, ds.GetPublicFeedParams{
		UserID: req.RequesterId.Int64(),
		Offset: req.Offset.Int32(),
		Limit:  req.Limit.Int32(),
	})
	if err != nil {
		return nil, err
	}

	posts := make([]models.Post, 0, len(rows))
	userIDs := make(ct.Ids, 0, len(rows))
	postImageIds := make(ct.Ids, 0, len(rows))

	for _, r := range rows {
		uid := r.CreatorID
		userIDs = append(userIDs, ct.Id(uid))

		posts = append(posts, models.Post{
			PostId: ct.Id(r.ID),
			Body:   ct.PostBody(r.PostBody),
			User: models.User{
				UserId: ct.Id(uid),
			},
			CommentsCount:   int(r.CommentsCount),
			ReactionsCount:  int(r.ReactionsCount),
			LastCommentedAt: ct.GenDateTime(r.LastCommentedAt.Time),
			CreatedAt:       ct.GenDateTime(r.CreatedAt.Time),
			UpdatedAt:       ct.GenDateTime(r.UpdatedAt.Time),
			LikedByUser:     r.LikedByUser,
			ImageId:         ct.Id(r.Image),
		})
		if r.Image > 0 {
			postImageIds = append(postImageIds, ct.Id(r.Image))
		}

	}
	if len(posts) == 0 {
		tele.Warn(ctx, "GetPublicFeed: no posts in feed")
		return posts, nil
	}

	userMap, err := s.userRetriever.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	var imageMap map[int64]string
	tele.Info(ctx, "GetPublicFeed needs these images", "image ids", postImageIds)
	if len(postImageIds) > 0 {
		imageMap, _, err = s.mediaRetriever.GetImages(ctx, postImageIds, media.FileVariant_MEDIUM)
	}

	for i := range posts {
		uid := posts[i].User.UserId
		if u, ok := userMap[uid]; ok {
			posts[i].User = u
		}
		posts[i].ImageUrl = imageMap[posts[i].ImageId.Int64()]
	}

	return posts, nil
}

func (s *Application) GetUserPostsPaginated(ctx context.Context, req models.GetUserPostsReq) ([]models.Post, error) {

	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}

	isFollowing, err := s.clients.IsFollowing(ctx, req.RequesterId.Int64(), int64(req.CreatorId))
	if err != nil {
		return nil, err
	}

	rows, err := s.db.GetUserPostsPaginated(ctx, ds.GetUserPostsPaginatedParams{
		CreatorID: req.CreatorId.Int64(),
		UserID:    req.RequesterId.Int64(),
		Column3:   isFollowing,
		Limit:     req.Limit.Int32(),
		Offset:    req.Offset.Int32(),
	})
	if err != nil {
		return nil, err
	}

	posts := make([]models.Post, 0, len(rows))
	userIDs := make(ct.Ids, 0, len(rows))
	PostImageIds := make(ct.Ids, 0, len(rows))

	for _, r := range rows {
		uid := r.CreatorID
		userIDs = append(userIDs, ct.Id(uid))

		posts = append(posts, models.Post{
			PostId: ct.Id(r.ID),
			Body:   ct.PostBody(r.PostBody),
			User: models.User{
				UserId: ct.Id(uid),
			},
			CommentsCount:   int(r.CommentsCount),
			ReactionsCount:  int(r.ReactionsCount),
			LastCommentedAt: ct.GenDateTime(r.LastCommentedAt.Time),
			CreatedAt:       ct.GenDateTime(r.CreatedAt.Time),
			UpdatedAt:       ct.GenDateTime(r.UpdatedAt.Time),
			LikedByUser:     r.LikedByUser,
			ImageId:         ct.Id(r.Image),
		})
		if r.Image > 0 {
			PostImageIds = append(PostImageIds, ct.Id(r.Image))
		}

	}

	if len(posts) == 0 {
		return posts, nil
	}

	userMap, err := s.userRetriever.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	var imageMap map[int64]string
	if len(PostImageIds) > 0 {
		imageMap, _, err = s.clients.GetImages(ctx, PostImageIds, media.FileVariant_MEDIUM)
	}

	for i := range posts {
		uid := posts[i].User.UserId
		if u, ok := userMap[uid]; ok {
			posts[i].User = u
		}
		posts[i].ImageUrl = imageMap[posts[i].ImageId.Int64()]
	}

	return posts, nil
}

func (s *Application) GetGroupPostsPaginated(ctx context.Context, req models.GetGroupPostsReq) ([]models.Post, error) {

	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}

	var groupId pgtype.Int8
	groupId.Int64 = req.GroupId.Int64()
	if req.GroupId == 0 {
		return nil, ErrNoGroupIdGiven
	}
	groupId.Valid = true

	isMember, err := s.clients.IsGroupMember(ctx, req.RequesterId.Int64(), req.GroupId.Int64())
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotAllowed
	}

	rows, err := s.db.GetGroupPostsPaginated(ctx, ds.GetGroupPostsPaginatedParams{
		GroupID: groupId,
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrNotFound
	}
	posts := make([]models.Post, 0, len(rows))
	userIDs := make(ct.Ids, 0, len(rows))
	PostImageIds := make(ct.Ids, 0, len(rows))

	for _, r := range rows {
		uid := r.CreatorID
		userIDs = append(userIDs, ct.Id(uid))

		posts = append(posts, models.Post{
			PostId: ct.Id(r.ID),
			Body:   ct.PostBody(r.PostBody),
			User: models.User{
				UserId: ct.Id(uid),
			},
			GroupId:         req.GroupId,
			Audience:        ct.Audience(r.Audience),
			CommentsCount:   int(r.CommentsCount),
			ReactionsCount:  int(r.ReactionsCount),
			LastCommentedAt: ct.GenDateTime(r.LastCommentedAt.Time),
			CreatedAt:       ct.GenDateTime(r.CreatedAt.Time),
			UpdatedAt:       ct.GenDateTime(r.UpdatedAt.Time),
			LikedByUser:     r.LikedByUser,
			ImageId:         ct.Id(r.Image),
		})

		if r.Image > 0 {
			PostImageIds = append(PostImageIds, ct.Id(r.Image))
		}
	}
	if len(posts) == 0 {
		return posts, nil
	}

	userMap, err := s.userRetriever.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	var imageMap map[int64]string
	if len(PostImageIds) > 0 {
		imageMap, _, err = s.clients.GetImages(ctx, PostImageIds, media.FileVariant_MEDIUM)
	}

	for i := range posts {
		uid := posts[i].User.UserId
		if u, ok := userMap[uid]; ok {
			posts[i].User = u
		}
		posts[i].ImageUrl = imageMap[posts[i].ImageId.Int64()]
	}

	return posts, nil
}

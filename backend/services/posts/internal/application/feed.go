package application

import (
	"context"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"

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

	rows, err := s.db.GetPersonalizedFeed(ctx, sqlc.GetPersonalizedFeedParams{
		UserID:  req.RequesterId.Int64(),
		Column2: idsRequesterFollows,
		Offset:  req.Offset.Int32(),
		Limit:   req.Limit.Int32(),
	})
	if err != nil {
		return nil, err
	}
	posts := make([]models.Post, 0, len(rows))
	for _, r := range rows {
		posts = append(posts, models.Post{
			PostId: ct.Id(r.ID),
			Body:   ct.PostBody(r.PostBody),
			User: models.User{
				UserId: ct.Id(r.CreatorID),
			},
			CommentsCount:   int(r.CommentsCount),
			ReactionsCount:  int(r.ReactionsCount),
			LastCommentedAt: r.LastCommentedAt.Time,
			CreatedAt:       r.CreatedAt.Time,
			UpdatedAt:       r.UpdatedAt.Time,
			LikedByUser:     r.LikedByUser,
			Image:           ct.Id(r.Image),
		})
	}
	if err := s.hydratePosts(ctx, posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Application) GetPublicFeed(ctx context.Context, req models.EntityIdPaginatedReq) ([]models.Post, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}
	rows, err := s.db.GetPublicFeed(ctx, sqlc.GetPublicFeedParams{
		UserID: req.RequesterId.Int64(),
		Offset: req.Offset.Int32(),
		Limit:  req.Limit.Int32(),
	})
	if err != nil {
		return nil, err
	}
	posts := make([]models.Post, 0, len(rows))
	for _, r := range rows {
		posts = append(posts, models.Post{
			PostId: ct.Id(r.ID),
			Body:   ct.PostBody(r.PostBody),
			User: models.User{
				UserId: ct.Id(r.CreatorID),
			},
			CommentsCount:   int(r.CommentsCount),
			ReactionsCount:  int(r.ReactionsCount),
			LastCommentedAt: r.LastCommentedAt.Time,
			CreatedAt:       r.CreatedAt.Time,
			UpdatedAt:       r.UpdatedAt.Time,
			LikedByUser:     r.LikedByUser,
			Image:           ct.Id(r.Image),
		})
	}

	if err := s.hydratePosts(ctx, posts); err != nil {
		return nil, err
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

	rows, err := s.db.GetUserPostsPaginated(ctx, sqlc.GetUserPostsPaginatedParams{
		CreatorID: req.CreatorId.Int64(),
		UserID:    req.RequesterId.Int64(),
		Column3:   isFollowing,
		Limit:     req.Limit.Int32(),
		Offset:    req.Offset.Int32(),
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrNotFound
	}
	posts := make([]models.Post, 0, len(rows))
	for _, r := range rows {
		posts = append(posts, models.Post{
			PostId: ct.Id(r.ID),
			Body:   ct.PostBody(r.PostBody),
			User: models.User{
				UserId: ct.Id(r.CreatorID),
			},
			CommentsCount:   int(r.CommentsCount),
			ReactionsCount:  int(r.ReactionsCount),
			LastCommentedAt: r.LastCommentedAt.Time,
			CreatedAt:       r.CreatedAt.Time,
			UpdatedAt:       r.UpdatedAt.Time,
			LikedByUser:     r.LikedByUser,
			Image:           ct.Id(r.Image),
		})

	}
	if err := s.hydratePosts(ctx, posts); err != nil {
		return nil, err
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

	rows, err := s.db.GetGroupPostsPaginated(ctx, sqlc.GetGroupPostsPaginatedParams{
		GroupID: groupId,
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrNotFound
	}
	posts := make([]models.Post, 0, len(rows))
	for _, r := range rows {
		posts = append(posts, models.Post{
			PostId: ct.Id(r.ID),
			Body:   ct.PostBody(r.PostBody),
			User: models.User{
				UserId: ct.Id(r.CreatorID),
			},
			GroupId:         req.GroupId,
			Audience:        ct.Audience(r.Audience),
			CommentsCount:   int(r.CommentsCount),
			ReactionsCount:  int(r.ReactionsCount),
			LastCommentedAt: r.LastCommentedAt.Time,
			CreatedAt:       r.CreatedAt.Time,
			UpdatedAt:       r.UpdatedAt.Time,
			LikedByUser:     r.LikedByUser,
			Image:           ct.Id(r.Image),
		})
	}

	if err := s.hydratePosts(ctx, posts); err != nil {
		return nil, err
	}

	return posts, nil
}

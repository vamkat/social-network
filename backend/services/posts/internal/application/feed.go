package application

import (
	"context"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
)

func (s *Application) GetPersonalizedFeed(ctx context.Context, req GetPersonalizedFeedReq) ([]Post, error) {
	//HANDLER needs to provide list of ids the requester follows
	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}
	rows, err := s.db.GetPersonalizedFeed(ctx, sqlc.GetPersonalizedFeedParams{
		UserID:  req.RequesterId.Int64(),
		Column2: req.RequesterFollowsIds.Int64(),
		Offset:  req.Offset.Int32(),
		Limit:   req.Limit.Int32(),
	})
	if err != nil {
		return nil, err
	}
	posts := make([]Post, 0, len(rows))
	for _, r := range rows {
		posts = append(posts, Post{
			PostId:          ct.Id(r.ID),
			Body:            ct.PostBody(r.PostBody),
			CreatorId:       ct.Id(r.CreatorID),
			CommentsCount:   int(r.CommentsCount),
			ReactionsCount:  int(r.ReactionsCount),
			LastCommentedAt: r.LastCommentedAt.Time,
			CreatedAt:       r.CreatedAt.Time,
			UpdatedAt:       r.UpdatedAt.Time,
			LikedByUser:     r.LikedByUser,
			Image:           ct.Id(r.Image),
			LatestComment: Comment{
				CommentId:      ct.Id(r.LatestCommentID),
				ParentId:       ct.Id(r.ID),
				Body:           ct.CommentBody(r.LatestCommentBody),
				CreatorId:      ct.Id(r.LatestCommentCreatorID),
				ReactionsCount: int(r.LatestCommentReactionsCount),
				CreatedAt:      r.LatestCommentCreatedAt.Time,
				UpdatedAt:      r.LatestCommentUpdatedAt.Time,
				LikedByUser:    r.LatestCommentLikedByUser,
				Image:          ct.Id(r.LatestCommentImage),
			},
		})
	}
	return posts, nil
}

func (s *Application) GetPublicFeed(ctx context.Context, req GenericPaginatedReq) ([]Post, error) {
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
	posts := make([]Post, 0, len(rows))
	for _, r := range rows {
		posts = append(posts, Post{
			PostId:          ct.Id(r.ID),
			Body:            ct.PostBody(r.PostBody),
			CreatorId:       ct.Id(r.CreatorID),
			CommentsCount:   int(r.CommentsCount),
			ReactionsCount:  int(r.ReactionsCount),
			LastCommentedAt: r.LastCommentedAt.Time,
			CreatedAt:       r.CreatedAt.Time,
			UpdatedAt:       r.UpdatedAt.Time,
			LikedByUser:     r.LikedByUser,
			Image:           ct.Id(r.Image),
			LatestComment: Comment{
				CommentId:      ct.Id(r.LatestCommentID),
				ParentId:       ct.Id(r.ID),
				Body:           ct.CommentBody(r.LatestCommentBody),
				CreatorId:      ct.Id(r.LatestCommentCreatorID),
				ReactionsCount: int(r.LatestCommentReactionsCount),
				CreatedAt:      r.LatestCommentCreatedAt.Time,
				UpdatedAt:      r.LatestCommentUpdatedAt.Time,
				LikedByUser:    r.LatestCommentLikedByUser,
				Image:          ct.Id(r.LatestCommentImage),
			},
		})
	}
	return posts, nil
}

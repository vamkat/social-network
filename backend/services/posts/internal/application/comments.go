package application

import (
	"context"
	ds "social-network/services/posts/internal/db/dbservice"
	"social-network/shared/gen-go/media"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"
)

func (s *Application) CreateComment(ctx context.Context, req models.CreateCommentReq) (err error) {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	accessCtx := accessContext{
		requesterId: req.CreatorId.Int64(),
		entityId:    req.ParentId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrNotAllowed
	}
	return s.txRunner.RunTx(ctx, func(q *ds.Queries) error {
		commentId, err := q.CreateComment(ctx, ds.CreateCommentParams{
			CommentCreatorID: req.CreatorId.Int64(),
			ParentID:         req.ParentId.Int64(),
			CommentBody:      req.Body.String(),
		})

		if err != nil {
			return err
		}

		if req.ImageId != 0 {
			err = q.UpsertImage(ctx, ds.UpsertImageParams{
				ID:       req.ImageId.Int64(),
				ParentID: commentId,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

}

func (s *Application) EditComment(ctx context.Context, req models.EditCommentReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	// accessCtx := accessContext{
	// 	requesterId: req.CreatorId.Int64(),
	// 	entityId:    req.CommentId.Int64(),
	// }

	// hasAccess, err := s.hasRightToView(ctx, accessCtx)
	// if err != nil {
	// 	return err
	// }
	// if !hasAccess {
	// 	return ErrNotAllowed
	// }

	return s.txRunner.RunTx(ctx, func(q *ds.Queries) error {
		rowsAffected, err := q.EditComment(ctx, ds.EditCommentParams{
			CommentBody:      req.Body.String(),
			ID:               req.CommentId.Int64(),
			CommentCreatorID: req.CreatorId.Int64(),
		})
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			return ErrNotFound
		}
		if req.ImageId > 0 {
			err := q.UpsertImage(ctx, ds.UpsertImageParams{
				ID:       req.ImageId.Int64(),
				ParentID: req.CommentId.Int64(),
			})
			if err != nil {
				return err
			}
		}
		if req.DeleteImage {
			rowsAffected, err := q.DeleteImage(ctx, req.CommentId.Int64())
			if err != nil {
				return err
			}
			if rowsAffected != 1 {
				tele.Warn(ctx, "EditComment: image to be deleted not found", "request", req)
			}
		}
		return nil
	})

}

func (s *Application) DeleteComment(ctx context.Context, req models.GenericReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	// accessCtx := accessContext{
	// 	requesterId: req.RequesterId.Int64(),
	// 	entityId:    req.EntityId.Int64(),
	// }

	// hasAccess, err := s.hasRightToView(ctx, accessCtx)
	// if err != nil {
	// 	return err
	// }
	// if !hasAccess {
	// 	return ErrNotAllowed
	// }

	rowsAffected, err := s.db.DeleteComment(ctx, ds.DeleteCommentParams{
		ID:               req.EntityId.Int64(),
		CommentCreatorID: req.RequesterId.Int64(),
	})
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrNotFound
	}
	return nil
}

func (s *Application) GetCommentsByParentId(ctx context.Context, req models.EntityIdPaginatedReq) ([]models.Comment, error) {

	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, ErrNotAllowed
	}

	rows, err := s.db.GetCommentsByPostId(ctx, ds.GetCommentsByPostIdParams{
		ParentID: req.EntityId.Int64(),
		UserID:   req.RequesterId.Int64(),
		Limit:    req.Limit.Int32(),
		Offset:   req.Offset.Int32(),
	})
	if err != nil {
		return nil, err
	}
	comments := make([]models.Comment, 0, len(rows))
	userIDs := make(ct.Ids, 0, len(rows))
	CommentImageIds := make(ct.Ids, 0, len(rows))

	for _, r := range rows {
		uid := ct.Id(r.CommentCreatorID)
		userIDs = append(userIDs, uid)

		comments = append(comments, models.Comment{
			CommentId: ct.Id(r.ID),
			ParentId:  req.EntityId,
			Body:      ct.CommentBody(r.CommentBody),
			User: models.User{
				UserId: ct.Id(r.CommentCreatorID),
			},
			ReactionsCount: int(r.ReactionsCount),
			CreatedAt:      ct.GenDateTime(r.CreatedAt.Time),
			UpdatedAt:      ct.GenDateTime(r.UpdatedAt.Time),
			LikedByUser:    r.LikedByUser,
			ImageId:        ct.Id(r.Image),
		})
		if r.Image > 0 {
			CommentImageIds = append(CommentImageIds, ct.Id(r.Image))
		}
	}

	if len(comments) == 0 {
		return comments, nil
	}

	userMap, err := s.userRetriever.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	var imageMap map[int64]string
	if len(CommentImageIds) > 0 {
		imageMap, _, err = s.mediaRetriever.GetImages(ctx, CommentImageIds, media.FileVariant_MEDIUM)
	}

	for i := range comments {
		uid := comments[i].User.UserId
		if u, ok := userMap[uid]; ok {
			comments[i].User = u
		}
		comments[i].ImageUrl = imageMap[comments[i].ImageId.Int64()]
	}

	return comments, nil
}

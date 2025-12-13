package application

import (
	"context"
	"fmt"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
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
	return s.txRunner.RunTx(ctx, func(q sqlc.Querier) error {
		err = q.CreateComment(ctx, sqlc.CreateCommentParams{
			CommentCreatorID: req.CreatorId.Int64(),
			ParentID:         req.ParentId.Int64(),
			CommentBody:      req.Body.String(),
		})

		if err != nil {
			return err
		}

		if req.Image != 0 {
			err = q.UpsertImage(ctx, sqlc.UpsertImageParams{
				ID:       req.Image.Int64(),
				ParentID: req.ParentId.Int64(),
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

	accessCtx := accessContext{
		requesterId: req.CreatorId.Int64(),
		entityId:    req.CommentId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrNotAllowed
	}

	return s.txRunner.RunTx(ctx, func(q sqlc.Querier) error {
		rowsAffected, err := q.EditComment(ctx, sqlc.EditCommentParams{
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
		if req.Image > 0 {
			err := q.UpsertImage(ctx, sqlc.UpsertImageParams{
				ID:       req.Image.Int64(),
				ParentID: req.CommentId.Int64(),
			})
			if err != nil {
				return err
			}
		} else {
			rowsAffected, err := q.DeleteImage(ctx, req.Image.Int64())
			if err != nil {
				return err
			}
			if rowsAffected != 1 {
				fmt.Println("image not found")
			}
		}
		return nil
	})

}

func (s *Application) DeleteComment(ctx context.Context, req models.GenericReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrNotAllowed
	}

	rowsAffected, err := s.db.DeleteComment(ctx, sqlc.DeleteCommentParams{
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

	rows, err := s.db.GetCommentsByPostId(ctx, sqlc.GetCommentsByPostIdParams{
		ParentID: req.EntityId.Int64(),
		UserID:   req.RequesterId.Int64(),
		Limit:    req.Limit.Int32(),
		Offset:   req.Offset.Int32(),
	})
	if err != nil {
		return nil, err
	}
	comments := make([]models.Comment, 0, len(rows))
	userIDs := make([]int64, 0, len(rows))

	for _, r := range rows {
		uid := r.CommentCreatorID
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
			Image:          ct.Id(r.Image),
		})
	}

	if len(comments) == 0 {
		return comments, nil
	}

	userMap, err := s.hydrator.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		uid := comments[i].User.UserId.Int64()
		if u, ok := userMap[uid]; ok {
			comments[i].User = u
		}
	}
	// if err := s.hydrateComments(ctx, comments); err != nil {
	// 	return nil, err
	// }

	return comments, nil
}

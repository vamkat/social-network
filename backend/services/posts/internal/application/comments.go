package application

import (
	"context"
	"fmt"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
)

func (s *Application) CreateComment(ctx context.Context, req CreateCommentReq) (err error) {
	// check requester can actually view parent entity? (probably not needed?)
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	err = s.db.CreateComment(ctx, sqlc.CreateCommentParams{
		CommentCreatorID: req.CreatorId.Int64(),
		ParentID:         req.ParentId.Int64(),
		CommentBody:      req.Body.String(),
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *Application) EditComment(ctx context.Context, req EditCommentReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	err := s.runTx(ctx, func(q *sqlc.Queries) error {
		rowsAffected, err := s.db.EditComment(ctx, sqlc.EditCommentParams{
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
			err := s.db.UpsertImage(ctx, sqlc.UpsertImageParams{
				ID:       req.Image.Int64(),
				ParentID: req.CommentId.Int64(),
			})
			if err != nil {
				return err
			}
		} else {
			rowsAffected, err := s.db.DeleteImage(ctx, req.Image.Int64())
			if err != nil {
				return err
			}
			if rowsAffected != 1 {
				fmt.Println("image not found")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Application) DeleteComment(ctx context.Context, req GenericReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
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

func (s *Application) GetCommentsByParentId(ctx context.Context, req GenericPaginatedReq) ([]Comment, error) {
	// check requester can actually view parent entity?
	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
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
	comments := make([]Comment, 0, len(rows))
	for _, r := range rows {
		comments = append(comments, Comment{
			CommentId:      ct.Id(r.ID),
			ParentId:       req.EntityId,
			Body:           ct.CommentBody(r.CommentBody),
			CreatorId:      ct.Id(r.CommentCreatorID),
			ReactionsCount: int(r.ReactionsCount),
			CreatedAt:      r.CreatedAt.Time,
			UpdatedAt:      r.UpdatedAt.Time,
			LikedByUser:    r.LikedByUser,
			Image:          ct.Id(r.Image),
		})
	}
	return comments, nil
}

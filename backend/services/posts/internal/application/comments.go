package application

import (
	"context"
	"fmt"
	ds "social-network/services/posts/internal/db/dbservice"
	"social-network/shared/gen-go/media"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"
)

const genericPublic = "posts service error"

func (s *Application) CreateComment(ctx context.Context, req models.CreateCommentReq) (err error) {
	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed").WithPublic("invalid data received")
	}

	accessCtx := accessContext{
		requesterId: req.CreatorId.Int64(),
		entityId:    req.ParentId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil { //err := "invalid parten id", "CreateComment>CreateComment>HasAccess: invalid parent id"
		return ce.Wrap(ce.ErrInternal, err).WithPublic(genericPublic)
	}
	if !hasAccess {
		return ce.Wrap(ce.ErrPermissionDenied, fmt.Errorf("user has no permission to view or edit entity: %v", req.ParentId)).WithPublic("permission denied")
	}
	var commentId int64
	err = s.txRunner.RunTx(ctx, func(q *ds.Queries) error {
		commentId, err = q.CreateComment(ctx, ds.CreateCommentParams{
			CommentCreatorID: req.CreatorId.Int64(),
			ParentID:         req.ParentId.Int64(),
			CommentBody:      req.Body.String(),
		})

		if err != nil {
			return ce.Wrap(ce.ErrInternal, err).WithPublic(genericPublic)
		}

		if req.ImageId != 0 {
			err = q.UpsertImage(ctx, ds.UpsertImageParams{
				ID:       req.ImageId.Int64(),
				ParentID: commentId,
			})
			if err != nil {
				return ce.Wrap(ce.ErrInternal, err).WithPublic(genericPublic)
			}
		}
		return nil
	})
	if err != nil {
		return ce.Wrap(nil, err)
	}

	//create notification
	userMap, err := s.userRetriever.GetUsers(ctx, ct.Ids{req.CreatorId})
	if err != nil {
		tele.Error(ctx, "failed to GetUsers for @1: @2 ", "userId", req.CreatorId, "error", err.Error())
		return nil //return with no error but without creating non-essential notif
	}
	var commenterUsername string
	if u, ok := userMap[req.CreatorId]; ok {
		commenterUsername = u.Username.String()
	}
	basicPost, err := s.db.GetBasicPostByID(ctx, req.ParentId.Int64())
	if err != nil {
		tele.Error(ctx, "get GetBasicPostByID failed for post @1: @2 ", "post id", req.ParentId, "error", err.Error())
		return nil //return with no error but without creating non-essential notif
	}
	err = s.clients.CreatePostComment(ctx, basicPost.CreatorID, req.CreatorId.Int64(), req.ParentId.Int64(), commenterUsername, req.Body.String())
	if err != nil {
		tele.Error(ctx, "CreatePostComment notification failed for comment @1: @2", "comment id", commentId, "error", err.Error())
		return nil //return with no error but without creating non-essential notif
	}
	return nil
}

func (s *Application) EditComment(ctx context.Context, req models.EditCommentReq) error {
	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed").WithPublic("invalid data received")
	}

	accessCtx := accessContext{
		requesterId: req.CreatorId.Int64(),
		entityId:    req.CommentId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return ce.Wrap(ce.ErrInternal, err).WithPublic(genericPublic)
	}
	if !hasAccess {
		return ce.Wrap(ce.ErrPermissionDenied, fmt.Errorf("user has no permission to view or edit this entity")).WithPublic("permission denied")
	}

	err = s.txRunner.RunTx(ctx, func(q *ds.Queries) error {
		rowsAffected, err := q.EditComment(ctx, ds.EditCommentParams{
			CommentBody:      req.Body.String(),
			ID:               req.CommentId.Int64(),
			CommentCreatorID: req.CreatorId.Int64(),
		})
		if err != nil {
			return ce.Wrap(ce.ErrInternal, err).WithPublic(genericPublic)
		}
		if rowsAffected != 1 {
			return ce.Wrap(ce.ErrNotFound, fmt.Errorf("comment not found or not owned by user")).WithPublic("permission denied")
		}

		if req.ImageId > 0 {
			err := q.UpsertImage(ctx, ds.UpsertImageParams{
				ID:       req.ImageId.Int64(),
				ParentID: req.CommentId.Int64(),
			})
			if err != nil {
				return ce.Wrap(ce.ErrInternal, err).WithPublic(genericPublic)
			}
		}

		if req.DeleteImage {
			rowsAffected, err := q.DeleteImage(ctx, req.CommentId.Int64())
			if err != nil {
				return ce.Wrap(ce.ErrInternal, err).WithPublic(genericPublic)
			}
			if rowsAffected != 1 {
				tele.Warn(ctx, "EditComment: image @1 for comment @2 could not be deleted: not found.", "image id", req.ImageId, "comment id", req.CommentId)
			}
		}
		return nil
	})
	if err != nil {
		return ce.Wrap(nil, err)
	}

	return nil
}
func (s *Application) DeleteComment(ctx context.Context, req models.GenericReq) error {
	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed").WithPublic("invalid data received")
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return ce.Wrap(ce.ErrInternal, err, "s.hasRightToView failed with accessCtx: %#v", ct.ReqActionDetails).WithPublic(genericPublic)
	}
	if !hasAccess {
		return ce.Wrap(ce.ErrPermissionDenied, fmt.Errorf("user has no permission to view or edit this entity")).WithPublic("permission denied")
	}

	rowsAffected, err := s.db.DeleteComment(ctx, ds.DeleteCommentParams{
		ID:               req.EntityId.Int64(),
		CommentCreatorID: req.RequesterId.Int64(),
	})
	if err != nil {
		return ce.Wrap(ce.ErrInternal, err, "db.DeleteComment failed, with entity").WithPublic(genericPublic)
	}
	if rowsAffected != 1 {
		tele.Warn(ctx, "DeleteComment: comment @1 could not be deleted: not found.", "comment id", req.EntityId)
	}

	return nil
}

func (s *Application) GetCommentsByParentId(ctx context.Context, req models.EntityIdPaginatedReq) ([]models.Comment, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return nil, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed").WithPublic("invalid data received")
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return nil, ce.Wrap(ce.ErrInternal, err).WithPublic(genericPublic)
	}
	if !hasAccess {
		return nil, ce.Wrap(ce.ErrPermissionDenied, fmt.Errorf("user has no permission to view comments of this entity")).WithPublic("permission denied")
	}

	rows, err := s.db.GetCommentsByPostId(ctx, ds.GetCommentsByPostIdParams{
		ParentID: req.EntityId.Int64(),
		UserID:   req.RequesterId.Int64(),
		Limit:    req.Limit.Int32(),
		Offset:   req.Offset.Int32(),
	})
	if err != nil {
		return nil, ce.Wrap(ce.ErrInternal, err).WithPublic(genericPublic)
	}

	if len(rows) == 0 {
		return []models.Comment{}, nil
	}

	comments := make([]models.Comment, 0, len(rows))
	userIDs := make(ct.Ids, 0, len(rows))
	commentImageIds := make(ct.Ids, 0, len(rows))

	for _, r := range rows {
		userIDs = append(userIDs, ct.Id(r.CommentCreatorID))

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
			commentImageIds = append(commentImageIds, ct.Id(r.Image))
		}
	}

	userMap, err := s.userRetriever.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, ce.Wrap(ce.ErrInternal, err).WithPublic("error retrieving users info")
	}

	var imageMap map[int64]string
	if len(commentImageIds) > 0 {
		imageMap, _, err = s.mediaRetriever.GetImages(ctx, commentImageIds, media.FileVariant_MEDIUM)
	}
	if err != nil {
		return nil, ce.Wrap(ce.ErrInternal, err).WithPublic("error retrieving images")
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

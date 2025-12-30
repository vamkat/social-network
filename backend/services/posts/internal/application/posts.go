package application

import (
	"context"
	"database/sql"
	"errors"
	ds "social-network/services/posts/internal/db/dbservice"
	"social-network/shared/gen-go/media"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Application) CreatePost(ctx context.Context, req models.CreatePostReq) (err error) {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	var groupId pgtype.Int8
	groupId.Int64 = req.GroupId.Int64()
	if req.GroupId == 0 {
		groupId.Valid = false
	}

	audience := ds.IntendedAudience(req.Audience.String())

	if !groupId.Valid && audience == "group" {
		return ErrNoGroupIdGiven
	}

	if groupId.Valid {
		isMember, err := s.clients.IsGroupMember(ctx, req.CreatorId.Int64(), req.GroupId.Int64())
		if err != nil {
			return err
		}
		if !isMember {
			return ErrNotAllowed
		}
	}
	return s.txRunner.RunTx(ctx, func(q *ds.Queries) error {

		postId, err := q.CreatePost(ctx, ds.CreatePostParams{
			PostBody:  req.Body.String(),
			CreatorID: req.CreatorId.Int64(),
			GroupID:   groupId,
			Audience:  audience,
		})
		if err != nil {
			return err
		}

		if audience == "selected" {
			if len(req.AudienceIds) < 1 {
				return ErrNoAudienceSelected
			}
			rowsAffected, err := q.InsertPostAudience(ctx, ds.InsertPostAudienceParams{
				PostID:         postId,
				AllowedUserIds: req.AudienceIds.Int64(), //does nil work here? TODO test
			})
			if err != nil {
				return err
			}
			if rowsAffected < 1 {
				return ErrNoAudienceSelected
			}
		}

		if req.ImageId != 0 {
			err = q.UpsertImage(ctx, ds.UpsertImageParams{
				ID:       req.ImageId.Int64(),
				ParentID: postId,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

}

func (s *Application) DeletePost(ctx context.Context, req models.GenericReq) error {
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

	rowsAffected, err := s.db.DeletePost(ctx, ds.DeletePostParams{
		ID:        int64(req.EntityId),
		CreatorID: req.RequesterId.Int64(),
	})
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrNotFound
	}
	return nil
}

func (s *Application) EditPost(ctx context.Context, req models.EditPostReq) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.PostId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrNotAllowed
	}

	return s.txRunner.RunTx(ctx, func(q *ds.Queries) error {
		//edit content
		if len(req.NewBody) > 0 {
			rowsAffected, err := q.EditPostContent(ctx, ds.EditPostContentParams{
				PostBody:  req.NewBody.String(),
				ID:        req.PostId.Int64(),
				CreatorID: req.RequesterId.Int64(),
			})
			if err != nil {
				return err
			}
			if rowsAffected != 1 {
				return ErrNotFound
			}
		}

		//edit image
		if req.ImageId > 0 {
			err := q.UpsertImage(ctx, ds.UpsertImageParams{
				ID:       req.ImageId.Int64(),
				ParentID: req.PostId.Int64(),
			})
			if err != nil {
				return err
			}
		}
		//delete image
		if req.DeleteImage {
			rowsAffected, err := q.DeleteImage(ctx, req.PostId.Int64())
			if err != nil {
				return err
			}
			if rowsAffected != 1 {
				tele.Warn(ctx, "EditPost: image to be deleted not found. @1", "request", req)
			}
		}
		// edit audience
		rowsAffected, err := q.UpdatePostAudience(ctx, ds.UpdatePostAudienceParams{
			ID:        req.PostId.Int64(),
			CreatorID: req.RequesterId.Int64(),
			Audience:  ds.IntendedAudience(req.Audience),
		})
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			tele.Warn(ctx, "EditPost: no audience change. @1", "request", req)
		}

		// edit audience ids
		if req.Audience == "selected" && rowsAffected == 1 {
			if len(req.AudienceIds) < 1 {
				return ErrNoAudienceSelected
			}
			//delete previous audience ids
			err := q.ClearPostAudience(ctx, req.PostId.Int64())
			if err != nil {
				return err
			}
			//insert new ids
			rowsAffected, err := q.InsertPostAudience(ctx, ds.InsertPostAudienceParams{
				PostID:         req.PostId.Int64(),
				AllowedUserIds: req.AudienceIds.Int64(),
			})
			if err != nil {
				return err
			}
			if rowsAffected < 1 {
				return ErrNoAudienceSelected
			}
		}

		return nil
	})

}

func (s *Application) GetMostPopularPostInGroup(ctx context.Context, req models.SimpleIdReq) (models.Post, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return models.Post{}, err
	}

	var groupId pgtype.Int8
	groupId.Int64 = req.Id.Int64()
	groupId.Valid = true
	p, err := s.db.GetMostPopularPostInGroup(ctx, groupId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, ErrNotFound
		}
		return models.Post{}, err
	}

	userMap, err := s.userRetriever.GetUsers(ctx, ct.Ids{ct.Id(p.CreatorID)})
	if err != nil {
		return models.Post{}, err
	}

	post := models.Post{
		PostId:          ct.Id(p.ID),
		Body:            ct.PostBody(p.PostBody),
		User:            userMap[ct.Id(p.CreatorID)],
		GroupId:         ct.Id(req.Id.Int64()),
		Audience:        ct.Audience(p.Audience),
		CommentsCount:   int(p.CommentsCount),
		ReactionsCount:  int(p.ReactionsCount),
		LastCommentedAt: ct.GenDateTime(p.LastCommentedAt.Time),
		CreatedAt:       ct.GenDateTime(p.CreatedAt.Time),
		UpdatedAt:       ct.GenDateTime(p.UpdatedAt.Time),
		ImageId:         ct.Id(p.Image),
	}

	if post.ImageId > 0 {
		imageUrl, err := s.mediaRetriever.GetImage(ctx, p.Image, media.FileVariant_MEDIUM)
		if err != nil {
			return models.Post{}, err
		}

		post.ImageUrl = imageUrl
	}

	return post, nil
}

func (s *Application) GetPostById(ctx context.Context, req models.GenericReq) (models.Post, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return models.Post{}, err
	}
	userCanView, err := s.hasRightToView(ctx, accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	})
	if err != nil {
		return models.Post{}, err
	}
	if !userCanView {
		return models.Post{}, ErrNotAllowed
	}

	p, err := s.db.GetPostByID(ctx, ds.GetPostByIDParams{
		UserID: req.RequesterId.Int64(),
		ID:     req.EntityId.Int64(),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, ErrNotFound
		}
		return models.Post{}, err
	}

	userMap, err := s.userRetriever.GetUsers(ctx, ct.Ids{ct.Id(p.CreatorID)})
	if err != nil {
		return models.Post{}, err
	}

	post := models.Post{
		PostId:          ct.Id(p.ID),
		Body:            ct.PostBody(p.PostBody),
		User:            userMap[ct.Id(p.CreatorID)],
		GroupId:         ct.Id(p.GroupID),
		Audience:        ct.Audience(p.Audience),
		CommentsCount:   int(p.CommentsCount),
		ReactionsCount:  int(p.ReactionsCount),
		LastCommentedAt: ct.GenDateTime(p.LastCommentedAt.Time),
		CreatedAt:       ct.GenDateTime(p.CreatedAt.Time),
		UpdatedAt:       ct.GenDateTime(p.UpdatedAt.Time),
		LikedByUser:     p.LikedByUser,
		ImageId:         ct.Id(p.Image),
	}

	if post.ImageId > 0 {
		imageUrl, err := s.mediaRetriever.GetImage(ctx, p.Image, media.FileVariant_MEDIUM)
		if err != nil {
			return models.Post{}, err
		}

		post.ImageUrl = imageUrl
	}

	return post, nil
}

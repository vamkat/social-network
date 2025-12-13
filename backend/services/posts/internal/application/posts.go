package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"

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

	audience := sqlc.IntendedAudience(req.Audience.String())

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
	return s.txRunner.RunTx(ctx, func(q sqlc.Querier) error {

		postId, err := q.CreatePost(ctx, sqlc.CreatePostParams{
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
			rowsAffected, err := q.InsertPostAudience(ctx, sqlc.InsertPostAudienceParams{
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

		if req.Image != 0 {
			err = q.UpsertImage(ctx, sqlc.UpsertImageParams{
				ID:       req.Image.Int64(),
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

	rowsAffected, err := s.db.DeletePost(ctx, sqlc.DeletePostParams{
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

	return s.txRunner.RunTx(ctx, func(q sqlc.Querier) error {
		//edit content
		if len(req.NewBody) > 0 {
			rowsAffected, err := q.EditPostContent(ctx, sqlc.EditPostContentParams{
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
		//edit image //does this also delete the image?
		if req.Image > 0 {
			err := q.UpsertImage(ctx, sqlc.UpsertImageParams{
				ID:       req.Image.Int64(),
				ParentID: req.PostId.Int64(),
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
		// edit audience
		rowsAffected, err := q.UpdatePostAudience(ctx, sqlc.UpdatePostAudienceParams{
			ID:        req.PostId.Int64(),
			CreatorID: req.RequesterId.Int64(),
			Audience:  sqlc.IntendedAudience(req.Audience),
		})
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			fmt.Println("no audience change")
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
			rowsAffected, err := q.InsertPostAudience(ctx, sqlc.InsertPostAudienceParams{
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

	userMap, err := s.hydrator.GetUsers(ctx, []int64{p.CreatorID})
	if err != nil {
		return models.Post{}, err
	}

	post := models.Post{
		PostId:          ct.Id(p.ID),
		Body:            ct.PostBody(p.PostBody),
		User:            userMap[p.CreatorID],
		GroupId:         ct.Id(req.Id.Int64()),
		Audience:        ct.Audience(p.Audience),
		CommentsCount:   int(p.CommentsCount),
		ReactionsCount:  int(p.ReactionsCount),
		LastCommentedAt: ct.GenDateTime(p.LastCommentedAt.Time),
		CreatedAt:       ct.GenDateTime(p.CreatedAt.Time),
		UpdatedAt:       ct.GenDateTime(p.UpdatedAt.Time),
		Image:           ct.Id(p.Image),
	}

	// if err := s.hydratePost(ctx, &post); err != nil {
	// 	return models.Post{}, err
	// }

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

	p, err := s.db.GetPostByID(ctx, sqlc.GetPostByIDParams{})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, ErrNotFound
		}
		return models.Post{}, err
	}

	userMap, err := s.hydrator.GetUsers(ctx, []int64{p.CreatorID})
	if err != nil {
		return models.Post{}, err
	}

	var groupId pgtype.Int8
	groupId.Int64 = p.GroupID.Int64
	groupId.Valid = true
	post := models.Post{
		PostId:          ct.Id(p.ID),
		Body:            ct.PostBody(p.PostBody),
		User:            userMap[p.CreatorID],
		GroupId:         ct.Id(groupId.Int64),
		Audience:        ct.Audience(p.Audience),
		CommentsCount:   int(p.CommentsCount),
		ReactionsCount:  int(p.ReactionsCount),
		LastCommentedAt: ct.GenDateTime(p.LastCommentedAt.Time),
		CreatedAt:       ct.GenDateTime(p.CreatedAt.Time),
		UpdatedAt:       ct.GenDateTime(p.UpdatedAt.Time),
		Image:           ct.Id(p.Image),
	}

	// if err := s.hydratePost(ctx, &post); err != nil {
	// 	return models.Post{}, err
	// }

	return post, nil
}

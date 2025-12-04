package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"

	"github.com/jackc/pgx/v5/pgtype"
)

// GENERAL NOTE For every response that includes a userId, actual basic user info will be retrieved by Gateway from Users

func (s *Application) CreatePost(ctx context.Context, req CreatePostReq) (err error) {
	// if group post, check creator is a member (HANDLER)

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	err = s.runTx(ctx, func(q *sqlc.Queries) error {

		var groupId pgtype.Int8
		groupId.Int64 = req.GroupId.Int64()
		if req.GroupId == 0 {
			groupId.Valid = false
		}

		audience := sqlc.IntendedAudience(req.Audience.String())

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
	if err != nil {
		return err
	}
	return nil
}

func (s *Application) DeletePost(ctx context.Context, req GenericReq) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
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

func (s *Application) EditPost(ctx context.Context, req EditPostReq) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	err := s.runTx(ctx, func(q *sqlc.Queries) error {
		//edit content
		if len(req.NewBody) > 0 {
			rowsAffected, err := s.db.EditPostContent(ctx, sqlc.EditPostContentParams{
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
			err := s.db.UpsertImage(ctx, sqlc.UpsertImageParams{
				ID:       req.Image.Int64(),
				ParentID: req.PostId.Int64(),
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
		// edit audience
		rowsAffected, err := s.db.UpdatePostAudience(ctx, sqlc.UpdatePostAudienceParams{
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
			err := s.db.ClearPostAudience(ctx, req.PostId.Int64())
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
	if err != nil {
		return err
	}
	return nil
}

func (s *Application) GetGroupPostsPaginated(ctx context.Context, req GenericPaginatedReq) ([]Post, error) {
	//check requester is group member (HANDLER)
	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}

	var groupId pgtype.Int8
	groupId.Int64 = req.EntityId.Int64()
	if req.EntityId == 0 {
		return nil, ErrNoGroupIdGiven
	}
	groupId.Valid = true

	rows, err := s.db.GetGroupPostsPaginated(ctx, sqlc.GetGroupPostsPaginatedParams{
		GroupID: groupId,
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrNotFound
	}
	posts := make([]Post, 0, len(rows))
	for _, r := range rows {
		posts = append(posts, Post{
			PostId:          ct.Id(r.ID),
			Body:            ct.PostBody(r.PostBody),
			CreatorId:       ct.Id(r.CreatorID),
			GroupId:         req.EntityId,
			Audience:        ct.Audience(r.Audience),
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
				UpdatedAt:      r.UpdatedAt.Time,
				LikedByUser:    r.LatestCommentLikedByUser,
				Image:          ct.Id(r.LatestCommentImage),
			},
		})
	}

	return posts, nil
}

func (s *Application) GetMostPopularPostInGroup(ctx context.Context, groupID ct.Id) (Post, error) {
	if err := ct.ValidateStruct(groupID); err != nil {
		return Post{}, err
	}

	var groupId pgtype.Int8
	groupId.Int64 = groupID.Int64()
	if groupID == 0 {
		return Post{}, ErrNoGroupIdGiven
	}
	groupId.Valid = true
	p, err := s.db.GetMostPopularPostInGroup(ctx, groupId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Post{}, ErrNotFound
		}
		return Post{}, err
	}
	post := Post{
		PostId:          ct.Id(p.ID),
		Body:            ct.PostBody(p.PostBody),
		CreatorId:       ct.Id(p.CreatorID),
		GroupId:         ct.Id(groupID.Int64()),
		Audience:        ct.Audience(p.Audience),
		CommentsCount:   int(p.CommentsCount),
		ReactionsCount:  int(p.ReactionsCount),
		LastCommentedAt: p.LastCommentedAt.Time,
		CreatedAt:       p.CreatedAt.Time,
		UpdatedAt:       p.UpdatedAt.Time,
		Image:           ct.Id(p.Image),
		//no latest comment or liked by user needed here
	}
	return post, nil
}

func (s *Application) GetUserPostsPaginated(ctx context.Context, req GetUserPostsReq) ([]Post, error) {
	// other than followers, rest of checks happen in query
	// HANDLER needs to get FOLLOWERS LIST for creatorId from users

	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}

	rows, err := s.db.GetUserPostsPaginated(ctx, sqlc.GetUserPostsPaginatedParams{
		CreatorID: req.CreatorId.Int64(),
		UserID:    req.RequesterId.Int64(),
		Column3:   req.CreatorFollowers.Int64(),
		Limit:     req.Limit.Int32(),
		Offset:    req.Offset.Int32(),
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrNotFound
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
				UpdatedAt:      r.UpdatedAt.Time,
				LikedByUser:    r.LatestCommentLikedByUser,
				Image:          ct.Id(r.LatestCommentImage),
			},
		})

	}

	return posts, nil
}

// NOT CURRENTLY NEEDED
func (s *Application) GetPostById(ctx context.Context, req GenericReq) (Post, error) {
	//check requester is allowed to view post, dependes on post audience:
	//everyone: any requester can see
	//followers: API GATEWAY(?) needs to get FOLLOWERS LIST for creatorId from users
	//selected: check can happen in posts service
	//group: API GATEWAY(?) needs to check requester is member of group
	return Post{}, nil
}

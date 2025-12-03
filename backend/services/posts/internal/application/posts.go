package application

import (
	"context"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"

	"github.com/jackc/pgx/v5/pgtype"
)

// GENERAL NOTE For every response that includes a userId, actual basic user info will be retrieved by Gateway from Users

func (s *PostsService) CreatePost(ctx context.Context, req CreatePostReq) (err error) {
	// if group post, check creator is a member (GATEWAY)

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	return s.runTx(ctx, func(q sqlc.Querier) error {

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
			err = q.InsertPostAudience(ctx, sqlc.InsertPostAudienceParams{
				PostID:  postId,
				Column2: req.AudienceIds.Int64(), //does nil work here? And do I allow it? TODO test
			})
			if err != nil {
				return err
			}
		}

		if req.Image != 0 {
			err = q.InsertImage(ctx, sqlc.InsertImageParams{
				Column1: req.Image.Int64(),
				Column2: postId,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

}

// NOT GRPC
func (s *PostsService) insertPostAudience(ctx context.Context, audienceIds []int64) error {
	//always in trasaction (change signature to accept transaction?)
	return nil
}

// NOT GRPC
func (s *PostsService) clearPostAudience(ctx context.Context, postId int64) error {
	//always in trasaction (change signature to accept transaction?)
	return nil
}

func (s *PostsService) DeletePost(ctx context.Context, req GenericReq) error {
	//check requester is post creator
	return nil
}

// FRONT: do you prefer just error instead of full post?
func (s *PostsService) EditPostContent(ctx context.Context, req EditPostContentReq) error {
	// check requester is post creator
	return nil
}

// // FRONT: do you prefer full post instead of just error?
// func (s *PostsService) EditPostAudience(ctx context.Context, req EditPostAudienceReq) error {
// 	// check requester is post creator
// 	//run in transaction
// 	//s.ClearPostAudience to remove previous audience
// 	//s.db.UpdatePostAudience
// 	// if audience=selected, s.InsertPostAudience
// 	return nil
// }

func (s *PostsService) GetGroupPostsPaginated(ctx context.Context, req GenericPaginatedReq) ([]Post, error) {
	//check requester is group member (cross service, API gateway should do it?)
	return nil, nil
}

// FRONT: If you also want to display liked by user information I also need the requester id- NOT liked by user
func (s *PostsService) GetMostPopularPostInGroup(ctx context.Context, groupId int64) (Post, error) {
	//anyone can see this
	// if no post ErrNoPost
	return Post{}, nil
}

// FRONT: if we only have one image per post it's likely you'll never need this? - LEAVE FOR LATER
func (s *PostsService) GetPostById(ctx context.Context, req GenericReq) (Post, error) {
	//check requester is allowed to view post, dependes on post audience:
	//everyone: any requester can see
	//followers: API GATEWAY(?) needs to get FOLLOWERS LIST for creatorId from users
	//selected: check can happen in posts service
	//group: API GATEWAY(?) needs to check requester is member of group
	return Post{}, nil
}

func (s *PostsService) GetUserPostsPaginated(ctx context.Context, req GetUserPostsReq) ([]Post, error) {
	// other than followers, rest of checks happen in query
	// API GATEWAY(?) needs to get FOLLOWERS LIST for creatorId from users
	return nil, nil
}

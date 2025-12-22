package handlers

import (
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/posts"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
)

func (h *Handlers) getPublicFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getPublicFeed handler called")

		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GenericPaginatedReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.GenericPaginatedReq{
			RequesterId: claims.UserId,
			Limit:       body.Limit.Int32(),
			Offset:      body.Offset.Int32(),
		}

		grpcResp, err := h.PostsService.GetPublicFeed(ctx, &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to get public feed: "+err.Error())
			return
		}

		fmt.Println("retrieved public feed: ", grpcResp)

		var postsResponse []models.Post
		for _, p := range grpcResp.Posts {
			post := models.Post{
				PostId: ct.Id(p.PostId),
				Body:   ct.PostBody(p.PostBody),
				User: models.User{
					UserId:    ct.Id(p.User.UserId),
					Username:  ct.Username(p.User.Username),
					AvatarId:  ct.Id(p.User.Avatar),
					AvatarURL: p.User.AvatarUrl,
				},
				GroupId:         ct.Id(p.GroupId),
				Audience:        ct.Audience(p.Audience),
				CommentsCount:   int(p.CommentsCount),
				ReactionsCount:  int(p.ReactionsCount),
				LastCommentedAt: ct.GenDateTime(p.LastCommentedAt.AsTime()),
				CreatedAt:       ct.GenDateTime(p.CreatedAt.AsTime()),
				UpdatedAt:       ct.GenDateTime(p.UpdatedAt.AsTime()),
				LikedByUser:     p.LikedByUser,
				Image:           ct.Id(p.Image),
			}
			postsResponse = append(postsResponse, post)
		}

		err = utils.WriteJSON(w, http.StatusOK, postsResponse)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send public feed")
			return
		}

	}
}

func (h *Handlers) getPersonalizedFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getPersonalizedFeed handler called")

		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GetPersonalizedFeedReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.GetPersonalizedFeedReq{
			RequesterId: claims.UserId,
			Limit:       body.Limit.Int32(),
			Offset:      body.Offset.Int32(),
		}

		grpcResp, err := h.PostsService.GetPersonalizedFeed(ctx, &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to get personalized feed: "+err.Error())
			return
		}

		fmt.Println("retrieved personalized feed: ", grpcResp)

		var postsResponse []models.Post
		for _, p := range grpcResp.Posts {
			post := models.Post{
				PostId: ct.Id(p.PostId),
				Body:   ct.PostBody(p.PostBody),
				User: models.User{
					UserId:    ct.Id(p.User.UserId),
					Username:  ct.Username(p.User.Username),
					AvatarId:  ct.Id(p.User.Avatar),
					AvatarURL: p.User.AvatarUrl,
				},
				GroupId:         ct.Id(p.GroupId),
				Audience:        ct.Audience(p.Audience),
				CommentsCount:   int(p.CommentsCount),
				ReactionsCount:  int(p.ReactionsCount),
				LastCommentedAt: ct.GenDateTime(p.LastCommentedAt.AsTime()),
				CreatedAt:       ct.GenDateTime(p.CreatedAt.AsTime()),
				UpdatedAt:       ct.GenDateTime(p.UpdatedAt.AsTime()),
				LikedByUser:     p.LikedByUser,
				Image:           ct.Id(p.Image),
			}
			postsResponse = append(postsResponse, post)
		}

		err = utils.WriteJSON(w, http.StatusOK, postsResponse)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send public feed")
			return
		}

	}
}

func (h *Handlers) getUserPostsPaginated() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getUserPostsPaginated handler called")

		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GetUserPostsReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.GetUserPostsReq{
			RequesterId: claims.UserId,
			CreatorId:   body.CreatorId.Int64(),
			Limit:       body.Limit.Int32(),
			Offset:      body.Offset.Int32(),
		}

		grpcResp, err := h.PostsService.GetUserPostsPaginated(ctx, &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to get personalized feed: "+err.Error())
			return
		}

		fmt.Println("retrieved personalized feed: ", grpcResp)

		var postsResponse []models.Post
		for _, p := range grpcResp.Posts {
			post := models.Post{
				PostId: ct.Id(p.PostId),
				Body:   ct.PostBody(p.PostBody),
				User: models.User{
					UserId:    ct.Id(p.User.UserId),
					Username:  ct.Username(p.User.Username),
					AvatarId:  ct.Id(p.User.Avatar),
					AvatarURL: p.User.AvatarUrl,
				},
				GroupId:         ct.Id(p.GroupId),
				Audience:        ct.Audience(p.Audience),
				CommentsCount:   int(p.CommentsCount),
				ReactionsCount:  int(p.ReactionsCount),
				LastCommentedAt: ct.GenDateTime(p.LastCommentedAt.AsTime()),
				CreatedAt:       ct.GenDateTime(p.CreatedAt.AsTime()),
				UpdatedAt:       ct.GenDateTime(p.UpdatedAt.AsTime()),
				LikedByUser:     p.LikedByUser,
				Image:           ct.Id(p.Image),
			}
			postsResponse = append(postsResponse, post)
		}

		err = utils.WriteJSON(w, http.StatusOK, postsResponse)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to send user %v posts: %v", body.CreatorId, err.Error()))
			return
		}

	}
}

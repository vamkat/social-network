package handlers

import (
	"fmt"
	"net/http"
	"social-network/shared/gen-go/posts"
	ct "social-network/shared/go/ct"
	utils "social-network/shared/go/http-utils"
	"social-network/shared/go/jwt"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"
)

func (h *Handlers) getPublicFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tele.Info(ctx, "getPublicFeed handler called")

		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GenericPaginatedReq{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.GenericPaginatedReq{
			RequesterId: claims.UserId,
			Limit:       body.Limit.Int32(),
			Offset:      body.Offset.Int32(),
		}

		grpcResp, err := h.PostsService.GetPublicFeed(ctx, &grpcReq)
		if err != nil {
			utils.ReturnHttpError(ctx, w, err)
			//utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to get public feed: "+err.Error())
			return
		}

		tele.Info(ctx, "retrieved public feed. @1", "grpcResp", grpcResp)

		postsResponse := []models.Post{}
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
				ImageId:         ct.Id(p.ImageId),
				ImageUrl:        p.ImageUrl,
			}
			postsResponse = append(postsResponse, post)
		}

		err = utils.WriteJSON(ctx, w, http.StatusOK, postsResponse)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to send public feed")
			return
		}

	}
}

func (h *Handlers) getPersonalizedFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tele.Info(ctx, "getPersonalizedFeed handler called")

		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GetPersonalizedFeedReq{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.GetPersonalizedFeedReq{
			RequesterId: claims.UserId,
			Limit:       body.Limit.Int32(),
			Offset:      body.Offset.Int32(),
		}

		grpcResp, err := h.PostsService.GetPersonalizedFeed(ctx, &grpcReq)
		if err != nil {
			utils.ReturnHttpError(ctx, w, err)
			//utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to get personalized feed: "+err.Error())
			return
		}

		tele.Info(ctx, "retrieved personalized feed. @1", "grpcResp", grpcResp)

		postsResponse := []models.Post{}
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
				ImageId:         ct.Id(p.ImageId),
				ImageUrl:        p.ImageUrl,
			}
			postsResponse = append(postsResponse, post)
		}

		err = utils.WriteJSON(ctx, w, http.StatusOK, postsResponse)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to send public feed")
			return
		}

	}
}

func (h *Handlers) getUserPostsPaginated() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tele.Info(ctx, "getUserPostsPaginated handler called")

		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GetUserPostsReq{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
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
			utils.ReturnHttpError(ctx, w, err)
			//utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to get personalized feed: "+err.Error())
			return
		}

		tele.Info(ctx, "retrieved personalized feed. @1", "grpcResp", grpcResp)

		postsResponse := []models.Post{}
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
				ImageId:         ct.Id(p.ImageId),
				ImageUrl:        p.ImageUrl,
			}
			postsResponse = append(postsResponse, post)
		}

		err = utils.WriteJSON(ctx, w, http.StatusOK, postsResponse)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("failed to send user %v posts: %v", body.CreatorId, err.Error()))
			return
		}

	}
}

func (h *Handlers) getGroupPostsPaginated() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tele.Info(ctx, "getGroupPostsPaginated handler called")

		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GetGroupPostsReq{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.GetGroupPostsReq{
			RequesterId: claims.UserId,
			GroupId:     body.GroupId.Int64(),
			Limit:       body.Limit.Int32(),
			Offset:      body.Offset.Int32(),
		}

		grpcResp, err := h.PostsService.GetGroupPostsPaginated(ctx, &grpcReq)
		if err != nil {
			utils.ReturnHttpError(ctx, w, err)
			//utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to get group feed: "+err.Error())
			return
		}

		tele.Info(ctx, "retrieved group feed. @1", "grpcResp", grpcResp)

		postsResponse := []models.Post{}
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
				ImageId:         ct.Id(p.ImageId),
				ImageUrl:        p.ImageUrl,
			}
			postsResponse = append(postsResponse, post)
		}

		err = utils.WriteJSON(ctx, w, http.StatusOK, postsResponse)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("failed to send group %v posts: %v", body.GroupId, err.Error()))
			return
		}

	}
}

package handlers

import (
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/common"
	"social-network/shared/gen-go/posts"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
)

func (h *Handlers) getPostById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getPostById handler called")

		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GenericReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.GenericReq{
			RequesterId: int64(claims.UserId),
			EntityId:    body.EntityId.Int64(),
		}

		grpcResp, err := h.PostsService.GetPostById(r.Context(), &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to get post with id %v: %s", body.EntityId, err.Error()))
			return
		}

		fmt.Println("retrieved post by id: ", grpcResp)

		post := models.Post{
			PostId: ct.Id(grpcResp.PostId),
			Body:   ct.PostBody(grpcResp.PostBody),
			User: models.User{
				UserId:    ct.Id(grpcResp.User.UserId),
				Username:  ct.Username(grpcResp.User.Username),
				AvatarId:  ct.Id(grpcResp.User.Avatar),
				AvatarURL: grpcResp.User.AvatarUrl,
			},
			GroupId:         ct.Id(grpcResp.GroupId),
			Audience:        ct.Audience(grpcResp.Audience),
			CommentsCount:   int(grpcResp.CommentsCount),
			ReactionsCount:  int(grpcResp.ReactionsCount),
			LastCommentedAt: ct.GenDateTime(grpcResp.LastCommentedAt.AsTime()),
			CreatedAt:       ct.GenDateTime(grpcResp.CreatedAt.AsTime()),
			UpdatedAt:       ct.GenDateTime(grpcResp.UpdatedAt.AsTime()),
			LikedByUser:     grpcResp.LikedByUser,
			Image:           ct.Id(grpcResp.Image),
		}

		err = utils.WriteJSON(w, http.StatusOK, post)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send post by id")
			return
		}

	}
}

func (h *Handlers) createPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("createPost handler called")

		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type createPostReq struct {
			Body        string  `json:"post_body"`
			GroupId     int64   `json:"group_id"`
			Audience    string  `json:"audience"`
			AudienceIds []int64 `json:"audience_ids"`
			Image       int64   `json:"image"`
		}

		body, err := utils.JSON2Struct(&createPostReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.CreatePostReq{
			CreatorId: int64(claims.UserId),
			Body:      body.Body,
			GroupId:   body.GroupId,
			Audience:  body.Audience,
			AudienceIds: &common.UserIds{
				Values: body.AudienceIds,
			},
			Image: body.Image,
		}

		_, err = h.PostsService.CreatePost(r.Context(), &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to create post: %v", err.Error()))
			return
		}

	}
}

func (h *Handlers) deletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("deletePost handler called")

		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GenericReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.GenericReq{
			RequesterId: int64(claims.UserId),
			EntityId:    body.EntityId.Int64(),
		}

		_, err = h.PostsService.DeletePost(r.Context(), &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete post with id %v: %v", body.EntityId, err.Error()))
			return
		}

	}
}

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

func (h *Handlers) createComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("createComment handler called")

		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.CreateCommentReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.CreateCommentReq{
			CreatorId: int64(claims.UserId),
			ParentId:  body.ParentId.Int64(),
			Body:      body.Body.String(),
			Image:     int64(body.Image),
		}

		_, err = h.PostsService.CreateComment(r.Context(), &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to create comment: %v", err.Error()))
			return
		}

	}
}

func (h *Handlers) getCommentsByParentId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getCommentsByParentId handler called")

		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.EntityIdPaginatedReq{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.EntityIdPaginatedReq{
			RequesterId: claims.UserId,
			EntityId:    body.EntityId.Int64(),
			Limit:       body.Limit.Int32(),
			Offset:      body.Offset.Int32(),
		}

		grpcResp, err := h.PostsService.GetCommentsByParentId(ctx, &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to get comments for post id %v: %v: ", body.EntityId, err.Error()))
			return
		}

		fmt.Println("retrieved comments: ", grpcResp)

		var commentsResponse []models.Comment
		for _, c := range grpcResp.Comments {
			comment := models.Comment{
				CommentId: ct.Id(c.CommentId),
				ParentId:  ct.Id(c.ParentId),
				Body:      ct.CommentBody(c.Body),
				User: models.User{
					UserId:    ct.Id(c.User.UserId),
					Username:  ct.Username(c.User.Username),
					AvatarId:  ct.Id(c.User.Avatar),
					AvatarURL: c.User.AvatarUrl,
				},
				ReactionsCount: int(c.ReactionsCount),
				CreatedAt:      ct.GenDateTime(c.CreatedAt.AsTime()),
				UpdatedAt:      ct.GenDateTime(c.UpdatedAt.AsTime()),
				LikedByUser:    c.LikedByUser,
				Image:          ct.Id(c.Image),
			}
			commentsResponse = append(commentsResponse, comment)
		}

		err = utils.WriteJSON(w, http.StatusOK, commentsResponse)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to send comments for post %v : %v", body.EntityId, err.Error()))
			return
		}

	}
}

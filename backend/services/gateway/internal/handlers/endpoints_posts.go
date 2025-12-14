package handlers

import (
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/posts"
	ct "social-network/shared/go/customtypes"
	"time"
)

func (h *Handlers) getPublicFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getPublicFeed handler called")

		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}
		requesterId := int64(claims.UserId)

		grpcReq := posts.GenericPaginatedReq{
			RequesterId: requesterId,
			Limit:       10, //hardcoded for now, TODO make dynamic
			Offset:      0,  ////hardcoded for now, TODO make dynamic
		}

		grpcResp, err := h.PostsService.GetPublicFeed(r.Context(), &grpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to get public feed: "+err.Error())
			return
		}

		fmt.Println("retrieved public feed: ", grpcResp)

		type Post struct {
			PostId          ct.Id       `json:"post_id"`
			Body            ct.PostBody `json:"post_body"`
			CreatorId       ct.Id       `json:"post_creator_id"`
			CreatorUsername ct.Username `json:"post_creator_username"`
			CreatorAvater   ct.Id       `json:"post_creator_avatar,omitempty"`
			CommentsCount   int         `json:"comments_count"`
			ReactionsCount  int         `json:"reactions_count"`
			LastCommentedAt time.Time   `json:"last_created_at"`
			CreatedAt       time.Time   `json:"created_at"`
			UpdatedAt       time.Time   `json:"updated_at,omitempty"`
			LikedByUser     bool        `json:"liked_by_user"`
			Image           ct.Id       `json:"image,omitempty"`
		}

		var postsResponse []Post
		for _, p := range grpcResp.Posts {
			post := Post{
				PostId:          ct.Id(p.PostId),
				Body:            ct.PostBody(p.PostBody),
				CreatorId:       ct.Id(p.User.UserId),
				CreatorUsername: ct.Username(p.User.Username),
				CreatorAvater:   ct.Id(p.User.Avatar),
				CommentsCount:   int(p.CommentsCount),
				ReactionsCount:  int(p.ReactionsCount),
				LastCommentedAt: p.LastCommentedAt.AsTime(),
				CreatedAt:       p.CreatedAt.AsTime(),
				UpdatedAt:       p.UpdatedAt.AsTime(),
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

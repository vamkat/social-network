package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/media"
	"social-network/shared/gen-go/posts"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"
	"time"
)

func (h *Handlers) createEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tele.Info(ctx, "createEvent handler called")

		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type CreatePostJSONRequest struct {
			Title     ct.Title         `json:"event_title"`
			Body      ct.EventBody     `json:"event_body"`
			GroupId   ct.Id            `json:"group_id"`
			EventDate ct.EventDateTime `json:"event_date"`

			ImageName string `json:"image_name"`
			ImageSize int64  `json:"image_size"`
			ImageType string `json:"image_type"`
		}

		httpReq := CreatePostJSONRequest{}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		if err := ct.ValidateStruct(httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		var ImageId ct.Id
		var uploadURL string
		if httpReq.ImageSize != 0 {
			exp := time.Duration(10 * time.Minute).Seconds()
			mediaRes, err := h.MediaService.UploadImage(ctx, &media.UploadImageRequest{
				Filename:          httpReq.ImageName,
				MimeType:          httpReq.ImageType,
				SizeBytes:         httpReq.ImageSize,
				Visibility:        media.FileVisibility_PUBLIC,
				Variants:          []media.FileVariant{media.FileVariant_MEDIUM},
				ExpirationSeconds: int64(exp),
			})
			if err != nil {
				utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
				return
			}
			ImageId = ct.Id(mediaRes.FileId)
			uploadURL = mediaRes.GetUploadUrl()
		}

		grpcReq := posts.CreateEventReq{
			Title:     httpReq.Title.String(),
			Body:      httpReq.Body.String(),
			CreatorId: int64(claims.UserId),
			GroupId:   httpReq.GroupId.Int64(),
			ImageId:   ImageId.Int64(),
			EventDate: httpReq.EventDate.ToProto(),
		}

		_, err := h.PostsService.CreateEvent(ctx, &grpcReq)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("failed to create post: %v", err.Error()))
			return
		}
		type httpResponse struct {
			UserId    ct.Id
			FileId    ct.Id
			UploadUrl string
		}
		httpResp := httpResponse{
			UserId:    ct.Id(claims.UserId),
			FileId:    ImageId,
			UploadUrl: uploadURL}

		utils.WriteJSON(ctx, w, http.StatusOK, httpResp)

	}
}

func (h *Handlers) editEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tele.Info(ctx, "editEvent handler called")

		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type EditEventJSONRequest struct {
			EventId     ct.Id            `json:"event_id"`
			Title       ct.Title         `json:"event_title"`
			Body        ct.EventBody     `json:"event_body"`
			EventDate   ct.EventDateTime `json:"event_date"`
			DeleteImage bool             `json:"delete_image"`

			ImageName string `json:"image_name"`
			ImageSize int64  `json:"image_size"`
			ImageType string `json:"image_type"`
		}

		httpReq := EditEventJSONRequest{}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		if err := ct.ValidateStruct(httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		var ImageId ct.Id
		var uploadURL string
		if httpReq.ImageSize != 0 {
			exp := time.Duration(10 * time.Minute).Seconds()
			mediaRes, err := h.MediaService.UploadImage(ctx, &media.UploadImageRequest{
				Filename:          httpReq.ImageName,
				MimeType:          httpReq.ImageType,
				SizeBytes:         httpReq.ImageSize,
				Visibility:        media.FileVisibility_PUBLIC,
				Variants:          []media.FileVariant{media.FileVariant_MEDIUM},
				ExpirationSeconds: int64(exp),
			})
			if err != nil {
				utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
				return
			}
			ImageId = ct.Id(mediaRes.FileId)
			uploadURL = mediaRes.GetUploadUrl()
		}

		grpcReq := posts.EditEventReq{
			EventId:     int64(httpReq.EventId),
			RequesterId: int64(claims.UserId),
			Title:       httpReq.Title.String(),
			Body:        httpReq.Body.String(),
			EventDate:   httpReq.EventDate.ToProto(),
			ImageId:     ImageId.Int64(),
			DeleteImage: httpReq.DeleteImage,
		}

		_, err := h.PostsService.EditEvent(ctx, &grpcReq)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("failed to create post: %v", err.Error()))
			return
		}
		type httpResponse struct {
			UserId    ct.Id
			FileId    ct.Id
			UploadUrl string
		}
		httpResp := httpResponse{
			UserId:    ct.Id(claims.UserId),
			FileId:    ImageId,
			UploadUrl: uploadURL}

		utils.WriteJSON(ctx, w, http.StatusOK, httpResp)

	}
}

func (h *Handlers) deleteEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tele.Info(ctx, "deleteEvent handler called")

		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.GenericReq{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.GenericReq{
			RequesterId: int64(claims.UserId),
			EntityId:    body.EntityId.Int64(),
		}

		_, err = h.PostsService.DeleteEvent(ctx, &grpcReq)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("failed to delete post with id %v: %v", body.EntityId, err.Error()))
			return
		}

	}
}

func (h *Handlers) getEventsByGroupId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tele.Info(ctx, "getEventsByGroupId handler called")

		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.EntityIdPaginatedReq{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		grpcReq := posts.EntityIdPaginatedReq{
			RequesterId: claims.UserId,
			EntityId:    body.EntityId.Int64(),
			Limit:       body.Limit.Int32(),
			Offset:      body.Offset.Int32(),
		}

		grpcResp, err := h.PostsService.GetEventsByGroupId(ctx, &grpcReq)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("failed to get events for group %v: %v ", body.EntityId, err.Error()))
			return
		}

		eventsResponse := []models.Event{}
		for _, e := range grpcResp.Events {
			event := models.Event{
				EventId: ct.Id(e.EventId),
				Title:   ct.Title(e.Title),
				Body:    ct.EventBody(e.Body),
				User: models.User{
					UserId:    ct.Id(e.User.UserId),
					Username:  ct.Username(e.User.Username),
					AvatarId:  ct.Id(e.User.Avatar),
					AvatarURL: e.User.AvatarUrl,
				},
				GroupId:       ct.Id(e.GroupId),
				EventDate:     ct.EventDateTime(e.EventDate.AsTime()),
				GoingCount:    int(e.GoingCount),
				NotGoingCount: int(e.NotGoingCount),
				ImageId:       ct.Id(e.ImageId),
				ImageUrl:      e.ImageUrl,
				CreatedAt:     ct.GenDateTime(e.CreatedAt.AsTime()),
				UpdatedAt:     ct.GenDateTime(e.UpdatedAt.AsTime()),
				UserResponse:  &e.UserResponse.Value,
			}
			eventsResponse = append(eventsResponse, event)
		}

		err = utils.WriteJSON(ctx, w, http.StatusOK, eventsResponse)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("failed to get events for group %v: %v ", body.EntityId, err.Error()))
			return
		}

	}
}

func (s *Handlers) respondToEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.RespondToEventReq{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := posts.RespondToEventReq{
			EventId:     int64(body.EventId),
			ResponderId: claims.UserId,
			Going:       body.Going,
		}

		_, err = s.PostsService.RespondToEvent(ctx, &req)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("Could not respond to event with id %v: %v ", body.EventId, err.Error()))
			return
		}

		utils.WriteJSON(ctx, w, http.StatusOK, nil)
	}
}

func (s *Handlers) RemoveEventResponse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		body, err := utils.JSON2Struct(&models.RespondToEventReq{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := posts.GenericReq{
			RequesterId: claims.UserId,
			EntityId:    int64(body.EventId),
		}

		_, err = s.PostsService.RemoveEventResponse(ctx, &req)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("Could not remove response from event with id %v: %v ", body.EventId, err.Error()))
			return
		}

		utils.WriteJSON(ctx, w, http.StatusOK, nil)
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/media"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/mapping"
	tele "social-network/shared/go/telemetry"
)

func (h *Handlers) validateFileUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		type validateUploadReq struct {
			FileId    ct.Id `json:"file_id"`
			ReturnURL bool  `json:"return_url"`
		}
		httpReq := validateUploadReq{}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.MediaService.ValidateUpload(
			r.Context(),
			&media.ValidateUploadRequest{
				FileId:    httpReq.FileId.Int64(),
				ReturnUrl: httpReq.ReturnURL},
		)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		tele.Info(ctx, fmt.Sprintf("Gateway: Successfully validated file upload for FileId: %v", httpReq.FileId), "FileId", httpReq.FileId)

		type httpResp struct {
			DownloadUrl string `json:"download_url"`
		}

		if err := utils.WriteJSON(ctx, w, http.StatusCreated, &httpResp{DownloadUrl: res.GetDownloadUrl()}); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, fmt.Sprintf("failed to validate file %v", httpReq.FileId))
			return
		}
	}
}

func (h *Handlers) getImageUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type getImageReq struct {
			ImageId ct.Id          `json:"image_id"`
			Variant ct.FileVariant `json:"variant"`
		}

		httpReq := getImageReq{}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.MediaService.GetImage(r.Context(), &media.GetImageRequest{
			ImageId: httpReq.ImageId.Int64(),
			Variant: mapping.CtToPbFileVariant(httpReq.Variant),
		})
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		type httpResp struct {
			DownloadUrl string `json:"download_url"`
		}

		httpRes := &httpResp{DownloadUrl: res.DownloadUrl}
		if err := utils.WriteJSON(ctx, w, http.StatusOK, httpRes); err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to send registration ACK")
			return
		}
	}
}

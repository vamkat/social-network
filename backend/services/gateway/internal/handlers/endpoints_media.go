package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/media"
	ct "social-network/shared/go/customtypes"
)

func (h *Handlers) validateFileUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type validateUploadReq struct {
			FileId ct.Id `json:"file_id"`
		}
		httpReq := validateUploadReq{}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err := h.MediaService.ValidateUpload(r.Context(),
			&media.ValidateUploadRequest{FileId: httpReq.FileId.Int64()},
		)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		log.Printf("Gateway: Successfully validated file upload for FileId: %v", httpReq.FileId)

		if err := utils.WriteJSON(w, http.StatusCreated, nil); err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send registration ACK")
			return
		}
	}
}

func (h *Handlers) GetImageUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type getImageReq struct {
			ImageId ct.Id          `json:"image_id"`
			Variant ct.FileVariant `json:"variant"`
		}

		httpReq := getImageReq{}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.MediaService.GetImage(r.Context(), &media.GetImageRequest{
			ImageId: httpReq.ImageId.Int64(),
			Variant: ProtoToFileVariant(httpReq.Variant),
		})
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		type httpResp struct {
			DownloadUrl string `json:"download_url"`
		}

		httpRes := &httpResp{DownloadUrl: res.DownloadUrl}
		if err := utils.WriteJSON(w, http.StatusOK, httpRes); err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send registration ACK")
			return
		}
	}
}

func ProtoToFileVariant(pv ct.FileVariant) media.FileVariant {
	switch pv {
	case ct.ImgThumbnail:
		return media.FileVariant_THUMBNAIL
	case ct.ImgSmall:
		return media.FileVariant_SMALL
	case ct.ImgMedium:
		return media.FileVariant_MEDIUM
	case ct.ImgLarge:
		return media.FileVariant_LARGE
	case ct.Original:
		return media.FileVariant_ORIGINAL
	default:
		return media.FileVariant_IMG_VARIANT_UNSPECIFIED
	}
}

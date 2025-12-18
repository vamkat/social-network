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

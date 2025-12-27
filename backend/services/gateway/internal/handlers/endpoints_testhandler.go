package handlers

import (
	"net/http"
	"social-network/services/gateway/internal/utils"
	ct "social-network/shared/go/ct"
	tele "social-network/shared/go/telemetry"
)

func (h *Handlers) testHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tele.Info(ctx, "test handler called")

		err := utils.WriteJSON(ctx, w, http.StatusOK, map[string]string{
			"message": "this request id is: " + r.Context().Value(ct.ReqID).(string),
		})

		if err != nil {
			tele.Warn(ctx, "failed to send test ACK: ", "error", err.Error())
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to send logout ACK")
			return
		}
	}
}

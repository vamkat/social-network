package handlers

import (
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/utils"
	ct "social-network/shared/go/ct"
)

func (h *Handlers) testHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test handler called")

		err := utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "this request id is: " + r.Context().Value(ct.ReqID).(string),
		})

		if err != nil {
			fmt.Println("failed to send test ACK: ", err.Error())
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send logout ACK")
			return
		}
	}
}

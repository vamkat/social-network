package handlers

import (
	"fmt"
	"net/http"
	utils "social-network/shared/go/http-utils"
)

func (h *Handlers) NoMid(w http.ResponseWriter, r *http.Request) {
	if err := utils.WriteJSON(r.Context(), w, 200, "Hello"); err != nil {
		fmt.Println(err)
		return
	}
}
func (h *Handlers) WithMid(w http.ResponseWriter, r *http.Request) {
	if err := utils.WriteJSON(r.Context(), w, 200, "Hello"); err != nil {
		fmt.Println(err)
		return
	}
}

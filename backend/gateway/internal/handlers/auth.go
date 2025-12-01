package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/gateway/internal/security"
	"social-network/gateway/internal/utils"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"
	"time"
)

func (h *Handlers) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login handler called")
		//READ REQUEST BODY
		type loginHttpRequest struct {
			Identifier ct.Identifier `json:"identifier"`
			Password   ct.Password   `json:"password"`
		}

		httpReq := loginHttpRequest{}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&httpReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//VALIDATE INPUT
		if err := ct.ValidateStruct(httpReq); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}

		//MAKE GRPC REQUEST
		gRpcReq := users.LoginRequest{
			Identifier: httpReq.Identifier.String(),
			Password:   httpReq.Password.String(),
		}

		user, err := h.Services.Users.LoginUser(r.Context(), &gRpcReq)
		if err != nil {
			//TODO: distinguish error types
			utils.ErrorJSON(w, http.StatusInternalServerError, "login failed")
			return
		}

		//PREPARE SUCCESS RESPONSE
		now := time.Now().Unix()
		exp := time.Now().AddDate(0, 6, 0).Unix() // six months from now

		claims := security.Claims{
			UserId: int64(user.UserId),
			Iat:    now,
			Exp:    exp,
		}

		token, err := security.CreateToken(claims)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "token generation failed")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    token,
			Path:     "/",
			Expires:  time.Unix(exp, 0),
			HttpOnly: true,
			Secure:   false, //TODO: set to true in production
			SameSite: http.SameSiteLaxMode,
		})

		//SEND RESPONSE
		err = utils.WriteJSON(w, http.StatusOK, user)
		if err != nil {
			utils.ErrorJSON(w, http.StatusUnauthorized, "failed to send login ACK")
			return
		}
	}
}

func (h *Handlers) registerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("register handler called, with: ", r.Body)
		//READ REQUEST BODY
		type registerHttpRequest struct {
			Username    ct.Username    `json:"username,omitempty"`
			FirstName   ct.Name        `json:"first_name,omitempty"`
			LastName    ct.Name        `json:"last_name,omitempty"`
			DateOfBirth ct.DateOfBirth `json:"date_of_birth,omitempty"`
			Avatar      string         `json:"avatar,omitempty"`
			About       ct.About       `json:"about,omitempty"`
			Public      bool           `json:"public,omitempty"`
			Email       ct.Email       `json:"email,omitempty"`
			Password    ct.Password    `json:"password,omitempty"`
		}

		httpReq := registerHttpRequest{}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&httpReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := ct.ValidateStruct(httpReq); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		}

		//VALIDATE INPUT
		// if httpReq.Username == "" || httpReq.Email == "" || httpReq.Password == "" {
		// 	utils.ErrorJSON(w, http.StatusBadRequest, "missing required fields")
		// 	return
		// }

		//MAKE GRPC REQUEST
		gRpcReq := users.RegisterUserRequest{
			Username:    string(httpReq.Username),
			FirstName:   string(httpReq.FirstName),
			LastName:    string(httpReq.LastName),
			DateOfBirth: httpReq.DateOfBirth.ToProto(),
			Avatar:      httpReq.Avatar,
			About:       string(httpReq.About),
			Public:      httpReq.Public,
			Email:       string(httpReq.Email),
			Password:    string(httpReq.Password),
		}

		_, err := h.Services.Users.RegisterUser(r.Context(), &gRpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "registration failed")
			return
		}

		//SEND RESPONSE
		if err := utils.WriteJSON(w, http.StatusCreated, "registered successfully"); err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send registration ACK")
			return
		}
	}
}

func (h *Handlers) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("logout handler called")
		//CLEAR COOKIE
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   false, //TODO: set to true in production
			SameSite: http.SameSiteLaxMode,
		})

		//SEND RESPONSE
		if err := utils.WriteJSON(w, http.StatusOK, "logged out successfully"); err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send logout ACK")
			return
		}
	}
}

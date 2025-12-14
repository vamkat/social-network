package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"
	"time"

	"github.com/minio/minio-go/v7"
)

func (h *Handlers) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login handler called with request_id:", r.Context().Value(ct.ReqID), " userid:", r.Context().Value(ct.UserId), " tracerId:", r.Context().Value(ct.TraceId))
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
			utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}

		httpReq.Password, err = httpReq.Password.Hash()
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "could not hash password")
			return
		}

		//VALIDATE INPUT
		if err := ct.ValidateStruct(httpReq); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		fmt.Println("password", httpReq.Password.String())
		//MAKE GRPC REQUEST
		gRpcReq := users.LoginRequest{
			Identifier: httpReq.Identifier.String(),
			Password:   httpReq.Password.String(),
		}

		fmt.Println(httpReq.Password.String())

		resp, err := h.UsersService.LoginUser(r.Context(), &gRpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		//PREPARE SUCCESS RESPONSE
		now := time.Now().Unix()
		exp := time.Now().AddDate(0, 6, 0).Unix() // six months from now

		claims := security.Claims{
			UserId: resp.UserId,
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

		type httpResponse struct {
			UserId ct.Id
		}

		httpResp := httpResponse{UserId: ct.Id(resp.UserId)}

		//SEND RESPONSE
		err = utils.WriteJSON(w, http.StatusCreated, httpResp)
		if err != nil {
			utils.ErrorJSON(w, http.StatusUnauthorized, "failed to send login ACK")
			return
		}
	}
}

func (h *Handlers) registerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user already logged in
		cookie, _ := r.Cookie("jwt")
		if cookie != nil {
			_, err := security.ParseAndValidate(cookie.Value)
			if err == nil {
				utils.ErrorJSON(w, http.StatusForbidden, "Already logged in. Log out to register.")
				return
			}
		}

		fmt.Println("register handler called, with: ", r.Body)
		//READ REQUEST BODY
		type registerHttpRequest struct {
			Username    ct.Username    `json:"username,omitempty" validate:"nullable"`
			FirstName   ct.Name        `json:"first_name,omitempty"`
			LastName    ct.Name        `json:"last_name,omitempty"`
			DateOfBirth ct.DateOfBirth `json:"date_of_birth,omitempty"`
			Avatar      ct.Id          `json:"avatar,omitempty" validate:"nullable"`
			About       ct.About       `json:"about,omitempty" validate:"nullable"`
			Public      bool           `json:"public,omitempty"`
			Email       ct.Email       `json:"email,omitempty"`
			Password    ct.Password    `json:"password,omitempty"`
		}

		// Parse up to 20MB file+fields
		err := r.ParseMultipartForm(20 << 20)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "invalid multipart form: "+err.Error())
			return
		}

		dob, err := ct.ParseDateOfBirth(r.FormValue("date_of_birth"))
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("Pass on register", ct.Password(r.FormValue("password")))

		password, err := ct.Password(r.FormValue("password")).Hash()
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "could not hash password")
			return
		}

		// Extract form fields
		httpReq := registerHttpRequest{
			Username:    ct.Username(r.FormValue("username")),
			FirstName:   ct.Name(r.FormValue("first_name")),
			LastName:    ct.Name(r.FormValue("last_name")),
			DateOfBirth: dob,
			About:       ct.About(r.FormValue("about")),
			Public:      r.FormValue("public") == "true",
			Email:       ct.Email(r.FormValue("email")),
			Password:    password,
		}

		if err := ct.ValidateStruct(httpReq); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		var uploadInfo minio.UploadInfo
		// Extract and upload avatar (optional)
		file, header, err := r.FormFile("avatar")
		if err == http.ErrMissingFile {
			// no avatar uploaded → fine
		} else if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "avatar upload error: "+err.Error())
			return
		} else {
			defer file.Close()
			_, err := utils.CheckImage(file, header)
			if err != nil {
				utils.ErrorJSON(w, http.StatusBadRequest, "avatar upload error: "+err.Error())
				return
			}

			//deprecated
			// uploadInfo, err = remoteservices.UploadToMinIO(r.Context(), h.MinIOClient, file, header, "images", fileType)
			// if err != nil {
			// 	utils.ErrorJSON(w, http.StatusInternalServerError, "failed to upload avatar: "+err.Error())
			// 	return
			// }
		}

		_ = uploadInfo

		//MAKE GRPC REQUEST
		gRpcReq := users.RegisterUserRequest{
			Username:    string(httpReq.Username),
			FirstName:   string(httpReq.FirstName),
			LastName:    string(httpReq.LastName),
			DateOfBirth: httpReq.DateOfBirth.ToProto(),
			About:       string(httpReq.About),
			Public:      httpReq.Public,
			Email:       string(httpReq.Email),
			Password:    string(httpReq.Password),
			// Avatar:      httpReq.Avatar,
		}

		resp, err := h.UsersService.RegisterUser(r.Context(), &gRpcReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		//PREPARE SUCCESS RESPONSE
		now := time.Now().Unix()
		exp := time.Now().AddDate(0, 6, 0).Unix() // six months from now

		claims := security.Claims{
			UserId: resp.UserId,
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

		type httpResponse struct {
			UserId ct.Id
		}
		httpResp := httpResponse{UserId: ct.Id(resp.UserId)}

		//SEND RESPONSE
		if err := utils.WriteJSON(w, http.StatusCreated, httpResp); err != nil {
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

// Returns status ok if passed Auth
func (h *Handlers) authStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, "user is logged in")
	}
}

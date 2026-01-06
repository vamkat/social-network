package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/shared/gen-go/media"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/ct"
	utils "social-network/shared/go/http-utils"
	"social-network/shared/go/jwt"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"
	"strings"
	"time"
)

func (h *Handlers) getUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tele.Info(ctx, "getUserProfile handler called")

		pathParts := strings.Split(r.URL.Path, "/")
		if pathParts[len(pathParts)-1] == "" {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "missing user_id in URL path")
			return
		}

		userId, err := ct.DecryptId(pathParts[len(pathParts)-1])
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "invalid user_id query param")
			return
		}

		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}
		requesterId := int64(claims.UserId)

		grpcReq := users.GetUserProfileRequest{
			UserId:      userId.Int64(),
			RequesterId: requesterId,
		}

		grpcResp, err := h.UsersService.GetUserProfile(r.Context(), &grpcReq)
		if err != nil {
			utils.ReturnHttpError(ctx, w, err)
			//utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to get user info: "+err.Error())
			return
		}

		tele.Info(ctx, "retrieved user profile. @1", "grpcResp", grpcResp)

		type userProfile struct {
			UserId            ct.Id          `json:"user_id"`
			Username          ct.Username    `json:"username"`
			FirstName         ct.Name        `json:"first_name"`
			LastName          ct.Name        `json:"last_name"`
			DateOfBirth       ct.DateOfBirth `json:"date_of_birth"`
			Avatar            ct.Id          `json:"avatar,omitempty"`
			AvatarURL         string         `json:"avatar_url,omitempty"`
			About             ct.About       `json:"about,omitempty"`
			Public            bool           `json:"public"`
			CreatedAt         time.Time      `json:"created_at"`
			Email             fmt.Stringer   `json:"email"`
			FollowersCount    int64          `json:"followers_count"`
			FollowingCount    int64          `json:"following_count"`
			GroupsCount       int64          `json:"groups_count"`
			OwnedGroupsCount  int64          `json:"owned_groups_count"`
			ViewerIsFollowing bool           `json:"viewer_is_following"`
			OwnProfile        bool           `json:"own_profile"`
			IsPending         bool           `json:"is_pending"`
		}

		userProfileResponse := userProfile{
			UserId:            ct.Id(grpcResp.UserId),
			Username:          ct.Username(grpcResp.Username),
			FirstName:         ct.Name(grpcResp.FirstName),
			LastName:          ct.Name(grpcResp.LastName),
			DateOfBirth:       ct.DateOfBirth(grpcResp.DateOfBirth.AsTime()),
			Avatar:            ct.Id(grpcResp.Avatar),
			AvatarURL:         grpcResp.AvatarUrl,
			About:             ct.About(grpcResp.About),
			Public:            grpcResp.Public,
			CreatedAt:         grpcResp.CreatedAt.AsTime(),
			Email:             ct.Email(grpcResp.Email),
			FollowersCount:    grpcResp.FollowersCount,
			FollowingCount:    grpcResp.FollowingCount,
			GroupsCount:       grpcResp.GroupsCount,
			OwnedGroupsCount:  grpcResp.OwnedGroupsCount,
			ViewerIsFollowing: grpcResp.ViewerIsFollowing,
			OwnProfile:        grpcResp.OwnProfile,
			IsPending:         grpcResp.IsPending,
		}

		tele.Info(ctx, "transformed profile struct. @1", "response", userProfileResponse)

		err = utils.WriteJSON(ctx, w, http.StatusOK, userProfileResponse)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to send user info")
			return
		}

	}
}

func (s *Handlers) searchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type reqBody struct {
			Query string `json:"query"`
			Limit int32  `json:"limit"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.UserSearchRequest{
			SearchTerm: body.Query,
			Limit:      body.Limit,
		}

		grpcResp, err := s.UsersService.SearchUsers(ctx, req)
		if err != nil {
			utils.ReturnHttpError(ctx, w, err)
			//utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "Could not search users: "+err.Error())
			return
		}

		resp := models.Users{
			Users: make([]models.User, 0, len(grpcResp.Users)),
		}

		for _, user := range grpcResp.Users {
			newUser := models.User{
				UserId:    ct.Id(user.UserId),
				Username:  ct.Username(user.Username),
				AvatarId:  ct.Id(user.Avatar),
				AvatarURL: user.AvatarUrl,
			}

			resp.Users = append(resp.Users, newUser)
		}

		utils.WriteJSON(ctx, w, http.StatusOK, resp)
	}
}

func (s *Handlers) updateProfilePrivacy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			Public bool `json:"public"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(ctx, w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := &users.UpdateProfilePrivacyRequest{
			UserId: claims.UserId,
			Public: body.Public,
		}

		_, err = s.UsersService.UpdateProfilePrivacy(ctx, req)
		if err != nil {
			utils.ReturnHttpError(ctx, w, err)
			//utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "Could not update privacy: "+err.Error())
			return
		}

		utils.WriteJSON(ctx, w, http.StatusOK, nil)
	}
}

func (s *Handlers) updateUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type UpdateProfileJSONRequest struct {
			Username    ct.Username    `json:"username"`
			FirstName   ct.Name        `json:"first_name"`
			LastName    ct.Name        `json:"last_name"`
			DateOfBirth ct.DateOfBirth `json:"date_of_birth"`
			About       ct.About       `json:"about" validate:"nullable"`

			AvatarName string `json:"avatar_name"`
			AvatarSize int64  `json:"avatar_size"`
			AvatarType string `json:"avatar_type"`
		}

		httpReq := UpdateProfileJSONRequest{}

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

		var AvatarId ct.Id
		var uploadURL string
		if httpReq.AvatarSize != 0 {
			exp := time.Duration(10 * time.Minute).Seconds()
			mediaRes, err := s.MediaService.UploadImage(r.Context(), &media.UploadImageRequest{
				Filename:          httpReq.AvatarName,
				MimeType:          httpReq.AvatarType,
				SizeBytes:         httpReq.AvatarSize,
				Visibility:        media.FileVisibility_PUBLIC,
				Variants:          []media.FileVariant{media.FileVariant_THUMBNAIL},
				ExpirationSeconds: int64(exp),
			})
			if err != nil {
				utils.ErrorJSON(ctx, w, http.StatusInternalServerError, err.Error())
				return
			}
			AvatarId = ct.Id(mediaRes.FileId)
			uploadURL = mediaRes.GetUploadUrl()
		}

		//MAKE GRPC REQUEST
		grpcRequest := &users.UpdateProfileRequest{
			UserId:      claims.UserId,
			Username:    httpReq.Username.String(),
			FirstName:   httpReq.FirstName.String(),
			LastName:    httpReq.LastName.String(),
			DateOfBirth: httpReq.DateOfBirth.ToProto(),
			Avatar:      AvatarId.Int64(),
			About:       httpReq.About.String(),
		}

		grpcResp, err := s.UsersService.UpdateUserProfile(ctx, grpcRequest)
		if err != nil {
			utils.ReturnHttpError(ctx, w, err)
			//utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "Could not update profile: "+err.Error())
			return
		}

		type httpResponse struct {
			UserId    ct.Id
			FileId    ct.Id
			UploadUrl string
		}
		httpResp := httpResponse{
			UserId:    ct.Id(grpcResp.UserId),
			FileId:    AvatarId,
			UploadUrl: uploadURL}

		utils.WriteJSON(ctx, w, http.StatusOK, httpResp)
	}
}

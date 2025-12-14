package handlers

import (
	"net/http"
	"social-network/services/gateway/internal/security"
	"social-network/services/gateway/internal/utils"
	"social-network/shared/gen-go/users"
	ct "social-network/shared/go/customtypes"
)

func (s *Handlers) CreateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type createGroupData struct {
			GroupTitle       string `json:"group_title"`
			GroupDescription string `json:"group_description"`
			GroupImage       string `json:"group_image,omitempty"`
		}

		createGroupDataRequest, err := utils.JSON2Struct(&createGroupData{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		createGroupRequest := users.CreateGroupRequest{
			OwnerId:          claims.UserId,
			GroupTitle:       createGroupDataRequest.GroupTitle,
			GroupDescription: createGroupDataRequest.GroupDescription,
			GroupImage:       createGroupDataRequest.GroupImage,
		}

		groupId, err := s.UsersService.CreateGroup(ctx, &createGroupRequest)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not create group: "+err.Error())
			return
		}

		type createGroupDataResponse struct {
			GroupId ct.Id `json:"group_id"`
		}

		resp := createGroupDataResponse{
			GroupId: ct.Id(groupId.Value),
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (s *Handlers) GetAllGroupsPaginated() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := utils.GetValue[security.Claims](r, ct.ClaimsKey)
		if !ok {
			panic(1)
		}

		type reqBody struct {
			Limit  int32 `json:"limit"`
			Offset int32 `json:"offset"`
		}

		body, err := utils.JSON2Struct(&reqBody{}, r)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad JSON data received")
			return
		}

		req := users.Pagination{
			UserId: claims.UserId,
			Limit:  body.Limit,
			Offset: body.Offset,
		}

		out, err := s.UsersService.GetAllGroupsPaginated(ctx, &req)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not fetch groups: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, out)
	}
}

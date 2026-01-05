package application

import (
	"context"
	"fmt"
	ds "social-network/services/users/internal/db/dbservice"
	"social-network/shared/gen-go/media"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
)

func (s *Application) GetAllGroupsPaginated(ctx context.Context, req models.Pagination) ([]models.Group, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return []models.Group{}, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	//paginated (sorting by most members first)
	rows, err := s.db.GetAllGroups(ctx, ds.GetAllGroupsParams{
		Offset: req.Offset.Int32(),
		Limit:  req.Limit.Int32(),
	})
	if err != nil {
		return nil, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}

	groups := make([]models.Group, 0, len(rows))
	var imageIds ct.Ids

	for _, r := range rows {
		userInfo, err := s.userInRelationToGroup(ctx, models.GeneralGroupReq{
			GroupId: ct.Id(r.ID),
			UserId:  req.UserId,
		})
		if err != nil {
			return nil, err
		}

		groups = append(groups, models.Group{
			GroupId:          ct.Id(r.ID),
			GroupOwnerId:     ct.Id(r.GroupOwner),
			GroupTitle:       ct.Title(r.GroupTitle),
			GroupDescription: ct.About(r.GroupDescription),
			GroupImage:       ct.Id(r.GroupImageID),
			MembersCount:     r.MembersCount,
			IsMember:         userInfo.isMember,
			IsOwner:          userInfo.isOwner,
			IsPending:        userInfo.isPending,
		})
		if r.GroupImageID > 0 {
			imageIds = append(imageIds, ct.Id(r.GroupImageID))
		}

	}

	//get image urls
	if len(imageIds) > 0 {
		imageMap, _, err := s.mediaRetriever.GetImages(ctx, imageIds, media.FileVariant(1)) //TODO delete failed
		if err != nil {
			return nil, ce.Wrap(nil, err, input).WithPublic("error retrieving group image")
		}
		for i := range groups {
			groups[i].GroupImageURL = imageMap[groups[i].GroupImage.Int64()]
		}
	}

	return groups, nil
}

func (s *Application) GetUserGroupsPaginated(ctx context.Context, req models.Pagination) ([]models.Group, error) {
	input := fmt.Sprintf("%#v", req)

	//paginated (joined latest first)
	if err := ct.ValidateStruct(req); err != nil {
		return []models.Group{}, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	rows, err := s.db.GetUserGroups(ctx, ds.GetUserGroupsParams{
		GroupOwner: req.UserId.Int64(),
		Limit:      req.Limit.Int32(),
		Offset:     req.Offset.Int32(),
	})
	if err != nil {
		return nil, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}

	groups := make([]models.Group, 0, len(rows))
	var imageIds ct.Ids

	for _, r := range rows {
		isPending, err := s.isGroupMembershipPending(ctx, models.GeneralGroupReq{
			GroupId: ct.Id(r.GroupID),
			UserId:  req.UserId,
		})
		if err != nil {
			return nil, ce.Wrap(nil, err)
		}
		groups = append(groups, models.Group{
			GroupId:          ct.Id(r.GroupID),
			GroupOwnerId:     ct.Id(r.GroupOwner),
			GroupTitle:       ct.Title(r.GroupTitle),
			GroupDescription: ct.About(r.GroupDescription),
			GroupImage:       ct.Id(r.GroupImageID),
			MembersCount:     r.MembersCount,
			IsMember:         r.IsMember,
			IsOwner:          r.IsOwner,
			IsPending:        isPending,
		})
		if r.GroupImageID > 0 {
			imageIds = append(imageIds, ct.Id(r.GroupImageID))
		}
	}

	//get image urls
	if len(imageIds) > 0 {
		imageMap, _, err := s.mediaRetriever.GetImages(ctx, imageIds, media.FileVariant(1)) //TODO delete failed
		if err != nil {
			return nil, ce.Wrap(nil, err, input).WithPublic("error retrieving images")
		}
		for i := range groups {
			groups[i].GroupImageURL = imageMap[groups[i].GroupImage.Int64()]
		}
	}

	return groups, nil
}

func (s *Application) GetGroupInfo(ctx context.Context, req models.GeneralGroupReq) (models.Group, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return models.Group{}, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	row, err := s.db.GetGroupInfo(ctx, req.GroupId.Int64())
	if err != nil {
		return models.Group{}, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	group := models.Group{
		GroupId:          ct.Id(row.ID),
		GroupOwnerId:     ct.Id(row.GroupOwner),
		GroupTitle:       ct.Title(row.GroupTitle),
		GroupDescription: ct.About(row.GroupDescription),
		GroupImage:       ct.Id(row.GroupImageID),
		MembersCount:     row.MembersCount,
	}
	userInfo, err := s.userInRelationToGroup(ctx, models.GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.UserId,
	})
	if err != nil {
		return models.Group{}, ce.Wrap(nil, err)
	}
	group.IsMember = userInfo.isMember
	group.IsOwner = userInfo.isOwner
	group.IsPending = userInfo.isPending

	if group.GroupImage > 0 {
		imageUrl, err := s.mediaRetriever.GetImage(ctx, group.GroupImage.Int64(), media.FileVariant(1))
		if err != nil {
			return models.Group{}, ce.Wrap(nil, err, input).WithPublic("error retrieving group image")
		}

		group.GroupImageURL = imageUrl
	}

	return group, nil

}

func (s *Application) GetGroupMembers(ctx context.Context, req models.GroupMembersReq) ([]models.GroupUser, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return []models.GroupUser{}, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	//check request comes from member
	isMember, err := s.IsGroupMember(ctx, models.GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.UserId,
	})
	if err != nil {
		return nil, ce.Wrap(nil, err)
	}
	if !isMember {
		return nil, ce.New(ce.ErrPermissionDenied, fmt.Errorf("user %v is not a member of group %v", req.UserId, req.GroupId), input).WithPublic("permission denied")
	}

	//paginated (newest first)
	rows, err := s.db.GetGroupMembers(ctx, ds.GetGroupMembersParams{
		GroupID: req.GroupId.Int64(),
		Limit:   req.Limit.Int32(),
		Offset:  req.Offset.Int32(),
	})
	if err != nil {
		return nil, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	members := make([]models.GroupUser, 0, len(rows))
	var imageIds ct.Ids

	for _, r := range rows {
		var role string
		if r.Role.Valid {
			role = string(r.Role.GroupRole)
		}

		members = append(members, models.GroupUser{
			UserId:    ct.Id(r.ID),
			Username:  ct.Username(r.Username),
			AvatarId:  ct.Id(r.AvatarID),
			GroupRole: role,
		})
		if r.AvatarID > 0 {
			imageIds = append(imageIds, ct.Id(r.AvatarID))
		}
	}

	//get avatar urls
	if len(imageIds) > 0 {
		avatarMap, _, err := s.mediaRetriever.GetImages(ctx, imageIds, media.FileVariant(1)) //TODO delete failed
		if err != nil {
			return []models.GroupUser{}, ce.Wrap(nil, err, input).WithPublic("error retrieving images")
		}
		for i := range members {
			members[i].AvatarUrl = avatarMap[members[i].AvatarId.Int64()]
		}
	}

	return members, nil
}

// intentionally doesn't return avatar urls to avoid redundant calls
func (s *Application) GetAllGroupMemberIds(ctx context.Context, req models.GroupId) (ct.Ids, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateBatch(ct.Id(req)); err != nil {
		return nil, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	//paginated (newest first)
	rows, err := s.db.GetAllGroupMemberIds(ctx, ds.GetAllGroupMemberIdsParams{
		GroupID: ct.Id(req).Int64(),
	})
	if err != nil {
		return nil, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	members := make([]ct.Id, 0, len(rows))

	for _, r := range rows {
		members = append(members, ct.Id(r.ID))
	}

	return members, nil
}

func (s *Application) SearchGroups(ctx context.Context, req models.GroupSearchReq) ([]models.Group, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return []models.Group{}, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	//weighted (title more important than description)
	//paginated (most members first)
	rows, err := s.db.SearchGroups(ctx, ds.SearchGroupsParams{
		Query:  req.SearchTerm.String(),
		UserID: req.UserId.Int64(),
		Limit:  req.Limit.Int32(),
		Offset: req.Offset.Int32(),
	})
	if err != nil {
		return []models.Group{}, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}

	if len(rows) == 0 {
		return []models.Group{}, nil
	}

	groups := make([]models.Group, 0, len(rows))
	var imageIds ct.Ids

	for _, r := range rows {
		isPending, err := s.isGroupMembershipPending(ctx, models.GeneralGroupReq{
			GroupId: ct.Id(r.ID),
			UserId:  req.UserId,
		})
		if err != nil {
			return nil, ce.Wrap(nil, err)
		}
		groups = append(groups, models.Group{
			GroupId:          ct.Id(r.ID),
			GroupOwnerId:     ct.Id(r.GroupOwner),
			GroupTitle:       ct.Title(r.GroupTitle),
			GroupDescription: ct.About(r.GroupDescription),
			GroupImage:       ct.Id(r.GroupImageID),
			MembersCount:     r.MembersCount,
			IsMember:         r.IsMember,
			IsOwner:          r.IsOwner,
			IsPending:        isPending,
		})
		if r.GroupImageID > 0 {
			imageIds = append(imageIds, ct.Id(r.GroupImageID))
		}
	}

	//get image urls
	if len(imageIds) > 0 {
		imageMap, _, err := s.mediaRetriever.GetImages(ctx, imageIds, media.FileVariant(1)) //TODO delete failed
		if err != nil {
			return nil, ce.Wrap(nil, err, input).WithPublic("error retrieving images")
		}
		for i := range groups {
			groups[i].GroupImageURL = imageMap[groups[i].GroupImage.Int64()]
		}
	}

	return groups, nil

}

func (s *Application) InviteToGroup(ctx context.Context, req models.InviteToGroupReq) error {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	//check request comes from member
	isMember, err := s.IsGroupMember(ctx, models.GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.InviterId,
	})
	if err != nil {
		return ce.Wrap(nil, err)
	}
	if !isMember {
		return ce.New(ce.ErrPermissionDenied, fmt.Errorf("user %v is not a member of group %v", req.InviterId, req.GroupId), input).WithPublic("permission denied")
	}

	err = s.db.SendGroupInvites(ctx, ds.SendGroupInvitesParams{
		GroupID:     req.GroupId.Int64(),
		SenderID:    req.InviterId.Int64(),
		ReceiverIDs: req.InvitedIds.Int64(),
	})
	if err != nil {
		return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	//create notification (waiting for batch notification)
	// inviter, err := s.GetBasicUserInfo(ctx, req.InviterId)
	// if err != nil {
	// 	//WHAT DO DO WITH ERROR HERE?
	// }
	// group, err := s.db.GetGroupBasicInfo(ctx, req.GroupId.Int64())
	// if err != nil {
	// 	//WHAT DO DO WITH ERROR HERE?
	// }
	// err = s.clients.CreateGroupInvite(ctx, req.InvitedId.Int64(), req.InviterId.Int64(), req.GroupId.Int64(), group.GroupTitle, inviter.Username.String())
	// if err != nil {
	// 	//WHAT DO DO WITH ERROR HERE?
	// }
	return nil
}

// SKIP GRPC FOR NOW
// func (s *Application) CancelInviteToGroup(ctx context.Context, req models.InviteToGroupReq) error {
// 	if err := ct.ValidateStruct(req); err != nil {
// 		return err
// 	}
// 	err := s.db.CancelGroupInvite(ctx, ds.CancelGroupInviteParams{
// 		GroupID:    req.GroupId.Int64(),
// 		ReceiverID: req.InvitedId.Int64(),
// 		SenderID:   req.InviterId.Int64(),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	//TODO REMOVE NOTIFICATION EVENT
// 	return nil
// }

func (s *Application) RequestJoinGroup(ctx context.Context, req models.GroupJoinRequest) error {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	err := s.db.SendGroupJoinRequest(ctx, ds.SendGroupJoinRequestParams{
		GroupID: req.GroupId.Int64(),
		UserID:  req.RequesterId.Int64(),
	})

	if err != nil {
		return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	//create notification
	requester, err := s.GetBasicUserInfo(ctx, req.RequesterId)
	if err != nil {
		//WHAT DO DO WITH ERROR HERE?
	}
	group, err := s.db.GetGroupBasicInfo(ctx, req.GroupId.Int64())
	if err != nil {
		//WHAT DO DO WITH ERROR HERE?
	}
	err = s.clients.CreateGroupJoinRequest(ctx, group.GroupOwner, int64(req.RequesterId.Int64()), req.GroupId.Int64(), group.GroupTitle, requester.Username.String())
	if err != nil {
		//WHAT DO DO WITH ERROR HERE?
	}
	return nil
}

// SKIP GRPC FOR NOW
func (s *Application) CancelJoinGroupRequest(ctx context.Context, req models.GroupJoinRequest) error {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	err := s.db.CancelGroupJoinRequest(ctx, ds.CancelGroupJoinRequestParams{
		GroupID: req.GroupId.Int64(),
		UserID:  req.RequesterId.Int64(),
	})
	if err != nil {
		return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	//TODO REMOVE NOTIFICATION EVENT
	return nil
}

func (s *Application) RespondToGroupInvite(ctx context.Context, req models.HandleGroupInviteRequest) error {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	inviterId, err := s.db.GetGroupInviterId(ctx, ds.GetGroupInviterIdParams{
		GroupID:    req.GroupId.Int64(),
		ReceiverID: req.InvitedId.Int64(),
	})
	if err != nil {
		return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}

	if req.Accepted {

		err := s.db.AcceptGroupInvite(ctx, ds.AcceptGroupInviteParams{
			GroupID:    req.GroupId.Int64(),
			ReceiverID: req.InvitedId.Int64(),
		})
		if err != nil {
			return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)

		}
		//create notification
		invited, err := s.GetBasicUserInfo(ctx, req.InvitedId)
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		group, err := s.db.GetGroupBasicInfo(ctx, req.GroupId.Int64())
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		err = s.clients.CreateGroupInviteAccepted(ctx, int64(req.InvitedId.Int64()), inviterId, req.GroupId.Int64(), group.GroupTitle, invited.Username.String())
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}

		// err = s.clients.AddMembersToGroupConversation(ctx, req.GroupId.Int64(), []int64{req.InvitedId.Int64()})
		// if err != nil {
		// 	tele.Info("could not add member to group conversation:", err)
		// }

	} else {
		err := s.db.DeclineGroupInvite(ctx, ds.DeclineGroupInviteParams{
			GroupID:    req.GroupId.Int64(),
			ReceiverID: req.InvitedId.Int64(),
		})
		if err != nil {
			return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
		}
		//create notification
		invited, err := s.GetBasicUserInfo(ctx, req.InvitedId)
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		group, err := s.db.GetGroupBasicInfo(ctx, req.GroupId.Int64())
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		err = s.clients.CreateGroupInviteRejected(ctx, int64(req.InvitedId.Int64()), inviterId, req.GroupId.Int64(), group.GroupTitle, invited.Username.String())
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
	}
	return nil
}

func (s *Application) HandleGroupJoinRequest(ctx context.Context, req models.HandleJoinRequest) error {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	isOwner, err := s.isGroupOwner(ctx, models.GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.OwnerId,
	})
	if err != nil {
		return ce.Wrap(nil, err)
	}
	if !isOwner {
		return ce.New(ce.ErrPermissionDenied, fmt.Errorf("user %v is not the owner of group %v", req.OwnerId, req.GroupId), input).WithPublic("permission denied")
	}

	if req.Accepted {
		err = s.db.AcceptGroupJoinRequest(ctx, ds.AcceptGroupJoinRequestParams{
			GroupID: req.GroupId.Int64(),
			UserID:  req.RequesterId.Int64(),
		})
		if err != nil {
			return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
		}

		//create notification
		group, err := s.db.GetGroupBasicInfo(ctx, req.GroupId.Int64())
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		err = s.clients.CreateGroupJoinRequestAccepted(ctx, int64(req.RequesterId.Int64()), group.GroupOwner, req.GroupId.Int64(), group.GroupTitle)
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}

		// err = s.clients.AddMembersToGroupConversation(ctx, req.GroupId.Int64(), []int64{req.RequesterId.Int64()})
		// if err != nil {
		// 	tele.Info("could not add member to group conversation:", err)
		// }

	} else {
		err = s.db.RejectGroupJoinRequest(ctx, ds.RejectGroupJoinRequestParams{
			GroupID: req.GroupId.Int64(),
			UserID:  req.RequesterId.Int64(),
		})
		if err != nil {
			return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
		}
		//create notification
		group, err := s.db.GetGroupBasicInfo(ctx, req.GroupId.Int64())
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
		err = s.clients.CreateGroupJoinRequestRejected(ctx, int64(req.RequesterId.Int64()), group.GroupOwner, req.GroupId.Int64(), group.GroupTitle)
		if err != nil {
			//WHAT DO DO WITH ERROR HERE?
		}
	}

	return nil
}

func (s *Application) LeaveGroup(ctx context.Context, req models.GeneralGroupReq) error {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	//check request comes from member
	isMember, err := s.IsGroupMember(ctx, models.GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.UserId,
	})
	if err != nil {
		return ce.Wrap(nil, err)
	}
	if !isMember {
		return ce.New(ce.ErrPermissionDenied, fmt.Errorf("user %v is not a member of group %v", req.UserId, req.GroupId), input).WithPublic("permission denied")
	}

	err = s.db.LeaveGroup(ctx, ds.LeaveGroupParams{
		GroupID: req.GroupId.Int64(),
		UserID:  req.UserId.Int64(),
	})
	if err != nil {
		return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	return nil
}

// SKIP GRPC FOR NOW
func (s *Application) RemoveFromGroup(ctx context.Context, req models.RemoveFromGroupRequest) error {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	//check owner has indeed the owner role
	isOwner, err := s.isGroupOwner(ctx, models.GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.OwnerId,
	})
	if err != nil {
		return ce.Wrap(nil, err)
	}
	if !isOwner {
		return ce.New(ce.ErrPermissionDenied, fmt.Errorf("user %v is not the owner of group %v", req.OwnerId, req.GroupId), input).WithPublic("permission denied")
	}

	err = s.LeaveGroup(ctx, models.GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.MemberId,
	})
	if err != nil {
		return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	return nil
}

func (s *Application) CreateGroup(ctx context.Context, req *models.CreateGroupRequest) (models.GroupId, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return 0, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	groupId, err := s.db.CreateGroup(ctx, ds.CreateGroupParams{
		GroupOwner:       req.OwnerId.Int64(),
		GroupTitle:       req.GroupTitle.String(),
		GroupDescription: req.GroupDescription.String(),
		GroupImageID:     req.GroupImage.Int64(),
	})
	if err != nil {
		return 0, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}

	//call to chat service to create group conversation with owner as member
	// err = s.clients.CreateGroupConversation(ctx, groupId, req.OwnerId.Int64())
	// if err != nil {
	// 	tele.Info("group conversation couldn't be created", err)
	// }

	return models.GroupId(groupId), nil
}

func (s *Application) UpdateGroup(ctx context.Context, req *models.UpdateGroupRequest) error {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	//check requester is owner
	isOwner, err := s.isGroupOwner(ctx, models.GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.RequesterId,
	})
	if err != nil {
		return ce.Wrap(nil, err)
	}
	if !isOwner {
		return ce.New(ce.ErrPermissionDenied, fmt.Errorf("user %v is not the owner of group %v", req.RequesterId, req.GroupId), input).WithPublic("permission denied")
	}

	rowsAffected, err := s.db.UpdateGroup(ctx, ds.UpdateGroupParams{
		ID:               req.GroupId.Int64(),
		GroupTitle:       req.GroupTitle.String(),
		GroupDescription: req.GroupDescription.String(),
		GroupImageID:     req.GroupImage.Int64(),
	})

	if err != nil {
		return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}

	if rowsAffected != 1 {
		return ce.New(ce.ErrNotFound, fmt.Errorf("group %v was not found or has been deleted", req.GroupId), input).WithPublic("not found")
	}

	return nil

}

// NOT GRPC
func (s *Application) userInRelationToGroup(ctx context.Context, req models.GeneralGroupReq) (resp userInRelationToGroup, err error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return userInRelationToGroup{}, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	resp.isOwner, err = s.isGroupOwner(ctx, req)
	if err != nil {
		return userInRelationToGroup{}, ce.Wrap(nil, err)
	}
	resp.isMember, err = s.IsGroupMember(ctx, req)
	if err != nil {
		return userInRelationToGroup{}, ce.Wrap(nil, err)
	}
	resp.isPending, err = s.isGroupMembershipPending(ctx, req)
	if err != nil {
		return userInRelationToGroup{}, ce.Wrap(nil, err)
	}
	return resp, nil
}

// NOT GRPC
func (s *Application) isGroupOwner(ctx context.Context, req models.GeneralGroupReq) (bool, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return false, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}
	isOwner, err := s.db.IsUserGroupOwner(ctx, ds.IsUserGroupOwnerParams{
		ID:         req.GroupId.Int64(),
		GroupOwner: req.UserId.Int64(),
	})
	if err != nil {
		return false, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	if !isOwner {
		return false, nil
	}
	return true, nil
}

func (s *Application) IsGroupMember(ctx context.Context, req models.GeneralGroupReq) (bool, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return false, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	isMember, err := s.db.IsUserGroupMember(ctx, ds.IsUserGroupMemberParams{
		GroupID: req.GroupId.Int64(),
		UserID:  req.UserId.Int64(),
	})
	if err != nil {
		return false, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	if !isMember {
		return false, nil
	}
	return true, nil
}

// NOT GRPC
func (s *Application) isGroupMembershipPending(ctx context.Context, req models.GeneralGroupReq) (bool, error) {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return false, ce.Wrap(ce.ErrInvalidArgument, err, "request validation failed", input).WithPublic("invalid data received")
	}

	isPending, err := s.db.IsGroupMembershipPending(ctx, ds.IsGroupMembershipPendingParams{
		GroupID: req.GroupId.Int64(),
		UserID:  req.UserId.Int64(),
	})
	if err != nil {
		return false, ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	if !isPending.Valid { //should never happen
		return false, nil
	}
	return isPending.Bool, nil
}

// ---------------------------------------------------------------------
// low priority
// ---------------------------------------------------------------------
func DeleteGroup() {}

//called with group_id, owner_id
//returns success or error
//request needs to come from owner
//---------------------------------------------------------------------

//initiated by ownder
//SoftDeleteGroup

func TranferGroupOwnerShip() {}

//called with group_id,previous_owner_id, new_owner_id
//returns success or error
//request needs to come from previous owner (or admin - not implemented)
//---------------------------------------------------------------------

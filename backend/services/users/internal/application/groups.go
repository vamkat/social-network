package application

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
)

func (s *Application) GetAllGroupsPaginated(ctx context.Context, req Pagination) ([]Group, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return []Group{}, err
	}
	//paginated (sorting by most members first)
	rows, err := s.db.GetAllGroups(ctx, sqlc.GetAllGroupsParams{
		Offset: req.Offset.Int32(),
		Limit:  req.Limit.Int32(),
	})
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0, len(rows))

	for _, r := range rows {
		userInfo, err := s.userInRelationToGroup(ctx, GeneralGroupReq{
			GroupId: ct.Id(r.ID),
			UserId:  req.UserId,
		})
		if err != nil {
			return nil, err
		}

		groups = append(groups, Group{
			GroupId:          ct.Id(r.ID),
			GroupOwnerId:     ct.Id(r.GroupOwner),
			GroupTitle:       ct.Title(r.GroupTitle),
			GroupDescription: ct.About(r.GroupDescription),
			GroupImage:       r.GroupImage,
			MembersCount:     r.MembersCount,
			IsMember:         userInfo.isMember,
			IsOwner:          userInfo.isOwner,
			IsPending:        userInfo.isPending,
		})

	}

	return groups, nil
}

func (s *Application) GetUserGroupsPaginated(ctx context.Context, req Pagination) ([]Group, error) {
	//paginated (joined latest first)
	if err := ct.ValidateStruct(req); err != nil {
		return []Group{}, err
	}
	rows, err := s.db.GetUserGroups(ctx, sqlc.GetUserGroupsParams{
		GroupOwner: req.UserId.Int64(),
		Limit:      req.Limit.Int32(),
		Offset:     req.Offset.Int32(),
	})
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0, len(rows))
	for _, r := range rows {
		isPending, err := s.isGroupMembershipPending(ctx, GeneralGroupReq{
			GroupId: ct.Id(r.GroupID),
			UserId:  req.UserId,
		})
		if err != nil {
			return nil, err
		}
		groups = append(groups, Group{
			GroupId:          ct.Id(r.GroupID),
			GroupOwnerId:     ct.Id(r.GroupOwner),
			GroupTitle:       ct.Title(r.GroupTitle),
			GroupDescription: ct.About(r.GroupDescription),
			GroupImage:       r.GroupImage,
			MembersCount:     r.MembersCount,
			IsMember:         r.IsMember,
			IsOwner:          r.IsOwner,
			IsPending:        isPending,
		})
	}

	return groups, nil
}

// SKIP GRPC FOR NOW
func (s *Application) GetGroupInfo(ctx context.Context, req GeneralGroupReq) (Group, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return Group{}, err
	}
	row, err := s.db.GetGroupInfo(ctx, req.GroupId.Int64())
	if err != nil {
		return Group{}, nil
	}
	group := Group{
		GroupId:          ct.Id(row.ID),
		GroupOwnerId:     ct.Id(row.GroupOwner),
		GroupTitle:       ct.Title(row.GroupTitle),
		GroupDescription: ct.About(row.GroupDescription),
		GroupImage:       row.GroupImage,
		MembersCount:     row.MembersCount,
	}
	userInfo, err := s.userInRelationToGroup(ctx, GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.UserId,
	})
	if err != nil {
		return Group{}, err
	}
	group.IsMember = userInfo.isMember
	group.IsOwner = userInfo.isOwner
	group.IsPending = userInfo.isPending

	return group, nil

	//different calls for chat and posts (API GATEWAY)
}

func (s *Application) GetGroupMembers(ctx context.Context, req GroupMembersReq) ([]GroupUser, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return []GroupUser{}, err
	}
	//check request comes from member
	isMember, err := s.isGroupMember(ctx, GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.UserId,
	})
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotAuthorized
	}

	//paginated (newest first)
	rows, err := s.db.GetGroupMembers(ctx, sqlc.GetGroupMembersParams{
		GroupID: req.GroupId.Int64(),
		Limit:   req.Limit.Int32(),
		Offset:  req.Offset.Int32(),
	})
	if err != nil {
		return nil, err
	}
	members := make([]GroupUser, 0, len(rows))

	for _, r := range rows {
		var role string
		if r.Role.Valid {
			role = string(r.Role.GroupRole)
		}

		members = append(members, GroupUser{
			UserId:    ct.Id(r.ID),
			Username:  ct.Username(r.Username),
			Avatar:    r.Avatar,
			GroupRole: role,
		})
	}
	return members, nil
}

func (s *Application) SearchGroups(ctx context.Context, req GroupSearchReq) ([]Group, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return []Group{}, err
	}
	//weighted (title more important than description)
	//paginated (most members first)
	rows, err := s.db.SearchGroupsFuzzy(ctx, sqlc.SearchGroupsFuzzyParams{
		Similarity: req.SearchTerm.String(),
		GroupOwner: req.UserId.Int64(),
		Limit:      req.Limit.Int32(),
		Offset:     req.Offset.Int32(),
	})
	if err != nil {
		return []Group{}, err
	}
	groups := make([]Group, 0, len(rows))
	for _, r := range rows {
		isPending, err := s.isGroupMembershipPending(ctx, GeneralGroupReq{
			GroupId: ct.Id(r.ID),
			UserId:  req.UserId,
		})
		if err != nil {
			return nil, err
		}
		groups = append(groups, Group{
			GroupId:          ct.Id(r.ID),
			GroupOwnerId:     ct.Id(r.GroupOwner),
			GroupTitle:       ct.Title(r.GroupTitle),
			GroupDescription: ct.About(r.GroupDescription),
			GroupImage:       r.GroupImage,
			MembersCount:     r.MembersCount,
			IsMember:         r.IsMember,
			IsOwner:          r.IsOwner,
			IsPending:        isPending,
		})
	}

	return groups, nil

}

func (s *Application) InviteToGroup(ctx context.Context, req InviteToGroupReq) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	//check request comes from member
	isMember, err := s.isGroupMember(ctx, GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.InviterId,
	})
	if err != nil {
		return err
	}
	if !isMember {
		return ErrNotAuthorized
	}

	err = s.db.SendGroupInvite(ctx, sqlc.SendGroupInviteParams{
		GroupID:    req.GroupId.Int64(),
		SenderID:   req.InviterId.Int64(),
		ReceiverID: req.InvitedId.Int64(),
	})
	if err != nil {
		return err
	}

	return nil
}

// SKIP GRPC FOR NOW
func (s *Application) CancelInviteToGroup(ctx context.Context, req InviteToGroupReq) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	err := s.db.CancelGroupInvite(ctx, sqlc.CancelGroupInviteParams{
		GroupID:    req.GroupId.Int64(),
		ReceiverID: req.InvitedId.Int64(),
		SenderID:   req.InviterId.Int64(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Application) RequestJoinGroupOrCancel(ctx context.Context, req GroupJoinRequest) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	err := s.db.SendGroupJoinRequest(ctx, sqlc.SendGroupJoinRequestParams{
		GroupID: req.GroupId.Int64(),
		UserID:  req.RequesterId.Int64(),
	})

	if err != nil {
		return err
	}

	return nil
}

// SKIP GRPC FOR NOW
func (s *Application) CancelJoinGroupRequest(ctx context.Context, req GroupJoinRequest) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	err := s.db.CancelGroupJoinRequest(ctx, sqlc.CancelGroupJoinRequestParams{
		GroupID: req.GroupId.Int64(),
		UserID:  req.RequesterId.Int64(),
	})
	if err != nil {
		return err
	}
	return nil
}

// CHAT SERVICE EVENT add member to group conversation if accepted
func (s *Application) RespondToGroupInvite(ctx context.Context, req HandleGroupInviteRequest) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	if req.Accepted {

		err := s.db.AcceptGroupInvite(ctx, sqlc.AcceptGroupInviteParams{
			GroupID:    req.GroupId.Int64(),
			ReceiverID: req.InvitedId.Int64(),
		})
		if err != nil {
			return err

		}
	} else {
		err := s.db.DeclineGroupInvite(ctx, sqlc.DeclineGroupInviteParams{
			GroupID:    req.GroupId.Int64(),
			ReceiverID: req.InvitedId.Int64(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// CHAT SERVICE EVENT add member to group conversation if accepted
func (s *Application) HandleGroupJoinRequest(ctx context.Context, req HandleJoinRequest) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	isOwner, err := s.isGroupOwner(ctx, GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.OwnerId,
	})
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrNotAuthorized
	}

	if req.Accepted {
		err = s.db.AcceptGroupJoinRequest(ctx, sqlc.AcceptGroupJoinRequestParams{
			GroupID: req.GroupId.Int64(),
			UserID:  req.RequesterId.Int64(),
		})
	} else {
		err = s.db.RejectGroupJoinRequest(ctx, sqlc.RejectGroupJoinRequestParams{
			GroupID: req.GroupId.Int64(),
			UserID:  req.RequesterId.Int64(),
		})
	}
	if err != nil {
		return err
	}
	return nil
}

// CHAT SERVICE EVENT soft remove member from group conversation (keep history)
func (s *Application) LeaveGroup(ctx context.Context, req GeneralGroupReq) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	err := s.db.LeaveGroup(ctx, sqlc.LeaveGroupParams{
		GroupID: req.GroupId.Int64(),
		UserID:  req.UserId.Int64(),
	})
	if err != nil {
		return err
	}
	return nil
}

// SKIP GRPC FOR NOW
func (s *Application) RemoveFromGroup(ctx context.Context, req RemoveFromGroupRequest) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	//check owner has indeed the owner role
	isOwner, err := s.isGroupOwner(ctx, GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.OwnerId,
	})
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrNotAuthorized
	}

	err = s.LeaveGroup(ctx, GeneralGroupReq{
		GroupId: req.GroupId,
		UserId:  req.MemberId,
	})
	if err != nil {
		return err
	}
	return nil
}

// CHAT SERVICE EVENT create group conversation
func (s *Application) CreateGroup(ctx context.Context, req CreateGroupRequest) (GroupId, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return 0, err
	}

	groupId, err := s.db.CreateGroup(ctx, sqlc.CreateGroupParams{
		GroupOwner:       req.OwnerId.Int64(),
		GroupTitle:       req.GroupTitle.String(),
		GroupDescription: req.GroupDescription.String(),
		GroupImage:       req.GroupImage,
	})
	if err != nil {
		return 0, err
	}
	return GroupId(groupId), nil
}

// NOT GRPC
func (s *Application) userInRelationToGroup(ctx context.Context, req GeneralGroupReq) (resp UserInRelationToGroup, err error) {
	if err := ct.ValidateStruct(req); err != nil {
		return UserInRelationToGroup{}, err
	}
	resp.isOwner, err = s.isGroupOwner(ctx, req)
	if err != nil {
		return UserInRelationToGroup{}, err
	}
	resp.isMember, err = s.isGroupMember(ctx, req)
	if err != nil {
		return UserInRelationToGroup{}, err
	}
	resp.isPending, err = s.isGroupMembershipPending(ctx, req)
	if err != nil {
		return UserInRelationToGroup{}, err
	}
	return resp, nil
}

// NOT GRPC
func (s *Application) isGroupOwner(ctx context.Context, req GeneralGroupReq) (bool, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}
	isOwner, err := s.db.IsUserGroupOwner(ctx, sqlc.IsUserGroupOwnerParams{
		ID:         req.GroupId.Int64(),
		GroupOwner: req.UserId.Int64(),
	})
	if err != nil {
		return false, err
	}
	if !isOwner {
		return false, nil
	}
	return true, nil
}

// NOT GRPC
func (s *Application) isGroupMember(ctx context.Context, req GeneralGroupReq) (bool, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}
	isMember, err := s.db.IsUserGroupMember(ctx, sqlc.IsUserGroupMemberParams{
		GroupID: req.GroupId.Int64(),
		UserID:  req.UserId.Int64(),
	})
	if err != nil {
		return false, err
	}
	if !isMember {
		return false, nil
	}
	return true, nil
}

// NOT GRPC
func (s *Application) isGroupMembershipPending(ctx context.Context, req GeneralGroupReq) (bool, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}
	isPending, err := s.db.IsGroupMembershipPending(ctx, sqlc.IsGroupMembershipPendingParams{
		GroupID: req.GroupId.Int64(),
		UserID:  req.UserId.Int64(),
	})
	if err != nil {
		return false, err
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

package application

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
)

func (s *UserService) GetAllGroupsPaginated(ctx context.Context, req Pagination) ([]Group, error) {
	//paginated (sorting by most members first)
	rows, err := s.db.GetAllGroups(ctx, sqlc.GetAllGroupsParams{
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0, len(rows))

	for _, r := range rows {
		userInfo, err := s.userInRelationToGroup(ctx, GeneralGroupReq{
			GroupId: r.ID,
			UserId:  req.UserId,
		})
		if err != nil {
			return nil, err
		}

		groups = append(groups, Group{
			GroupId:          r.ID,
			GroupOwnerId:     r.GroupOwner,
			GroupTitle:       r.GroupTitle,
			GroupDescription: r.GroupDescription,
			GroupImage:       r.GroupImage,
			MembersCount:     r.MembersCount,
			IsMember:         userInfo.isMember,
			IsOwner:          userInfo.isOwner,
			IsPending:        userInfo.isPending,
		})

	}

	return groups, nil
}

func (s *UserService) GetUserGroupsPaginated(ctx context.Context, req Pagination) ([]Group, error) {
	//paginated (joined latest first)

	rows, err := s.db.GetUserGroups(ctx, sqlc.GetUserGroupsParams{
		GroupOwner: req.UserId,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0, len(rows))
	for _, r := range rows {
		isPending, err := s.isGroupMembershipPending(ctx, GeneralGroupReq{
			GroupId: r.GroupID,
			UserId:  req.UserId,
		})
		if err != nil {
			return nil, err
		}
		groups = append(groups, Group{
			GroupId:          r.GroupID,
			GroupOwnerId:     r.GroupOwner,
			GroupTitle:       r.GroupTitle,
			GroupDescription: r.GroupDescription,
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
func (s *UserService) GetGroupInfo(ctx context.Context, req GeneralGroupReq) (Group, error) {
	row, err := s.db.GetGroupInfo(ctx, req.GroupId)
	if err != nil {
		return Group{}, nil
	}
	group := Group{
		GroupId:          row.ID,
		GroupOwnerId:     row.GroupOwner,
		GroupTitle:       row.GroupTitle,
		GroupDescription: row.GroupDescription,
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

func (s *UserService) GetGroupMembers(ctx context.Context, req GroupMembersReq) ([]GroupUser, error) {
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
		GroupID: req.GroupId,
		Limit:   req.Limit,
		Offset:  req.Offset,
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
			UserId:    r.ID,
			Username:  r.Username,
			Avatar:    r.Avatar,
			GroupRole: role,
		})
	}
	return members, nil
}

func (s *UserService) SearchGroups(ctx context.Context, req GroupSearchReq) ([]Group, error) {
	//weighted (title more important than description)
	//paginated (most members first)
	rows, err := s.db.SearchGroupsFuzzy(ctx, sqlc.SearchGroupsFuzzyParams{
		Similarity: req.SearchTerm,
		GroupOwner: req.UserId,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})
	if err != nil {
		return []Group{}, err
	}
	groups := make([]Group, 0, len(rows))
	for _, r := range rows {
		isPending, err := s.isGroupMembershipPending(ctx, GeneralGroupReq{
			GroupId: r.ID,
			UserId:  req.UserId,
		})
		if err != nil {
			return nil, err
		}
		groups = append(groups, Group{
			GroupId:          r.ID,
			GroupOwnerId:     r.GroupOwner,
			GroupTitle:       r.GroupTitle,
			GroupDescription: r.GroupDescription,
			GroupImage:       r.GroupImage,
			MembersCount:     r.MembersCount,
			IsMember:         r.IsMember,
			IsOwner:          r.IsOwner,
			IsPending:        isPending,
		})
	}

	return groups, nil

}

func (s *UserService) InviteToGroup(ctx context.Context, req InviteToGroupReq) error {
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
		GroupID:    req.GroupId,
		SenderID:   req.InviterId,
		ReceiverID: req.InvitedId,
	})
	if err != nil {
		return err
	}

	return nil
}

// SKIP GRPC FOR NOW
func (s *UserService) CancelInviteToGroup(ctx context.Context, req InviteToGroupReq) error {

	err := s.db.CancelGroupInvite(ctx, sqlc.CancelGroupInviteParams{
		GroupID:    req.GroupId,
		ReceiverID: req.InvitedId,
		SenderID:   req.InviterId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) RequestJoinGroupOrCancel(ctx context.Context, req GroupJoinRequest) error {
	err := s.db.SendGroupJoinRequest(ctx, sqlc.SendGroupJoinRequestParams{
		GroupID: req.GroupId,
		UserID:  req.RequesterId,
	})

	if err != nil {
		return err
	}

	return nil
}

// SKIP GRPC FOR NOW
func (s *UserService) CancelJoinGroupRequest(ctx context.Context, req GroupJoinRequest) error {
	err := s.db.CancelGroupJoinRequest(ctx, sqlc.CancelGroupJoinRequestParams{
		GroupID: req.GroupId,
		UserID:  req.RequesterId,
	})
	if err != nil {
		return err
	}
	return nil
}

// CHAT SERVICE EVENT add member to group conversation if accepted
func (s *UserService) RespondToGroupInvite(ctx context.Context, req HandleGroupInviteRequest) error {

	if req.Accepted {

		err := s.db.AcceptGroupInvite(ctx, sqlc.AcceptGroupInviteParams{
			GroupID:    req.GroupId,
			ReceiverID: req.InvitedId,
		})
		if err != nil {
			return err

		}
	} else {
		err := s.db.DeclineGroupInvite(ctx, sqlc.DeclineGroupInviteParams{
			GroupID:    req.GroupId,
			ReceiverID: req.InvitedId,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// CHAT SERVICE EVENT add member to group conversation if accepted
func (s *UserService) HandleGroupJoinRequest(ctx context.Context, req HandleJoinRequest) error {

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
			GroupID: req.GroupId,
			UserID:  req.RequesterId,
		})
	} else {
		err = s.db.RejectGroupJoinRequest(ctx, sqlc.RejectGroupJoinRequestParams{
			GroupID: req.GroupId,
			UserID:  req.RequesterId,
		})
	}
	if err != nil {
		return err
	}
	return nil
}

// CHAT SERVICE EVENT soft remove member from group conversation (keep history)
func (s *UserService) LeaveGroup(ctx context.Context, req GeneralGroupReq) error {

	err := s.db.LeaveGroup(ctx, sqlc.LeaveGroupParams{
		GroupID: req.GroupId,
		UserID:  req.UserId,
	})
	if err != nil {
		return err
	}
	return nil
}

// SKIP GRPC FOR NOW
func (s *UserService) RemoveFromGroup(ctx context.Context, req RemoveFromGroupRequest) error {
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
func (s *UserService) CreateGroup(ctx context.Context, req CreateGroupRequest) (GroupId, error) {

	groupId, err := s.db.CreateGroup(ctx, sqlc.CreateGroupParams{
		GroupOwner:       req.OwnerId,
		GroupTitle:       req.GroupTitle,
		GroupDescription: req.GroupDescription,
		GroupImage:       req.GroupImage,
	})
	if err != nil {
		return 0, err
	}
	return GroupId(groupId), nil
}

// NOT GRPC
func (s *UserService) userInRelationToGroup(ctx context.Context, req GeneralGroupReq) (resp UserInRelationToGroup, err error) {
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
func (s *UserService) isGroupOwner(ctx context.Context, req GeneralGroupReq) (bool, error) {

	isOwner, err := s.db.IsUserGroupOwner(ctx, sqlc.IsUserGroupOwnerParams{
		ID:         req.GroupId,
		GroupOwner: req.UserId,
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
func (s *UserService) isGroupMember(ctx context.Context, req GeneralGroupReq) (bool, error) {

	isMember, err := s.db.IsUserGroupMember(ctx, sqlc.IsUserGroupMemberParams{
		GroupID: req.GroupId,
		UserID:  req.UserId,
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
func (s *UserService) isGroupMembershipPending(ctx context.Context, req GeneralGroupReq) (bool, error) {
	isPending, err := s.db.IsGroupMembershipPending(ctx, sqlc.IsGroupMembershipPendingParams{
		GroupID: req.GroupId,
		UserID:  req.UserId,
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

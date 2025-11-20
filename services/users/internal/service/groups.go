package userservice

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
)

func (s *UserService) GetAllGroupsPaginated(ctx context.Context) ([]Group, error) {
	//TODO add pagination (sort how?)
	rows, err := s.db.GetAllGroups(ctx)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0, len(rows))
	for _, r := range rows {
		groups = append(groups, Group{
			GroupId:          r.ID,
			GroupTitle:       r.GroupTitle,
			GroupDescription: r.GroupDescription,
			MembersCount:     r.MembersCount,
		})
	}

	return groups, nil
}

func (s *UserService) GetUserGroupsPaginated(ctx context.Context, userId int64) ([]Group, error) {
	//TODO add pagination (sort how?)
	rows, err := s.db.GetUserGroups(ctx, userId)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0, len(rows))
	for _, r := range rows {
		groups = append(groups, Group{
			GroupId:          r.GroupID,
			GroupTitle:       r.GroupTitle,
			GroupDescription: r.GroupDescription,
			MembersCount:     r.MembersCount,
			Role:             r.Role,
		})
	}

	return groups, nil
}

func (s *UserService) GetGroupInfo(ctx context.Context, groupId int64) (Group, error) {
	row, err := s.db.GetGroupInfo(ctx, groupId)
	if err != nil {
		return Group{}, nil
	}
	group := Group{
		GroupId:          row.ID,
		GroupTitle:       row.GroupTitle,
		GroupDescription: row.GroupDescription,
		MembersCount:     row.MembersCount,
	}

	return group, nil

	//different calls for chat and posts (API GATEWAY)
}

func (s *UserService) GetGroupMembers(ctx context.Context, groupId int64) ([]GroupUser, error) {
	//TODO add pagination (sort how?)
	rows, err := s.db.GetGroupMembers(ctx, groupId)
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
			UserId:    r.UserID,
			Username:  r.Username,
			Avatar:    r.Avatar,
			Public:    r.ProfilePublic,
			GroupRole: role,
		})
	}
	return members, nil
}

func (s *UserService) SearchGroups(ctx context.Context, searchTerm string) ([]Group, error) {
	//TODO add pagination? (sort how?)
	rows, err := s.db.SearchGroupsFuzzy(ctx, searchTerm)
	if err != nil {
		return []Group{}, err
	}
	groups := make([]Group, 0, len(rows))
	for _, r := range rows {
		groups = append(groups, Group{
			GroupId:          r.ID,
			GroupTitle:       r.GroupTitle,
			GroupDescription: r.GroupDescription,
			MembersCount:     r.MembersCount,
		})
	}

	return groups, nil
	//Do we want to include userId so that we also have the info if the user is a member, owner, or nothing?
}

func (s *UserService) InviteToGroupOrCancel(ctx context.Context, req InviteToGroupOrCancelRequest) error {

	if req.Cancel {
		err := s.db.CancelGroupInvite(ctx, sqlc.CancelGroupInviteParams{
			GroupID:    req.GroupId,
			ReceiverID: req.InvitedId,
			SenderID:   req.InviterId,
		})
		if err != nil {
			return err
		}
	} else {

		role, err := s.GetUserGroupRole(ctx, GroupRoleReq{
			GroupId: req.GroupId,
			UserId:  req.InviterId,
		})
		if err != nil {
			return err
		}
		if role == "" {
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
	}

	return nil
}

func (s *UserService) GetUserGroupRole(ctx context.Context, req GroupRoleReq) (GroupRole, error) {
	role, err := s.db.GetUserGroupRole(ctx, sqlc.GetUserGroupRoleParams{
		GroupID: req.GroupId,
		UserID:  req.UserId,
	})
	if err != nil {
		return "", err
	}
	if !role.Valid {
		return "", err
	}
	return GroupRole(role.GroupRole), nil
}

func (s *UserService) RequestJoinGroupOrCancel(ctx context.Context, req GroupJoinOrCancelRequest) error {
	var err error
	if req.Cancel {
		//userId needs to come from token, to make sure the user is not trying to cancel someone else's join request
		err = s.db.CancelGroupJoinRequest(ctx, sqlc.CancelGroupJoinRequestParams{
			GroupID: req.GroupId,
			UserID:  req.RequesterId,
		})

	} else {
		err = s.db.SendGroupJoinRequest(ctx, sqlc.SendGroupJoinRequestParams{
			GroupID: req.GroupId,
			UserID:  req.RequesterId,
		})

	}
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) RespondToGroupInvite(ctx context.Context, req HandleGroupInviteRequest) error {
	//userId needs to come from token, to make sure the user is not trying to answer invite not for them
	if req.Accepted {
		err := s.runTx(ctx, func(q *sqlc.Queries) error {
			err := q.AcceptGroupInvite(ctx, sqlc.AcceptGroupInviteParams{
				GroupID:    req.GroupId,
				ReceiverID: req.InvitedId,
			})
			if err != nil {
				return err
			}
			err = q.AddUserToGroup(ctx, sqlc.AddUserToGroupParams{
				GroupID: req.GroupId,
				UserID:  req.InvitedId,
			})
			if err != nil {
				return err
			}
			return nil
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

func HandleGroupJoinRequest() {
	//called with group_id,user_id (who requested to join),owner_id(who responds), bool (accept or decline)
	//returns success or error
	//request needs to come from group owner
	//---------------------------------------------------------------------

	//yes or no
	//AcceptGroupJoinRequest & addUserToGroup
	//RejectGroupJoinRequest
}

func LeaveGroup() {
	//called with group_id,user_id
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------

	//initiated by user
	//LeaveGroup
}

func RemoveFromGroup() {
	//called with group_id,user_id (who is removed), owner_id(making the request)
	//returns success or error
	//request needs to come from group owner
	//---------------------------------------------------------------------

	//initiated by owner
	//LeaveGroup
}

func CreateGroup() {
	//called with owner_id, group_title, group_description
	//returns group_id
	//---------------------------------------------------------------------

	//CreateGroup
	//AddGroupOwnerAsMember
}

func DeleteGroup() { //low priorty
	//called with group_id, owner_id
	//returns success or error
	//request needs to come from owner
	//---------------------------------------------------------------------

	//initiated by ownder
	//SoftDeleteGroup
}

func TranferGroupOwnerShip() { //low priority
	//called with group_id,previous_owner_id, new_owner_id
	//returns success or error
	//request needs to come from previous owner (or admin - not implemented)
	//---------------------------------------------------------------------

}

package userservice

import (
	"context"

	"social-network/services/users/internal/db/sqlc"

	"github.com/stretchr/testify/mock"
)

// MockQuerier is a mock implementation of sqlc.Querier
type MockQuerier struct {
	mock.Mock
}

// Auth-related methods
func (m *MockQuerier) InsertNewUser(ctx context.Context, arg sqlc.InsertNewUserParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) InsertNewUserAuth(ctx context.Context, arg sqlc.InsertNewUserAuthParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) GetUserForLogin(ctx context.Context, identifier string) (sqlc.GetUserForLoginRow, error) {
	args := m.Called(ctx, identifier)
	if args.Get(0) == nil {
		return sqlc.GetUserForLoginRow{}, args.Error(1)
	}
	return args.Get(0).(sqlc.GetUserForLoginRow), args.Error(1)
}

func (m *MockQuerier) IncrementFailedLoginAttempts(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockQuerier) ResetFailedLoginAttempts(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockQuerier) GetUserPassword(ctx context.Context, userID int64) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockQuerier) UpdateUserPassword(ctx context.Context, arg sqlc.UpdateUserPasswordParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdateUserEmail(ctx context.Context, arg sqlc.UpdateUserEmailParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// Profile-related methods
func (m *MockQuerier) GetUserBasic(ctx context.Context, id int64) (sqlc.GetUserBasicRow, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return sqlc.GetUserBasicRow{}, args.Error(1)
	}
	return args.Get(0).(sqlc.GetUserBasicRow), args.Error(1)
}

func (m *MockQuerier) GetUserProfile(ctx context.Context, id int64) (sqlc.GetUserProfileRow, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return sqlc.GetUserProfileRow{}, args.Error(1)
	}
	return args.Get(0).(sqlc.GetUserProfileRow), args.Error(1)
}

func (m *MockQuerier) GetFollowerCount(ctx context.Context, id int64) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) GetFollowingCount(ctx context.Context, id int64) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) GetUserGroups(ctx context.Context, id int64) ([]sqlc.GetUserGroupsRow, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sqlc.GetUserGroupsRow), args.Error(1)
}

func (m *MockQuerier) SearchUsers(ctx context.Context, arg sqlc.SearchUsersParams) ([]sqlc.SearchUsersRow, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sqlc.SearchUsersRow), args.Error(1)
}

func (m *MockQuerier) UpdateUserProfile(ctx context.Context, arg sqlc.UpdateUserProfileParams) (sqlc.User, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return sqlc.User{}, args.Error(1)
	}
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockQuerier) UpdateProfilePrivacy(ctx context.Context, arg sqlc.UpdateProfilePrivacyParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// Followers-related methods
func (m *MockQuerier) GetFollowers(ctx context.Context, arg sqlc.GetFollowersParams) ([]sqlc.GetFollowersRow, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sqlc.GetFollowersRow), args.Error(1)
}

func (m *MockQuerier) GetFollowing(ctx context.Context, arg sqlc.GetFollowingParams) ([]sqlc.GetFollowingRow, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sqlc.GetFollowingRow), args.Error(1)
}

func (m *MockQuerier) FollowUser(ctx context.Context, arg sqlc.FollowUserParams) (string, error) {
	args := m.Called(ctx, arg)
	return args.String(0), args.Error(1)
}

func (m *MockQuerier) UnfollowUser(ctx context.Context, arg sqlc.UnfollowUserParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) AcceptFollowRequest(ctx context.Context, arg sqlc.AcceptFollowRequestParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) RejectFollowRequest(ctx context.Context, arg sqlc.RejectFollowRequestParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) IsFollowing(ctx context.Context, arg sqlc.IsFollowingParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) GetMutualFollowers(ctx context.Context, arg sqlc.GetMutualFollowersParams) ([]sqlc.GetMutualFollowersRow, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sqlc.GetMutualFollowersRow), args.Error(1)
}

func (m *MockQuerier) IsFollowingEither(ctx context.Context, arg sqlc.IsFollowingEitherParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

// Groups-related methods
func (m *MockQuerier) GetAllGroups(ctx context.Context) ([]sqlc.GetAllGroupsRow, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sqlc.GetAllGroupsRow), args.Error(1)
}

func (m *MockQuerier) GetGroupInfo(ctx context.Context, groupID int64) (sqlc.GetGroupInfoRow, error) {
	args := m.Called(ctx, groupID)
	if args.Get(0) == nil {
		return sqlc.GetGroupInfoRow{}, args.Error(1)
	}
	return args.Get(0).(sqlc.GetGroupInfoRow), args.Error(1)
}

func (m *MockQuerier) GetGroupMembers(ctx context.Context, groupID int64) ([]sqlc.GetGroupMembersRow, error) {
	args := m.Called(ctx, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sqlc.GetGroupMembersRow), args.Error(1)
}

func (m *MockQuerier) SearchGroupsFuzzy(ctx context.Context, searchTerm string) ([]sqlc.SearchGroupsFuzzyRow, error) {
	args := m.Called(ctx, searchTerm)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sqlc.SearchGroupsFuzzyRow), args.Error(1)
}

func (m *MockQuerier) CreateGroup(ctx context.Context, arg sqlc.CreateGroupParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) AddUserToGroup(ctx context.Context, arg sqlc.AddUserToGroupParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) SendGroupInvite(ctx context.Context, arg sqlc.SendGroupInviteParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) CancelGroupInvite(ctx context.Context, arg sqlc.CancelGroupInviteParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) AcceptGroupInvite(ctx context.Context, arg sqlc.AcceptGroupInviteParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) SendGroupJoinRequest(ctx context.Context, arg sqlc.SendGroupJoinRequestParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) CancelGroupJoinRequest(ctx context.Context, arg sqlc.CancelGroupJoinRequestParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) AcceptGroupJoinRequest(ctx context.Context, arg sqlc.AcceptGroupJoinRequestParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) DeclineGroupInvite(ctx context.Context, arg sqlc.DeclineGroupInviteParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) RejectGroupJoinRequest(ctx context.Context, arg sqlc.RejectGroupJoinRequestParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) IsUserGroupMember(ctx context.Context, arg sqlc.IsUserGroupMemberParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) IsUserGroupOwner(ctx context.Context, arg sqlc.IsUserGroupOwnerParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) LeaveGroup(ctx context.Context, arg sqlc.LeaveGroupParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// Admin-related methods
func (m *MockQuerier) AddGroupOwnerAsMember(ctx context.Context, arg sqlc.AddGroupOwnerAsMemberParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) BanUser(ctx context.Context, arg sqlc.BanUserParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) SoftDeleteGroup(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) SoftDeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) TransferOwnership(ctx context.Context, arg sqlc.TransferOwnershipParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UnbanUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) GetUserGroupRole(ctx context.Context, arg sqlc.GetUserGroupRoleParams) (sqlc.NullGroupRole, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return sqlc.NullGroupRole{}, args.Error(1)
	}
	return args.Get(0).(sqlc.NullGroupRole), args.Error(1)
}

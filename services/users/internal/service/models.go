package userservice

import (
	"time"
)

type UserId int64

type GroupId int64

type GroupRole string

type LoginReq struct {
	Identifier string
	Password   string
}

// returned by login, basicUser,searchUsers([]), getFollowers([]), getFollowing([])
type User struct {
	UserId   int64
	Username string
	Avatar   *string
	Public   bool
}

// getGroupMembers([])
type GroupUser struct {
	UserId    int64
	Username  string
	Avatar    *string
	Public    bool
	GroupRole string //only applicable in group members
}

// returned by getAllGroups([]), getUserGroups([]),getGroupInfo, searchGroup([])
type Group struct {
	GroupId          int64
	GroupTitle       string
	GroupDescription string
	MembersCount     *int32
	Role             string
}

type RegisterUserRequest struct {
	Username    string
	FirstName   string
	LastName    string
	DateOfBirth string
	Avatar      *string
	About       *string
	Public      bool
	Email       string
	Password    string
}

type RegisterUserResponse struct {
	UserId string
}

type LoginRequest struct {
	Identifier string //username or email
	Password   string
}

type UserProfileRequest struct {
	UserId      int64
	RequesterId int64
}

// returned by getUserProfile, updateUserProfile
type UserProfileResponse struct {
	UserId         int64
	Username       string
	FirstName      string
	LastName       string
	DateOfBirth    time.Time
	Avatar         *string
	About          *string
	Public         bool
	FollowersCount int64
	FollowingCount int64
	Groups         []Group
}

type UserSearchReq struct {
	SearchTerm string
	Limit      int32
}

type UpdateProfileRequest struct {
	Username    string
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Avatar      string
	About       string
	Public      bool
}

type GetFollowersReq struct {
	FollowingID int64
	Limit       int32
	Offset      int32
}

type GetFollowingReq struct {
	FollowerID int64
	Limit      int32
	Offset     int32
}

type UpdatePasswordRequest struct {
	UserId      int64
	OldPassword string
	Password    string
}

type UpdateEmailRequest struct {
	UserId int64
	Email  string
}

type InviteToGroupOrCancelRequest struct {
	InviterId int64
	InvitedId int64
	GroupId   int64
	Cancel    bool
}

type HandleGroupInviteRequest struct {
	GroupId   int64
	InvitedId int64
	Accepted  bool
}

type GroupJoinOrCancelRequest struct {
	GroupId     int64
	RequesterId int64
	Cancel      bool
}

type HandleJoinRequest struct {
	GroupId     int64
	RequesterId int64
	OwnerId     int64
	Accepted    bool
}

type LeaveGroupRequest struct {
	GroupId  int64
	MemberId int64
	OwnerId  int64 //nil if initiated by member
}

type CreateGroupRequest struct {
	OwnerId          int64
	GroupTitle       string
	GroupDescription string
}

type GroupRoleReq struct {
	GroupId int64
	UserId  int64
}

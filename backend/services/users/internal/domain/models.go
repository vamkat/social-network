package userservice

import (
	"time"
)

//-------------------------------------------
// Auth
//-------------------------------------------

type RegisterUserRequest struct {
	Username    string
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Avatar      string
	About       string
	Public      bool
	Email       string
	Password    string
}

type RegisterUserResponse struct {
	UserId int64
}

type LoginRequest struct {
	Identifier string //username or email
	Password   string
}

type UpdatePasswordRequest struct {
	UserId      int64
	OldPassword string
	NewPassword string
}

type UpdateEmailRequest struct {
	UserId int64
	Email  string
}

//-------------------------------------------
// Profile
//-------------------------------------------

type UserId int64

type User struct {
	UserId   int64
	Username string
	Avatar   string
}

type UserSearchReq struct {
	SearchTerm string
	Limit      int32
}

type UserProfileRequest struct {
	UserId      int64
	RequesterId int64
}

type UserProfileResponse struct {
	UserId            int64
	Username          string
	FirstName         string
	LastName          string
	DateOfBirth       time.Time
	Avatar            string
	About             string
	Public            bool
	CreatedAt         time.Time
	FollowersCount    int64
	FollowingCount    int64
	GroupsCount       int64
	OwnedGroupsCount  int64
	ViewerIsFollowing bool
	OwnProfile        bool
}

type UpdateProfileRequest struct {
	UserId      int64
	Username    string
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Avatar      string
	About       string
}

type UpdateProfilePrivacyRequest struct {
	UserId int64
	Public bool
}

// -------------------------------------------
// Groups
// -------------------------------------------
type GroupId int64

type GroupRole string

type Pagination struct {
	Limit  int32
	Offset int32
}

type GroupMembersReq struct {
	GroupId int64
	Limit   int32
	Offset  int32
}

type UserGroupsPaginated struct {
	UserId int64
	Limit  int32
	Offset int32
}

type GroupUser struct {
	UserId    int64
	Username  string
	Avatar    string
	Public    bool
	GroupRole string
}

type GroupSearchReq struct {
	SearchTerm string
	UserId     int64
	Limit      int32
	Offset     int32
}

// add owner to group
type Group struct {
	GroupId          int64
	GroupOwnerId     int64
	GroupTitle       string
	GroupDescription string
	MembersCount     int32
	IsMember         bool
	IsOwner          bool
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

type GeneralGroupReq struct {
	GroupId int64
	UserId  int64
}

type RemoveFromGroupRequest struct {
	GroupId  int64
	MemberId int64
	OwnerId  int64
}

type CreateGroupRequest struct {
	OwnerId          int64
	GroupTitle       string
	GroupDescription string
}

// -------------------------------------------
// Followers
// -------------------------------------------

type FollowUserReq struct {
	FollowerId   int64
	TargetUserId int64
}

type HandleFollowRequestReq struct {
	UserId      int64
	RequesterId int64
	Accept      bool
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

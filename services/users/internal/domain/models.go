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
	UserId string
}

type LoginRequest struct {
	Identifier string //username or email
	Password   string
}

type UpdatePasswordRequest struct {
	UserId      string
	OldPassword string
	NewPassword string
}

type UpdateEmailRequest struct {
	UserId string
	Email  string
}

//-------------------------------------------
// Profile
//-------------------------------------------

type UserId string

type User struct {
	UserId   string
	Username string
	Avatar   string
}

type UserSearchReq struct {
	SearchTerm string
	Limit      int32
}

type UserProfileRequest struct {
	UserId      string
	RequesterId string
}

type UserProfileResponse struct {
	UserId            string
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
	UserId      string
	Username    string
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Avatar      string
	About       string
}

type UpdateProfilePrivacyRequest struct {
	UserId string
	Public bool
}

// -------------------------------------------
// Groups
// -------------------------------------------
type GroupId int64

type GroupRole string

type GroupUser struct {
	UserId    string
	Username  string
	Avatar    string
	Public    bool
	GroupRole string
}

// add owner to group
type Group struct {
	GroupId          int64
	GroupTitle       string
	GroupDescription string
	MembersCount     int32
	Role             string
}

type InviteToGroupOrCancelRequest struct {
	InviterId string
	InvitedId string
	GroupId   int64
	Cancel    bool
}

type HandleGroupInviteRequest struct {
	GroupId   int64
	InvitedId string
	Accepted  bool
}

type GroupJoinOrCancelRequest struct {
	GroupId     int64
	RequesterId string
	Cancel      bool
}

type HandleJoinRequest struct {
	GroupId     int64
	RequesterId string
	OwnerId     string
	Accepted    bool
}

type GeneralGroupReq struct {
	GroupId int64
	UserId  string
}

type RemoveFromGroupRequest struct {
	GroupId  int64
	MemberId string
	OwnerId  string
}

type CreateGroupRequest struct {
	OwnerId          string
	GroupTitle       string
	GroupDescription string
}

// -------------------------------------------
// Followers
// -------------------------------------------

type FollowUserReq struct {
	FollowerId   string
	TargetUserId string
}

type HandleFollowRequestReq struct {
	UserId      string
	RequesterId string
	Accept      bool
}

type GetFollowersReq struct {
	FollowingID string
	Limit       int32
	Offset      int32
}

type GetFollowingReq struct {
	FollowerID string
	Limit      int32
	Offset     int32
}

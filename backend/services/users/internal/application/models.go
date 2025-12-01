package application

import (
	ct "social-network/shared/go/customtypes"
	"time"
)

//-------------------------------------------
// Auth
//-------------------------------------------

type RegisterUserRequest struct {
	Username    ct.Username
	FirstName   ct.Name
	LastName    ct.Name
	DateOfBirth ct.DateOfBirth
	Avatar      string
	About       ct.About
	Public      bool
	Email       ct.Email
	Password    ct.Password
}

type RegisterUserResponse struct {
	UserId int64
}

type LoginRequest struct {
	Identifier ct.Identifier //username or email
	Password   ct.Password
}

type UpdatePasswordRequest struct {
	UserId      ct.Id
	OldPassword ct.Password
	NewPassword ct.Password
}

type UpdateEmailRequest struct {
	UserId ct.Id
	Email  ct.Email
}

//-------------------------------------------
// Profile
//-------------------------------------------

type UserId int64

type User struct {
	UserId   ct.Id
	Username ct.Username
	Avatar   string
}

type UserSearchReq struct {
	SearchTerm ct.SearchTerm
	Limit      ct.Limit
}

type UserProfileRequest struct {
	UserId      ct.Id
	RequesterId ct.Id
}

type UserProfileResponse struct {
	UserId            ct.Id
	Username          ct.Username
	FirstName         ct.Name
	LastName          ct.Name
	DateOfBirth       ct.DateOfBirth
	Avatar            string
	About             ct.About
	Public            bool
	CreatedAt         time.Time
	FollowersCount    int64
	FollowingCount    int64
	GroupsCount       int64
	OwnedGroupsCount  int64
	ViewerIsFollowing bool
	OwnProfile        bool
	IsPending         bool
}

type UpdateProfileRequest struct {
	UserId      ct.Id
	Username    ct.Username
	FirstName   ct.Name
	LastName    ct.Name
	DateOfBirth ct.DateOfBirth
	Avatar      string
	About       ct.About
}

type UpdateProfilePrivacyRequest struct {
	UserId ct.Id
	Public bool
}

// -------------------------------------------
// Groups
// -------------------------------------------
type GroupId int64

type GroupRole string

type GroupMembersReq struct {
	UserId  ct.Id
	GroupId ct.Id
	Limit   ct.Limit
	Offset  ct.Offset
}

type Pagination struct {
	UserId ct.Id
	Limit  ct.Limit
	Offset ct.Offset
}

type GroupUser struct {
	UserId    ct.Id
	Username  ct.Username
	Avatar    string
	GroupRole string
}

type GroupSearchReq struct {
	SearchTerm ct.SearchTerm
	UserId     ct.Id
	Limit      ct.Limit
	Offset     ct.Offset
}

// add owner to group
type Group struct {
	GroupId          ct.Id
	GroupOwnerId     ct.Id
	GroupTitle       ct.Title
	GroupDescription ct.About
	GroupImage       string
	MembersCount     int32
	IsMember         bool
	IsOwner          bool
	IsPending        bool
}

type InviteToGroupReq struct {
	InviterId ct.Id
	InvitedId ct.Id
	GroupId   ct.Id
}

type HandleGroupInviteRequest struct {
	GroupId   ct.Id
	InvitedId ct.Id
	Accepted  bool
}

type GroupJoinRequest struct {
	GroupId     ct.Id
	RequesterId ct.Id
}

type HandleJoinRequest struct {
	GroupId     ct.Id
	RequesterId ct.Id
	OwnerId     ct.Id
	Accepted    bool
}

type GeneralGroupReq struct {
	GroupId ct.Id
	UserId  ct.Id
}

type RemoveFromGroupRequest struct {
	GroupId  ct.Id
	MemberId ct.Id
	OwnerId  ct.Id
}

type CreateGroupRequest struct {
	OwnerId          ct.Id
	GroupTitle       ct.Title
	GroupDescription ct.About
	GroupImage       string
}

type UserInRelationToGroup struct {
	isOwner   bool
	isMember  bool
	isPending bool
}

// -------------------------------------------
// Followers
// -------------------------------------------

type FollowUserReq struct {
	FollowerId   ct.Id
	TargetUserId ct.Id
}

type FollowUserResp struct {
	IsPending         bool
	ViewerIsFollowing bool
}

type HandleFollowRequestReq struct {
	UserId      ct.Id
	RequesterId ct.Id
	Accept      bool
}

package userservice

import "time"

//returned by login, basicUser,getGroupMembers([]),searchUsers([]), getFollowers([]), getFollowing([])
type user struct {
	userId    int64
	username  string
	avatar    string
	public    bool
	groupRole string //only applicable in group members
}

//returned by getAllGroups([]), getUserGroups([]),getGroupInfo, searchGroup([])
type group struct {
	groupId          int64
	ownerId          int64
	groupTitle       string
	groupDescription string
	membersCount     int
}

type registerUserRequest struct {
	username    string
	firstName   string
	lastName    string
	dateOfBirth time.Time
	avatar      string
	about       string
	public      bool
	email       string
	password    string
}

type registerUserResponse struct {
	userId string
}

type loginRequest struct {
	identifier string //username or email
	password   string
}

type userProfileRequest struct {
	userId      int64
	requesterId int64
}

//returned by getUserProfile, updateUserProfile
type userProfileResponse struct {
	userId         int64
	username       string
	firstName      string
	lastName       string
	dateOfBirth    time.Time
	avatar         string
	about          string
	public         bool
	followersCount int
	followingCount int
	groups         []group
}

type updateProfileRequest struct {
	username    string
	firstName   string
	lastName    string
	dateOfBirth time.Time
	avatar      string
	about       string
	public      bool
}

type updatePasswordRequest struct {
	userId      int64
	oldPassword string
	password    string
}

type updateEmailRequest struct {
	userId int64
	email  string
}

type inviteToGroupOrCancelRequest struct {
	inviterId int64
	invitedId int64
	groupId   int64
	cancel    bool
}

type handleGroupInviteRequest struct {
	groupId   int64
	invitedId int64
	accepted  bool
}

type groupJoinOrCancelRequest struct {
	groupId     int64
	requesterId int64
	cancel      bool
}

type handleGroupJoinRequest struct {
	groupId     int64
	requesterId int64
	ownerId     int64
	accepted    bool
}

type leaveGroupRequest struct {
	groupId  int64
	memberId int64
	ownerId  int64 //nil if initiated by member
}

type createGroupRequest struct {
	ownerId          int64
	groupTitle       string
	groupDescription string
}

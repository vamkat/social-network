package userservice

import (
	"context"
)

func RegisterUser() {
	//called with: username, first_name, last_name, date_of_birth, avatar, about_me, profile_public(bool), email, password hash, salt
	//returns: user_id or error
	//---------------------------------------------------------------------
	//InsertNewUser, get id
	//hash password, add salt (is this the service's responsibility? Probably at gateway))
	//if no conflict insert InsertNewUserAuth
	//if no error LoginUser(id)
}

func LoginUser() {
	//called with: username (and/or email?), password
	//returns user_id or error
	//---------------------------------------------------------------------
	//BY EMAIL OR ONLY USERNAME???? (discuss with front)
	//GetUserForLogin (id, username, password hash, salt), check status=active (maybe also email TODO)
	//check password correct
	//if failed login IncrementFailedLoginAttempts
	//if success ResetFailedLoginAttempts
}

func DeleteUser() {
	//called with: user_id
	//returns success or error
	//request needs to come from same user or admin (not implemented)
	//---------------------------------------------------------------------
	//softDeleteUser(id)
}

func BanUser() {
	//called with: user_id
	//returns success or error
	//request needs to come from admin (not implemented)
	//---------------------------------------------------------------------
	//BanUser(id, liftBanDate)
	//TODO logic to automatically unban after liftBanDate
}

func UnbanUser() {
	//called with: user_id
	//returns success or error
	//request needs to come from admin (not implemented) or automatic on expiration
	//---------------------------------------------------------------------
	//UnbanUser(id)
}

type BasicUserInfo struct {
	UserName      string
	Avatar        string
	PublicProfile bool
}

func GetBasicUserInfo(ctx context.Context, userID int64) (resp BasicUserInfo, err error) {
	//called with: user_id
	//returns username, avatar, profile_public(bool)
	//---------------------------------------------------------------------
	// GetUserBasic(id)
	return BasicUserInfo{}, nil
}

func GetUserProfile() {
	//called with: user_id, viewer_id (to check permission to view)
	//returns id, username, first_name, last_name, date_of_birth, avatar, about_me, profile_public, number of followers, number of following, list of groups
	//---------------------------------------------------------------------
	// check if user has permission to see (public profile or isFollower)
	// GetUserProlife(id)

	// number of followers, following (TODO keep in profile with trigger?)
	// number of groups? (TODO keep in profile with trigger?)
	// which groups

	// THIS CAN BE HANDLED BY THE API GATEWAY (and different call from front):
	// from forum service get all posts paginated (and number of posts)
	// and within all posts check each one if viewer has permission
}

func UpdateUserProfile() {
	//called with user_id and any of: username (TODO), first_name, last_name, date_of_birth, avatar, about_me
	//returns full profile
	//request needs to come from same user
	//---------------------------------------------------------------------

	//UpdateUserProflie
	//TODO check how to not update all fields but only changes (nil pointers?)
}

func UpdateUserPassword() {
	//called with user_id, old password, new password_hash, salt
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------
	//UpdateUserPassword
}

func UpdateUserEmail() {
	//called with user_id, new email
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------
	//UpdateUserEmail
}

func GetAllGroups() {
	//called with nothing
	//returns list of groups containing group_id, group_title, group_description, members_count
	//---------------------------------------------------------------------
	//GetAllGroups
}

func GetUserGroups() {
	//called with user_id
	//returns list of groups containing group_id, group_title, group_description, members_count
	//---------------------------------------------------------------------

	//GetUserGroups
}

func GetGroupInfo() {
	//called with group_id
	//returns group_id, group_title, group_description, members_count (owner?)
	//---------------------------------------------------------------------

	//GetGroupInfo

	//different calls for chat and posts (API GATEWAY)
}

func GetGroupMembers() {
	//called with group_id
	//returns list of members containing user_id, username, avatar, profile_public(bool), role(owner or member), joined_at
	//---------------------------------------------------------------------

	//getGroupMembers
}

func SeachByUserrnameOrName() {
	//called with search term
	//returns list of users containing user_id, username, avatar, profile_public(bool)
	//---------------------------------------------------------------------

	//SearchUsers
}

func SearchGroup() {
	//called with search term
	//returns list of groups containing group_id, group_title, group_description, members_count
	//---------------------------------------------------------------------

	//SeachGroupsFuzzy
}

func GetFollowers() {
	//called with user_id
	//returns list of users containing user_id, username, avatar, profile_public(bool)
	//---------------------------------------------------------------------

	//GetFollowers() TODO FIX RETURNS
}

func GetFollowing() {
	//called with user_id
	//returns list of users containing user_id, username, avatar, profile_public(bool)
	//---------------------------------------------------------------------
	//GetFollowing()
}

func InviteToGroup() {
	//called with group_id,user_id
	//returns success or error
	//request needs to come from group owner or group member
	//---------------------------------------------------------------------

	//SendGroupInvite
}

func RequestJoinGroup() {
	//called with group_id,user_id
	//returns success or error
	//---------------------------------------------------------------------

	//SendGroupJoinRequest
}

func HandleGroupInvite() {
	//called with group_id,user_id, bool (accept or decline)
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------

	//yes or no
	//AcceptGroupInvite & addUserToGroup
	//DeclineGroupInvite
}

func CancelGroupInvite() {
	//called with group_id,user_id (who is invited), sender_id
	//returns success or error
	//request needs to come from sender
	//---------------------------------------------------------------------

	//CancelGroupInvite
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

func CancelGroupJoinRequest() {
	//called with group_id,user_id, bool (accept or decline)
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------

	//CancelGroupJoinRequest
}

func LeaveGroup() {
	//called with group_id,user_id, bool
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

func DeleteGroup() {
	//called with group_id, owner_id
	//returns success or error
	//request needs to come from owner
	//---------------------------------------------------------------------

	//initiated by ownder
	//SoftDeleteGroup
}

func TranferGroupOwnerShip() {
	//called with group_id,previous_owner_id, new_owner_id
	//returns success or error
	//request needs to come from previous owner (or admin - not implemented)
	//---------------------------------------------------------------------

}

func hashPasswordWithSalt() {

}

func comparePasswords() {

}

func issueToken() {

}

func checkToken() {

}

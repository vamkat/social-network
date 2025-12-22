package handlers

import ct "social-network/shared/go/customtypes"

type GroupT struct {
	GroupId          ct.Id  `json:"group_id"`
	GroupOwnerId     ct.Id  `json:"group_owner_id"`
	GroupTitle       string `json:"group_title"`
	GroupDescription string `json:"group_description"`
	GroupImage       ct.Id  `json:"group_image"`
	MembersCount     int32  `json:"members_count"`
	IsMember         bool   `json:"is_member"`
	IsOwner          bool   `json:"is_owner"`
	IsPending        bool   `json:"is_pending"`
}

type GroupsT struct {
	Groups []GroupT
}

type UserT struct {
	UserId   ct.Id       `json:"id"`
	Username ct.Username `json:"username"`
	AvatarId ct.Id       `json:"avatar_id"`
}

type GroupUserT struct {
	UserId    ct.Id       `json:"id"`
	Username  ct.Username `json:"username"`
	AvatarId  ct.Id       `json:"avatar_id"`
	GroupRole string      `json:"group_role"`
}

type UsersRespT struct {
	Users []UserT `json:"users"`
}

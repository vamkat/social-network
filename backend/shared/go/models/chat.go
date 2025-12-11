package models

import (
	ct "social-network/shared/go/customtypes"
)

type AddConversationMembersParams struct {
	ConversationId ct.Id
	UserIds        ct.Ids
}

type CreatePrivateConvParams struct {
	UserA ct.Id `json:"user_a"`
	UserB ct.Id `json:"user_b"`
}

type CreateGroupConvParams struct {
	GroupId ct.Id  `json:"group_id"`
	UserIds ct.Ids `json:"users_id"`
}

type ConversationDeleteResp struct {
	Id        ct.Id
	GroupId   ct.Id
	CreatedAt ct.GenDateTime
	UpdatedAt ct.GenDateTime
	DeletedAt ct.GenDateTime
}

type ConversationResponse struct {
	Id        ct.Id
	GroupId   ct.Id
	CreatedAt ct.GenDateTime
	UpdatedAt ct.GenDateTime `validation:"nullable"`
	DeletedAt ct.GenDateTime `validation:"nullable"`
}

type GetConversationMembersParams struct {
	ConversationID ct.Id
	UserID         ct.Id
}

type GetUserConversationsParams struct {
	UserId  ct.Id
	GroupId ct.Id
	Limit   ct.Limit
	Offset  ct.Offset
}

type GetUserConversationsRow struct {
	ConversationId       ct.Id
	CreatedAt            ct.GenDateTime
	UpdatedAt            ct.GenDateTime
	MemberIds            []int64
	UnreadCount          int64
	FirstUnreadMessageId *int64
}

type AddMembersToGroupConversationParams struct {
	GroupId ct.Id
	UserIds ct.Ids
}

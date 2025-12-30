package models

import (
	ct "social-network/shared/go/ct"
)

type AddConversationMembersParams struct {
	ConversationId ct.Id
	UserIds        ct.Ids
}

type AddMembersToGroupConversationParams struct {
	GroupId ct.Id
	UserIds ct.Ids
}

type CreatePrivateConvParams struct {
	UserA ct.Id `json:"user_a"`
	UserB ct.Id `json:"user_b"`
}

type CreateGroupConvParams struct {
	GroupId ct.Id  `json:"group_id"`
	UserIds ct.Ids `json:"users_id"`
}

type CreateMessageParams struct {
	ConversationId ct.Id
	SenderId       ct.Id
	MessageText    ct.MsgBody
}

type ConversationDeleteResp struct {
	Id        ct.Id
	GroupId   ct.Id
	CreatedAt ct.GenDateTime
	UpdatedAt ct.GenDateTime
	DeletedAt ct.GenDateTime
}

type ConversationResponse struct {
	Id             ct.Id
	GroupId        ct.Id
	LastMessageId  ct.Id
	FirstMessageId ct.Id
	CreatedAt      ct.GenDateTime
	UpdatedAt      ct.GenDateTime `validation:"nullable"`
	DeletedAt      ct.GenDateTime `validation:"nullable"`
}

type ConversationMember struct {
	ConversationId    ct.Id
	UserId            ct.Id
	LastReadMessageId ct.Id `validation:"nullable"`
	JoinedAt          ct.GenDateTime
	DeletedAt         ct.GenDateTime `validation:"nullable"`
}

// All fields are required except LastReadMessgeId
type ConversationMemberDeleted struct {
	ConversationId    ct.Id
	UserId            ct.Id
	LastReadMessageId ct.Id `validation:"nullable"`
	JoinedAt          ct.GenDateTime
	DeletedAt         ct.GenDateTime
}

type GetConversationMembersParams struct {
	ConversationId ct.Id
	UserID         ct.Id
}

type GetPrevMessagesParams struct {
	UserId            ct.Id
	ConversationId    ct.Id
	BoundaryMessageId ct.Id `validation:"nullable"`
	Limit             ct.Limit
	HydrateUsers      bool
}

type GetPrevMessagesResp struct {
	FirstMessageId ct.Id
	HaveMoreBefore bool
	Messages       []MessageResp
}

type GetNextMessageParams struct {
	BoundaryMessageId ct.Id
	ConversationId    ct.Id
	UserId            ct.Id
	Limit             ct.Limit
	RetrieveUsers     bool
}

type GetNextMessagesResp struct {
	LastMessageId ct.Id
	HaveMoreAfter bool
	Messages      []MessageResp
}

type GetUserConversationsParams struct {
	UserId       ct.Id
	GroupId      ct.Id `validation:"nullable"`
	Limit        ct.Limit
	Offset       ct.Offset
	HydrateUsers bool
}

type GetUserConversationsResp struct {
	ConversationId    ct.Id
	CreatedAt         ct.GenDateTime
	UpdatedAt         ct.GenDateTime
	Members           []User
	UnreadCount       int64
	LastReadMessageId ct.Id `validation:"nullable"`
}

// All fields are required except deleted at which in most cases is null.
type MessageResp struct {
	Id             ct.Id
	ConversationID ct.Id
	Sender         User
	MessageText    ct.MsgBody
	CreatedAt      ct.GenDateTime
	UpdatedAt      ct.GenDateTime
	DeletedAt      ct.GenDateTime `validation:"nullable"`
}

type DeleteConversationMemberParams struct {
	ConversationID ct.Id
	Owner          ct.Id
	ToDelete       ct.Id
}

// Last Read message is not nullable. If it is null then request is invalid.
type UpdateLastReadMessageParams struct {
	ConversationId    ct.Id
	UserID            ct.Id
	LastReadMessageId ct.Id
}

package models

import (
	ct "social-network/shared/go/ct"
)

// ================================
// PMs
// ================================

type CreateGroupConvReq struct {
	GroupId ct.Id  `json:"group_id"`
	UserIds ct.Ids `json:"user_ids"`
}

type CreateGroupMsgReq struct {
	GroupId     ct.Id      `json:"group_id"`
	SenderId    ct.Id      `json:"sender_id"`
	MessageText ct.MsgBody `json:"message_text"`
}

type GetGroupMsgsReq struct {
	GroupId           ct.Id    `json:"user_id"`
	MemberId          ct.Id    `json:"member_id"`
	BoundaryMessageId ct.Id    `json:"boundary_message_id" validation:"nullable"`
	Limit             ct.Limit `json:"limit"`
}

type GroupMsg struct {
	Id            ct.Id
	ConvesationId ct.Id
	GroupId       ct.Id
	Sender        User
	MessageText   ct.MsgBody
	CreatedAt     ct.GenDateTime `validation:"nullable"`
	UpdatedAt     ct.GenDateTime `validation:"nullable"`
	DeletedAt     ct.GenDateTime `validation:"nullable"`
}

type GetGetGroupMsgsResp struct {
	HaveMore bool
	Messages []GroupMsg
}

// ================================
// PMs
// ================================

type GetOrCreatePrivateConvReq struct {
	UserId            ct.Id `json:"user"`
	OtherUserId       ct.Id `json:"other_user"`
	RetrieveOtherUser bool  `json:"retrieve_other_user"`
}

type GetOrCreatePrivateConvResp struct {
	ConversationId  ct.Id
	OtherUser       User
	LastReadMessage ct.Id `validation:"nullable"`
	IsNew           bool
}

type CreatePrivatMsgReq struct {
	ConversationId ct.Id      `json:"conversation_id"`
	SenderId       ct.Id      `json:"sender_id"`
	MessageText    ct.MsgBody `json:"message_text"`
}

type GetPrivatMsgsReq struct {
	ConversationId    ct.Id    `json:"conversation_id"`
	UserId            ct.Id    `json:"user_id"`
	BoundaryMessageId ct.Id    `json:"boundary_message_id" validation:"nullable"`
	Limit             ct.Limit `json:"limit"`
	RetrieveUsers     bool     `json:"retrieve_users"`
}

type GetPrivateMsgsResp struct {
	HaveMore bool
	Messages []PrivateMsg
}

type GetPrivateConvsReq struct {
	UserId     ct.Id          `json:"user_id"`
	BeforeDate ct.GenDateTime `json:"before_date"`
	Limit      ct.Limit       `json:"limit"`
}

type PrivateConvsPreview struct {
	ConversationId ct.Id
	UpdatedAt      ct.GenDateTime
	OtherUser      User
	LastMessage    PrivateMsg
	UnreadCount    int
}

type PrivateMsg struct {
	Id             ct.Id
	ConversationID ct.Id
	Sender         User
	MessageText    ct.MsgBody
	CreatedAt      ct.GenDateTime `validation:"nullable"`
	UpdatedAt      ct.GenDateTime `validation:"nullable"`
	DeletedAt      ct.GenDateTime `validation:"nullable"`
}

type UpdateLastReadMsgParams struct {
	ConversationId    ct.Id `json:"conversation_id"`
	UserId            ct.Id `json:"user_id"`
	LastReadMessageId ct.Id `json:"last_read_message_id"`
}

type ConversationMember struct {
	ConversationId    ct.Id
	UserId            ct.Id
	LastReadMessageId ct.Id `validation:"nullable"`
	JoinedAt          ct.GenDateTime
	DeletedAt         ct.GenDateTime `validation:"nullable"`
}

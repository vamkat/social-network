package models

import (
	ct "social-network/shared/go/ct"
)

// ================================
// Group Conversations
// ================================

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

type GetGroupMsgsResp struct {
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

type CreatePrivateMsgReq struct {
	ConversationId ct.Id      `json:"conversation_id"`
	SenderId       ct.Id      `json:"sender_id"`
	MessageText    ct.MsgBody `json:"message_text"`
}

type GetPrivateMsgsReq struct {
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
	UserId            ct.Id          `json:"user_id"`
	BeforeDateUpdated ct.GenDateTime `json:"before_date_updated"`
	Limit             ct.Limit       `json:"limit"`
}

type PrivateConvsPreview struct {
	ConversationId ct.Id
	UpdatedAt      ct.GenDateTime
	OtherUser      User
	LastMessage    PrivateMsg
	UnreadCount    int
}

type PrivateMsg struct {
	Id             ct.Id          `json:"id"`
	ConversationId ct.Id          `json:"conversation_id"`
	Sender         User           `json:"sender"`
	ReceiverId     ct.Id          `json:"receiver_id,omitempty" validation:"nullable"`
	MessageText    ct.MsgBody     `json:"message_text"`
	CreatedAt      ct.GenDateTime `json:"created_at" validation:"nullable"`
	UpdatedAt      ct.GenDateTime `json:"updated_at" validation:"nullable"`
	DeletedAt      ct.GenDateTime `json:"deleted_at" validation:"nullable"`
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

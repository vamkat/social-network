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
	UserId            ct.Id    `json:"member_id"`
	BoundaryMessageId ct.Id    `json:"boundary_message_id" validate:"nullable"`
	Limit             ct.Limit `json:"limit"`
}

type GroupMsg struct {
	Id             ct.Id
	ConversationId ct.Id
	GroupId        ct.Id
	Sender         User
	MessageText    ct.MsgBody
	CreatedAt      ct.GenDateTime `validate:"nullable"`
	UpdatedAt      ct.GenDateTime `validate:"nullable"`
	DeletedAt      ct.GenDateTime `validate:"nullable"`
}

type GetGroupMsgsResp struct {
	HaveMore bool
	Messages []GroupMsg
}

// ================================
// PMs
// ================================

type GetOrCreatePrivateConvReq struct {
	UserId               ct.Id `json:"user"`
	InterlocutorId       ct.Id `json:"interlocutor"`
	RetrieveInterlocutor bool  `json:"retrieve_interlocutor"`
}

// DEPRECATED
type GetOrCreatePrivateConvResp struct {
	ConversationId  ct.Id
	Interlocutor    User
	LastReadMessage ct.Id `validate:"nullable"`
	IsNew           bool
}

type CreatePrivateMsgReq struct {
	SenderId       ct.Id      `json:"sender_id"`
	InterlocutorId ct.Id      `json:"interlocutor_id"`
	MessageText    ct.MsgBody `json:"message_text"`
}

type GetPrivateMsgsReq struct {
	UserId            ct.Id    `json:"user_id"`
	InterlocutorId    ct.Id    `json:"interlocutor_id"`
	BoundaryMessageId ct.Id    `json:"boundary_message_id" validate:"nullable"`
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
type GetPrivateConvByIdReq struct {
	UserId         ct.Id `json:"user_id"`
	ConversationId ct.Id `json:"conversation_id"`
	InterlocutorId ct.Id `json:"interlocutor_id"`
}

type PrivateConvsPreview struct {
	ConversationId ct.Id
	UpdatedAt      ct.GenDateTime
	Interlocutor   User
	LastMessage    PrivateMsg
	UnreadCount    int
}

type PrivateMsg struct {
	Id             ct.Id          `json:"id"`
	ConversationId ct.Id          `json:"conversation_id"`
	Sender         User           `json:"sender"`
	ReceiverId     ct.Id          `json:"receiver_id,omitempty" validate:"nullable"`
	MessageText    ct.MsgBody     `json:"message_text"`
	CreatedAt      ct.GenDateTime `json:"created_at" validate:"nullable"`
	UpdatedAt      ct.GenDateTime `json:"updated_at" validate:"nullable"`
	DeletedAt      ct.GenDateTime `json:"deleted_at" validate:"nullable"`
}

type UpdateLastReadMsgParams struct {
	ConversationId    ct.Id `json:"conversation_id"`
	UserId            ct.Id `json:"user_id"`
	LastReadMessageId ct.Id `json:"last_read_message_id"`
}

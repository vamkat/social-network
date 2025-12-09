package models

import (
	ct "social-network/shared/go/customtypes"
	"time"
)

type CreatePrivateConvParams struct {
	UserA ct.Id `json:"user_a"`
	UserB ct.Id `json:"user_b"`
}

type CreateGroupConvParams struct {
	GroupId ct.Id  `json:"group_id"`
	UserIds ct.Ids `json:"users_id"`
}

type Conversation struct {
	ID        ct.Id
	GroupID   ct.Id
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type AddMembersToGroupConversationParams struct {
	GroupID ct.Id
	UserIds ct.Ids
}

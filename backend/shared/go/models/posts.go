package models

import (
	ct "social-network/shared/go/customtypes"
)

type SimpleIdReq struct {
	Id ct.Id
}

type GenericReq struct {
	RequesterId ct.Id
	EntityId    ct.Id `json:"entity_id"`
}

type EntityIdPaginatedReq struct {
	RequesterId ct.Id
	EntityId    ct.Id     `json:"entity_id"`
	Limit       ct.Limit  `json:"limit"`
	Offset      ct.Offset `json:"offset"`
}

type GenericPaginatedReq struct {
	RequesterId ct.Id
	Limit       ct.Limit  `json:"limit"`
	Offset      ct.Offset `json:"offset"`
}

// -------------------------------------------
// Posts
// -------------------------------------------
type Post struct {
	PostId          ct.Id          `json:"post_id"`
	Body            ct.PostBody    `json:"post_body"`
	User            User           `json:"post_user"`
	GroupId         ct.Id          `json:"group_id" validate:"nullable"`
	Audience        ct.Audience    `json:"audience"`
	CommentsCount   int            `json:"comments_count"`
	ReactionsCount  int            `json:"reactions_count"`
	LastCommentedAt ct.GenDateTime `json:"last_commented_at"`
	CreatedAt       ct.GenDateTime `json:"created_at"`
	UpdatedAt       ct.GenDateTime `json:"updated_at" validate:"nullable"`
	LikedByUser     bool           `json:"liked_by_user"`
	ImageId         ct.Id          `json:"image" validate:"nullable"`
	ImageUrl        string         `json:"image_url"`
}

type CreatePostReq struct {
	CreatorId   ct.Id
	Body        ct.PostBody `json:"post_body"`
	GroupId     ct.Id       `json:"group_id" validate:"nullable"`
	Audience    ct.Audience `json:"audience"`
	AudienceIds ct.Ids      `json:"audience_ids" validate:"nullable"`
	ImageId     ct.Id       `json:"image" validate:"nullable"`
}

type EditPostReq struct {
	RequesterId ct.Id
	PostId      ct.Id
	NewBody     ct.PostBody `validate:"nullable"`
	ImageId     ct.Id       `validate:"nullable"`
	Audience    ct.Audience
	AudienceIds ct.Ids `validate:"nullable"`
}

type GetUserPostsReq struct {
	CreatorId   ct.Id `json:"creator_id"`
	RequesterId ct.Id
	Limit       ct.Limit  `json:"limit"`
	Offset      ct.Offset `json:"offset"`
}

type GetPersonalizedFeedReq struct {
	RequesterId ct.Id
	Limit       ct.Limit  `json:"limit"`
	Offset      ct.Offset `json:"offset"`
}

type GetGroupPostsReq struct {
	RequesterId ct.Id
	GroupId     ct.Id
	Limit       ct.Limit
	Offset      ct.Offset
}

//-------------------------------------------
// Comments
//-------------------------------------------

type Comment struct {
	CommentId      ct.Id          `json:"comment_id"`
	ParentId       ct.Id          `json:"parent_id"`
	Body           ct.CommentBody `json:"comment_body"`
	User           User           `json:"user"`
	ReactionsCount int            `json:"reactions_count"`
	CreatedAt      ct.GenDateTime `json:"created_at"`
	UpdatedAt      ct.GenDateTime `json:"updated_at"`
	LikedByUser    bool           `json:"liked_by_user"`
	ImageId        ct.Id          `json:"image" validate:"nullable"`
	ImageUrl       string         `json:"image_url"`
}

type CreateCommentReq struct {
	CreatorId ct.Id
	ParentId  ct.Id          `json:"parent_id"`
	Body      ct.CommentBody `json:"comment_body"`
	ImageId   ct.Id          `json:"image" validate:"nullable"`
}

type EditCommentReq struct {
	CreatorId ct.Id
	CommentId ct.Id          `json:"comment_id"`
	Body      ct.CommentBody `json:"comment_body" validate:"nullable"`
	ImageId   ct.Id          `json:"image" validate:"nullable"`
	ImageUrl  string         `json:"image_url"`
}

//-------------------------------------------
// Events
//-------------------------------------------

type Event struct {
	EventId       ct.Id
	Title         ct.Title
	Body          ct.EventBody
	User          User
	GroupId       ct.Id
	EventDate     ct.EventDateTime
	GoingCount    int
	NotGoingCount int
	ImageId       ct.Id  `validate:"nullable"`
	ImageUrl      string `json:"image_url"`
	CreatedAt     ct.GenDateTime
	UpdatedAt     ct.GenDateTime
	UserResponse  *bool
}

type CreateEventReq struct {
	Title     ct.Title
	Body      ct.EventBody
	CreatorId ct.Id
	GroupId   ct.Id
	ImageId   ct.Id `validate:"nullable"`
	EventDate ct.EventDateTime
}

type EditEventReq struct {
	EventId     ct.Id
	RequesterId ct.Id
	Title       ct.Title
	Body        ct.EventBody
	Image       ct.Id `validate:"nullable"`
	EventDate   ct.EventDateTime
}

type RespondToEventReq struct {
	EventId     ct.Id
	ResponderId ct.Id
	Going       bool
}

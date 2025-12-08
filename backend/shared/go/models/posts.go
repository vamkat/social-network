package models

import (
	ct "social-network/shared/go/customtypes"
	"time"
)

type SimpleIdReq struct {
	Id ct.Id
}

type GenericReq struct {
	RequesterId ct.Id
	EntityId    ct.Id
}

type EntityIdPaginatedReq struct {
	RequesterId ct.Id
	EntityId    ct.Id
	Limit       ct.Limit
	Offset      ct.Offset
}

type GenericPaginatedReq struct {
	RequesterId ct.Id
	Limit       ct.Limit
	Offset      ct.Offset
}

type HasUser interface {
	GetUserId() int64
	SetUser(User)
}

// -------------------------------------------
// Posts
// -------------------------------------------
type Post struct {
	PostId          ct.Id
	Body            ct.PostBody
	User            User
	GroupId         ct.Id `validate:"nullable"`
	Audience        ct.Audience
	CommentsCount   int
	ReactionsCount  int
	LastCommentedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LikedByUser     bool
	Image           ct.Id `validate:"nullable"`
}

func (p *Post) GetUserId() int64 { return p.User.UserId.Int64() }
func (p *Post) SetUser(u User)   { p.User = u }

type CreatePostReq struct {
	CreatorId   ct.Id
	Body        ct.PostBody
	GroupId     ct.Id `validate:"nullable"`
	Audience    ct.Audience
	AudienceIds ct.Ids `validate:"nullable"`
	Image       ct.Id  `validate:"nullable"`
}

type EditPostReq struct {
	RequesterId ct.Id
	PostId      ct.Id
	NewBody     ct.PostBody `validate:"nullable"`
	Image       ct.Id       `validate:"nullable"`
	Audience    ct.Audience
	AudienceIds ct.Ids `validate:"nullable"`
}

type GetUserPostsReq struct {
	CreatorId   ct.Id
	RequesterId ct.Id
	Limit       ct.Limit
	Offset      ct.Offset
}

type GetPersonalizedFeedReq struct {
	RequesterId ct.Id
	Limit       ct.Limit
	Offset      ct.Offset
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
	CommentId      ct.Id
	ParentId       ct.Id
	Body           ct.CommentBody
	User           User
	ReactionsCount int
	CreatedAt      time.Time
	UpdatedAt      time.Time //can be nil
	LikedByUser    bool
	Image          ct.Id `validate:"nullable"`
}

func (c *Comment) GetUserId() int64 { return c.User.UserId.Int64() }
func (c *Comment) SetUser(u User)   { c.User = u }

type CreateCommentReq struct {
	CreatorId ct.Id
	ParentId  ct.Id
	Body      ct.CommentBody
	Image     ct.Id `validate:"nullable"`
}

type EditCommentReq struct {
	CreatorId ct.Id
	CommentId ct.Id
	Body      ct.CommentBody `validate:"nullable"`
	Image     ct.Id          `validate:"nullable"`
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
	EventDate     ct.EventDate
	GoingCount    int
	NotGoingCount int
	Image         ct.Id `validate:"nullable"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserResponse  *bool
}

func (e *Event) GetUserId() int64 { return e.User.UserId.Int64() }
func (e *Event) SetUser(u User)   { e.User = u }

type CreateEventReq struct {
	Title     ct.Title
	Body      ct.EventBody
	CreatorId ct.Id
	GroupId   ct.Id
	Image     ct.Id `validate:"nullable"`
	EventDate ct.EventDate
}

type EditEventReq struct {
	EventId     ct.Id
	RequesterId ct.Id
	Title       ct.Title
	Body        ct.EventBody
	Image       ct.Id `validate:"nullable"`
	EventDate   ct.EventDate
}

type RespondToEventReq struct {
	EventId     ct.Id
	ResponderId ct.Id
	Going       bool
}

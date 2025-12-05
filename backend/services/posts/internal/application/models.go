package application

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

type accessContext struct {
	requesterId int64
	entityId    int64
}

// -------------------------------------------
// Posts
// -------------------------------------------
type Post struct {
	PostId          ct.Id
	Body            ct.PostBody
	CreatorId       ct.Id
	GroupId         ct.Id `validate:"nullable"` //add check that if audience=group it can't be nil
	Audience        ct.Audience
	CommentsCount   int
	ReactionsCount  int
	LastCommentedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LikedByUser     bool
	Image           ct.Id `validate:"nullable"`
	LatestComment   Comment
}

type CreatePostReq struct {
	CreatorId       ct.Id
	Body            ct.PostBody
	GroupId         ct.Id `validate:"nullable"`
	Audience        ct.Audience
	AudienceIds     ct.Ids `validate:"nullable"`
	Image           ct.Id  `validate:"nullable"`
	RequesterGroups ct.Ids `validate:"nullable"`
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
	CreatorId        ct.Id
	CreatorFollowers ct.Ids `validate:"nullable"`
	RequesterId      ct.Id
	Limit            ct.Limit
	Offset           ct.Offset
}

type GetPersonalizedFeedReq struct {
	RequesterId         ct.Id
	RequesterFollowsIds ct.Ids //from user service
	Limit               ct.Limit
	Offset              ct.Offset
}

type GetGroupPostsReq struct {
	RequesterId     ct.Id
	GroupId         ct.Id
	Limit           ct.Limit
	Offset          ct.Offset
	RequesterGroups ct.Ids `validate:"nullable"`
}

//-------------------------------------------
// Comments
//-------------------------------------------

type Comment struct {
	CommentId      ct.Id
	ParentId       ct.Id
	Body           ct.CommentBody
	CreatorId      ct.Id
	ReactionsCount int
	CreatedAt      time.Time
	UpdatedAt      time.Time //can be nil
	LikedByUser    bool
	Image          ct.Id `validate:"nullable"`
}

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
	CreatorId     ct.Id
	GroupId       ct.Id
	EventDate     ct.EventDate
	GoingCount    int
	NotGoingCount int
	Image         ct.Id `validate:"nullable"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserResponse  *bool
}

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

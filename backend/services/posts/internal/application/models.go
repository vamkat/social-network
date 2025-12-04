package application

import (
	ct "social-network/shared/go/customtypes"
	"time"
)

type GenericReq struct {
	RequesterId ct.Id
	EntityId    ct.Id
}

type GenericPaginatedReq struct { //two different ones with nullable and not?
	RequesterId ct.Id
	EntityId    ct.Id `validate:"nullable"`
	Limit       ct.Limit
	Offset      ct.Offset
}

type hasRightToView struct {
	RequesterId         ct.Id
	ParentEntityId      ct.Id
	RequesterFollowsIds ct.Ids `validate:"nullable"`
	RequesterGroups     ct.Ids `validate:"nullable"`
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

type insertPostAudienceReq struct {
	PostId      ct.Id
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

//-------------------------------------------
// Reactions
//-------------------------------------------

//-------------------------------------------
// Images
//-------------------------------------------

type ImageReq struct {
	RequesterId ct.Id
	PostId      ct.Id
	Image       ct.Id
}

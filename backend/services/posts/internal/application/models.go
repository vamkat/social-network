package application

import (
	ct "social-network/shared/go/customtypes"
	"time"
)

type GenericReq struct {
	RequesterId ct.Id
	EntityId    ct.Id
}

type GenericPaginatedReq struct {
	RequesterId ct.Id
	EntityId    ct.Id `validate:"nullable"`
	Limit       ct.Limit
	Offset      ct.Offset
}

type hasRightToView struct {
	RequesterId         ct.Id
	ParentEntityId      ct.Id
	RequesterFollowsIds []ct.Id
	RequesterGroups     []ct.Id
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
	ImagesCount     int
	LastCommentedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LikedByUser     bool
	FirstImage      string
}

type CreatePostReq struct {
	CreatorId   ct.Id
	Body        ct.PostBody
	GroupId     ct.Id `validate:"nullable"`
	Audience    ct.Audience
	AudienceIds []ct.Id `validate:"nullable"`
}

type EditPostContentReq struct {
	RequesterId ct.Id
	PostId      ct.Id
	NewBody     ct.PostBody
}

type EditPostAudienceReq struct {
	RequesterId ct.Id
	PostId      ct.Id
	Audience    ct.Audience
	AudienceIds []ct.Id `validate:"nullable"`
}

type GetUserPostsReq struct {
	CreatorId        ct.Id
	CreatorFollowers []ct.Id `validate:"nullable"`
	RequesterId      ct.Id
	Limit            ct.Limit
	Offset           ct.Offset
}

type GetPersonalizedFeedReq struct {
	RequesterId         ct.Id
	RequesterFollowsIds []ct.Id //from user service
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
	ImagesCount    int
	CreatedAt      time.Time
	UpdatedAt      time.Time //can be nil
	LikedByUser    bool
	FirstImage     string //can be nil (or "")
}

type CreateCommentReq struct {
	CreatorId ct.Id
	ParentId  ct.Id
	Body      ct.CommentBody
}

type EditCommentReq struct {
	CreatorId ct.Id
	CommentId ct.Id
	Body      ct.CommentBody
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
	StillValid    bool //still not sure this is needed
	GoingCount    int
	NotGoingCount int
	ImagesCount   int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateEventReq struct {
	Title     ct.Title
	Body      ct.EventBody
	CreatorId ct.Id
	GroupId   ct.Id
	EventDate ct.EventDate
}

type EditEventReq struct {
	EventId     ct.Id
	RequesterId ct.Id
	Title       ct.Title
	Body        ct.EventBody
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

type InsertImagesReq struct {
	PostId ct.Id
	Images []string
}

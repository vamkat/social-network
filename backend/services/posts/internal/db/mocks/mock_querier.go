package mocks

import (
	"context"
	"social-network/services/posts/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
)

type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CanUserSeeEntity(ctx context.Context, arg sqlc.CanUserSeeEntityParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) ClearPostAudience(ctx context.Context, postID int64) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockQuerier) CreateComment(ctx context.Context, arg sqlc.CreateCommentParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) CreateEvent(ctx context.Context, arg sqlc.CreateEventParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) CreatePost(ctx context.Context, arg sqlc.CreatePostParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) GetPostByID(ctx context.Context, arg sqlc.GetPostByIDParams) (sqlc.GetPostByIDRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.GetPostByIDRow), args.Error(1)
}

func (m *MockQuerier) DeleteComment(ctx context.Context, arg sqlc.DeleteCommentParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) DeleteEvent(ctx context.Context, arg sqlc.DeleteEventParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) DeleteEventResponse(ctx context.Context, arg sqlc.DeleteEventResponseParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) DeleteImage(ctx context.Context, id int64) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) DeletePost(ctx context.Context, arg sqlc.DeletePostParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) EditComment(ctx context.Context, arg sqlc.EditCommentParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) EditEvent(ctx context.Context, arg sqlc.EditEventParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) EditPostContent(ctx context.Context, arg sqlc.EditPostContentParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) GetCommentsByPostId(ctx context.Context, arg sqlc.GetCommentsByPostIdParams) ([]sqlc.GetCommentsByPostIdRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetCommentsByPostIdRow), args.Error(1)
}

func (m *MockQuerier) GetEntityCreatorAndGroup(ctx context.Context, id int64) (sqlc.GetEntityCreatorAndGroupRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlc.GetEntityCreatorAndGroupRow), args.Error(1)
}

func (m *MockQuerier) GetEventsByGroupId(ctx context.Context, arg sqlc.GetEventsByGroupIdParams) ([]sqlc.GetEventsByGroupIdRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetEventsByGroupIdRow), args.Error(1)
}

func (m *MockQuerier) GetGroupPostsPaginated(ctx context.Context, arg sqlc.GetGroupPostsPaginatedParams) ([]sqlc.GetGroupPostsPaginatedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetGroupPostsPaginatedRow), args.Error(1)
}

func (m *MockQuerier) GetImages(ctx context.Context, parentID int64) (int64, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) GetLatestCommentforPostId(ctx context.Context, arg sqlc.GetLatestCommentforPostIdParams) (sqlc.GetLatestCommentforPostIdRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.GetLatestCommentforPostIdRow), args.Error(1)
}

func (m *MockQuerier) GetMostPopularPostInGroup(ctx context.Context, groupID pgtype.Int8) (sqlc.GetMostPopularPostInGroupRow, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(sqlc.GetMostPopularPostInGroupRow), args.Error(1)
}

func (m *MockQuerier) GetPersonalizedFeed(ctx context.Context, arg sqlc.GetPersonalizedFeedParams) ([]sqlc.GetPersonalizedFeedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetPersonalizedFeedRow), args.Error(1)
}

func (m *MockQuerier) GetPostAudience(ctx context.Context, postID int64) ([]int64, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockQuerier) GetPublicFeed(ctx context.Context, arg sqlc.GetPublicFeedParams) ([]sqlc.GetPublicFeedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetPublicFeedRow), args.Error(1)
}

func (m *MockQuerier) GetUserPostsPaginated(ctx context.Context, arg sqlc.GetUserPostsPaginatedParams) ([]sqlc.GetUserPostsPaginatedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetUserPostsPaginatedRow), args.Error(1)
}

func (m *MockQuerier) GetWhoLikedEntityId(ctx context.Context, contentID int64) ([]int64, error) {
	args := m.Called(ctx, contentID)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockQuerier) InsertPostAudience(ctx context.Context, arg sqlc.InsertPostAudienceParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) SuggestUsersByPostActivity(ctx context.Context, creatorID int64) ([]int64, error) {
	args := m.Called(ctx, creatorID)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockQuerier) ToggleOrInsertReaction(ctx context.Context, arg sqlc.ToggleOrInsertReactionParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) UpdatePostAudience(ctx context.Context, arg sqlc.UpdatePostAudienceParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) UpsertEventResponse(ctx context.Context, arg sqlc.UpsertEventResponseParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) UpsertImage(ctx context.Context, arg sqlc.UpsertImageParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

package mocks

import (
	"context"
	"social-network/services/posts/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
)

type MockQueries struct {
	mock.Mock
}

func (m *MockQueries) CanUserSeeEntity(ctx context.Context, arg sqlc.CanUserSeeEntityParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueries) ClearPostAudience(ctx context.Context, postID int64) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockQueries) CreateComment(ctx context.Context, arg sqlc.CreateCommentParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQueries) CreateEvent(ctx context.Context, arg sqlc.CreateEventParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQueries) CreatePost(ctx context.Context, arg sqlc.CreatePostParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) GetPostByID(ctx context.Context, arg sqlc.GetPostByIDParams) (sqlc.GetPostByIDRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.GetPostByIDRow), args.Error(1)
}

func (m *MockQueries) DeleteComment(ctx context.Context, arg sqlc.DeleteCommentParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) DeleteEvent(ctx context.Context, arg sqlc.DeleteEventParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) DeleteEventResponse(ctx context.Context, arg sqlc.DeleteEventResponseParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) DeleteImage(ctx context.Context, id int64) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) DeletePost(ctx context.Context, arg sqlc.DeletePostParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) EditComment(ctx context.Context, arg sqlc.EditCommentParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) EditEvent(ctx context.Context, arg sqlc.EditEventParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) EditPostContent(ctx context.Context, arg sqlc.EditPostContentParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) GetCommentsByPostId(ctx context.Context, arg sqlc.GetCommentsByPostIdParams) ([]sqlc.GetCommentsByPostIdRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetCommentsByPostIdRow), args.Error(1)
}

func (m *MockQueries) GetEntityCreatorAndGroup(ctx context.Context, id int64) (sqlc.GetEntityCreatorAndGroupRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlc.GetEntityCreatorAndGroupRow), args.Error(1)
}

func (m *MockQueries) GetEventsByGroupId(ctx context.Context, arg sqlc.GetEventsByGroupIdParams) ([]sqlc.GetEventsByGroupIdRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetEventsByGroupIdRow), args.Error(1)
}

func (m *MockQueries) GetGroupPostsPaginated(ctx context.Context, arg sqlc.GetGroupPostsPaginatedParams) ([]sqlc.GetGroupPostsPaginatedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetGroupPostsPaginatedRow), args.Error(1)
}

func (m *MockQueries) GetImages(ctx context.Context, parentID int64) (int64, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) GetLatestCommentforPostId(ctx context.Context, arg sqlc.GetLatestCommentforPostIdParams) (sqlc.GetLatestCommentforPostIdRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.GetLatestCommentforPostIdRow), args.Error(1)
}

func (m *MockQueries) GetMostPopularPostInGroup(ctx context.Context, groupID pgtype.Int8) (sqlc.GetMostPopularPostInGroupRow, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(sqlc.GetMostPopularPostInGroupRow), args.Error(1)
}

func (m *MockQueries) GetPersonalizedFeed(ctx context.Context, arg sqlc.GetPersonalizedFeedParams) ([]sqlc.GetPersonalizedFeedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetPersonalizedFeedRow), args.Error(1)
}

func (m *MockQueries) GetPostAudience(ctx context.Context, postID int64) ([]int64, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockQueries) GetPublicFeed(ctx context.Context, arg sqlc.GetPublicFeedParams) ([]sqlc.GetPublicFeedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetPublicFeedRow), args.Error(1)
}

func (m *MockQueries) GetUserPostsPaginated(ctx context.Context, arg sqlc.GetUserPostsPaginatedParams) ([]sqlc.GetUserPostsPaginatedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetUserPostsPaginatedRow), args.Error(1)
}

func (m *MockQueries) GetWhoLikedEntityId(ctx context.Context, contentID int64) ([]int64, error) {
	args := m.Called(ctx, contentID)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockQueries) InsertPostAudience(ctx context.Context, arg sqlc.InsertPostAudienceParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) SuggestUsersByPostActivity(ctx context.Context, creatorID int64) ([]int64, error) {
	args := m.Called(ctx, creatorID)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockQueries) ToggleOrInsertReaction(ctx context.Context, arg sqlc.ToggleOrInsertReactionParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) UpdatePostAudience(ctx context.Context, arg sqlc.UpdatePostAudienceParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) UpsertEventResponse(ctx context.Context, arg sqlc.UpsertEventResponseParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) UpsertImage(ctx context.Context, arg sqlc.UpsertImageParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQueries) WithTx(ctx context.Context, fn func(*MockQueries) error) error {
	m.Called(ctx)
	return fn(m)
}

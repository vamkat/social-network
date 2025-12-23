package mocks

import (
	"context"
	ds "social-network/services/posts/internal/db/dbservice"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
)

type MockQueries struct {
	mock.Mock
}

func (m *MockQueries) CanUserSeeEntity(ctx context.Context, arg ds.CanUserSeeEntityParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueries) ClearPostAudience(ctx context.Context, postID int64) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockQueries) CreateComment(ctx context.Context, arg ds.CreateCommentParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQueries) CreateEvent(ctx context.Context, arg ds.CreateEventParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQueries) CreatePost(ctx context.Context, arg ds.CreatePostParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) GetPostByID(ctx context.Context, arg ds.GetPostByIDParams) (ds.GetPostByIDRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(ds.GetPostByIDRow), args.Error(1)
}

func (m *MockQueries) DeleteComment(ctx context.Context, arg ds.DeleteCommentParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) DeleteEvent(ctx context.Context, arg ds.DeleteEventParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) DeleteEventResponse(ctx context.Context, arg ds.DeleteEventResponseParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) DeleteImage(ctx context.Context, id int64) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) DeletePost(ctx context.Context, arg ds.DeletePostParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) EditComment(ctx context.Context, arg ds.EditCommentParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) EditEvent(ctx context.Context, arg ds.EditEventParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) EditPostContent(ctx context.Context, arg ds.EditPostContentParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) GetCommentsByPostId(ctx context.Context, arg ds.GetCommentsByPostIdParams) ([]ds.GetCommentsByPostIdRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]ds.GetCommentsByPostIdRow), args.Error(1)
}

func (m *MockQueries) GetEntityCreatorAndGroup(ctx context.Context, id int64) (ds.GetEntityCreatorAndGroupRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ds.GetEntityCreatorAndGroupRow), args.Error(1)
}

func (m *MockQueries) GetEventsByGroupId(ctx context.Context, arg ds.GetEventsByGroupIdParams) ([]ds.GetEventsByGroupIdRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]ds.GetEventsByGroupIdRow), args.Error(1)
}

func (m *MockQueries) GetGroupPostsPaginated(ctx context.Context, arg ds.GetGroupPostsPaginatedParams) ([]ds.GetGroupPostsPaginatedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]ds.GetGroupPostsPaginatedRow), args.Error(1)
}

func (m *MockQueries) GetImages(ctx context.Context, parentID int64) (int64, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) GetLatestCommentforPostId(ctx context.Context, arg ds.GetLatestCommentforPostIdParams) (ds.GetLatestCommentforPostIdRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(ds.GetLatestCommentforPostIdRow), args.Error(1)
}

func (m *MockQueries) GetMostPopularPostInGroup(ctx context.Context, groupID pgtype.Int8) (ds.GetMostPopularPostInGroupRow, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(ds.GetMostPopularPostInGroupRow), args.Error(1)
}

func (m *MockQueries) GetPersonalizedFeed(ctx context.Context, arg ds.GetPersonalizedFeedParams) ([]ds.GetPersonalizedFeedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]ds.GetPersonalizedFeedRow), args.Error(1)
}

func (m *MockQueries) GetPostAudience(ctx context.Context, postID int64) ([]int64, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockQueries) GetPublicFeed(ctx context.Context, arg ds.GetPublicFeedParams) ([]ds.GetPublicFeedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]ds.GetPublicFeedRow), args.Error(1)
}

func (m *MockQueries) GetUserPostsPaginated(ctx context.Context, arg ds.GetUserPostsPaginatedParams) ([]ds.GetUserPostsPaginatedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]ds.GetUserPostsPaginatedRow), args.Error(1)
}

func (m *MockQueries) GetWhoLikedEntityId(ctx context.Context, contentID int64) ([]int64, error) {
	args := m.Called(ctx, contentID)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockQueries) InsertPostAudience(ctx context.Context, arg ds.InsertPostAudienceParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) SuggestUsersByPostActivity(ctx context.Context, creatorID int64) ([]int64, error) {
	args := m.Called(ctx, creatorID)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockQueries) ToggleOrInsertReaction(ctx context.Context, arg ds.ToggleOrInsertReactionParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) UpdatePostAudience(ctx context.Context, arg ds.UpdatePostAudienceParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) UpsertEventResponse(ctx context.Context, arg ds.UpsertEventResponseParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueries) UpsertImage(ctx context.Context, arg ds.UpsertImageParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQueries) WithTx(ctx context.Context, fn func(*MockQueries) error) error {
	m.Called(ctx)
	return fn(m)
}

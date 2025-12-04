package application

import (
	"context"
	"database/sql"
	"errors"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// MOCKS
// ============================================

// MockQuerier is a mock implementation of sqlc.Querier
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

// MockTxRunner is a mock implementation of TxRunner
type MockTxRunner struct {
	mock.Mock
}

func (m *MockTxRunner) RunTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	// We call the mock but DON'T execute fn
	// This means we're testing that CreatePost calls RunTx correctly
	// but we're NOT testing what happens inside the transaction
	// (that would require integration tests or refactoring to pass Querier interface)
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// ============================================
// TEST SETUP HELPERS
// ============================================

func setupTestApp() (*Application, *MockQuerier) {
	mockDB := new(MockQuerier)
	app := &Application{
		db:       mockDB,
		txRunner: nil,
	}
	return app, mockDB
}

func setupTestAppWithTx() (*Application, *MockQuerier, *MockTxRunner) {
	mockDB := new(MockQuerier)
	mockTx := new(MockTxRunner)

	app := NewApplicationWithTxRunner(mockDB, mockTx)

	return app, mockDB, mockTx
}

// ============================================
// TESTS
// ============================================

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name          string
		req           CreatePostReq
		setupMock     func(*MockTxRunner)
		expectedError error
	}{
		{
			name: "successful post creation - transaction called",
			req: CreatePostReq{
				Body:            ct.PostBody("Test post content"),
				CreatorId:       ct.Id(1),
				GroupId:         ct.Id(10),
				Audience:        ct.Audience("everyone"),
				RequesterGroups: []ct.Id{10, 20},
			},
			setupMock: func(tx *MockTxRunner) {
				// Just verify RunTx is called and return success
				tx.On("RunTx", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "transaction returns error",
			req: CreatePostReq{
				Body:            ct.PostBody("Test post"),
				CreatorId:       ct.Id(1),
				GroupId:         ct.Id(10),
				Audience:        ct.Audience("everyone"),
				RequesterGroups: []ct.Id{10},
			},
			setupMock: func(tx *MockTxRunner) {
				tx.On("RunTx", mock.Anything, mock.Anything).Return(errors.New("transaction failed"))
			},
			expectedError: errors.New("transaction failed"),
		},
		{
			name: "user not member of group",
			req: CreatePostReq{
				Body:            ct.PostBody("Test post"),
				CreatorId:       ct.Id(1),
				GroupId:         ct.Id(10),
				Audience:        ct.Audience("everyone"),
				RequesterGroups: []ct.Id{20, 30},
			},
			setupMock:     func(tx *MockTxRunner) {},
			expectedError: ErrNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _, mockTx := setupTestAppWithTx()
			tt.setupMock(mockTx)

			err := app.CreatePost(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError.Error() != "" {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			mockTx.AssertExpectations(t)
		})
	}
}

func TestDeletePost(t *testing.T) {
	tests := []struct {
		name          string
		req           GenericReq
		setupMock     func(*MockQuerier)
		expectedError error
	}{
		{
			name: "successful deletion",
			req: GenericReq{
				EntityId:    ct.Id(100),
				RequesterId: ct.Id(1),
			},
			setupMock: func(m *MockQuerier) {
				m.On("DeletePost", mock.Anything, sqlc.DeletePostParams{
					ID:        100,
					CreatorID: 1,
				}).Return(int64(1), nil)
			},
			expectedError: nil,
		},
		{
			name: "post not found",
			req: GenericReq{
				EntityId:    ct.Id(999),
				RequesterId: ct.Id(1),
			},
			setupMock: func(m *MockQuerier) {
				m.On("DeletePost", mock.Anything, sqlc.DeletePostParams{
					ID:        999,
					CreatorID: 1,
				}).Return(int64(0), nil)
			},
			expectedError: ErrNotFound,
		},
		{
			name: "database error",
			req: GenericReq{
				EntityId:    ct.Id(100),
				RequesterId: ct.Id(1),
			},
			setupMock: func(m *MockQuerier) {
				m.On("DeletePost", mock.Anything, mock.Anything).
					Return(int64(0), errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockDB := setupTestApp()
			tt.setupMock(mockDB)

			err := app.DeletePost(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetGroupPostsPaginated(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		req           GetGroupPostsReq
		setupMock     func(*MockQuerier)
		expectedPosts int
		expectedError error
	}{
		{
			name: "successful retrieval",
			req: GetGroupPostsReq{
				GroupId:         ct.Id(10),
				RequesterId:     ct.Id(1),
				RequesterGroups: []ct.Id{10, 20},
			},
			setupMock: func(m *MockQuerier) {
				m.On("GetGroupPostsPaginated", mock.Anything, mock.MatchedBy(func(params sqlc.GetGroupPostsPaginatedParams) bool {
					return params.GroupID.Int64 == 10 && params.GroupID.Valid
				})).Return([]sqlc.GetGroupPostsPaginatedRow{
					{
						ID:                          100,
						PostBody:                    "Test post 1",
						CreatorID:                   1,
						Audience:                    "everyone",
						CommentsCount:               5,
						ReactionsCount:              10,
						LastCommentedAt:             pgtype.Timestamptz{Time: now, Valid: true},
						CreatedAt:                   pgtype.Timestamptz{Time: now, Valid: true},
						UpdatedAt:                   pgtype.Timestamptz{Time: now, Valid: true},
						LikedByUser:                 true,
						Image:                       200,
						LatestCommentID:             300,
						LatestCommentBody:           "Latest comment",
						LatestCommentCreatorID:      2,
						LatestCommentReactionsCount: 3,
						LatestCommentCreatedAt:      pgtype.Timestamptz{Time: now, Valid: true},
						LatestCommentLikedByUser:    false,
						LatestCommentImage:          0,
					},
				}, nil)
			},
			expectedPosts: 1,
			expectedError: nil,
		},
		{
			name: "user not member of group",
			req: GetGroupPostsReq{
				GroupId:         ct.Id(10),
				RequesterId:     ct.Id(1),
				RequesterGroups: []ct.Id{20, 30},
			},
			setupMock:     func(m *MockQuerier) {},
			expectedPosts: 0,
			expectedError: ErrNotAllowed,
		},
		{
			name: "no posts found",
			req: GetGroupPostsReq{
				GroupId:         ct.Id(10),
				RequesterId:     ct.Id(1),
				RequesterGroups: []ct.Id{10},
			},
			setupMock: func(m *MockQuerier) {
				m.On("GetGroupPostsPaginated", mock.Anything, mock.Anything).
					Return([]sqlc.GetGroupPostsPaginatedRow{}, nil)
			},
			expectedPosts: 0,
			expectedError: ErrNotFound,
		},
		{
			name: "no group ID provided",
			req: GetGroupPostsReq{
				GroupId:         ct.Id(0),
				RequesterId:     ct.Id(1),
				RequesterGroups: []ct.Id{10},
			},
			setupMock:     func(m *MockQuerier) {},
			expectedPosts: 0,
			expectedError: ErrNoGroupIdGiven,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockDB := setupTestApp()
			tt.setupMock(mockDB)

			posts, err := app.GetGroupPostsPaginated(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, posts)
			} else {
				assert.NoError(t, err)
				assert.Len(t, posts, tt.expectedPosts)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetMostPopularPostInGroup(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		groupID       ct.Id
		setupMock     func(*MockQuerier)
		expectPost    bool
		expectedError error
	}{
		{
			name:    "successful retrieval",
			groupID: ct.Id(10),
			setupMock: func(m *MockQuerier) {
				m.On("GetMostPopularPostInGroup", mock.Anything, mock.MatchedBy(func(id pgtype.Int8) bool {
					return id.Int64 == 10 && id.Valid
				})).Return(sqlc.GetMostPopularPostInGroupRow{
					ID:              100,
					PostBody:        "Most popular post",
					CreatorID:       1,
					Audience:        "everyone",
					CommentsCount:   50,
					ReactionsCount:  100,
					LastCommentedAt: pgtype.Timestamptz{Time: now, Valid: true},
					CreatedAt:       pgtype.Timestamptz{Time: now, Valid: true},
					UpdatedAt:       pgtype.Timestamptz{Time: now, Valid: true},
					Image:           200,
				}, nil)
			},
			expectPost:    true,
			expectedError: nil,
		},
		{
			name:    "no posts in group",
			groupID: ct.Id(10),
			setupMock: func(m *MockQuerier) {
				m.On("GetMostPopularPostInGroup", mock.Anything, mock.Anything).
					Return(sqlc.GetMostPopularPostInGroupRow{}, sql.ErrNoRows)
			},
			expectPost:    false,
			expectedError: ErrNotFound,
		},
		{
			name:          "no group ID provided",
			groupID:       ct.Id(0),
			setupMock:     func(m *MockQuerier) {},
			expectPost:    false,
			expectedError: ErrNoGroupIdGiven,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockDB := setupTestApp()
			tt.setupMock(mockDB)

			post, err := app.GetMostPopularPostInGroup(context.Background(), tt.groupID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				if tt.expectPost {
					assert.NotEmpty(t, post.PostId)
					assert.NotEmpty(t, post.Body)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetUserPostsPaginated(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		req           GetUserPostsReq
		setupMock     func(*MockQuerier)
		expectedPosts int
		expectedError error
	}{
		{
			name: "successful retrieval",
			req: GetUserPostsReq{
				CreatorId:        ct.Id(1),
				RequesterId:      ct.Id(2),
				CreatorFollowers: []ct.Id{2, 3},
				Limit:            ct.Limit(10),
				Offset:           ct.Offset(0),
			},
			setupMock: func(m *MockQuerier) {
				m.On("GetUserPostsPaginated", mock.Anything, mock.Anything).
					Return([]sqlc.GetUserPostsPaginatedRow{
						{
							ID:                          100,
							PostBody:                    "User post 1",
							CreatorID:                   1,
							CommentsCount:               5,
							ReactionsCount:              10,
							LastCommentedAt:             pgtype.Timestamptz{Time: now, Valid: true},
							CreatedAt:                   pgtype.Timestamptz{Time: now, Valid: true},
							UpdatedAt:                   pgtype.Timestamptz{Time: now, Valid: true},
							LikedByUser:                 true,
							Image:                       200,
							LatestCommentID:             300,
							LatestCommentBody:           "Latest comment",
							LatestCommentCreatorID:      2,
							LatestCommentReactionsCount: 3,
							LatestCommentCreatedAt:      pgtype.Timestamptz{Time: now, Valid: true},
							LatestCommentLikedByUser:    false,
							LatestCommentImage:          0,
						},
					}, nil)
			},
			expectedPosts: 1,
			expectedError: nil,
		},
		{
			name: "no posts found",
			req: GetUserPostsReq{
				CreatorId:        ct.Id(1),
				RequesterId:      ct.Id(2),
				CreatorFollowers: []ct.Id{2},
				Limit:            ct.Limit(10),
				Offset:           ct.Offset(0),
			},
			setupMock: func(m *MockQuerier) {
				m.On("GetUserPostsPaginated", mock.Anything, mock.Anything).
					Return([]sqlc.GetUserPostsPaginatedRow{}, nil)
			},
			expectedPosts: 0,
			expectedError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockDB := setupTestApp()
			tt.setupMock(mockDB)

			posts, err := app.GetUserPostsPaginated(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, posts)
			} else {
				assert.NoError(t, err)
				assert.Len(t, posts, tt.expectedPosts)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

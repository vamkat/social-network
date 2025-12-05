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

func (m *MockQuerier) GetEntityCreatorAndGroup(ctx context.Context, arg int64) (sqlc.GetEntityCreatorAndGroupRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.GetEntityCreatorAndGroupRow), args.Error(1)
}

// MockTxRunner is a mock implementation of TxRunner
type MockTxRunner struct {
	mock.Mock
	MockQuerier sqlc.Querier // The querier to pass to transaction functions
}

func (m *MockTxRunner) RunTx(ctx context.Context, fn func(sqlc.Querier) error) error {
	args := m.Called(ctx, fn)

	// Check if we should return an error before executing
	if err := args.Error(0); err != nil {
		return err
	}

	// Execute the function with our mock querier
	// Now this works because fn expects sqlc.Querier interface!
	return fn(m.MockQuerier)
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
	mockTx := &MockTxRunner{
		MockQuerier: mockDB, // Pass the mock querier to the transaction runner
	}

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
		setupMock     func(*MockQuerier, *MockTxRunner)
		expectedError error
	}{
		{
			name: "successful post creation with public audience",
			req: CreatePostReq{
				Body:            ct.PostBody("Test post content"),
				CreatorId:       ct.Id(1),
				GroupId:         ct.Id(10),
				Audience:        ct.Audience("everyone"),
				RequesterGroups: []ct.Id{10, 20},
			},
			setupMock: func(m *MockQuerier, tx *MockTxRunner) {
				// Mock RunTx to actually execute the function
				tx.On("RunTx", mock.Anything, mock.Anything).Return(nil)

				// Mock the database calls that happen inside the transaction
				m.On("CreatePost", mock.Anything, mock.MatchedBy(func(params sqlc.CreatePostParams) bool {
					return params.PostBody == "Test post content" &&
						params.CreatorID == 1 &&
						params.Audience == "everyone"
				})).Return(int64(100), nil)
			},
			expectedError: nil,
		},
		{
			name: "successful post creation with selected audience",
			req: CreatePostReq{
				Body:            ct.PostBody("Test post"),
				CreatorId:       ct.Id(1),
				GroupId:         ct.Id(10),
				Audience:        ct.Audience("selected"),
				AudienceIds:     []ct.Id{2, 3, 4},
				RequesterGroups: []ct.Id{10},
			},
			setupMock: func(m *MockQuerier, tx *MockTxRunner) {
				tx.On("RunTx", mock.Anything, mock.Anything).Return(nil)

				m.On("CreatePost", mock.Anything, mock.Anything).Return(int64(100), nil)
				m.On("InsertPostAudience", mock.Anything, mock.MatchedBy(func(params sqlc.InsertPostAudienceParams) bool {
					return params.PostID == 100 && len(params.AllowedUserIds) == 3
				})).Return(int64(3), nil)
			},
			expectedError: nil,
		},
		{
			name: "successful post creation with image",
			req: CreatePostReq{
				Body:            ct.PostBody("Post with image"),
				CreatorId:       ct.Id(1),
				GroupId:         ct.Id(10),
				Audience:        ct.Audience("everyone"),
				Image:           ct.Id(500),
				RequesterGroups: []ct.Id{10},
			},
			setupMock: func(m *MockQuerier, tx *MockTxRunner) {
				tx.On("RunTx", mock.Anything, mock.Anything).Return(nil)

				m.On("CreatePost", mock.Anything, mock.Anything).Return(int64(100), nil)
				m.On("UpsertImage", mock.Anything, mock.MatchedBy(func(params sqlc.UpsertImageParams) bool {
					return params.ID == 500 && params.ParentID == 100
				})).Return(nil)
			},
			expectedError: nil,
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
			setupMock:     func(m *MockQuerier, tx *MockTxRunner) {},
			expectedError: ErrNotAllowed,
		},
		{
			name: "selected audience with no audience IDs",
			req: CreatePostReq{
				Body:            ct.PostBody("Test post"),
				CreatorId:       ct.Id(1),
				GroupId:         ct.Id(10),
				Audience:        ct.Audience("selected"),
				AudienceIds:     []ct.Id{},
				RequesterGroups: []ct.Id{10},
			},
			setupMock: func(m *MockQuerier, tx *MockTxRunner) {
				tx.On("RunTx", mock.Anything, mock.Anything).Return(nil)

				m.On("CreatePost", mock.Anything, mock.Anything).Return(int64(100), nil)
			},
			expectedError: ErrNoAudienceSelected,
		},
		{
			name: "database error during CreatePost",
			req: CreatePostReq{
				Body:            ct.PostBody("Test post"),
				CreatorId:       ct.Id(1),
				GroupId:         ct.Id(10),
				Audience:        ct.Audience("everyone"),
				RequesterGroups: []ct.Id{10},
			},
			setupMock: func(m *MockQuerier, tx *MockTxRunner) {
				tx.On("RunTx", mock.Anything, mock.Anything).Return(nil)

				m.On("CreatePost", mock.Anything, mock.Anything).
					Return(int64(0), errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockDB, mockTx := setupTestAppWithTx()
			tt.setupMock(mockDB, mockTx)

			err := app.CreatePost(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if !errors.Is(tt.expectedError, ErrNotAllowed) && !errors.Is(tt.expectedError, ErrNoAudienceSelected) {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
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
		name              string
		req               GetGroupPostsReq
		setupMock         func(*MockQuerier)
		expectedPosts     int
		expectedError     error
		expectedErrSubstr string
	}{
		{
			name: "successful retrieval",
			req: GetGroupPostsReq{
				GroupId:         ct.Id(10),
				RequesterId:     ct.Id(1),
				RequesterGroups: []ct.Id{10, 20},
				Limit:           ct.Limit(10), // Added - was missing
				Offset:          ct.Offset(0), // Added - was missing
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
				Limit:           ct.Limit(10), // Added
				Offset:          ct.Offset(0), // Added
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
				Limit:           ct.Limit(10), // Added
				Offset:          ct.Offset(0), // Added
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
				Limit:           ct.Limit(10), // Added
				Offset:          ct.Offset(0), // Added
			},
			setupMock:         func(m *MockQuerier) {},
			expectedPosts:     0,
			expectedError:     nil, // expected via substring, not ErrorIs
			expectedErrSubstr: "required field missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockDB := setupTestApp()
			tt.setupMock(mockDB)

			posts, err := app.GetGroupPostsPaginated(context.Background(), tt.req)

			if tt.expectedErrSubstr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErrSubstr)
			} else if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
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
		name              string
		req               SimpleIdReq
		setupMock         func(*MockQuerier)
		expectPost        bool
		expectedError     error
		expectedErrSubstr string
	}{
		{
			name: "successful retrieval",
			req: SimpleIdReq{
				Id: ct.Id(10),
			},
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
			name: "no posts in group",
			req: SimpleIdReq{
				Id: ct.Id(10),
			},
			setupMock: func(m *MockQuerier) {
				m.On("GetMostPopularPostInGroup", mock.Anything, mock.Anything).
					Return(sqlc.GetMostPopularPostInGroupRow{}, sql.ErrNoRows)
			},
			expectPost:    false,
			expectedError: ErrNotFound,
		},
		{
			name: "no group ID provided",
			req: SimpleIdReq{
				Id: ct.Id(0),
			},
			setupMock:         func(m *MockQuerier) {},
			expectPost:        false,
			expectedError:     nil, // expected via substring, not ErrorIs
			expectedErrSubstr: "required field missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockDB := setupTestApp()
			tt.setupMock(mockDB)

			post, err := app.GetMostPopularPostInGroup(context.Background(), tt.req)

			if tt.expectedErrSubstr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErrSubstr)
			} else if tt.expectedError != nil {
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

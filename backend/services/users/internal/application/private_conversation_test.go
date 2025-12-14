package application

import (
	"context"
	"testing"

	"social-network/services/users/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClients is a testify mock implementing ClientsInterface for tests.
type MockClients struct {
	mock.Mock
}

func (m *MockClients) CreateGroupConversation(ctx context.Context, groupId int64, ownerId int64) error {
	args := m.Called(ctx, groupId, ownerId)
	return args.Error(0)
}

func (m *MockClients) CreatePrivateConversation(ctx context.Context, userId1, userId2 int64) error {
	args := m.Called(ctx, userId1, userId2)
	return args.Error(0)
}

func (m *MockClients) AddMembersToGroupConversation(ctx context.Context, groupId int64, userIds []int64) error {
	args := m.Called(ctx, groupId, userIds)
	return args.Error(0)
}

func (m *MockClients) DeleteConversationByExactMembers(ctx context.Context, userIds []int64) error {
	args := m.Called(ctx, userIds)
	return args.Error(0)
}

// func TestCreatePrivateConversation_OneWayFollows_CallsClient(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	mockClients := new(MockClients)
// 	service := NewApplicationWithMocks(mockDB, mockClients)

// 	ctx := context.Background()
// 	req := models.FollowUserReq{
// 		FollowerId:   ct.Id(1),
// 		TargetUserId: ct.Id(2),
// 	}

// 	// One-way follow: user1 follows user2 but not vice-versa
// 	mockDB.On("AreFollowingEachOther", ctx, sqlc.AreFollowingEachOtherParams{
// 		FollowerID:  1,
// 		FollowingID: 2,
// 	}).Return(sqlc.AreFollowingEachOtherRow{User1FollowsUser2: true, User2FollowsUser1: false}, nil)

// 	mockClients.On("CreatePrivateConversation", ctx, int64(1), int64(2)).Return(nil)

// 	err := service.createPrivateConversation(ctx, req)

// 	assert.NoError(t, err)
// 	mockClients.AssertExpectations(t)
// 	mockDB.AssertExpectations(t)
// }

// func TestCreatePrivateConversation_NeitherFollows_DoesNotCallClient(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	mockClients := new(MockClients)
// 	service := NewApplicationWithMocks(mockDB, mockClients)

// 	ctx := context.Background()
// 	req := models.FollowUserReq{
// 		FollowerId:   ct.Id(1),
// 		TargetUserId: ct.Id(2),
// 	}

// 	// Neither follows the other => AreFollowingEachOther returns empty/zero row
// 	mockDB.On("AreFollowingEachOther", ctx, sqlc.AreFollowingEachOtherParams{
// 		FollowerID:  1,
// 		FollowingID: 2,
// 	}).Return(sqlc.AreFollowingEachOtherRow{}, nil)

// 	err := service.createPrivateConversation(ctx, req)

// 	assert.NoError(t, err)
// 	// Ensure client was not called
// 	mockClients.AssertNotCalled(t, "CreatePrivateConversation", mock.Anything, mock.Anything, mock.Anything)
// 	mockDB.AssertExpectations(t)
// }

// func TestCreatePrivateConversation_BothFollow_DoesNotCallClient(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	mockClients := new(MockClients)
// 	service := NewApplicationWithMocks(mockDB, mockClients)

// 	ctx := context.Background()
// 	req := models.FollowUserReq{
// 		FollowerId:   ct.Id(1),
// 		TargetUserId: ct.Id(2),
// 	}

// 	// Both follow each other => no private conversation creation
// 	mockDB.On("AreFollowingEachOther", ctx, sqlc.AreFollowingEachOtherParams{
// 		FollowerID:  1,
// 		FollowingID: 2,
// 	}).Return(sqlc.AreFollowingEachOtherRow{User1FollowsUser2: true, User2FollowsUser1: true}, nil)

// 	err := service.createPrivateConversation(ctx, req)

// 	assert.NoError(t, err)
// 	mockClients.AssertNotCalled(t, "CreatePrivateConversation", mock.Anything, mock.Anything, mock.Anything)
// 	mockDB.AssertExpectations(t)
// }

// func TestCreatePrivateConversation_DbError_Propagates(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	mockClients := new(MockClients)
// 	service := NewApplicationWithMocks(mockDB, mockClients)

// 	ctx := context.Background()
// 	req := models.FollowUserReq{
// 		FollowerId:   ct.Id(1),
// 		TargetUserId: ct.Id(2),
// 	}

// 	expectedErr := errors.New("db failure")
// 	mockDB.On("AreFollowingEachOther", ctx, sqlc.AreFollowingEachOtherParams{
// 		FollowerID:  1,
// 		FollowingID: 2,
// 	}).Return(sqlc.AreFollowingEachOtherRow{}, expectedErr)

// 	err := service.createPrivateConversation(ctx, req)

// 	assert.Equal(t, expectedErr, err)
// 	mockClients.AssertNotCalled(t, "CreatePrivateConversation", mock.Anything, mock.Anything, mock.Anything)
// 	mockDB.AssertExpectations(t)
// }

func TestAreFollowingEachOther_VariousCases(t *testing.T) {
	mockDB := new(MockQuerier)
	mockClients := new(MockClients)
	service := NewApplicationWithMocks(mockDB, mockClients)

	ctx := context.Background()
	req := models.FollowUserReq{
		FollowerId:   ct.Id(1),
		TargetUserId: ct.Id(2),
	}

	// Case: neither follows => expect nil pointer
	mockDB.On("AreFollowingEachOther", ctx, sqlc.AreFollowingEachOtherParams{FollowerID: 1, FollowingID: 2}).Return(sqlc.AreFollowingEachOtherRow{}, nil)
	res, err := service.AreFollowingEachOther(ctx, req)
	assert.NoError(t, err)
	assert.Nil(t, res)
	mockDB.AssertExpectations(t)

	// Case: one-way follows => expect pointer to false
	mockDB = new(MockQuerier)
	service = NewApplicationWithMocks(mockDB, mockClients)
	mockDB.On("AreFollowingEachOther", ctx, sqlc.AreFollowingEachOtherParams{FollowerID: 1, FollowingID: 2}).Return(sqlc.AreFollowingEachOtherRow{User1FollowsUser2: true, User2FollowsUser1: false}, nil)
	res2, err := service.AreFollowingEachOther(ctx, req)
	assert.NoError(t, err)
	if assert.NotNil(t, res2) {
		assert.False(t, *res2)
	}
	mockDB.AssertExpectations(t)

	// Case: both follow => expect pointer to true
	mockDB = new(MockQuerier)
	service = NewApplicationWithMocks(mockDB, mockClients)
	mockDB.On("AreFollowingEachOther", ctx, sqlc.AreFollowingEachOtherParams{FollowerID: 1, FollowingID: 2}).Return(sqlc.AreFollowingEachOtherRow{User1FollowsUser2: true, User2FollowsUser1: true}, nil)
	res3, err := service.AreFollowingEachOther(ctx, req)
	assert.NoError(t, err)
	if assert.NotNil(t, res3) {
		assert.True(t, *res3)
	}
	mockDB.AssertExpectations(t)
}

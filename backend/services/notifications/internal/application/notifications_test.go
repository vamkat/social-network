package application

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"social-network/services/notifications/internal/db/sqlc"
)

// MockDB is a mock implementation of the database queries
type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateNotification(ctx context.Context, arg sqlc.CreateNotificationParams) (sqlc.Notification, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.Notification), args.Error(1)
}

func (m *MockDB) GetNotificationByID(ctx context.Context, id int64) (sqlc.Notification, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlc.Notification), args.Error(1)
}

func (m *MockDB) GetUserNotifications(ctx context.Context, arg sqlc.GetUserNotificationsParams) ([]sqlc.Notification, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.Notification), args.Error(1)
}

func (m *MockDB) GetUserNotificationsCount(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDB) GetUserUnreadNotificationsCount(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDB) MarkNotificationAsRead(ctx context.Context, arg sqlc.MarkNotificationAsReadParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockDB) MarkAllAsRead(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockDB) DeleteNotification(ctx context.Context, arg sqlc.DeleteNotificationParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockDB) CreateNotificationType(ctx context.Context, arg sqlc.CreateNotificationTypeParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockDB) GetNotificationType(ctx context.Context, notifType string) (sqlc.NotificationType, error) {
	args := m.Called(ctx, notifType)
	return args.Get(0).(sqlc.NotificationType), args.Error(1)
}

// Test CreateNotification function
func TestCreateNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	userID := int64(1)
	notifType := FollowRequest
	title := "Follow Request"
	message := "User wants to follow you"
	sourceService := "users"
	sourceEntityID := int64(2)
	needsAction := true
	payload := map[string]string{"requester_id": "2", "requester_name": "testuser"}

	payloadBytes, _ := json.Marshal(payload)

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         userID,
		NotifType:      string(notifType),
		SourceService:  sourceService,
		SourceEntityID: pgtype.Int8{Int64: sourceEntityID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: needsAction, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	notification, err := app.CreateNotification(ctx, userID, notifType, title, message, sourceService, sourceEntityID, needsAction, payload)

	assert.NoError(t, err)
	assert.NotNil(t, notification)
	assert.Equal(t, userID, notification.UserID)
	assert.Equal(t, notifType, notification.Type)
	assert.Equal(t, title, notification.Title)
	assert.Equal(t, message, notification.Message)
	assert.Equal(t, sourceService, notification.SourceService)
	assert.Equal(t, sourceEntityID, notification.SourceEntityID)
	assert.Equal(t, needsAction, notification.NeedsAction)
	assert.Equal(t, payload, notification.Payload)

	mockDB.AssertExpectations(t)
}

// Test GetNotification function
func TestGetNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	notificationID := int64(1)
	userID := int64(10)
	payloadBytes, _ := json.Marshal(map[string]string{
		"requester_id": "2",
		"requester_name": "testuser",
	})

	expectedDBNotification := sqlc.Notification{
		ID:             notificationID,
		UserID:         userID,
		NotifType:      "follow_request",
		SourceService:  "users",
		SourceEntityID: pgtype.Int8{Int64: 2, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: true, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
	}

	mockDB.On("GetNotificationByID", ctx, notificationID).Return(expectedDBNotification, nil)

	notification, err := app.GetNotification(ctx, notificationID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, notification)
	assert.Equal(t, notificationID, notification.ID)
	assert.Equal(t, userID, notification.UserID)
	assert.Equal(t, FollowRequest, notification.Type)
	assert.Equal(t, "users", notification.SourceService)
	assert.Equal(t, int64(2), notification.SourceEntityID)
	assert.False(t, notification.Seen)
	assert.True(t, notification.NeedsAction)
	assert.False(t, notification.Acted)

	mockDB.AssertExpectations(t)
}

// Test GetUserNotifications function
func TestGetUserNotifications(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	userID := int64(1)
	limit := int32(10)
	offset := int32(0)

	payloadBytes, _ := json.Marshal(map[string]string{
		"requester_id": "2",
		"requester_name": "testuser",
	})

	expectedDBNotifications := []sqlc.Notification{
		{
			ID:             1,
			UserID:         userID,
			NotifType:      "follow_request",
			SourceService:  "users",
			SourceEntityID: pgtype.Int8{Int64: 2, Valid: true},
			Seen:           pgtype.Bool{Bool: false, Valid: true},
			NeedsAction:    pgtype.Bool{Bool: true, Valid: true},
			Acted:          pgtype.Bool{Bool: false, Valid: true},
			CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
			ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
			DeletedAt:      pgtype.Timestamptz{Valid: false},
			Payload:        payloadBytes,
		},
	}

	mockDB.On("GetUserNotifications", ctx, mock.AnythingOfType("sqlc.GetUserNotificationsParams")).Return(expectedDBNotifications, nil)

	notifications, err := app.GetUserNotifications(ctx, userID, limit, offset)

	assert.NoError(t, err)
	assert.Len(t, notifications, 1)
	assert.Equal(t, int64(1), notifications[0].ID)
	assert.Equal(t, userID, notifications[0].UserID)

	mockDB.AssertExpectations(t)
}

// Test GetUserNotificationsCount function
func TestGetUserNotificationsCount(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	userID := int64(1)
	expectedCount := int64(5)

	mockDB.On("GetUserNotificationsCount", ctx, userID).Return(expectedCount, nil)

	count, err := app.GetUserNotificationsCount(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)

	mockDB.AssertExpectations(t)
}

// Test GetUserUnreadNotificationsCount function
func TestGetUserUnreadNotificationsCount(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	userID := int64(1)
	expectedCount := int64(3)

	mockDB.On("GetUserUnreadNotificationsCount", ctx, userID).Return(expectedCount, nil)

	count, err := app.GetUserUnreadNotificationsCount(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)

	mockDB.AssertExpectations(t)
}

// Test MarkNotificationAsRead function
func TestMarkNotificationAsRead(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	notificationID := int64(1)
	userID := int64(10)

	mockDB.On("MarkNotificationAsRead", ctx, mock.AnythingOfType("sqlc.MarkNotificationAsReadParams")).Return(nil)

	err := app.MarkNotificationAsRead(ctx, notificationID, userID)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test MarkAllAsRead function
func TestMarkAllAsRead(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	userID := int64(10)

	mockDB.On("MarkAllAsRead", ctx, userID).Return(nil)

	err := app.MarkAllAsRead(ctx, userID)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test DeleteNotification function
func TestDeleteNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	notificationID := int64(1)
	userID := int64(10)

	mockDB.On("DeleteNotification", ctx, mock.AnythingOfType("sqlc.DeleteNotificationParams")).Return(nil)

	err := app.DeleteNotification(ctx, notificationID, userID)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test CreateNotificationType function
func TestCreateNotificationType(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()

	mockDB.On("CreateNotificationType", ctx, mock.AnythingOfType("sqlc.CreateNotificationTypeParams")).Return(nil)

	err := app.CreateDefaultNotificationTypes(ctx)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test the specific notification trigger functions
func TestCreateFollowRequestNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	targetUserID := int64(1)
	requesterUserID := int64(2)
	requesterUsername := "testuser"

	payloadBytes, _ := json.Marshal(map[string]string{
		"requester_id":   "2",
		"requester_name": "testuser",
	})

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         targetUserID,
		NotifType:      string(FollowRequest),
		SourceService:  "users",
		SourceEntityID: pgtype.Int8{Int64: requesterUserID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: true, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateFollowRequestNotification(ctx, targetUserID, requesterUserID, requesterUsername)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test CreateGroupInviteNotification function
func TestCreateGroupInviteNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	invitedUserID := int64(1)
	inviterUserID := int64(2)
	groupID := int64(100)
	groupName := "Test Group"
	inviterUsername := "testuser"

	payloadBytes, _ := json.Marshal(map[string]string{
		"inviter_id":   "2",
		"inviter_name": "testuser",
		"group_id":     "100",
		"group_name":   "Test Group",
		"action":       "accept_or_decline",
	})

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         invitedUserID,
		NotifType:      string(GroupInvite),
		SourceService:  "users",
		SourceEntityID: pgtype.Int8{Int64: groupID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: true, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateGroupInviteNotification(ctx, invitedUserID, inviterUserID, groupID, groupName, inviterUsername)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test CreateGroupJoinRequestNotification function
func TestCreateGroupJoinRequestNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	groupOwnerID := int64(1)
	requesterID := int64(2)
	groupID := int64(100)
	groupName := "Test Group"
	requesterUsername := "testuser"

	payloadBytes, _ := json.Marshal(map[string]string{
		"requester_id":  "2",
		"requester_name": "testuser",
		"group_id":      "100",
		"group_name":    "Test Group",
		"action":        "accept_or_decline",
	})

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         groupOwnerID,
		NotifType:      string(GroupJoinRequest),
		SourceService:  "users",
		SourceEntityID: pgtype.Int8{Int64: groupID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: true, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateGroupJoinRequestNotification(ctx, groupOwnerID, requesterID, groupID, groupName, requesterUsername)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test CreateNewEventNotification function
func TestCreateNewEventNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	userID := int64(1)
	groupID := int64(100)
	eventID := int64(200)
	groupName := "Test Group"
	eventTitle := "Test Event"

	payloadBytes, _ := json.Marshal(map[string]string{
		"group_id":      "100",
		"group_name":    "Test Group",
		"event_id":      "200",
		"event_title":   "Test Event",
		"action":        "view_event",
	})

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         userID,
		NotifType:      string(NewEvent),
		SourceService:  "posts",
		SourceEntityID: pgtype.Int8{Int64: eventID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateNewEventNotification(ctx, userID, groupID, eventID, groupName, eventTitle)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}
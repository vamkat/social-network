package application

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"social-network/services/notifications/internal/db/sqlc"
)

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
		Count:          pgtype.Int4{Int32: 1, Valid: true}, // Add the count field
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
		Count:          pgtype.Int4{Int32: 1, Valid: true}, // Add the count field
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
			Count:          pgtype.Int4{Int32: 1, Valid: true}, // Add the count field
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
		Count:          pgtype.Int4{Int32: 1, Valid: true}, // Add the count field
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
		Count:          pgtype.Int4{Int32: 1, Valid: true}, // Add the count field
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
		Count:          pgtype.Int4{Int32: 1, Valid: true}, // Add the count field
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
		Count:          pgtype.Int4{Int32: 1, Valid: true}, // Add the count field
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateNewEventNotification(ctx, userID, groupID, eventID, groupName, eventTitle)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test CreateNotificationWithAggregation function - no aggregation path
func TestCreateNotificationWithAggregationNoExisting(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	userID := int64(1)
	notifType := PostLike
	title := "Post Liked"
	message := "User liked your post"
	sourceService := "posts"
	sourceEntityID := int64(2)
	needsAction := false
	payload := map[string]string{"liker_id": "2", "liker_name": "testuser"}
	aggregate := true

	payloadBytes, _ := json.Marshal(payload)

	// Expect GetUnreadNotificationByTypeAndEntity to return error (no existing notification)
	mockDB.On("GetUnreadNotificationByTypeAndEntity", ctx, mock.AnythingOfType("sqlc.GetUnreadNotificationByTypeAndEntityParams")).Return(sqlc.Notification{}, fmt.Errorf("sql: no rows in result set"))

	// Expect CreateNotification to be called
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
		Count:          pgtype.Int4{Int32: 1, Valid: true}, // Add the count field
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	notification, err := app.CreateNotificationWithAggregation(ctx, userID, notifType, title, message, sourceService, sourceEntityID, needsAction, payload, aggregate)

	assert.NoError(t, err)
	assert.NotNil(t, notification)
	assert.Equal(t, int32(1), notification.Count)  // Should be 1 since no existing notification was found
	assert.Equal(t, userID, notification.UserID)
	assert.Equal(t, notifType, notification.Type)

	mockDB.AssertExpectations(t)
}

// Test CreateNotificationWithAggregation function - aggregation path
func TestCreateNotificationWithAggregationExisting(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	userID := int64(1)
	notifType := PostLike
	title := "Post Liked"
	message := "User liked your post"
	sourceService := "posts"
	sourceEntityID := int64(2)
	needsAction := false
	payload := map[string]string{"liker_id": "2", "liker_name": "testuser"}
	aggregate := true

	payloadBytes, _ := json.Marshal(payload)

	// First, expect GetUnreadNotificationByTypeAndEntity to return an existing notification with count=2
	existingNotification := sqlc.Notification{
		ID:             10,
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
		Count:          pgtype.Int4{Int32: 2, Valid: true}, // Existing count is 2
	}

	mockDB.On("GetUnreadNotificationByTypeAndEntity", ctx, mock.AnythingOfType("sqlc.GetUnreadNotificationByTypeAndEntityParams")).Return(existingNotification, nil)

	// Then expect UpdateNotificationCount to be called to increment the count to 3
	mockDB.On("UpdateNotificationCount", ctx, mock.AnythingOfType("sqlc.UpdateNotificationCountParams")).Return(nil)

	// Finally expect GetNotificationByID to fetch the updated notification
	updatedNotification := sqlc.Notification{
		ID:             10,
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
		Count:          pgtype.Int4{Int32: 3, Valid: true}, // Updated count is 3
	}

	mockDB.On("GetNotificationByID", ctx, int64(10)).Return(updatedNotification, nil)

	notification, err := app.CreateNotificationWithAggregation(ctx, userID, notifType, title, message, sourceService, sourceEntityID, needsAction, payload, aggregate)

	assert.NoError(t, err)
	assert.NotNil(t, notification)
	assert.Equal(t, int32(3), notification.Count)  // Should be 3 (existing count 2 + 1)
	assert.Equal(t, userID, notification.UserID)
	assert.Equal(t, notifType, notification.Type)

	mockDB.AssertExpectations(t)
}

// Test CreateNotificationWithAggregation with aggregation disabled
func TestCreateNotificationWithAggregationDisabled(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	userID := int64(1)
	notifType := PostLike
	title := "Post Liked"
	message := "User liked your post"
	sourceService := "posts"
	sourceEntityID := int64(2)
	needsAction := false
	payload := map[string]string{"liker_id": "2", "liker_name": "testuser"}
	aggregate := false // Aggregation disabled

	payloadBytes, _ := json.Marshal(payload)

	// When aggregation is disabled, CreateNotification should be called directly
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
		Count:          pgtype.Int4{Int32: 1, Valid: true}, // Add the count field
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	notification, err := app.CreateNotificationWithAggregation(ctx, userID, notifType, title, message, sourceService, sourceEntityID, needsAction, payload, aggregate)

	assert.NoError(t, err)
	assert.NotNil(t, notification)
	assert.Equal(t, int32(1), notification.Count)  // Should be 1 since aggregation is disabled
	assert.Equal(t, userID, notification.UserID)
	assert.Equal(t, notifType, notification.Type)

	mockDB.AssertExpectations(t)
}

// Test CreateFollowRequestAcceptedNotification function
func TestCreateFollowRequestAcceptedNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	requesterUserID := int64(1)
	targetUserID := int64(2)
	targetUsername := "targetuser"

	payloadBytes, _ := json.Marshal(map[string]string{
		"target_id":   "2",
		"target_name": "targetuser",
	})

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         requesterUserID,
		NotifType:      string(FollowRequestAccepted),
		SourceService:  "users",
		SourceEntityID: pgtype.Int8{Int64: targetUserID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
		Count:          pgtype.Int4{Int32: 1, Valid: true},
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateFollowRequestAcceptedNotification(ctx, requesterUserID, targetUserID, targetUsername)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test CreateFollowRequestRejectedNotification function
func TestCreateFollowRequestRejectedNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	requesterUserID := int64(1)
	targetUserID := int64(2)
	targetUsername := "targetuser"

	payloadBytes, _ := json.Marshal(map[string]string{
		"target_id":   "2",
		"target_name": "targetuser",
	})

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         requesterUserID,
		NotifType:      string(FollowRequestRejected),
		SourceService:  "users",
		SourceEntityID: pgtype.Int8{Int64: targetUserID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
		Count:          pgtype.Int4{Int32: 1, Valid: true},
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateFollowRequestRejectedNotification(ctx, requesterUserID, targetUserID, targetUsername)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test CreateGroupInviteAcceptedNotification function
func TestCreateGroupInviteAcceptedNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	inviterUserID := int64(1)
	invitedUserID := int64(2)
	groupID := int64(100)
	invitedUsername := "inviteduser"
	groupName := "Test Group"

	payloadBytes, _ := json.Marshal(map[string]string{
		"invited_id":    "2",
		"invited_name":  "inviteduser",
		"group_id":      "100",
		"group_name":    "Test Group",
	})

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         inviterUserID,
		NotifType:      string(GroupInviteAccepted),
		SourceService:  "users",
		SourceEntityID: pgtype.Int8{Int64: groupID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
		Count:          pgtype.Int4{Int32: 1, Valid: true},
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateGroupInviteAcceptedNotification(ctx, inviterUserID, invitedUserID, groupID, invitedUsername, groupName)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test CreateGroupInviteRejectedNotification function
func TestCreateGroupInviteRejectedNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	inviterUserID := int64(1)
	invitedUserID := int64(2)
	groupID := int64(100)
	invitedUsername := "inviteduser"
	groupName := "Test Group"

	payloadBytes, _ := json.Marshal(map[string]string{
		"invited_id":    "2",
		"invited_name":  "inviteduser",
		"group_id":      "100",
		"group_name":    "Test Group",
	})

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         inviterUserID,
		NotifType:      string(GroupInviteRejected),
		SourceService:  "users",
		SourceEntityID: pgtype.Int8{Int64: groupID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
		Count:          pgtype.Int4{Int32: 1, Valid: true},
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateGroupInviteRejectedNotification(ctx, inviterUserID, invitedUserID, groupID, invitedUsername, groupName)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test CreateGroupJoinRequestAcceptedNotification function
func TestCreateGroupJoinRequestAcceptedNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	requesterUserID := int64(1)
	groupOwnerID := int64(2)
	groupID := int64(100)
	groupName := "Test Group"

	payloadBytes, _ := json.Marshal(map[string]string{
		"group_owner_id": "2",
		"group_id":       "100",
		"group_name":     "Test Group",
	})

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         requesterUserID,
		NotifType:      string(GroupJoinRequestAccepted),
		SourceService:  "users",
		SourceEntityID: pgtype.Int8{Int64: groupID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
		Count:          pgtype.Int4{Int32: 1, Valid: true},
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateGroupJoinRequestAcceptedNotification(ctx, requesterUserID, groupOwnerID, groupID, groupName)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

// Test CreateGroupJoinRequestRejectedNotification function
func TestCreateGroupJoinRequestRejectedNotification(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	requesterUserID := int64(1)
	groupOwnerID := int64(2)
	groupID := int64(100)
	groupName := "Test Group"

	payloadBytes, _ := json.Marshal(map[string]string{
		"group_owner_id": "2",
		"group_id":       "100",
		"group_name":     "Test Group",
	})

	expectedNotification := sqlc.Notification{
		ID:             1,
		UserID:         requesterUserID,
		NotifType:      string(GroupJoinRequestRejected),
		SourceService:  "users",
		SourceEntityID: pgtype.Int8{Int64: groupID, Valid: true},
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		DeletedAt:      pgtype.Timestamptz{Valid: false},
		Payload:        payloadBytes,
		Count:          pgtype.Int4{Int32: 1, Valid: true},
	}

	mockDB.On("CreateNotification", ctx, mock.AnythingOfType("sqlc.CreateNotificationParams")).Return(expectedNotification, nil)

	err := app.CreateGroupJoinRequestRejectedNotification(ctx, requesterUserID, groupOwnerID, groupID, groupName)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}
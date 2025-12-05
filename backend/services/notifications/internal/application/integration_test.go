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

// Test that all notification trigger functions work
func TestNotificationTriggerFunctions(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()

	// Test CreateGroupInviteNotification
	t.Run("CreateGroupInviteNotification", func(t *testing.T) {
		invitedUserID := int64(1)
		inviterUserID := int64(2)
		groupID := int64(3)
		groupName := "Test Group"
		inviterUsername := "testuser"

		payloadBytes, _ := json.Marshal(map[string]string{
			"inviter_id":   "2",
			"inviter_name": "testuser",
			"group_id":     "3",
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
	})

	// Test CreateGroupJoinRequestNotification
	t.Run("CreateGroupJoinRequestNotification", func(t *testing.T) {
		groupOwnerID := int64(1)
		requesterID := int64(2)
		groupID := int64(3)
		groupName := "Test Group"
		requesterUsername := "testuser"

		payloadBytes, _ := json.Marshal(map[string]string{
			"requester_id":  "2",
			"requester_name": "testuser",
			"group_id":      "3",
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
	})

	// Test CreateNewEventNotification
	t.Run("CreateNewEventNotification", func(t *testing.T) {
		userID := int64(1)
		groupID := int64(2)
		eventID := int64(3)
		groupName := "Test Group"
		eventTitle := "Test Event"

		payloadBytes, _ := json.Marshal(map[string]string{
			"group_id":      "2",
			"group_name":    "Test Group",
			"event_id":      "3",
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
	})
}

// Test notification queries work
func TestNotificationQueries(t *testing.T) {
	mockDB := new(MockDB)
	app := NewApplication(mockDB)

	ctx := context.Background()
	userID := int64(1)

	// Test GetUserUnreadNotificationsCount
	t.Run("GetUserUnreadNotificationsCount", func(t *testing.T) {
		mockDB.On("GetUserUnreadNotificationsCount", ctx, userID).Return(int64(5), nil)

		count, err := app.GetUserUnreadNotificationsCount(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
		mockDB.AssertExpectations(t)
	})

	// Test MarkNotificationAsRead
	t.Run("MarkNotificationAsRead", func(t *testing.T) {
		notificationID := int64(123)
		mockDB.On("MarkNotificationAsRead", ctx, mock.AnythingOfType("sqlc.MarkNotificationAsReadParams")).Return(nil)

		err := app.MarkNotificationAsRead(ctx, notificationID, userID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}
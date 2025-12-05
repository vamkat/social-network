package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "social-network/services/notifications/internal/db/sqlc"
)

// CreateNotification creates a new notification
func (a *Application) CreateNotification(ctx context.Context, userID int64, notifType NotificationType, title, message, sourceService string, sourceEntityID int64, needsAction bool, payload map[string]string) (*Notification, error) {
	// Prepare the JSON payload
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Calculate expiration time (default 30 days if not specified)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	// Create the notification in the database
	dbNotification, err := a.DB.CreateNotification(ctx, db.CreateNotificationParams{
		UserID:         userID,
		NotifType:      string(notifType),
		SourceService:  sourceService,
		SourceEntityID: pgtype.Int8{Int64: sourceEntityID, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: needsAction, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true}, // New notifications haven't been acted upon yet
		Payload:        payloadJSON,
		ExpiresAt:      pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Convert database notification to our model
	notification := &Notification{
		ID:    dbNotification.ID,
		UserID: dbNotification.UserID,
		Type:  NotificationType(dbNotification.NotifType),
		SourceService: dbNotification.SourceService,
		Title:          title,
		Message:        message,
	}

	// Handle optional fields with proper type conversion
	if dbNotification.SourceEntityID.Valid {
		notification.SourceEntityID = dbNotification.SourceEntityID.Int64
	}
	notification.Seen = dbNotification.Seen.Bool
	notification.NeedsAction = dbNotification.NeedsAction.Bool
	notification.Acted = dbNotification.Acted.Bool

	if dbNotification.CreatedAt.Valid {
		notification.CreatedAt = dbNotification.CreatedAt.Time
	}
	if dbNotification.ExpiresAt.Valid {
		notification.ExpiresAt = &dbNotification.ExpiresAt.Time
	}
	if dbNotification.DeletedAt.Valid {
		notification.DeletedAt = &dbNotification.DeletedAt.Time
	}

	// Parse the payload JSON back to map
	if len(dbNotification.Payload) > 0 {
		err = json.Unmarshal(dbNotification.Payload, &notification.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}
	}

	return notification, nil
}

// CreateNotifications creates multiple notifications in a batch
func (a *Application) CreateNotifications(ctx context.Context, notifications []struct {
	UserID         int64
	Type           NotificationType
	Title          string
	Message        string
	SourceService  string
	SourceEntityID int64
	NeedsAction    bool
	Payload        map[string]string
}) ([]*Notification, error) {
	result := make([]*Notification, 0, len(notifications))

	for _, n := range notifications {
		notification, err := a.CreateNotification(ctx, n.UserID, n.Type, n.Title, n.Message, n.SourceService, n.SourceEntityID, n.NeedsAction, n.Payload)
		if err != nil {
			return nil, err
		}
		result = append(result, notification)
	}

	return result, nil
}

// GetNotification retrieves a single notification by ID
func (a *Application) GetNotification(ctx context.Context, notificationID, userID int64) (*Notification, error) {
	dbNotification, err := a.DB.GetNotificationByID(ctx, notificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	// Only return the notification if it belongs to the user
	if dbNotification.UserID != userID {
		return nil, fmt.Errorf("notification not found")
	}

	notification := &Notification{
		ID:             dbNotification.ID,
		UserID:         dbNotification.UserID,
		Type:           NotificationType(dbNotification.NotifType),
		SourceService:  dbNotification.SourceService,
	}

	// Handle optional fields with proper type conversion
	if dbNotification.SourceEntityID.Valid {
		notification.SourceEntityID = dbNotification.SourceEntityID.Int64
	}
	notification.Seen = dbNotification.Seen.Bool
	notification.NeedsAction = dbNotification.NeedsAction.Bool
	notification.Acted = dbNotification.Acted.Bool

	if dbNotification.CreatedAt.Valid {
		notification.CreatedAt = dbNotification.CreatedAt.Time
	}
	if dbNotification.ExpiresAt.Valid {
		notification.ExpiresAt = &dbNotification.ExpiresAt.Time
	}
	if dbNotification.DeletedAt.Valid {
		notification.DeletedAt = &dbNotification.DeletedAt.Time
	}

	// Parse the payload JSON back to map
	if len(dbNotification.Payload) > 0 {
		err = json.Unmarshal(dbNotification.Payload, &notification.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}
	}

	return notification, nil
}

// GetUserNotifications retrieves notifications for a user
func (a *Application) GetUserNotifications(ctx context.Context, userID int64, limit, offset int32) ([]*Notification, error) {
	dbNotifications, err := a.DB.GetUserNotifications(ctx, db.GetUserNotificationsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}

	notifications := make([]*Notification, len(dbNotifications))
	for i, dbNotif := range dbNotifications {
		notification := &Notification{
			ID:             dbNotif.ID,
			UserID:         dbNotif.UserID,
			Type:           NotificationType(dbNotif.NotifType),
			SourceService:  dbNotif.SourceService,
		}

		// Handle optional fields with proper type conversion
		if dbNotif.SourceEntityID.Valid {
			notification.SourceEntityID = dbNotif.SourceEntityID.Int64
		}
		notification.Seen = dbNotif.Seen.Bool
		notification.NeedsAction = dbNotif.NeedsAction.Bool
		notification.Acted = dbNotif.Acted.Bool

		if dbNotif.CreatedAt.Valid {
			notification.CreatedAt = dbNotif.CreatedAt.Time
		}
		if dbNotif.ExpiresAt.Valid {
			notification.ExpiresAt = &dbNotif.ExpiresAt.Time
		}
		if dbNotif.DeletedAt.Valid {
			notification.DeletedAt = &dbNotif.DeletedAt.Time
		}

		// Parse the payload JSON back to map
		if len(dbNotif.Payload) > 0 {
			err = json.Unmarshal(dbNotif.Payload, &notification.Payload)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
			}
		}

		notifications[i] = notification
	}

	return notifications, nil
}

// GetUserNotificationsCount gets the total count of notifications for a user
func (a *Application) GetUserNotificationsCount(ctx context.Context, userID int64) (int64, error) {
	count, err := a.DB.GetUserNotificationsCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get notifications count: %w", err)
	}
	return count, nil
}

// GetUserUnreadNotificationsCount gets the count of unread notifications for a user
func (a *Application) GetUserUnreadNotificationsCount(ctx context.Context, userID int64) (int64, error) {
	count, err := a.DB.GetUserUnreadNotificationsCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread notifications count: %w", err)
	}
	return count, nil
}

// MarkNotificationAsRead marks a notification as read
func (a *Application) MarkNotificationAsRead(ctx context.Context, notificationID, userID int64) error {
	err := a.DB.MarkNotificationAsRead(ctx, db.MarkNotificationAsReadParams{
		ID:     notificationID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

// MarkAllAsRead marks all notifications for a user as read
func (a *Application) MarkAllAsRead(ctx context.Context, userID int64) error {
	err := a.DB.MarkAllAsRead(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}

// DeleteNotification soft deletes a notification
func (a *Application) DeleteNotification(ctx context.Context, notificationID, userID int64) error {
	err := a.DB.DeleteNotification(ctx, db.DeleteNotificationParams{
		ID:     notificationID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// CreateDefaultNotificationTypes ensures default notification types are in the database
func (a *Application) CreateDefaultNotificationTypes(ctx context.Context) error {
	defaultTypes := []struct {
		Type          string
		Category      string
		DefaultEnabled bool
	}{
		{string(FollowRequest), "social", true},
		{string(NewFollower), "social", true},
		{string(GroupInvite), "group", true},
		{string(GroupJoinRequest), "group", true},
		{string(NewEvent), "group", true},
		{string(PostLike), "posts", true},
		{string(PostComment), "posts", true},
		{string(Mention), "posts", true},
	}

	for _, nt := range defaultTypes {
		err := a.DB.CreateNotificationType(ctx, db.CreateNotificationTypeParams{
			NotifType:      nt.Type,
			Category:       pgtype.Text{String: nt.Category, Valid: true},
			DefaultEnabled: pgtype.Bool{Bool: nt.DefaultEnabled, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create notification type %s: %w", nt.Type, err)
		}
	}

	return nil
}
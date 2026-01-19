package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	db "social-network/services/notifications/internal/db/sqlc"
	ct "social-network/shared/go/ct"
	tele "social-network/shared/go/telemetry"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateNotification creates a new notification
func (a *Application) CreateNotification(ctx context.Context, userID int64, notifType NotificationType, title, message, sourceService string, sourceEntityID int64, needsAction bool, payload map[string]string) (*Notification, error) {
	return a.CreateNotificationWithAggregation(ctx, userID, notifType, title, message, sourceService, sourceEntityID, needsAction, payload, false)
}

// CreateNotificationWithAggregation creates a new notification or aggregates with an existing one if applicable
func (a *Application) CreateNotificationWithAggregation(ctx context.Context, userID int64, notifType NotificationType, title, message, sourceService string, sourceEntityID int64, needsAction bool, payload map[string]string, aggregate bool) (*Notification, error) {
	if !aggregate {
		// If aggregation is disabled, create a new notification as before
		return a.createNotification(ctx, userID, notifType, title, message, sourceService, sourceEntityID, needsAction, payload, 1)
	}

	// If aggregation is enabled, first check for an existing unread notification of same type and entity
	existingNotification, err := a.DB.GetUnreadNotificationByTypeAndEntity(ctx, db.GetUnreadNotificationByTypeAndEntityParams{
		UserID:         userID,
		NotifType:      string(notifType),
		SourceEntityID: pgtype.Int8{Int64: sourceEntityID, Valid: true},
	})

	if err != nil {
		// If no existing notification found (which is normal), create a new one
		if errors.Is(err, pgx.ErrNoRows) {
			return a.createNotification(ctx, userID, notifType, title, message, sourceService, sourceEntityID, needsAction, payload, 1)
		}
		return nil, fmt.Errorf("failed to check for existing notification: %w", err)
	}

	// If an existing unread notification is found, increment its count and update it
	newCount := existingNotification.Count.Int32 + 1
	err = a.DB.UpdateNotificationCount(ctx, db.UpdateNotificationCountParams{
		Count:  pgtype.Int4{Int32: newCount, Valid: true},
		ID:     existingNotification.ID,
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update notification count: %w", err)
	}

	// Fetch and return the updated notification
	updatedNotification, err := a.DB.GetNotificationByID(ctx, existingNotification.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated notification: %w", err)
	}

	// Convert database notification to our model
	notification := &Notification{
		ID:            updatedNotification.ID,
		UserID:        updatedNotification.UserID,
		Type:          NotificationType(updatedNotification.NotifType),
		SourceService: updatedNotification.SourceService,
		Title:         a.formatAggregatedTitle(title, int64(newCount)),
		Message:       a.formatAggregatedMessage(message, int64(newCount)),
		Count:         newCount,
	}

	// Handle optional fields with proper type conversion
	if updatedNotification.SourceEntityID.Valid {
		notification.SourceEntityID = updatedNotification.SourceEntityID.Int64
	}
	notification.Seen = updatedNotification.Seen.Bool
	notification.NeedsAction = updatedNotification.NeedsAction.Bool
	notification.Acted = updatedNotification.Acted.Bool

	if updatedNotification.CreatedAt.Valid {
		notification.CreatedAt = updatedNotification.CreatedAt.Time
	}
	if updatedNotification.ExpiresAt.Valid {
		notification.ExpiresAt = &updatedNotification.ExpiresAt.Time
	}
	if updatedNotification.DeletedAt.Valid {
		notification.DeletedAt = &updatedNotification.DeletedAt.Time
	}

	// Parse the payload JSON back to map
	if len(updatedNotification.Payload) > 0 {
		err = json.Unmarshal(updatedNotification.Payload, &notification.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}
	}

	// Publish the notification to NATS for real-time delivery to the live service
	// We do this asynchronously to not block the notification creation
	go func() {
		// Create a background context for the NATS publish operation
		natsCtx := context.Background()
		if err := a.publishNotificationToNATS(natsCtx, notification); err != nil {
			// Log the error but don't fail the notification creation
			tele.Error(natsCtx, "failed to publish notification to nats in background: @1", "error", err.Error())
		}
	}()

	return notification, nil
}

// createNotification is a helper function that creates a notification with a specific count
func (a *Application) createNotification(ctx context.Context, userID int64, notifType NotificationType, title, message, sourceService string, sourceEntityID int64, needsAction bool, payload map[string]string, count int32) (*Notification, error) {
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
		Count:          pgtype.Int4{Int32: count, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Convert database notification to our model
	notification := &Notification{
		ID:            dbNotification.ID,
		UserID:        dbNotification.UserID,
		Type:          NotificationType(dbNotification.NotifType),
		SourceService: dbNotification.SourceService,
		Title:         title,
		Message:       message,
	}

	// Handle optional fields with proper type conversion
	if dbNotification.SourceEntityID.Valid {
		notification.SourceEntityID = dbNotification.SourceEntityID.Int64
	}
	notification.Seen = dbNotification.Seen.Bool
	notification.NeedsAction = dbNotification.NeedsAction.Bool
	notification.Acted = dbNotification.Acted.Bool
	notification.Count = dbNotification.Count.Int32 // Set the count from the database response

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

	// Publish the notification to NATS for real-time delivery to the live service
	// We do this asynchronously to not block the notification creation
	go func() {
		// Create a background context for the NATS publish operation
		natsCtx := context.Background()
		if err := a.publishNotificationToNATS(natsCtx, notification); err != nil {
			// Log the error but don't fail the notification creation
			tele.Error(natsCtx, "failed to publish notification to nats in background: @1", "error", err.Error())
		}
	}()

	return notification, nil
}

// formatAggregatedTitle formats notification titles when notifications are aggregated
func (a *Application) formatAggregatedTitle(originalTitle string, count int64) string {
	if count <= 1 {
		return originalTitle
	}

	// For now, we'll handle a few common cases, but this could be extended based on notification type
	switch originalTitle {
	case "Post Liked":
		return fmt.Sprintf("%d People Liked Your Post", count)
	case "New Comment":
		return fmt.Sprintf("%d People Commented On Your Post", count)
	case "New Follower":
		return fmt.Sprintf("%d New Followers", count)
	case "New Message":
		return fmt.Sprintf("%d New Messages", count)
	default:
		return fmt.Sprintf("%d Notifications", count)
	}
}

// formatAggregatedMessage formats notification messages when notifications are aggregated
func (a *Application) formatAggregatedMessage(originalMessage string, count int64) string {
	if count <= 1 {
		return originalMessage
	}

	// For now, we'll handle a few common cases, but this could be extended based on notification type
	switch {
	case strings.Contains(originalMessage, "liked your post"):
		return fmt.Sprintf("%d people liked your post", count)
	case strings.Contains(originalMessage, "commented on your post"):
		return fmt.Sprintf("%d people commented on your post", count)
	case strings.Contains(originalMessage, "is now following you"):
		return fmt.Sprintf("%d people are now following you", count)
	case strings.Contains(originalMessage, "sent you a message"):
		return fmt.Sprintf("%d people sent you a message", count)
	default:
		return fmt.Sprintf("You have %d notifications", count)
	}
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
		ID:            dbNotification.ID,
		UserID:        dbNotification.UserID,
		Type:          NotificationType(dbNotification.NotifType),
		SourceService: dbNotification.SourceService,
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
			ID:            dbNotif.ID,
			UserID:        dbNotif.UserID,
			Type:          NotificationType(dbNotif.NotifType),
			SourceService: dbNotif.SourceService,
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
		Type           string
		Category       string
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
		{string(NewMessage), "chat", true},
		{string(FollowRequestAccepted), "social", true},
		{string(FollowRequestRejected), "social", true},
		{string(GroupInviteAccepted), "group", true},
		{string(GroupInviteRejected), "group", true},
		{string(GroupJoinRequestAccepted), "group", true},
		{string(GroupJoinRequestRejected), "group", true},
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

// publishNotificationToNATS publishes a notification to NATS for real-time delivery to the live service
func (a *Application) publishNotificationToNATS(ctx context.Context, notification *Notification) error {
	if a.NatsConn == nil {
		tele.Warn(ctx, "NATS connection is nil, skipping notification publish")
		return nil
	}

	// Marshal the notification to JSON format (similar to chat service)
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification to JSON: %w", err)
	}

	// Publish to the user-specific NATS subject
	subject := ct.NotificationKey(notification.UserID)
	err = a.NatsConn.Publish(subject, notificationJSON)
	if err != nil {
		// Log the error but don't fail the notification creation
		tele.Error(ctx, "failed to publish notification to nats: @1", "error", err.Error())
		return fmt.Errorf("failed to publish notification to nats: %w", err)
	}

	// Flush to ensure the message is sent
	err = a.NatsConn.Flush()
	if err != nil {
		tele.Error(ctx, "failed to flush nats connection: @1", "error", err.Error())
	}

	tele.Info(ctx, "Published notification to nats for user @1", "userId", notification.UserID)
	return nil
}

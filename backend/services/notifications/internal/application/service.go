package application

import (
	"context"

	"social-network/services/notifications/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

// Service implements the NotificationService interface
type Service struct {
	db DBQuerier
}

// NewService creates a new notification service
func NewService(db DBQuerier) *Service {
	return &Service{
		db: db,
	}
}

// GetNotifications returns notifications for a user based on the specified status
func (s *Service) GetNotifications(ctx context.Context, req *GetNotificationsRequest) (*GetNotificationsResponse, error) {
	var notifications []interface{}
	var totalCount int64

	if req.Status == nil {
		// Get all notifications (read and unread)
		params := sqlc.GetNotificationsByUserIdAllParams{
			Column1: pgtype.Int8{Int64: req.UserID, Valid: true},
			Column2: pgtype.Int8{Int64: int64(req.Limit), Valid: true},
			Column3: pgtype.Int8{Int64: int64(req.Offset), Valid: true},
		}
		notificationsRaw, err := s.db.GetNotificationsByUserIdAll(ctx, params)
		if err != nil {
			return nil, err
		}

		// Convert to interface{}
		notifications = make([]interface{}, len(notificationsRaw))
		for i, n := range notificationsRaw {
			notifications[i] = n
		}

		count, err := s.db.GetNotificationCountAll(ctx, pgtype.Int8{Int64: req.UserID, Valid: true})
		if err != nil {
			return nil, err
		}
		totalCount = count.Int64
	} else {
		// Get notifications with a specific status (read or unread)
		status := pgtype.Bool{Bool: *req.Status, Valid: true}
		params := sqlc.GetNotificationsByUserIdFilteredByStatusParams{
			Column1: pgtype.Int8{Int64: req.UserID, Valid: true},
			Column2: status,
			Column3: pgtype.Int8{Int64: int64(req.Limit), Valid: true},
			Column4: pgtype.Int8{Int64: int64(req.Offset), Valid: true},
		}
		notificationsRaw, err := s.db.GetNotificationsByUserIdFilteredByStatus(ctx, params)
		if err != nil {
			return nil, err
		}

		// Convert to interface{}
		notifications = make([]interface{}, len(notificationsRaw))
		for i, n := range notificationsRaw {
			notifications[i] = n
		}

		countParams := sqlc.GetNotificationCountFilteredByStatusParams{
			Column1: pgtype.Int8{Int64: req.UserID, Valid: true},
			Column2: status,
		}
		count, err := s.db.GetNotificationCountFilteredByStatus(ctx, countParams)
		if err != nil {
			return nil, err
		}
		totalCount = count.Int64
	}

	result := make([]Notification, len(notifications))
	for i, n := range notifications {
		result[i] = ConvertDBNotification(n)
	}

	return &GetNotificationsResponse{
		Notifications: result,
		TotalCount:    totalCount,
	}, nil
}

// MarkAsRead marks the given notifications as read
func (s *Service) MarkAsRead(ctx context.Context, req *MarkAsReadRequest) (*MarkAsReadResponse, error) {
	params := sqlc.MarkNotificationsAsReadParams{
		Column1: pgtype.Int8{Int64: req.UserID, Valid: true},
		Column2: req.NotificationIDs,
	}

	err := s.db.MarkNotificationsAsRead(ctx, params)
	if err != nil {
		return nil, err
	}

	// Return the number of notifications actually updated
	// For simplicity, we return the count of IDs passed in
	return &MarkAsReadResponse{
		Success:      true,
		UpdatedCount: int32(len(req.NotificationIDs)),
	}, nil
}

// CreateNotification creates a new notification
func (s *Service) CreateNotification(ctx context.Context, req *CreateNotificationRequest) (*CreateNotificationResponse, error) {
	params, err := ConvertCreateNotificationParams(req)
	if err != nil {
		return nil, err
	}

	dbNotification, err := s.db.CreateNotification(ctx, params)
	if err != nil {
		return nil, err
	}

	return &CreateNotificationResponse{
		ID:      dbNotification.ID,
		Success: true,
	}, nil
}

// MarkNotificationAsActed marks a notification as having been acted upon
func (s *Service) MarkNotificationAsActed(ctx context.Context, notificationID int64) error {
	return s.db.MarkNotificationAsActed(ctx, pgtype.Int8{Int64: notificationID, Valid: true})
}
package application

import (
	"context"
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NewFollower      NotificationType = "new_follower"
	FollowRequest    NotificationType = "follow_request"
	GroupInvite      NotificationType = "group_invite"
	GroupJoinRequest NotificationType = "group_join_request"
	NewEvent         NotificationType = "new_event"
	NewMessage       NotificationType = "new_message"
	PostReply        NotificationType = "post_reply"
	Like             NotificationType = "like"
)

// SourceService represents the service that generated the notification
type SourceService string

const (
	UsersService SourceService = "users"
	ChatService  SourceService = "chat"
	PostsService SourceService = "posts"
)

// Notification represents a notification entity
type Notification struct {
	ID             int64             `json:"id"`
	UserID         int64             `json:"user_id"`
	Type           NotificationType  `json:"type"`
	SourceService  SourceService     `json:"source_service"`
	SourceEntityID int64             `json:"source_entity_id"`
	Seen           bool              `json:"seen"`
	NeedsAction    bool              `json:"needs_action"`
	Acted          bool              `json:"acted"`
	Payload        map[string]string `json:"payload"`
	CreatedAt      time.Time         `json:"created_at"`
	ExpiresAt      time.Time         `json:"expires_at"`
}

// GetNotificationsRequest represents the request to get notifications
type GetNotificationsRequest struct {
	UserID int64  `json:"user_id"`
	Status *bool  `json:"status,omitempty"` // nil = all, true = read, false = unread
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

// GetNotificationsResponse represents the response for getting notifications
type GetNotificationsResponse struct {
	Notifications []Notification `json:"notifications"`
	TotalCount    int64          `json:"total_count"`
}

// MarkAsReadRequest represents the request to mark notifications as read
type MarkAsReadRequest struct {
	UserID          int64   `json:"user_id"`
	NotificationIDs []int64 `json:"notification_ids"`
}

// MarkAsReadResponse represents the response for marking notifications as read
type MarkAsReadResponse struct {
	Success      bool `json:"success"`
	UpdatedCount int32 `json:"updated_count"`
}

// CreateNotificationRequest represents the request to create a notification
type CreateNotificationRequest struct {
	UserID         int64             `json:"user_id"`
	Type           NotificationType  `json:"type"`
	SourceService  SourceService     `json:"source_service"`
	SourceEntityID int64             `json:"source_entity_id"`
	Payload        map[string]string `json:"payload"`
	NeedsAction    bool              `json:"needs_action"`
}

// CreateNotificationResponse represents the response for creating a notification
type CreateNotificationResponse struct {
	ID      int64 `json:"id"`
	Success bool  `json:"success"`
}

// NotificationService interface defines the notifications service operations
type NotificationService interface {
	GetNotifications(ctx context.Context, req *GetNotificationsRequest) (*GetNotificationsResponse, error)
	MarkAsRead(ctx context.Context, req *MarkAsReadRequest) (*MarkAsReadResponse, error)
	CreateNotification(ctx context.Context, req *CreateNotificationRequest) (*CreateNotificationResponse, error)
	MarkNotificationAsActed(ctx context.Context, notificationID int64) error
}
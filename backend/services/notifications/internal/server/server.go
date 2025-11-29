package server

import (
	"context"

	"social-network/services/notifications/internal/application"

	notificationpb "social-network/shared/gen/notifications"
)

// NotificationServer implements the gRPC server for notifications
type NotificationServer struct {
	notificationpb.UnimplementedNotificationServiceServer
	service application.NotificationService
}

// NewNotificationServer creates a new notification server
func NewNotificationServer(service application.NotificationService) *NotificationServer {
	return &NotificationServer{
		service: service,
	}
}

// GetNotifications handles the get notifications request
func (s *NotificationServer) GetNotifications(ctx context.Context, req *notificationpb.GetNotificationsRequest) (*notificationpb.GetNotificationsResponse, error) {
	// Convert gRPC request to application request
	appReq := &application.GetNotificationsRequest{
		UserID: req.UserId,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	// Set status filter if provided
	if req.Status != nil {
		status := req.GetStatus()
		appReq.Status = &status
	}

	appResp, err := s.service.GetNotifications(ctx, appReq)
	if err != nil {
		return nil, err
	}

	// Convert application response to gRPC response
	grpcResp := &notificationpb.GetNotificationsResponse{
		Notifications: make([]*notificationpb.Notification, len(appResp.Notifications)),
		TotalCount:    appResp.TotalCount,
	}

	for i, notif := range appResp.Notifications {
		var notificationType notificationpb.NotificationType
		switch application.NotificationType(notif.Type) {
		case application.NewFollower:
			notificationType = notificationpb.NotificationType_NEW_FOLLOWER
		case application.FollowRequest:
			notificationType = notificationpb.NotificationType_FOLLOW_REQUEST
		case application.GroupInvite:
			notificationType = notificationpb.NotificationType_GROUP_INVITE
		case application.GroupJoinRequest:
			notificationType = notificationpb.NotificationType_GROUP_JOIN_REQUEST
		case application.NewEvent:
			notificationType = notificationpb.NotificationType_NEW_EVENT
		case application.NewMessage:
			notificationType = notificationpb.NotificationType_NEW_MESSAGE
		case application.PostReply:
			notificationType = notificationpb.NotificationType_POST_REPLY
		case application.Like:
			notificationType = notificationpb.NotificationType_LIKE
		default:
			notificationType = notificationpb.NotificationType_NEW_FOLLOWER // default
		}

		var sourceService notificationpb.SourceService
		switch application.SourceService(notif.SourceService) {
		case application.UsersService:
			sourceService = notificationpb.SourceService_USERS
		case application.ChatService:
			sourceService = notificationpb.SourceService_CHAT
		case application.PostsService:
			sourceService = notificationpb.SourceService_POSTS
		default:
			sourceService = notificationpb.SourceService_USERS // default
		}

		grpcResp.Notifications[i] = &notificationpb.Notification{
			Id:             notif.ID,
			UserId:         notif.UserID,
			Type:           notificationType,
			SourceService:  sourceService,
			SourceEntityId: notif.SourceEntityID,
			Seen:           notif.Seen,
			NeedsAction:    notif.NeedsAction,
			Acted:          notif.Acted,
			Payload:        notif.Payload,
			CreatedAt:      notif.CreatedAt.Unix(),  // Convert to Unix timestamp
			ExpiresAt:      notif.ExpiresAt.Unix(),  // Convert to Unix timestamp
		}
	}

	return grpcResp, nil
}

// MarkAsRead handles the mark as read request
func (s *NotificationServer) MarkAsRead(ctx context.Context, req *notificationpb.MarkAsReadRequest) (*notificationpb.MarkAsReadResponse, error) {
	appReq := &application.MarkAsReadRequest{
		UserID:          req.UserId,
		NotificationIDs: req.NotificationIds,
	}

	appResp, err := s.service.MarkAsRead(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &notificationpb.MarkAsReadResponse{
		Success:      appResp.Success,
		UpdatedCount: appResp.UpdatedCount,
	}, nil
}

// CreateNotification handles the create notification request
func (s *NotificationServer) CreateNotification(ctx context.Context, req *notificationpb.CreateNotificationRequest) (*notificationpb.CreateNotificationResponse, error) {
	var notificationType application.NotificationType
	switch req.Type {
	case notificationpb.NotificationType_NEW_FOLLOWER:
		notificationType = application.NewFollower
	case notificationpb.NotificationType_FOLLOW_REQUEST:
		notificationType = application.FollowRequest
	case notificationpb.NotificationType_GROUP_INVITE:
		notificationType = application.GroupInvite
	case notificationpb.NotificationType_GROUP_JOIN_REQUEST:
		notificationType = application.GroupJoinRequest
	case notificationpb.NotificationType_NEW_EVENT:
		notificationType = application.NewEvent
	case notificationpb.NotificationType_NEW_MESSAGE:
		notificationType = application.NewMessage
	case notificationpb.NotificationType_POST_REPLY:
		notificationType = application.PostReply
	case notificationpb.NotificationType_LIKE:
		notificationType = application.Like
	default:
		notificationType = application.NewFollower // default
	}

	var sourceService application.SourceService
	switch req.SourceService {
	case notificationpb.SourceService_USERS:
		sourceService = application.UsersService
	case notificationpb.SourceService_CHAT:
		sourceService = application.ChatService
	case notificationpb.SourceService_POSTS:
		sourceService = application.PostsService
	default:
		sourceService = application.UsersService // default
	}

	appReq := &application.CreateNotificationRequest{
		UserID:         req.UserId,
		Type:           notificationType,
		SourceService:  sourceService,
		SourceEntityID: req.SourceEntityId,
		Payload:        req.Payload,
		NeedsAction:    req.NeedsAction,
	}

	appResp, err := s.service.CreateNotification(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &notificationpb.CreateNotificationResponse{
		Id:      appResp.ID,
		Success: appResp.Success,
	}, nil
}

// MarkNotificationAsActed handles marking a notification as acted upon
func (s *NotificationServer) MarkNotificationAsActed(ctx context.Context, req *notificationpb.MarkNotificationAsActedRequest) (*notificationpb.MarkAsActedResponse, error) {
	err := s.service.MarkNotificationAsActed(ctx, req.NotificationId)
	if err != nil {
		return nil, err
	}

	return &notificationpb.MarkAsActedResponse{
		Success: true,
	}, nil
}
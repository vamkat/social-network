package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"social-network/services/notifications/internal/application"
	pb "social-network/shared/gen-go/notifications"
)

// CreateNotification creates a new notification
func (s *Server) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.Notification, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	payload := make(map[string]string)
	for k, v := range req.Payload {
		payload[k] = v
	}

	// For now, use aggregation based on notification type (like, comment, follower, message)
	// When protobuf is regenerated with the Aggregate field, this can be controlled from the request
	aggregate := shouldAggregateNotification(req.Type)

	notification, err := s.Application.CreateNotificationWithAggregation(
		ctx,
		req.UserId,
		convertProtoNotificationTypeToApplication(req.Type),
		req.Title,
		req.Message,
		req.SourceService,
		req.SourceEntityId,
		req.NeedsAction,
		payload,
		aggregate,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create notification: %v", err)
	}

	return s.convertToProtoNotification(notification), nil
}

// CreateNotifications creates multiple notifications
func (s *Server) CreateNotifications(ctx context.Context, req *pb.CreateNotificationsRequest) (*pb.CreateNotificationsResponse, error) {
	// Create notifications individually to allow for aggregation control
	createdNotifications := make([]*application.Notification, 0, len(req.Notifications))

	for _, n := range req.Notifications {
		payload := make(map[string]string)
		for k, v := range n.Payload {
			payload[k] = v
		}

		// For now, use aggregation based on notification type (like, comment, follower, message)
		// When protobuf is regenerated with the Aggregate field, this can be controlled from the request
		aggregate := shouldAggregateNotification(n.Type)

		notification, err := s.Application.CreateNotificationWithAggregation(
			ctx,
			n.UserId,
			convertProtoNotificationTypeToApplication(n.Type),
			n.Title,
			n.Message,
			n.SourceService,
			n.SourceEntityId,
			n.NeedsAction,
			payload,
			aggregate,  // Use the aggregate flag determined by type
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create notification: %v", err)
		}

		createdNotifications = append(createdNotifications, notification)
	}

	pbNotifications := make([]*pb.Notification, len(createdNotifications))
	for i, n := range createdNotifications {
		pbNotifications[i] = s.convertToProtoNotification(n)
	}

	return &pb.CreateNotificationsResponse{
		CreatedNotifications: pbNotifications,
	}, nil
}

// GetUserNotifications retrieves notifications for a user
func (s *Server) GetUserNotifications(ctx context.Context, req *pb.GetUserNotificationsRequest) (*pb.GetUserNotificationsResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.Limit == 0 {
		req.Limit = 20 // default limit
	}

	notifications, err := s.Application.GetUserNotifications(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user notifications: %v", err)
	}

	// Convert to protobuf notifications
	pbNotifications := make([]*pb.Notification, len(notifications))
	for i, n := range notifications {
		pbNotifications[i] = s.convertToProtoNotification(n)
	}

	// Get total count
	totalCount, err := s.Application.GetUserNotificationsCount(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get notifications count: %v", err)
	}

	// Get unread count
	unreadCount, err := s.Application.GetUserUnreadNotificationsCount(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get unread notifications count: %v", err)
	}

	return &pb.GetUserNotificationsResponse{
		Notifications: pbNotifications,
		TotalCount:    int32(totalCount),
		UnreadCount:   int32(unreadCount),
	}, nil
}

// GetUnreadNotificationsCount returns the count of unread notifications for a user
func (s *Server) GetUnreadNotificationsCount(ctx context.Context, req *wrapperspb.Int64Value) (*wrapperspb.Int64Value, error) {
	if req.Value == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	count, err := s.Application.GetUserUnreadNotificationsCount(ctx, req.Value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get unread notifications count: %v", err)
	}

	return &wrapperspb.Int64Value{
		Value: count,
	}, nil
}

// MarkNotificationAsRead marks a notification as read
func (s *Server) MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*emptypb.Empty, error) {
	if req.NotificationId == 0 || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "notification_id and user_id are required")
	}

	err := s.Application.MarkNotificationAsRead(ctx, req.NotificationId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mark notification as read: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// MarkAllAsRead marks all notifications for a user as read
func (s *Server) MarkAllAsRead(ctx context.Context, req *wrapperspb.Int64Value) (*emptypb.Empty, error) {
	if req.Value == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	err := s.Application.MarkAllAsRead(ctx, req.Value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mark all notifications as read: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// DeleteNotification deletes a notification
func (s *Server) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*emptypb.Empty, error) {
	if req.NotificationId == 0 || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "notification_id and user_id are required")
	}

	err := s.Application.DeleteNotification(ctx, req.NotificationId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete notification: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetNotificationPreferences returns notification preferences for a user
func (s *Server) GetNotificationPreferences(ctx context.Context, req *wrapperspb.Int64Value) (*pb.NotificationPreferences, error) {
	// For now, return default preferences
	// In a real implementation, this would fetch from a user preferences table
	defaultPrefs := make(map[string]bool)
	for _, notifType := range pb.NotificationType_name {
		defaultPrefs[notifType] = true
	}

	return &pb.NotificationPreferences{
		UserId:       req.Value,
		Preferences:  defaultPrefs,
	}, nil
}

// UpdateNotificationPreferences updates notification preferences for a user
func (s *Server) UpdateNotificationPreferences(ctx context.Context, req *pb.UpdateNotificationPreferencesRequest) (*emptypb.Empty, error) {
	// For now, just return success
	// In a real implementation, this would update a user preferences table
	return &emptypb.Empty{}, nil
}

// convertToProtoNotification converts our internal notification model to protobuf format
func (s *Server) convertToProtoNotification(notification *application.Notification) *pb.Notification {
	// Convert map[string]string to map[string]string (which protobuf handles as map<string, string>)
	payload := make(map[string]string)
	for k, v := range notification.Payload {
		payload[k] = v
	}

	var createdAt *timestamppb.Timestamp
	if !notification.CreatedAt.IsZero() {
		createdAt = timestamppb.New(notification.CreatedAt)
	}

	var expiresAt *timestamppb.Timestamp
	if notification.ExpiresAt != nil {
		expiresAt = timestamppb.New(*notification.ExpiresAt)
	}

	var status pb.NotificationStatus
	if notification.Seen {
		status = pb.NotificationStatus_NOTIFICATION_STATUS_READ
	} else {
		status = pb.NotificationStatus_NOTIFICATION_STATUS_UNREAD
	}

	// If deleted_at is not nil, the notification is deleted
	if notification.DeletedAt != nil {
		status = pb.NotificationStatus_NOTIFICATION_STATUS_DELETED
	}

	return &pb.Notification{
		Id:             notification.ID,
		UserId:         notification.UserID,
		Type:           convertApplicationNotificationTypeToProto(notification.Type),
		Title:          notification.Title,
		Message:        notification.Message,
		SourceService:  notification.SourceService,
		SourceEntityId: notification.SourceEntityID,
		NeedsAction:    notification.NeedsAction,
		Acted:          notification.Acted,
		Count:          notification.Count,
		Payload:        payload,
		CreatedAt:      createdAt,
		ExpiresAt:      expiresAt,
		Status:         status,
	}
}

// Helper function to determine if a notification type should be aggregated
func shouldAggregateNotification(notificationType pb.NotificationType) bool {
	switch notificationType {
	case pb.NotificationType_NOTIFICATION_TYPE_POST_LIKE:
		return true
	case pb.NotificationType_NOTIFICATION_TYPE_POST_COMMENT:
		return true
	case pb.NotificationType_NOTIFICATION_TYPE_NEW_FOLLOWER:
		return true
	case pb.NotificationType_NOTIFICATION_TYPE_NEW_MESSAGE:
		return true
	case pb.NotificationType_NOTIFICATION_TYPE_FOLLOW_REQUEST_ACCEPTED:
		return false  // Follow request responses are specific to each request
	case pb.NotificationType_NOTIFICATION_TYPE_FOLLOW_REQUEST_REJECTED:
		return false  // Follow request responses are specific to each request
	case pb.NotificationType_NOTIFICATION_TYPE_GROUP_INVITE_ACCEPTED:
		return false  // Group invite responses are specific to each invitation
	case pb.NotificationType_NOTIFICATION_TYPE_GROUP_INVITE_REJECTED:
		return false  // Group invite responses are specific to each invitation
	case pb.NotificationType_NOTIFICATION_TYPE_GROUP_JOIN_REQUEST_ACCEPTED:
		return false  // Group join request responses are specific to each request
	case pb.NotificationType_NOTIFICATION_TYPE_GROUP_JOIN_REQUEST_REJECTED:
		return false  // Group join request responses are specific to each request
	default:
		return false
	}
}

// convertApplicationNotificationTypeToProto converts application notification type to protobuf notification type
func convertApplicationNotificationTypeToProto(appType application.NotificationType) pb.NotificationType {
	switch appType {
	case application.FollowRequest:
		return pb.NotificationType_NOTIFICATION_TYPE_FOLLOW_REQUEST
	case application.NewFollower:
		return pb.NotificationType_NOTIFICATION_TYPE_NEW_FOLLOWER
	case application.GroupInvite:
		return pb.NotificationType_NOTIFICATION_TYPE_GROUP_INVITE
	case application.GroupJoinRequest:
		return pb.NotificationType_NOTIFICATION_TYPE_GROUP_JOIN_REQUEST
	case application.NewEvent:
		return pb.NotificationType_NOTIFICATION_TYPE_NEW_EVENT
	case application.PostLike:
		return pb.NotificationType_NOTIFICATION_TYPE_POST_LIKE
	case application.PostComment:
		return pb.NotificationType_NOTIFICATION_TYPE_POST_COMMENT
	case application.Mention:
		return pb.NotificationType_NOTIFICATION_TYPE_MENTION
	case application.FollowRequestAccepted:
		return pb.NotificationType_NOTIFICATION_TYPE_FOLLOW_REQUEST_ACCEPTED
	case application.FollowRequestRejected:
		return pb.NotificationType_NOTIFICATION_TYPE_FOLLOW_REQUEST_REJECTED
	case application.GroupInviteAccepted:
		return pb.NotificationType_NOTIFICATION_TYPE_GROUP_INVITE_ACCEPTED
	case application.GroupInviteRejected:
		return pb.NotificationType_NOTIFICATION_TYPE_GROUP_INVITE_REJECTED
	case application.GroupJoinRequestAccepted:
		return pb.NotificationType_NOTIFICATION_TYPE_GROUP_JOIN_REQUEST_ACCEPTED
	case application.GroupJoinRequestRejected:
		return pb.NotificationType_NOTIFICATION_TYPE_GROUP_JOIN_REQUEST_REJECTED
	default:
		return pb.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED
	}
}

// convertProtoNotificationTypeToApplication converts protobuf notification type to application notification type
func convertProtoNotificationTypeToApplication(protoType pb.NotificationType) application.NotificationType {
	switch protoType {
	case pb.NotificationType_NOTIFICATION_TYPE_FOLLOW_REQUEST:
		return application.FollowRequest
	case pb.NotificationType_NOTIFICATION_TYPE_NEW_FOLLOWER:
		return application.NewFollower
	case pb.NotificationType_NOTIFICATION_TYPE_GROUP_INVITE:
		return application.GroupInvite
	case pb.NotificationType_NOTIFICATION_TYPE_GROUP_JOIN_REQUEST:
		return application.GroupJoinRequest
	case pb.NotificationType_NOTIFICATION_TYPE_NEW_EVENT:
		return application.NewEvent
	case pb.NotificationType_NOTIFICATION_TYPE_POST_LIKE:
		return application.PostLike
	case pb.NotificationType_NOTIFICATION_TYPE_POST_COMMENT:
		return application.PostComment
	case pb.NotificationType_NOTIFICATION_TYPE_MENTION:
		return application.Mention
	case pb.NotificationType_NOTIFICATION_TYPE_FOLLOW_REQUEST_ACCEPTED:
		return application.FollowRequestAccepted
	case pb.NotificationType_NOTIFICATION_TYPE_FOLLOW_REQUEST_REJECTED:
		return application.FollowRequestRejected
	case pb.NotificationType_NOTIFICATION_TYPE_GROUP_INVITE_ACCEPTED:
		return application.GroupInviteAccepted
	case pb.NotificationType_NOTIFICATION_TYPE_GROUP_INVITE_REJECTED:
		return application.GroupInviteRejected
	case pb.NotificationType_NOTIFICATION_TYPE_GROUP_JOIN_REQUEST_ACCEPTED:
		return application.GroupJoinRequestAccepted
	case pb.NotificationType_NOTIFICATION_TYPE_GROUP_JOIN_REQUEST_REJECTED:
		return application.GroupJoinRequestRejected
	default:
		return application.NotificationType("")
	}
}
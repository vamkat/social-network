package events

import (
	"context"
	"testing"
	"social-network/services/notifications/internal/application"
	pb "social-network/shared/gen-go/notifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockApplication is a mock that implements the Application interface
type MockApplication struct {
	mock.Mock
}

// Implement all the methods needed for testing
func (m *MockApplication) CreatePostCommentNotification(ctx context.Context, userID, commenterID, postID int64, commenterUsername, commentContent string, aggregate bool) error {
	args := m.Called(ctx, userID, commenterID, postID, commenterUsername, commentContent, aggregate)
	return args.Error(0)
}

func (m *MockApplication) CreatePostLikeNotification(ctx context.Context, userID, likerID, postID int64, likerUsername string, aggregate bool) error {
	args := m.Called(ctx, userID, likerID, postID, likerUsername, aggregate)
	return args.Error(0)
}

func (m *MockApplication) CreateFollowRequestNotification(ctx context.Context, targetUserID, requesterUserID int64, requesterUsername string) error {
	args := m.Called(ctx, targetUserID, requesterUserID, requesterUsername)
	return args.Error(0)
}

func (m *MockApplication) CreateNewFollowerNotification(ctx context.Context, targetUserID, followerUserID int64, followerUsername string, aggregate bool) error {
	args := m.Called(ctx, targetUserID, followerUserID, followerUsername, aggregate)
	return args.Error(0)
}

func (m *MockApplication) CreateGroupInviteNotification(ctx context.Context, invitedUserID, inviterUserID, groupID int64, groupName, inviterUsername string) error {
	args := m.Called(ctx, invitedUserID, inviterUserID, groupID, groupName, inviterUsername)
	return args.Error(0)
}

func (m *MockApplication) CreateGroupJoinRequestNotification(ctx context.Context, groupOwnerID, requesterID, groupID int64, groupName, requesterUsername string) error {
	args := m.Called(ctx, groupOwnerID, requesterID, groupID, groupName, requesterUsername)
	return args.Error(0)
}

func (m *MockApplication) CreateNewEventNotification(ctx context.Context, userID, groupID, eventID int64, groupName, eventTitle string) error {
	args := m.Called(ctx, userID, groupID, eventID, groupName, eventTitle)
	return args.Error(0)
}

func (m *MockApplication) CreateMentionNotification(ctx context.Context, userID, mentionerID, postID int64, mentionerUsername, postContent, mentionText string) error {
	args := m.Called(ctx, userID, mentionerID, postID, mentionerUsername, postContent, mentionText)
	return args.Error(0)
}

func (m *MockApplication) CreateNewMessageNotification(ctx context.Context, userID, senderID, chatID int64, senderUsername, messageContent string, aggregate bool) error {
	args := m.Called(ctx, userID, senderID, chatID, senderUsername, messageContent, aggregate)
	return args.Error(0)
}

func (m *MockApplication) CreateFollowRequestAcceptedNotification(ctx context.Context, requesterUserID, targetUserID int64, targetUsername string) error {
	args := m.Called(ctx, requesterUserID, targetUserID, targetUsername)
	return args.Error(0)
}

func (m *MockApplication) CreateFollowRequestRejectedNotification(ctx context.Context, requesterUserID, targetUserID int64, targetUsername string) error {
	args := m.Called(ctx, requesterUserID, targetUserID, targetUsername)
	return args.Error(0)
}

func (m *MockApplication) CreateGroupInviteAcceptedNotification(ctx context.Context, inviterUserID, invitedUserID, groupID int64, invitedUsername, groupName string) error {
	args := m.Called(ctx, inviterUserID, invitedUserID, groupID, invitedUsername, groupName)
	return args.Error(0)
}

func (m *MockApplication) CreateGroupInviteRejectedNotification(ctx context.Context, inviterUserID, invitedUserID, groupID int64, invitedUsername, groupName string) error {
	args := m.Called(ctx, inviterUserID, invitedUserID, groupID, invitedUsername, groupName)
	return args.Error(0)
}

func (m *MockApplication) CreateGroupJoinRequestAcceptedNotification(ctx context.Context, requesterUserID, groupOwnerID, groupID int64, groupName string) error {
	args := m.Called(ctx, requesterUserID, groupOwnerID, groupID, groupName)
	return args.Error(0)
}

func (m *MockApplication) CreateGroupJoinRequestRejectedNotification(ctx context.Context, requesterUserID, groupOwnerID, groupID int64, groupName string) error {
	args := m.Called(ctx, requesterUserID, groupOwnerID, groupID, groupName)
	return args.Error(0)
}

func (m *MockApplication) CreateDefaultNotificationTypes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockApplication) GetNotification(ctx context.Context, notificationID, userID int64) (*application.Notification, error) {
	args := m.Called(ctx, notificationID, userID)
	return args.Get(0).(*application.Notification), args.Error(1)
}

func (m *MockApplication) GetUserNotifications(ctx context.Context, userID int64, limit, offset int32) ([]*application.Notification, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*application.Notification), args.Error(1)
}

func (m *MockApplication) GetUserNotificationsCount(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockApplication) GetUserUnreadNotificationsCount(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockApplication) MarkNotificationAsRead(ctx context.Context, notificationID, userID int64) error {
	args := m.Called(ctx, notificationID, userID)
	return args.Error(0)
}

func (m *MockApplication) MarkAllAsRead(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockApplication) DeleteNotification(ctx context.Context, notificationID, userID int64) error {
	args := m.Called(ctx, notificationID, userID)
	return args.Error(0)
}

func (m *MockApplication) CreateNotification(ctx context.Context, userID int64, notifType application.NotificationType, title, message, sourceService string, sourceEntityID int64, needsAction bool, payload map[string]string) (*application.Notification, error) {
	args := m.Called(ctx, userID, notifType, title, message, sourceService, sourceEntityID, needsAction, payload)
	return args.Get(0).(*application.Notification), args.Error(1)
}

func (m *MockApplication) CreateNotificationWithAggregation(ctx context.Context, userID int64, notifType application.NotificationType, title, message, sourceService string, sourceEntityID int64, needsAction bool, payload map[string]string, aggregate bool) (*application.Notification, error) {
	args := m.Called(ctx, userID, notifType, title, message, sourceService, sourceEntityID, needsAction, payload, aggregate)
	return args.Get(0).(*application.Notification), args.Error(1)
}

func (m *MockApplication) CreateNotifications(ctx context.Context, notifications []struct {
	UserID         int64
	Type           application.NotificationType
	Title          string
	Message        string
	SourceService  string
	SourceEntityID int64
	NeedsAction    bool
	Payload        map[string]string
}) ([]*application.Notification, error) {
	args := m.Called(ctx, notifications)
	return args.Get(0).([]*application.Notification), args.Error(1)
}

// Unit tests for each event handler
func TestEventHandler_HandlePostCommentCreated(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_POST_COMMENT_CREATED,
		Payload: &pb.NotificationEvent_PostCommentCreated{
			PostCommentCreated: &pb.PostCommentCreated{
				PostId:             123,
				CommentId:          456,
				CommenterUserId:    789,
				CommenterUsername:  "test_user",
				Body:               "This is a test comment",
				Aggregate:          true,
			},
		},
	}

	// Set up expectations
	mockApp.On("CreatePostCommentNotification", 
		mock.Anything, 
		int64(123),      // userID (post owner)
		int64(789),      // commenterID
		int64(123),      // postID
		"test_user",     // commenterUsername
		"This is a test comment", // commentContent
		true,            // aggregate
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandlePostLiked(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_POST_LIKED,
		Payload: &pb.NotificationEvent_PostLiked{
			PostLiked: &pb.PostLiked{
				PostId:        123,
				LikerUserId:   789,
				LikerUsername: "test_user",
				Aggregate:     true,
			},
		},
	}

	// Set up expectations
	mockApp.On("CreatePostLikeNotification", 
		mock.Anything, 
		int64(123),      // userID (post owner)
		int64(789),      // likerID
		int64(123),      // postID
		"test_user",     // likerUsername
		true,            // aggregate
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleFollowRequestCreated(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_FOLLOW_REQUEST_CREATED,
		Payload: &pb.NotificationEvent_FollowRequestCreated{
			FollowRequestCreated: &pb.FollowRequestCreated{
				TargetUserId:      123,
				RequesterUserId:   789,
				RequesterUsername: "test_user",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateFollowRequestNotification", 
		mock.Anything, 
		int64(123),      // targetUserID
		int64(789),      // requesterUserID
		"test_user",     // requesterUsername
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleNewFollowerCreated(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_NEW_FOLLOWER_CREATED,
		Payload: &pb.NotificationEvent_NewFollowerCreated{
			NewFollowerCreated: &pb.NewFollowerCreated{
				TargetUserId:      123,
				FollowerUserId:    789,
				FollowerUsername:  "test_user",
				Aggregate:         true,
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateNewFollowerNotification", 
		mock.Anything, 
		int64(123),      // targetUserID
		int64(789),      // followerUserID
		"test_user",     // followerUsername
		true,            // aggregate
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleGroupInviteCreated(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_GROUP_INVITE_CREATED,
		Payload: &pb.NotificationEvent_GroupInviteCreated{
			GroupInviteCreated: &pb.GroupInviteCreated{
				InvitedUserId:     123,
				InviterUserId:     789,
				GroupId:           456,
				GroupName:         "test_group",
				InviterUsername:   "inviter_user",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateGroupInviteNotification", 
		mock.Anything, 
		int64(123),      // invitedUserID
		int64(789),      // inviterUserID
		int64(456),      // groupID
		"test_group",    // groupName
		"inviter_user",  // inviterUsername
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleGroupJoinRequestCreated(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_GROUP_JOIN_REQUEST_CREATED,
		Payload: &pb.NotificationEvent_GroupJoinRequestCreated{
			GroupJoinRequestCreated: &pb.GroupJoinRequestCreated{
				GroupOwnerId:      123,
				RequesterUserId:   789,
				GroupId:           456,
				GroupName:         "test_group",
				RequesterUsername: "requester_user",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateGroupJoinRequestNotification", 
		mock.Anything, 
		int64(123),      // groupOwnerID
		int64(789),      // requesterID
		int64(456),      // groupID
		"test_group",    // groupName
		"requester_user", // requesterUsername
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleNewEventCreated(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_NEW_EVENT_CREATED,
		Payload: &pb.NotificationEvent_NewEventCreated{
			NewEventCreated: &pb.NewEventCreated{
				UserId:      123,
				GroupId:     456,
				EventId:     789,
				GroupName:   "test_group",
				EventTitle:  "Test Event",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateNewEventNotification", 
		mock.Anything, 
		int64(123),      // userID
		int64(456),      // groupID
		int64(789),      // eventID
		"test_group",    // groupName
		"Test Event",    // eventTitle
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleMentionCreated(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_MENTION_CREATED,
		Payload: &pb.NotificationEvent_MentionCreated{
			MentionCreated: &pb.MentionCreated{
				MentionedUserId:   123,
				MentionerUserId:   789,
				PostId:            456,
				MentionerUsername: "mentioner_user",
				PostContent:       "This is a post with @mention",
				MentionText:       "@mention",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateMentionNotification", 
		mock.Anything, 
		int64(123),      // userID
		int64(789),      // mentionerID
		int64(456),      // postID
		"mentioner_user", // mentionerUsername
		"This is a post with @mention", // postContent
		"@mention",      // mentionText
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleNewMessageCreated(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_NEW_MESSAGE_CREATED,
		Payload: &pb.NotificationEvent_NewMessageCreated{
			NewMessageCreated: &pb.NewMessageCreated{
				UserId:           123,
				SenderUserId:     789,
				ChatId:           456,
				SenderUsername:   "sender_user",
				MessageContent:   "Hello, this is a test message!",
				Aggregate:        true,
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateNewMessageNotification", 
		mock.Anything, 
		int64(123),      // userID
		int64(789),      // senderID
		int64(456),      // chatID
		"sender_user",   // senderUsername
		"Hello, this is a test message!", // messageContent
		true,            // aggregate
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleFollowRequestAccepted(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_FOLLOW_REQUEST_ACCEPTED,
		Payload: &pb.NotificationEvent_FollowRequestAccepted{
			FollowRequestAccepted: &pb.FollowRequestAccepted{
				RequesterUserId: 123,
				TargetUserId:    789,
				TargetUsername:  "target_user",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateFollowRequestAcceptedNotification", 
		mock.Anything, 
		int64(123),      // requesterUserID
		int64(789),      // targetUserID
		"target_user",   // targetUsername
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleFollowRequestRejected(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_FOLLOW_REQUEST_REJECTED,
		Payload: &pb.NotificationEvent_FollowRequestRejected{
			FollowRequestRejected: &pb.FollowRequestRejected{
				RequesterUserId: 123,
				TargetUserId:    789,
				TargetUsername:  "target_user",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateFollowRequestRejectedNotification", 
		mock.Anything, 
		int64(123),      // requesterUserID
		int64(789),      // targetUserID
		"target_user",   // targetUsername
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleGroupInviteAccepted(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_GROUP_INVITE_ACCEPTED,
		Payload: &pb.NotificationEvent_GroupInviteAccepted{
			GroupInviteAccepted: &pb.GroupInviteAccepted{
				InviterUserId:    123,
				InvitedUserId:    789,
				GroupId:          456,
				InvitedUsername:  "invited_user",
				GroupName:        "test_group",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateGroupInviteAcceptedNotification", 
		mock.Anything, 
		int64(123),      // inviterUserID
		int64(789),      // invitedUserID
		int64(456),      // groupID
		"invited_user",  // invitedUsername
		"test_group",    // groupName
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleGroupInviteRejected(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_GROUP_INVITE_REJECTED,
		Payload: &pb.NotificationEvent_GroupInviteRejected{
			GroupInviteRejected: &pb.GroupInviteRejected{
				InviterUserId:    123,
				InvitedUserId:    789,
				GroupId:          456,
				InvitedUsername:  "invited_user",
				GroupName:        "test_group",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateGroupInviteRejectedNotification", 
		mock.Anything, 
		int64(123),      // inviterUserID
		int64(789),      // invitedUserID
		int64(456),      // groupID
		"invited_user",  // invitedUsername
		"test_group",    // groupName
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleGroupJoinRequestAccepted(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_GROUP_JOIN_REQUEST_ACCEPTED,
		Payload: &pb.NotificationEvent_GroupJoinRequestAccepted{
			GroupJoinRequestAccepted: &pb.GroupJoinRequestAccepted{
				RequesterUserId: 123,
				GroupOwnerId:    789,
				GroupId:         456,
				GroupName:       "test_group",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateGroupJoinRequestAcceptedNotification", 
		mock.Anything, 
		int64(123),      // requesterUserID
		int64(789),      // groupOwnerID
		int64(456),      // groupID
		"test_group",    // groupName
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleGroupJoinRequestRejected(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_GROUP_JOIN_REQUEST_REJECTED,
		Payload: &pb.NotificationEvent_GroupJoinRequestRejected{
			GroupJoinRequestRejected: &pb.GroupJoinRequestRejected{
				RequesterUserId: 123,
				GroupOwnerId:    789,
				GroupId:         456,
				GroupName:       "test_group",
			},
		},
	}

	// Set up expectations
	mockApp.On("CreateGroupJoinRequestRejectedNotification", 
		mock.Anything, 
		int64(123),      // requesterUserID
		int64(789),      // groupOwnerID
		int64(456),      // groupID
		"test_group",    // groupName
	).Return(nil)

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestEventHandler_HandleUnknownEvent(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	// Create an event with no payload
	event := &pb.NotificationEvent{
		EventId:    "test-event-id",
		EventType:  pb.EventType_EVENT_TYPE_UNSPECIFIED,
		Payload:    nil,
	}

	// Execute
	err := eventHandler.Handle(context.Background(), event)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown notification event payload type")
}

// Integration test: test the event handler with all event types
func TestEventHandler_Integration(t *testing.T) {
	mockApp := new(MockApplication)
	eventHandler := &EventHandler{App: mockApp}

	// Test all event types in sequence
	events := []*pb.NotificationEvent{
		{
			EventId:    "test-1",
			EventType:  pb.EventType_POST_COMMENT_CREATED,
			Payload: &pb.NotificationEvent_PostCommentCreated{
				PostCommentCreated: &pb.PostCommentCreated{
					PostId:             123,
					CommenterUserId:    456,
					CommenterUsername:  "user1",
					Body:               "comment body",
					Aggregate:          true,
				},
			},
		},
		{
			EventId:    "test-2",
			EventType:  pb.EventType_POST_LIKED,
			Payload: &pb.NotificationEvent_PostLiked{
				PostLiked: &pb.PostLiked{
					PostId:        123,
					LikerUserId:   456,
					LikerUsername: "user2",
					Aggregate:     true,
				},
			},
		},
		{
			EventId:    "test-3",
			EventType:  pb.EventType_FOLLOW_REQUEST_CREATED,
			Payload: &pb.NotificationEvent_FollowRequestCreated{
				FollowRequestCreated: &pb.FollowRequestCreated{
					TargetUserId:      123,
					RequesterUserId:   456,
					RequesterUsername: "user3",
				},
			},
		},
	}

	// Set up expectations for all events
	mockApp.On("CreatePostCommentNotification", mock.Anything, int64(123), int64(456), int64(123), "user1", "comment body", true).Return(nil)
	mockApp.On("CreatePostLikeNotification", mock.Anything, int64(123), int64(456), int64(123), "user2", true).Return(nil)
	mockApp.On("CreateFollowRequestNotification", mock.Anything, int64(123), int64(456), "user3").Return(nil)

	// Execute all events
	for _, event := range events {
		err := eventHandler.Handle(context.Background(), event)
		assert.NoError(t, err)
	}

	// Assert all expectations were met
	mockApp.AssertExpectations(t)
}
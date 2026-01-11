package events

import (
	"context"
	"testing"
	"time"
	pb "social-network/shared/gen-go/notifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite runs integration tests for the event handler
type IntegrationTestSuite struct {
	suite.Suite
	eventHandler *EventHandler
	mockApp      *MockApplication
}

// SetupTest runs before each test
func (suite *IntegrationTestSuite) SetupTest() {
	suite.mockApp = new(MockApplication)
	suite.eventHandler = &EventHandler{App: suite.mockApp}
}

// Test that all event types can be processed successfully
func (suite *IntegrationTestSuite) TestProcessEventTypesIntegration() {
	ctx := context.Background()

	// Test PostCommentCreated
	suite.mockApp.On("CreatePostCommentNotification", mock.Anything, int64(123), int64(456), int64(123), "user1", "comment body", true).Return(nil)
	event1 := &pb.NotificationEvent{
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
	}
	err := suite.eventHandler.Handle(ctx, event1)
	assert.NoError(suite.T(), err)

	// Test PostLiked
	suite.mockApp.On("CreatePostLikeNotification", mock.Anything, int64(123), int64(456), int64(123), "user2", true).Return(nil)
	event2 := &pb.NotificationEvent{
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
	}
	err = suite.eventHandler.Handle(ctx, event2)
	assert.NoError(suite.T(), err)

	// Test FollowRequestCreated
	suite.mockApp.On("CreateFollowRequestNotification", mock.Anything, int64(123), int64(456), "user3").Return(nil)
	event3 := &pb.NotificationEvent{
		EventId:    "test-3",
		EventType:  pb.EventType_FOLLOW_REQUEST_CREATED,
		Payload: &pb.NotificationEvent_FollowRequestCreated{
			FollowRequestCreated: &pb.FollowRequestCreated{
				TargetUserId:      123,
				RequesterUserId:   456,
				RequesterUsername: "user3",
			},
		},
	}
	err = suite.eventHandler.Handle(ctx, event3)
	assert.NoError(suite.T(), err)

	// Test NewFollowerCreated
	suite.mockApp.On("CreateNewFollowerNotification", mock.Anything, int64(123), int64(456), "user4", true).Return(nil)
	event4 := &pb.NotificationEvent{
		EventId:    "test-4",
		EventType:  pb.EventType_NEW_FOLLOWER_CREATED,
		Payload: &pb.NotificationEvent_NewFollowerCreated{
			NewFollowerCreated: &pb.NewFollowerCreated{
				TargetUserId:      123,
				FollowerUserId:    456,
				FollowerUsername:  "user4",
				Aggregate:         true,
			},
		},
	}
	err = suite.eventHandler.Handle(ctx, event4)
	assert.NoError(suite.T(), err)

	// Test GroupInviteCreated
	suite.mockApp.On("CreateGroupInviteNotification", mock.Anything, int64(123), int64(456), int64(789), "group1", "user5").Return(nil)
	event5 := &pb.NotificationEvent{
		EventId:    "test-5",
		EventType:  pb.EventType_GROUP_INVITE_CREATED,
		Payload: &pb.NotificationEvent_GroupInviteCreated{
			GroupInviteCreated: &pb.GroupInviteCreated{
				InvitedUserId:     123,
				InviterUserId:     456,
				GroupId:           789,
				GroupName:         "group1",
				InviterUsername:   "user5",
			},
		},
	}
	err = suite.eventHandler.Handle(ctx, event5)
	assert.NoError(suite.T(), err)

	// Verify all expectations were met
	suite.mockApp.AssertExpectations(suite.T())
}

// Test error handling when application layer returns an error
func (suite *IntegrationTestSuite) TestErrorHandlingIntegration() {
	ctx := context.Background()

	// Set up expectation to return an error
	expectedErr := assert.AnError
	suite.mockApp.On("CreatePostCommentNotification", mock.Anything, int64(123), int64(456), int64(123), "user1", "comment body", true).Return(expectedErr)

	event := &pb.NotificationEvent{
		EventId:    "error-test",
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
	}

	err := suite.eventHandler.Handle(ctx, event)
	assert.Equal(suite.T(), expectedErr, err)
	suite.mockApp.AssertExpectations(suite.T())
}

// Test processing multiple events in sequence
func (suite *IntegrationTestSuite) TestSequentialProcessingIntegration() {
	ctx := context.Background()

	events := []*pb.NotificationEvent{
		{
			EventId:    "seq-1",
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
			EventId:    "seq-2",
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
			EventId:    "seq-3",
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
	suite.mockApp.On("CreatePostCommentNotification", mock.Anything, int64(123), int64(456), int64(123), "user1", "comment body", true).Return(nil)
	suite.mockApp.On("CreatePostLikeNotification", mock.Anything, int64(123), int64(456), int64(123), "user2", true).Return(nil)
	suite.mockApp.On("CreateFollowRequestNotification", mock.Anything, int64(123), int64(456), "user3").Return(nil)

	// Process all events sequentially
	for _, event := range events {
		err := suite.eventHandler.Handle(ctx, event)
		assert.NoError(suite.T(), err)
	}

	// Verify all expectations were met
	suite.mockApp.AssertExpectations(suite.T())
}

// Test concurrent processing of events (simulated)
func (suite *IntegrationTestSuite) TestConcurrentProcessingIntegration() {
	ctx := context.Background()

	// Simulate concurrent processing by processing multiple events rapidly
	events := []*pb.NotificationEvent{
		{
			EventId:    "concurrent-1",
			EventType:  pb.EventType_NEW_FOLLOWER_CREATED,
			Payload: &pb.NotificationEvent_NewFollowerCreated{
				NewFollowerCreated: &pb.NewFollowerCreated{
					TargetUserId:      123,
					FollowerUserId:    456,
					FollowerUsername:  "user1",
					Aggregate:         true,
				},
			},
		},
		{
			EventId:    "concurrent-2",
			EventType:  pb.EventType_NEW_MESSAGE_CREATED,
			Payload: &pb.NotificationEvent_NewMessageCreated{
				NewMessageCreated: &pb.NewMessageCreated{
					UserId:           123,
					SenderUserId:     456,
					ChatId:           789,
					SenderUsername:   "user2",
					MessageContent:   "hello",
					Aggregate:        true,
				},
			},
		},
		{
			EventId:    "concurrent-3",
			EventType:  pb.EventType_MENTION_CREATED,
			Payload: &pb.NotificationEvent_MentionCreated{
				MentionCreated: &pb.MentionCreated{
					MentionedUserId:   123,
					MentionerUserId:   456,
					PostId:            789,
					MentionerUsername: "user3",
					PostContent:       "post content",
					MentionText:       "@user",
				},
			},
		},
	}

	// Set up expectations
	suite.mockApp.On("CreateNewFollowerNotification", mock.Anything, int64(123), int64(456), "user1", true).Return(nil)
	suite.mockApp.On("CreateNewMessageNotification", mock.Anything, int64(123), int64(456), int64(789), "user2", "hello", true).Return(nil)
	suite.mockApp.On("CreateMentionNotification", mock.Anything, int64(123), int64(456), int64(789), "user3", "post content", "@user").Return(nil)

	// Process events (simulating concurrent processing)
	for _, event := range events {
		err := suite.eventHandler.Handle(ctx, event)
		assert.NoError(suite.T(), err)
	}

	// Verify all expectations were met
	suite.mockApp.AssertExpectations(suite.T())
}

// Test aggregation behavior with multiple similar events
func (suite *IntegrationTestSuite) TestAggregationBehaviorIntegration() {
	ctx := context.Background()

	// Test multiple post likes (should aggregate)
	events := []*pb.NotificationEvent{
		{
			EventId:    "agg-1",
			EventType:  pb.EventType_POST_LIKED,
			Payload: &pb.NotificationEvent_PostLiked{
				PostLiked: &pb.PostLiked{
					PostId:        123,
					LikerUserId:   456,
					LikerUsername: "user1",
					Aggregate:     true,
				},
			},
		},
		{
			EventId:    "agg-2",
			EventType:  pb.EventType_POST_LIKED,
			Payload: &pb.NotificationEvent_PostLiked{
				PostLiked: &pb.PostLiked{
					PostId:        123,
					LikerUserId:   789,
					LikerUsername: "user2",
					Aggregate:     true,
				},
			},
		},
		{
			EventId:    "agg-3",
			EventType:  pb.EventType_POST_LIKED,
			Payload: &pb.NotificationEvent_PostLiked{
				PostLiked: &pb.PostLiked{
					PostId:        123,
					LikerUserId:   101,
					LikerUsername: "user3",
					Aggregate:     true,
				},
			},
		},
	}

	// Set up expectations for each event
	suite.mockApp.On("CreatePostLikeNotification", mock.Anything, int64(123), int64(456), int64(123), "user1", true).Return(nil)
	suite.mockApp.On("CreatePostLikeNotification", mock.Anything, int64(123), int64(789), int64(123), "user2", true).Return(nil)
	suite.mockApp.On("CreatePostLikeNotification", mock.Anything, int64(123), int64(101), int64(123), "user3", true).Return(nil)

	// Process all aggregation events
	for _, event := range events {
		err := suite.eventHandler.Handle(ctx, event)
		assert.NoError(suite.T(), err)
	}

	// Verify all expectations were met
	suite.mockApp.AssertExpectations(suite.T())
}

// Test non-aggregatable events (should not aggregate)
func (suite *IntegrationTestSuite) TestNonAggregatableEventsIntegration() {
	ctx := context.Background()

	// Test follow request events (should not aggregate)
	events := []*pb.NotificationEvent{
		{
			EventId:    "nonagg-1",
			EventType:  pb.EventType_FOLLOW_REQUEST_CREATED,
			Payload: &pb.NotificationEvent_FollowRequestCreated{
				FollowRequestCreated: &pb.FollowRequestCreated{
					TargetUserId:      123,
					RequesterUserId:   456,
					RequesterUsername: "user1",
				},
			},
		},
		{
			EventId:    "nonagg-2",
			EventType:  pb.EventType_FOLLOW_REQUEST_CREATED,
			Payload: &pb.NotificationEvent_FollowRequestCreated{
				FollowRequestCreated: &pb.FollowRequestCreated{
					TargetUserId:      123,
					RequesterUserId:   789,
					RequesterUsername: "user2",
				},
			},
		},
	}

	// Set up expectations for each event
	suite.mockApp.On("CreateFollowRequestNotification", mock.Anything, int64(123), int64(456), "user1").Return(nil)
	suite.mockApp.On("CreateFollowRequestNotification", mock.Anything, int64(123), int64(789), "user2").Return(nil)

	// Process all non-aggregatable events
	for _, event := range events {
		err := suite.eventHandler.Handle(ctx, event)
		assert.NoError(suite.T(), err)
	}

	// Verify all expectations were met
	suite.mockApp.AssertExpectations(suite.T())
}

// Test mixed event types processing
func (suite *IntegrationTestSuite) TestMixedEventTypesIntegration() {
	ctx := context.Background()

	// Mix of different event types
	events := []*pb.NotificationEvent{
		{
			EventId:    "mixed-1",
			EventType:  pb.EventType_POST_COMMENT_CREATED,
			Payload: &pb.NotificationEvent_PostCommentCreated{
				PostCommentCreated: &pb.PostCommentCreated{
					PostId:             123,
					CommenterUserId:    456,
					CommenterUsername:  "user1",
					Body:               "comment 1",
					Aggregate:          true,
				},
			},
		},
		{
			EventId:    "mixed-2",
			EventType:  pb.EventType_FOLLOW_REQUEST_CREATED,
			Payload: &pb.NotificationEvent_FollowRequestCreated{
				FollowRequestCreated: &pb.FollowRequestCreated{
					TargetUserId:      123,
					RequesterUserId:   456,
					RequesterUsername: "user2",
				},
			},
		},
		{
			EventId:    "mixed-3",
			EventType:  pb.EventType_NEW_MESSAGE_CREATED,
			Payload: &pb.NotificationEvent_NewMessageCreated{
				NewMessageCreated: &pb.NewMessageCreated{
					UserId:           123,
					SenderUserId:     456,
					ChatId:           789,
					SenderUsername:   "user3",
					MessageContent:   "message 1",
					Aggregate:        true,
				},
			},
		},
		{
			EventId:    "mixed-4",
			EventType:  pb.EventType_GROUP_INVITE_CREATED,
			Payload: &pb.NotificationEvent_GroupInviteCreated{
				GroupInviteCreated: &pb.GroupInviteCreated{
					InvitedUserId:     123,
					InviterUserId:     456,
					GroupId:           789,
					GroupName:         "group1",
					InviterUsername:   "user4",
				},
			},
		},
	}

	// Set up expectations for all events
	suite.mockApp.On("CreatePostCommentNotification", mock.Anything, int64(123), int64(456), int64(123), "user1", "comment 1", true).Return(nil)
	suite.mockApp.On("CreateFollowRequestNotification", mock.Anything, int64(123), int64(456), "user2").Return(nil)
	suite.mockApp.On("CreateNewMessageNotification", mock.Anything, int64(123), int64(456), int64(789), "user3", "message 1", true).Return(nil)
	suite.mockApp.On("CreateGroupInviteNotification", mock.Anything, int64(123), int64(456), int64(789), "group1", "user4").Return(nil)

	// Process mixed event types
	for _, event := range events {
		err := suite.eventHandler.Handle(ctx, event)
		assert.NoError(suite.T(), err)
	}

	// Verify all expectations were met
	suite.mockApp.AssertExpectations(suite.T())
}

// Test performance with many events
func (suite *IntegrationTestSuite) TestPerformanceIntegration() {
	ctx := context.Background()

	const numEvents = 100
	events := make([]*pb.NotificationEvent, numEvents)

	// Create many similar events
	for i := 0; i < numEvents; i++ {
		events[i] = &pb.NotificationEvent{
			EventId:    "perf-" + string(rune('0'+i%10)),
			EventType:  pb.EventType_POST_LIKED,
			Payload: &pb.NotificationEvent_PostLiked{
				PostLiked: &pb.PostLiked{
					PostId:        123,
					LikerUserId:   int64(456 + i),
					LikerUsername: "user" + string(rune('0'+i%10)),
					Aggregate:     true,
				},
			},
		}
	}

	// Set up expectations for all events
	for i := 0; i < numEvents; i++ {
		suite.mockApp.On("CreatePostLikeNotification", mock.Anything, int64(123), int64(456+i), int64(123), "user"+string(rune('0'+i%10)), true).Return(nil)
	}

	// Measure processing time
	start := time.Now()
	for _, event := range events {
		err := suite.eventHandler.Handle(ctx, event)
		assert.NoError(suite.T(), err)
	}
	duration := time.Since(start)

	// Verify performance (should be fast)
	assert.True(suite.T(), duration < 5*time.Second, "Processing took too long: %v", duration)

	// Verify all expectations were met
	suite.mockApp.AssertExpectations(suite.T())
}

// Test that unknown event types return appropriate errors
func (suite *IntegrationTestSuite) TestUnknownEventTypeIntegration() {
	ctx := context.Background()

	// Test unknown event type
	event := &pb.NotificationEvent{
		EventId:    "unknown-test",
		EventType:  pb.EventType_EVENT_TYPE_UNSPECIFIED,
		Payload:    nil,
	}

	err := suite.eventHandler.Handle(ctx, event)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "unknown notification event payload type")
}

// Run the test suite
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
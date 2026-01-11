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

// EndToEndIntegrationTestSuite tests the complete flow from event reception to processing
type EndToEndIntegrationTestSuite struct {
	suite.Suite
	eventHandler *EventHandler
	mockApp      *MockApplication
	ctx          context.Context
}

// SetupTest runs before each test
func (suite *EndToEndIntegrationTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockApp = new(MockApplication)
	suite.eventHandler = &EventHandler{App: suite.mockApp}
}

// TestCompleteEventFlow tests the complete flow of receiving and processing events
func (suite *EndToEndIntegrationTestSuite) TestCompleteEventFlow() {
	// Simulate a realistic sequence of events that might happen in the system
	events := []*pb.NotificationEvent{
		// User A follows User B
		{
			EventId:    "event-1",
			EventType:  pb.EventType_FOLLOW_REQUEST_CREATED,
			Metadata:   map[string]string{"source": "users-service"},
			Payload: &pb.NotificationEvent_FollowRequestCreated{
				FollowRequestCreated: &pb.FollowRequestCreated{
					TargetUserId:      1001,
					RequesterUserId:   2001,
					RequesterUsername: "user_a",
				},
			},
		},
		// User B posts something
		{
			EventId:    "event-2",
			EventType:  pb.EventType_NEW_EVENT_CREATED,
			Metadata:   map[string]string{"source": "posts-service"},
			Payload: &pb.NotificationEvent_NewEventCreated{
				NewEventCreated: &pb.NewEventCreated{
					UserId:      1001,
					GroupId:     3001,
					EventId:     4001,
					GroupName:   "Tech Enthusiasts",
					EventTitle:  "Weekly Tech Meetup",
				},
			},
		},
		// Someone comments on User B's post
		{
			EventId:    "event-3",
			EventType:  pb.EventType_POST_COMMENT_CREATED,
			Metadata:   map[string]string{"source": "posts-service"},
			Payload: &pb.NotificationEvent_PostCommentCreated{
				PostCommentCreated: &pb.PostCommentCreated{
					PostId:             5001,
					CommentId:          6001,
					CommenterUserId:    3001,
					CommenterUsername:  "user_c",
					Body:               "Great post! Thanks for sharing.",
					Aggregate:          true,
				},
			},
		},
		// Someone likes User B's post
		{
			EventId:    "event-4",
			EventType:  pb.EventType_POST_LIKED,
			Metadata:   map[string]string{"source": "posts-service"},
			Payload: &pb.NotificationEvent_PostLiked{
				PostLiked: &pb.PostLiked{
					PostId:        5001,
					LikerUserId:   4001,
					LikerUsername: "user_d",
					Aggregate:     true,
				},
			},
		},
		// User B gets a message
		{
			EventId:    "event-5",
			EventType:  pb.EventType_NEW_MESSAGE_CREATED,
			Metadata:   map[string]string{"source": "chat-service"},
			Payload: &pb.NotificationEvent_NewMessageCreated{
				NewMessageCreated: &pb.NewMessageCreated{
					UserId:           1001,
					SenderUserId:     5001,
					ChatId:           7001,
					SenderUsername:   "user_e",
					MessageContent:   "Hey, are we still meeting tomorrow?",
					Aggregate:        true,
				},
			},
		},
	}

	// Set up expectations for all events
	suite.mockApp.On("CreateFollowRequestNotification", 
		mock.Anything, int64(1001), int64(2001), "user_a").Return(nil)
	suite.mockApp.On("CreateNewEventNotification", 
		mock.Anything, int64(1001), int64(3001), int64(4001), "Tech Enthusiasts", "Weekly Tech Meetup").Return(nil)
	suite.mockApp.On("CreatePostCommentNotification", 
		mock.Anything, int64(5001), int64(3001), int64(5001), "user_c", "Great post! Thanks for sharing.", true).Return(nil)
	suite.mockApp.On("CreatePostLikeNotification", 
		mock.Anything, int64(5001), int64(4001), int64(5001), "user_d", true).Return(nil)
	suite.mockApp.On("CreateNewMessageNotification", 
		mock.Anything, int64(1001), int64(5001), int64(7001), "user_e", "Hey, are we still meeting tomorrow?", true).Return(nil)

	// Process all events in sequence (simulating real Kafka consumption)
	for i, event := range events {
		err := suite.eventHandler.Handle(suite.ctx, event)
		assert.NoError(suite.T(), err, "Event %d should be processed successfully", i+1)
	}

	// Verify all expectations were met
	suite.mockApp.AssertExpectations(suite.T())
}

// TestRealWorldScenario tests a realistic scenario with mixed event types
func (suite *EndToEndIntegrationTestSuite) TestRealWorldScenario() {
	// Simulate a busy day on the social network
	scenarioEvents := []*pb.NotificationEvent{
		// Morning: Several likes on yesterday's post
		{
			EventId:   "morning-1",
			EventType: pb.EventType_POST_LIKED,
			Payload: &pb.NotificationEvent_PostLiked{
				PostLiked: &pb.PostLiked{
					PostId:        1001,
					LikerUserId:   2001,
					LikerUsername: "alice",
					Aggregate:     true,
				},
			},
		},
		{
			EventId:   "morning-2",
			EventType: pb.EventType_POST_LIKED,
			Payload: &pb.NotificationEvent_PostLiked{
				PostLiked: &pb.PostLiked{
					PostId:        1001,
					LikerUserId:   2002,
					LikerUsername: "bob",
					Aggregate:     true,
				},
			},
		},
		// Mid-day: New followers
		{
			EventId:   "midday-1",
			EventType: pb.EventType_NEW_FOLLOWER_CREATED,
			Payload: &pb.NotificationEvent_NewFollowerCreated{
				NewFollowerCreated: &pb.NewFollowerCreated{
					TargetUserId:     1001,
					FollowerUserId:   3001,
					FollowerUsername: "charlie",
					Aggregate:        true,
				},
			},
		},
		// Afternoon: Group activity
		{
			EventId:   "afternoon-1",
			EventType: pb.EventType_NEW_EVENT_CREATED,
			Payload: &pb.NotificationEvent_NewEventCreated{
				NewEventCreated: &pb.NewEventCreated{
					UserId:     1001,
					GroupId:    4001,
					EventId:    5001,
					GroupName:  "Photography Club",
					EventTitle: "Sunset Photography Session",
				},
			},
		},
		// Evening: Chat messages
		{
			EventId:   "evening-1",
			EventType: pb.EventType_NEW_MESSAGE_CREATED,
			Payload: &pb.NotificationEvent_NewMessageCreated{
				NewMessageCreated: &pb.NewMessageCreated{
					UserId:           1001,
					SenderUserId:     6001,
					ChatId:           7001,
					SenderUsername:   "diana",
					MessageContent:   "Did you see the new camera gear?",
					Aggregate:        true,
				},
			},
		},
	}

	// Set up expectations
	suite.mockApp.On("CreatePostLikeNotification", mock.Anything, int64(1001), int64(2001), int64(1001), "alice", true).Return(nil)
	suite.mockApp.On("CreatePostLikeNotification", mock.Anything, int64(1001), int64(2002), int64(1001), "bob", true).Return(nil)
	suite.mockApp.On("CreateNewFollowerNotification", mock.Anything, int64(1001), int64(3001), "charlie", true).Return(nil)
	suite.mockApp.On("CreateNewEventNotification", mock.Anything, int64(1001), int64(4001), int64(5001), "Photography Club", "Sunset Photography Session").Return(nil)
	suite.mockApp.On("CreateNewMessageNotification", mock.Anything, int64(1001), int64(6001), int64(7001), "diana", "Did you see the new camera gear?", true).Return(nil)

	// Process the scenario
	for _, event := range scenarioEvents {
		err := suite.eventHandler.Handle(suite.ctx, event)
		assert.NoError(suite.T(), err, "Event %s should be processed successfully", event.EventId)
	}

	suite.mockApp.AssertExpectations(suite.T())
}

// TestErrorRecovery tests how the system handles errors and recovers
func (suite *EndToEndIntegrationTestSuite) TestErrorRecovery() {
	// Set up an event that will cause an error
	suite.mockApp.On("CreatePostCommentNotification", 
		mock.Anything, int64(123), int64(456), int64(123), "user1", "comment body", true).Return(assert.AnError)

	// This event should fail
	failingEvent := &pb.NotificationEvent{
		EventId:   "failing-event",
		EventType: pb.EventType_POST_COMMENT_CREATED,
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

	err := suite.eventHandler.Handle(suite.ctx, failingEvent)
	assert.Error(suite.T(), err)

	// Now test that the system can still process valid events after an error
	suite.mockApp.ExpectedCalls = nil // Clear previous expectations
	suite.mockApp.On("CreatePostLikeNotification", 
		mock.Anything, int64(123), int64(456), int64(123), "user2", true).Return(nil)

	successfulEvent := &pb.NotificationEvent{
		EventId:   "successful-event",
		EventType: pb.EventType_POST_LIKED,
		Payload: &pb.NotificationEvent_PostLiked{
			PostLiked: &pb.PostLiked{
				PostId:        123,
				LikerUserId:   456,
				LikerUsername: "user2",
				Aggregate:     true,
			},
		},
	}

	err = suite.eventHandler.Handle(suite.ctx, successfulEvent)
	assert.NoError(suite.T(), err)

	suite.mockApp.AssertExpectations(suite.T())
}

// TestHighThroughput tests processing many events quickly
func (suite *EndToEndIntegrationTestSuite) TestHighThroughput() {
	const numEvents = 50

	// Create many events of different types
	events := make([]*pb.NotificationEvent, numEvents)
	for i := 0; i < numEvents; i++ {
		eventType := i % 5 // Cycle through 5 different event types
		eventId := "ht-" + string(rune('0'+i%10)) + "-" + string(rune('0'+i/10))

		switch eventType {
		case 0:
			events[i] = &pb.NotificationEvent{
				EventId:   eventId,
				EventType: pb.EventType_POST_LIKED,
				Payload: &pb.NotificationEvent_PostLiked{
					PostLiked: &pb.PostLiked{
						PostId:        int64(100 + i),
						LikerUserId:   int64(200 + i),
						LikerUsername: "user_" + string(rune('a'+i%26)),
						Aggregate:     true,
					},
				},
			}
			suite.mockApp.On("CreatePostLikeNotification", 
				mock.Anything, int64(100+i), int64(200+i), int64(100+i), "user_"+string(rune('a'+i%26)), true).Return(nil)
		case 1:
			events[i] = &pb.NotificationEvent{
				EventId:   eventId,
				EventType: pb.EventType_POST_COMMENT_CREATED,
				Payload: &pb.NotificationEvent_PostCommentCreated{
					PostCommentCreated: &pb.PostCommentCreated{
						PostId:             int64(100 + i),
						CommenterUserId:    int64(200 + i),
						CommenterUsername:  "user_" + string(rune('a'+i%26)),
						Body:               "comment " + string(rune('0'+i%10)),
						Aggregate:          true,
					},
				},
			}
			suite.mockApp.On("CreatePostCommentNotification", 
				mock.Anything, int64(100+i), int64(200+i), int64(100+i), "user_"+string(rune('a'+i%26)), "comment "+string(rune('0'+i%10)), true).Return(nil)
		case 2:
			events[i] = &pb.NotificationEvent{
				EventId:   eventId,
				EventType: pb.EventType_NEW_FOLLOWER_CREATED,
				Payload: &pb.NotificationEvent_NewFollowerCreated{
					NewFollowerCreated: &pb.NewFollowerCreated{
						TargetUserId:     int64(100 + i),
						FollowerUserId:   int64(200 + i),
						FollowerUsername: "user_" + string(rune('a'+i%26)),
						Aggregate:        true,
					},
				},
			}
			suite.mockApp.On("CreateNewFollowerNotification", 
				mock.Anything, int64(100+i), int64(200+i), "user_"+string(rune('a'+i%26)), true).Return(nil)
		case 3:
			events[i] = &pb.NotificationEvent{
				EventId:   eventId,
				EventType: pb.EventType_NEW_MESSAGE_CREATED,
				Payload: &pb.NotificationEvent_NewMessageCreated{
					NewMessageCreated: &pb.NewMessageCreated{
						UserId:           int64(100 + i),
						SenderUserId:     int64(200 + i),
						ChatId:           int64(300 + i),
						SenderUsername:   "user_" + string(rune('a'+i%26)),
						MessageContent:   "message " + string(rune('0'+i%10)),
						Aggregate:        true,
					},
				},
			}
			suite.mockApp.On("CreateNewMessageNotification", 
				mock.Anything, int64(100+i), int64(200+i), int64(300+i), "user_"+string(rune('a'+i%26)), "message "+string(rune('0'+i%10)), true).Return(nil)
		case 4:
			events[i] = &pb.NotificationEvent{
				EventId:   eventId,
				EventType: pb.EventType_MENTION_CREATED,
				Payload: &pb.NotificationEvent_MentionCreated{
					MentionCreated: &pb.MentionCreated{
						MentionedUserId:   int64(100 + i),
						MentionerUserId:   int64(200 + i),
						PostId:            int64(300 + i),
						MentionerUsername: "user_" + string(rune('a'+i%26)),
						PostContent:       "post content " + string(rune('0'+i%10)),
						MentionText:       "@user",
					},
				},
			}
			suite.mockApp.On("CreateMentionNotification", 
				mock.Anything, int64(100+i), int64(200+i), int64(300+i), "user_"+string(rune('a'+i%26)), "post content "+string(rune('0'+i%10)), "@user").Return(nil)
		}
	}

	// Measure processing time
	start := time.Now()
	
	// Process all events
	for _, event := range events {
		err := suite.eventHandler.Handle(suite.ctx, event)
		assert.NoError(suite.T(), err)
	}

	duration := time.Since(start)
	
	// Should process 50 events quickly
	assert.True(suite.T(), duration < 2*time.Second, "Processing 50 events took too long: %v", duration)

	// Verify all expectations were met
	suite.mockApp.AssertExpectations(suite.T())
}

// TestEventRoutingCorrectness tests that events are routed to correct handlers
func (suite *EndToEndIntegrationTestSuite) TestEventRoutingCorrectness() {
	// Test that each event type goes to the correct handler method
	testCases := []struct {
		name     string
		event    *pb.NotificationEvent
		expectedCall string
	}{
		{
			name: "PostLike",
			event: &pb.NotificationEvent{
				EventId:   "route-1",
				EventType: pb.EventType_POST_LIKED,
				Payload: &pb.NotificationEvent_PostLiked{
					PostLiked: &pb.PostLiked{
						PostId:        123,
						LikerUserId:   456,
						LikerUsername: "user",
						Aggregate:     true,
					},
				},
			},
			expectedCall: "CreatePostLikeNotification",
		},
		{
			name: "FollowRequest",
			event: &pb.NotificationEvent{
				EventId:   "route-2",
				EventType: pb.EventType_FOLLOW_REQUEST_CREATED,
				Payload: &pb.NotificationEvent_FollowRequestCreated{
					FollowRequestCreated: &pb.FollowRequestCreated{
						TargetUserId:      123,
						RequesterUserId:   456,
						RequesterUsername: "user",
					},
				},
			},
			expectedCall: "CreateFollowRequestNotification",
		},
		{
			name: "GroupInvite",
			event: &pb.NotificationEvent{
				EventId:   "route-3",
				EventType: pb.EventType_GROUP_INVITE_CREATED,
				Payload: &pb.NotificationEvent_GroupInviteCreated{
					GroupInviteCreated: &pb.GroupInviteCreated{
						InvitedUserId:   123,
						InviterUserId:   456,
						GroupId:         789,
						GroupName:       "group",
						InviterUsername: "user",
					},
				},
			},
			expectedCall: "CreateGroupInviteNotification",
		},
	}

	// Set up expectations for each test case
	suite.mockApp.On("CreatePostLikeNotification", mock.Anything, int64(123), int64(456), int64(123), "user", true).Return(nil)
	suite.mockApp.On("CreateFollowRequestNotification", mock.Anything, int64(123), int64(456), "user").Return(nil)
	suite.mockApp.On("CreateGroupInviteNotification", mock.Anything, int64(123), int64(456), int64(789), "group", "user").Return(nil)

	// Process each test case
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.eventHandler.Handle(suite.ctx, tc.event)
			assert.NoError(t, err)
		})
	}

	// Verify all expectations were met
	suite.mockApp.AssertExpectations(suite.T())
}

// TestUnknownEventTypeHandling tests handling of unknown event types
func (suite *EndToEndIntegrationTestSuite) TestUnknownEventTypeHandling() {
	// Test with nil payload
	unknownEvent1 := &pb.NotificationEvent{
		EventId:   "unknown-1",
		EventType: pb.EventType_EVENT_TYPE_UNSPECIFIED,
		Payload:   nil,
	}

	err := suite.eventHandler.Handle(suite.ctx, unknownEvent1)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "unknown notification event payload type")

	// Test with unsupported payload type (though this shouldn't happen in practice)
	// The system should gracefully handle unknown event types
}

// Run the end-to-end integration test suite
func TestEndToEndIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(EndToEndIntegrationTestSuite))
}
package application

import (
	"context"
	"fmt"
)

// CreateFollowRequestNotification creates a notification when a user sends a follow request to a private account
func (a *Application) CreateFollowRequestNotification(ctx context.Context, targetUserID, requesterUserID int64, requesterUsername string) error {
	title := "New Follow Request"
	message := fmt.Sprintf("%s wants to follow you", requesterUsername)
	
	payload := map[string]string{
		"requester_id":   fmt.Sprintf("%d", requesterUserID),
		"requester_name": requesterUsername,
	}

	_, err := a.CreateNotification(
		ctx,
		targetUserID,           // recipient
		FollowRequest,          // type
		title,                  // title
		message,                // message
		"users",                // source service
		requesterUserID,        // source entity ID
		true,                   // needs action
		payload,                // payload
	)
	if err != nil {
		return fmt.Errorf("failed to create follow request notification: %w", err)
	}

	return nil
}

// CreateNewFollowerNotification creates a notification when someone follows a user
func (a *Application) CreateNewFollowerNotification(ctx context.Context, targetUserID, followerUserID int64, followerUsername string) error {
	title := "New Follower"
	message := fmt.Sprintf("%s is now following you", followerUsername)
	
	payload := map[string]string{
		"follower_id":   fmt.Sprintf("%d", followerUserID),
		"follower_name": followerUsername,
	}

	_, err := a.CreateNotification(
		ctx,
		targetUserID,           // recipient
		NewFollower,            // type
		title,                  // title
		message,                // message
		"users",                // source service
		followerUserID,         // source entity ID
		false,                  // doesn't need action
		payload,                // payload
	)
	if err != nil {
		return fmt.Errorf("failed to create new follower notification: %w", err)
	}

	return nil
}

// CreateGroupInviteNotification creates a notification when a user is invited to join a group
func (a *Application) CreateGroupInviteNotification(ctx context.Context, invitedUserID, inviterUserID, groupID int64, groupName, inviterUsername string) error {
	title := "Group Invitation"
	message := fmt.Sprintf("%s invited you to join the group \"%s\"", inviterUsername, groupName)
	
	payload := map[string]string{
		"inviter_id":    fmt.Sprintf("%d", inviterUserID),
		"inviter_name":  inviterUsername,
		"group_id":      fmt.Sprintf("%d", groupID),
		"group_name":    groupName,
		"action":        "accept_or_decline",
	}

	_, err := a.CreateNotification(
		ctx,
		invitedUserID,          // recipient
		GroupInvite,            // type
		title,                  // title
		message,                // message
		"users",                // source service
		groupID,                // source entity ID (the group)
		true,                   // needs action
		payload,                // payload
	)
	if err != nil {
		return fmt.Errorf("failed to create group invite notification: %w", err)
	}

	return nil
}

// CreateGroupJoinRequestNotification creates a notification when someone requests to join a group
func (a *Application) CreateGroupJoinRequestNotification(ctx context.Context, groupOwnerID, requesterID, groupID int64, groupName, requesterUsername string) error {
	title := "New Group Join Request"
	message := fmt.Sprintf("%s wants to join your group \"%s\"", requesterUsername, groupName)
	
	payload := map[string]string{
		"requester_id":  fmt.Sprintf("%d", requesterID),
		"requester_name": requesterUsername,
		"group_id":      fmt.Sprintf("%d", groupID),
		"group_name":    groupName,
		"action":        "accept_or_decline",
	}

	_, err := a.CreateNotification(
		ctx,
		groupOwnerID,           // recipient (group owner)
		GroupJoinRequest,       // type
		title,                  // title
		message,                // message
		"users",                // source service
		groupID,                // source entity ID (the group)
		true,                   // needs action
		payload,                // payload
	)
	if err != nil {
		return fmt.Errorf("failed to create group join request notification: %w", err)
	}

	return nil
}

// CreateNewEventNotification creates a notification when a new event is created in a group the user is part of
func (a *Application) CreateNewEventNotification(ctx context.Context, userID, groupID, eventID int64, groupName, eventTitle string) error {
	title := "New Event in Group"
	message := fmt.Sprintf("New event \"%s\" was created in group \"%s\"", eventTitle, groupName)
	
	payload := map[string]string{
		"group_id":      fmt.Sprintf("%d", groupID),
		"group_name":    groupName,
		"event_id":      fmt.Sprintf("%d", eventID),
		"event_title":   eventTitle,
		"action":        "view_event",
	}

	_, err := a.CreateNotification(
		ctx,
		userID,                 // recipient
		NewEvent,               // type
		title,                  // title
		message,                // message
		"posts",                // source service
		eventID,                // source entity ID (the event)
		false,                  // doesn't need action (just informational)
		payload,                // payload
	)
	if err != nil {
		return fmt.Errorf("failed to create new event notification: %w", err)
	}

	return nil
}

// Additional notification types for extended functionality

// CreatePostLikeNotification creates a notification when someone likes a user's post
func (a *Application) CreatePostLikeNotification(ctx context.Context, userID, likerID, postID int64, likerUsername string) error {
	title := "Post Liked"
	message := fmt.Sprintf("%s liked your post", likerUsername)
	
	payload := map[string]string{
		"liker_id":     fmt.Sprintf("%d", likerID),
		"liker_name":   likerUsername,
		"post_id":      fmt.Sprintf("%d", postID),
		"action":       "view_post",
	}

	_, err := a.CreateNotification(
		ctx,
		userID,                 // recipient
		PostLike,               // type
		title,                  // title
		message,                // message
		"posts",                // source service
		postID,                 // source entity ID
		false,                  // doesn't need action
		payload,                // payload
	)
	if err != nil {
		return fmt.Errorf("failed to create post like notification: %w", err)
	}

	return nil
}

// CreatePostCommentNotification creates a notification when someone comments on a user's post
func (a *Application) CreatePostCommentNotification(ctx context.Context, userID, commenterID, postID int64, commenterUsername, commentContent string) error {
	title := "New Comment"
	message := fmt.Sprintf("%s commented on your post", commenterUsername)
	
	payload := map[string]string{
		"commenter_id":    fmt.Sprintf("%d", commenterID),
		"commenter_name":  commenterUsername,
		"post_id":         fmt.Sprintf("%d", postID),
		"comment_content": commentContent,
		"action":          "view_post",
	}

	_, err := a.CreateNotification(
		ctx,
		userID,                 // recipient
		PostComment,            // type
		title,                  // title
		message,                // message
		"posts",                // source service
		postID,                 // source entity ID
		false,                  // doesn't need action
		payload,                // payload
	)
	if err != nil {
		return fmt.Errorf("failed to create post comment notification: %w", err)
	}

	return nil
}

// CreateMentionNotification creates a notification when a user is mentioned in a post or comment
func (a *Application) CreateMentionNotification(ctx context.Context, userID, mentionerID, postID int64, mentionerUsername, postContent string) error {
	title := "You were mentioned"
	message := fmt.Sprintf("%s mentioned you in a post", mentionerUsername)
	
	payload := map[string]string{
		"mentioner_id":    fmt.Sprintf("%d", mentionerID),
		"mentioner_name":  mentionerUsername,
		"post_id":         fmt.Sprintf("%d", postID),
		"post_content":    postContent,
		"action":          "view_post",
	}

	_, err := a.CreateNotification(
		ctx,
		userID,                 // recipient
		Mention,                // type
		title,                  // title
		message,                // message
		"posts",                // source service
		postID,                 // source entity ID
		false,                  // doesn't need action
		payload,                // payload
	)
	if err != nil {
		return fmt.Errorf("failed to create mention notification: %w", err)
	}

	return nil
}
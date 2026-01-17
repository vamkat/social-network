package application

import (
	"context"
	"fmt"
	ds "social-network/services/posts/internal/db/dbservice"
	notifpb "social-network/shared/gen-go/notifications"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"
)

func (s *Application) ToggleOrInsertReaction(ctx context.Context, req models.GenericReq) error {
	input := fmt.Sprintf("%#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		return ce.Wrap(ce.ErrInvalidArgument, err, input).WithPublic("invalid data received")

	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return ce.Wrap(ce.ErrInternal, err, fmt.Sprintf("%#v", accessCtx)).WithPublic(genericPublic)
	}
	if !hasAccess {
		return ce.New(ce.ErrPermissionDenied, fmt.Errorf("user has no permission to react to entity %v", req.EntityId), input).WithPublic("permission denied")
	}

	action, err := s.db.ToggleOrInsertReaction(ctx, ds.ToggleOrInsertReactionParams{
		ContentID: req.EntityId.Int64(),
		UserID:    req.RequesterId.Int64(),
	})
	if err != nil {
		return ce.New(ce.ErrInternal, err, input).WithPublic(genericPublic)
	}
	if action == "added" {
		//create notification
		liker, err := s.userRetriever.GetUser(ctx, ct.Id(req.RequesterId))
		if err != nil {
			//log error
		}

		row, err := s.db.GetEntityCreatorAndGroup(ctx, req.EntityId.Int64())
		if err != nil {
			//log and don't proceed to notif
		}

		//if liker is same as entity creator, return without creating a notification
		if row.CreatorID == int64(liker.UserId) {
			return nil
		}

		// build the notification event
		event := &notifpb.NotificationEvent{
			EventType: notifpb.EventType_POST_LIKED,
			Payload: &notifpb.NotificationEvent_PostLiked{
				PostLiked: &notifpb.PostLiked{
					EntityCreatorId: row.CreatorID,
					PostId:          req.EntityId.Int64(),
					LikerUserId:     req.RequesterId.Int64(),
					LikerUsername:   liker.Username.String(),
					Aggregate:       true,
				},
			},
		}

		if err := s.eventProducer.CreateAndSendNotificationEvent(ctx, event); err != nil {
			tele.Error(ctx, "failed to send new reaction notification: @1", "error", err.Error())
		}
		tele.Info(ctx, "new reaction notification event created")

	}
	return nil
}

// SKIP FOR NOW
func (s *Application) GetWhoLikedEntityId(ctx context.Context, req models.GenericReq) ([]int64, error) {
	return nil, nil
}

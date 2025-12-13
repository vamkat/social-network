package application

import (
	"context"
	"fmt"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Application) CreateEvent(ctx context.Context, req models.CreateEventReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	isMember, err := s.clients.IsGroupMember(ctx, req.CreatorId.Int64(), req.GroupId.Int64())
	if err != nil {
		return err
	}
	if !isMember {
		return ErrNotAllowed
	}

	// convert date
	eventDate := pgtype.Date{
		Time:  req.EventDate.Time(),
		Valid: true,
	}

	err = s.db.CreateEvent(ctx, sqlc.CreateEventParams{
		EventTitle:     req.Title.String(),
		EventBody:      req.Body.String(),
		EventCreatorID: req.CreatorId.Int64(),
		GroupID:        req.GroupId.Int64(),
		EventDate:      eventDate,
	})
	if err != nil {
		return err
	}

	//TODO CREATE NOTIFICATION EVENT

	return nil
}

func (s *Application) DeleteEvent(ctx context.Context, req models.GenericReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrNotAllowed
	}

	rowsAffected, err := s.db.DeleteEvent(ctx, sqlc.DeleteEventParams{
		ID:             req.EntityId.Int64(),
		EventCreatorID: req.RequesterId.Int64(),
	})
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrNotFound
	}

	return nil
}

func (s *Application) EditEvent(ctx context.Context, req models.EditEventReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EventId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrNotAllowed
	}

	return s.txRunner.RunTx(ctx, func(q sqlc.Querier) error {
		// convert date
		eventDate := pgtype.Date{
			Time:  req.EventDate.Time(),
			Valid: true,
		}
		rowsAffected, err := q.EditEvent(ctx, sqlc.EditEventParams{
			EventTitle:     req.Title.String(),
			EventBody:      req.Body.String(),
			EventDate:      eventDate,
			ID:             req.EventId.Int64(),
			EventCreatorID: req.RequesterId.Int64(),
		})
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			return ErrNotFound
		}
		if req.Image > 0 {
			err := q.UpsertImage(ctx, sqlc.UpsertImageParams{
				ID:       req.Image.Int64(),
				ParentID: req.EventId.Int64(),
			})
			if err != nil {
				return err
			}
		} else {
			rowsAffected, err := q.DeleteImage(ctx, req.Image.Int64())
			if err != nil {
				return err
			}
			if rowsAffected != 1 {
				fmt.Println("image not found")
			}
		}
		return nil
	})

}

func (s *Application) GetEventsByGroupId(ctx context.Context, req models.EntityIdPaginatedReq) ([]models.Event, error) {

	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, ErrNotAllowed
	}
	rows, err := s.db.GetEventsByGroupId(ctx, sqlc.GetEventsByGroupIdParams{
		GroupID: req.EntityId.Int64(),
		Offset:  req.Offset.Int32(),
		Limit:   req.Limit.Int32(),
		UserID:  req.RequesterId.Int64(),
	})
	if err != nil {
		return nil, nil
	}
	events := make([]models.Event, 0, len(rows))
	userIDs := make([]int64, 0, len(rows))

	for _, r := range rows {
		uid := r.EventCreatorID
		userIDs = append(userIDs, uid)

		events = append(events, models.Event{
			EventId: ct.Id(r.ID),
			Title:   ct.Title(r.EventTitle),
			Body:    ct.EventBody(r.EventBody),
			User: models.User{
				UserId: ct.Id(uid),
			},
			GroupId:       ct.Id(r.GroupID),
			EventDate:     ct.EventDateTime(r.EventDate.Time),
			GoingCount:    int(r.GoingCount),
			NotGoingCount: int(r.NotGoingCount),
			Image:         ct.Id(r.Image),
			CreatedAt:     ct.GenDateTime(r.CreatedAt.Time),
			UpdatedAt:     ct.GenDateTime(r.UpdatedAt.Time),
			UserResponse:  &r.UserResponse.Bool,
		})
	}

	if len(events) == 0 {
		return events, nil
	}

	userMap, err := s.hydrator.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	for i := range events {
		uid := events[i].User.UserId.Int64()
		if u, ok := userMap[uid]; ok {
			events[i].User = u
		}
	}

	// if err := s.hydrateEvents(ctx, events); err != nil {
	// 	return nil, err
	// }

	return events, nil
}

func (s *Application) RespondToEvent(ctx context.Context, req models.RespondToEventReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	accessCtx := accessContext{
		requesterId: req.ResponderId.Int64(),
		entityId:    req.EventId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrNotAllowed
	}

	rowsAffected, err := s.db.UpsertEventResponse(ctx, sqlc.UpsertEventResponseParams{
		EventID: req.EventId.Int64(),
		UserID:  req.ResponderId.Int64(),
		Going:   req.Going,
	})
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrNotFound
	}
	return nil
}

func (s *Application) RemoveEventResponse(ctx context.Context, req models.GenericReq) error {

	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrNotAllowed
	}

	rowsAffected, err := s.db.DeleteEventResponse(ctx, sqlc.DeleteEventResponseParams{
		EventID: req.EntityId.Int64(),
		UserID:  req.RequesterId.Int64(),
	})
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrNotFound
	}
	return nil
}

package application

import (
	"context"
	"fmt"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Application) CreateEvent(ctx context.Context, req CreateEventReq) error {
	//check creator is member of group (HANDLER)
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	// convert date
	eventDate := pgtype.Date{
		Time:  req.EventDate.Time(),
		Valid: true,
	}

	err := s.db.CreateEvent(ctx, sqlc.CreateEventParams{
		EventTitle:     req.Title.String(),
		EventBody:      req.Body.String(),
		EventCreatorID: req.CreatorId.Int64(),
		GroupID:        req.GroupId.Int64(),
		EventDate:      eventDate,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Application) DeleteEvent(ctx context.Context, req GenericReq) error {
	//check requester is member of the group?(HANDLER)
	if err := ct.ValidateStruct(req); err != nil {
		return err
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

func (s *Application) EditEvent(ctx context.Context, req EditEventReq) error {
	//check requester is creator of event (and member of the group? what happens if they're not any more?)
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}
	err := s.runTx(ctx, func(q *sqlc.Queries) error {
		// convert date
		eventDate := pgtype.Date{
			Time:  req.EventDate.Time(),
			Valid: true,
		}
		rowsAffected, err := s.db.EditEvent(ctx, sqlc.EditEventParams{
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
			err := s.db.UpsertImage(ctx, sqlc.UpsertImageParams{
				ID:       req.Image.Int64(),
				ParentID: req.EventId.Int64(),
			})
			if err != nil {
				return err
			}
		} else {
			rowsAffected, err := s.db.DeleteImage(ctx, req.Image.Int64())
			if err != nil {
				return err
			}
			if rowsAffected != 1 {
				fmt.Println("image not found")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Application) GetEventsByGroupId(ctx context.Context, req GenericPaginatedReq) ([]Event, error) {
	//check requester is member of group (HANDLER)
	if err := ct.ValidateStruct(req); err != nil {
		return nil, err
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
	events := make([]Event, 0, len(rows))
	for _, r := range rows {

		events = append(events, Event{
			EventId:       ct.Id(r.ID),
			Title:         ct.Title(r.EventTitle),
			Body:          ct.EventBody(r.EventBody),
			CreatorId:     ct.Id(r.EventCreatorID),
			GroupId:       ct.Id(r.GroupID),
			EventDate:     ct.EventDate(r.EventDate.Time),
			GoingCount:    int(r.GoingCount),
			NotGoingCount: int(r.NotGoingCount),
			Image:         ct.Id(r.Image),
			CreatedAt:     r.CreatedAt.Time,
			UpdatedAt:     r.UpdatedAt.Time,
			UserResponse:  &r.UserResponse.Bool,
		})
	}

	return events, nil
}

func (s *Application) RespondToEvent(ctx context.Context, req RespondToEventReq) error {
	//check requester is member of group (HANDLER)
	if err := ct.ValidateStruct(req); err != nil {
		return err
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

func (s *Application) RemoveEventResponse(ctx context.Context, req GenericReq) error {
	//check requester is member of group (HANDLER)
	if err := ct.ValidateStruct(req); err != nil {
		return err
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

package application

import (
	"context"
	"fmt"
	"social-network/services/posts/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
)

// runTx runs a function inside a database transaction.
// If fn returns an error, the tx is rolled back.
func (s *Application) runTx(ctx context.Context, fn func(q *sqlc.Queries) error) error {
	// start tx
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// db must be a *sqlc.Queries to use WithTx
	base, ok := s.db.(*sqlc.Queries)
	if !ok {
		return fmt.Errorf("UserService.db must be *sqlc.Queries for transactions")
	}

	qtx := base.WithTx(tx)

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// group and post audience=group: only members can see
// post audience=everyone: everyone can see (can we check this before all the fetches from users?)
// post audience=followers: requester can see if they follow creator
// post audience=selected: requester can see if they are in post audience table
func (s *Application) hasRightToView(ctx context.Context, req AccessContext) (bool, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return false, err
	}
	canSee, err := s.db.CanUserSeeEntity(ctx, sqlc.CanUserSeeEntityParams{
		UserID:       req.RequesterId.Int64(),
		FollowingIds: req.RequesterFollowsIds.Int64(),
		GroupIds:     req.RequesterGroups.Int64(),
		EntityID:     req.ParentEntityId.Int64(),
	})
	if err != nil {
		return false, err
	}
	return canSee, nil
}

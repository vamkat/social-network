package application

import (
	"context"
	"fmt"
	"social-network/services/posts/internal/db/sqlc"
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

func (s *Application) hasRightToView(ctx context.Context, req hasRightToView) (bool, error) {
	// get the requester id, the parent entity id (so post or event even if the request is for comments)
	// user ids the requester follows and group ids the requester belongs to
	// group and post audience=group: only members can see
	// post audience=everyone: everyone can see (can we check this before all the fetches from users?)
	// post audience=followers: requester can see if they follow creator
	// post audience=selected: requester can see if they are in post audience table
	return false, nil
}

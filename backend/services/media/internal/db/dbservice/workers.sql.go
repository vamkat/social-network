package dbservice

import (
	"context"
	"fmt"
	tele "social-network/shared/go/telemetry"
	"time"
)

func (w *Workers) StartStaleFilesWorker(ctx context.Context, period time.Duration) {
	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := w.db.MarkStaleFilesFailed(ctx); err != nil {
					tele.Error(ctx, fmt.Sprintf("Error processing pending variants: %v", err))
				}
			case <-ctx.Done():
				tele.Warn(ctx, "Stale files worker stopped")
				return
			}
		}
	}()
}

func (q *Queries) MarkStaleFilesFailed(ctx context.Context) error {
	tag, err := q.db.Exec(ctx, `
		UPDATE files
		SET status = 'failed',
		    updated_at = now()
		WHERE status = 'pending'
		  AND created_at < now() - interval '24 hours'
	`)
	if err != nil {
		return err
	}
	tele.Warn(ctx, "Number of files marked as failed, tag: "+fmt.Sprint(tag), "tag", tag)
	return nil
}

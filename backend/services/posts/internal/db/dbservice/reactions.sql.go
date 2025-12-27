package dbservice

import (
	"context"
)

const getWhoLikedEntityId = `-- name: GetWhoLikedEntityId :many
SELECT user_id
FROM reactions
WHERE content_id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetWhoLikedEntityId(ctx context.Context, contentID int64) ([]int64, error) {
	rows, err := q.db.Query(ctx, getWhoLikedEntityId, contentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var user_id int64
		if err := rows.Scan(&user_id); err != nil {
			return nil, err
		}
		items = append(items, user_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const toggleOrInsertReaction = `-- name: ToggleReaction :execrows
WITH deleted AS (
    DELETE FROM reactions
    WHERE content_id = $1
      AND user_id = $2
    RETURNING 1
)
INSERT INTO reactions (content_id, user_id)
SELECT $1, $2
WHERE NOT EXISTS (SELECT 1 FROM deleted);
`

type ToggleOrInsertReactionParams struct {
	ContentID int64
	UserID    int64
}

func (q *Queries) ToggleOrInsertReaction(ctx context.Context, arg ToggleOrInsertReactionParams) (int64, error) {
	result, err := q.db.Exec(ctx, toggleOrInsertReaction, arg.ContentID, arg.UserID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

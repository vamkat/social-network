-- name: ToggleReactionIfExists :one
UPDATE reactions
SET deleted_at = CASE
                     WHEN deleted_at IS NULL THEN NOW()   -- make inactive
                     ELSE NULL                           -- restore
                 END,
    updated_at = NOW()
WHERE content_id = $1
  AND user_id = $2
RETURNING id, content_id, user_id, deleted_at, updated_at;

-- name: InsertReaction :one
INSERT INTO reactions (content_id, user_id)
VALUES ($1, $2)
RETURNING id, content_id, user_id, deleted_at, updated_at;

-- name: GetWhoLikedEntityId :many
SELECT user_id
FROM reactions
WHERE content_id = $1 AND deleted_at IS NULL;
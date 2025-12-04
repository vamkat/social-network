-- name: ToggleOrInsertReaction :execrows
INSERT INTO reactions (content_id, user_id)
VALUES ($1, $2)
ON CONFLICT (content_id, user_id) DO UPDATE
SET deleted_at = CASE
                     WHEN reactions.deleted_at IS NULL THEN NOW()
                     ELSE NULL
                 END;

-- name: GetWhoLikedEntityId :many
SELECT user_id
FROM reactions
WHERE content_id = $1 AND deleted_at IS NULL;
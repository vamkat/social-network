-- name: GetMessages :many
SELECT m.*
FROM messages m
JOIN conversation_members cm 
  ON cm.conversation_id = m.conversation_id
WHERE m.conversation_id = $1
  AND cm.user_id = $2
  AND m.deleted_at IS NULL
ORDER BY m.created_at ASC
LIMIT $3 OFFSET $4;


-- name: CreateMessage :one
INSERT INTO messages (conversation_id, sender_id, message_text)
SELECT $1, $2, $3
FROM conversation_members
WHERE conversation_id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING *;


-- name: UpdateLastReadMessage :one
UPDATE conversation_members cm
SET last_read_message_id = $3
WHERE cm.conversation_id = $1
  AND cm.user_id = $2
  AND cm.deleted_at IS NULL
RETURNING *;







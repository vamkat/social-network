-- name: CreatePrivateConv :one
-- Creates a Conversation if and only if a DM between the same 2 users does not exist.
-- Returns NULL if a duplicate DM exists (sqlc will error if RETURNING finds no rows).
WITH existing AS (
    SELECT c.id
    FROM conversations c
    JOIN conversation_members cm1 ON cm1.conversation_id = c.id AND cm1.user_id = $1
    JOIN conversation_members cm2 ON cm2.conversation_id = c.id AND cm2.user_id = $2
    WHERE c.group_id IS NULL
)
INSERT INTO conversations (group_id)
SELECT NULL
WHERE NOT EXISTS (SELECT 1 FROM existing)
RETURNING id;


-- name: CreateGroupConv :one
INSERT INTO conversations (group_id)
VALUES ($1)
RETURNING id;

-- name: GetUserConversations :many
SELECT 
    c.id AS conversation_id,
    c.group_id,
    c.created_at,
    c.updated_at,
    cm2.user_id AS member_id
FROM conversations c
JOIN conversation_members cm1
    ON cm1.conversation_id = c.id
    AND cm1.user_id = $1
    AND cm1.deleted_at IS NULL
JOIN conversation_members cm2
    ON cm2.conversation_id = c.id
    AND cm2.user_id <> $1
    AND cm2.deleted_at IS NULL
WHERE c.deleted_at IS NULL
AND (
    ($2 IS NULL AND c.group_id IS NULL)
    OR ($2 IS NOT NULL AND c.group_id = $2)
)
ORDER BY c.updated_at DESC
LIMIT $3 OFFSET $4;


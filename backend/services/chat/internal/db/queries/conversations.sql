-- name: CreatePrivateConv :one
-- Creates a Conversation if and only if a conversation between the same 2 users does not exist.
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
-- Get all conversations paginated by user id excluding group conversations.
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

-- name: DeleteConversationByExactMembers :one
-- Delete a conversation only if its members exactly match the provided list.
-- Returns 0 rows if conversation doesn't exist, members don’t match exactly, conversation has extra or missing members.
WITH target_members AS (
    SELECT unnest(@member_ids::bigint[]) AS user_id
),
matched_conversation AS (
    SELECT cm.conversation_id
    FROM conversation_members cm
    JOIN target_members tm ON tm.user_id = cm.user_id
    GROUP BY cm.conversation_id
    HAVING 
        -- same count of overlapping members
        COUNT(*) = (SELECT COUNT(*) FROM target_members)
        -- and the conversation has no extra members
        AND COUNT(*) = (
            SELECT COUNT(*) 
            FROM conversation_members cm2 
            WHERE cm2.conversation_id = cm.conversation_id
              AND cm2.deleted_at IS NULL
        )
)
UPDATE conversations c
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE c.id = (SELECT conversation_id FROM matched_conversation)
RETURNING *;



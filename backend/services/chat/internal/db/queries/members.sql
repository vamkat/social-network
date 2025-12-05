-- name: AddConversationMembers :exec
INSERT INTO conversation_members (conversation_id, user_id, last_read_message_id)
SELECT sqlc.arg(conversation_id), UNNEST(sqlc.arg(user_ids)::bigint[]), NULL;

-- name: GetConversationMembers :many
SELECT cm2.user_id
FROM conversation_members cm1
JOIN conversation_members cm2
  ON cm2.conversation_id = cm1.conversation_id
WHERE cm1.user_id = $2
  AND cm2.conversation_id = $1
  AND cm2.user_id <> $2
  AND cm2.deleted_at IS NULL;

-- name: SoftDeleteConversationMember :one
UPDATE conversation_members cm_target
SET deleted_at = NOW()
FROM conversation_members cm_actor
WHERE cm_target.conversation_id = $1
  AND cm_target.user_id = $2
  AND cm_target.deleted_at IS NULL
  AND cm_actor.conversation_id = $1
  AND cm_actor.user_id = $3
  AND cm_actor.deleted_at IS NULL
RETURNING cm_target.*;
-- name: AddMembersToGroupConversation :one
-- Find a conversation by group_id and insert the given user_ids into conversation_members.
-- existing members are ignored, new members are added.
-- Returns:
--   BIGINT          -- the conversation id
WITH convo AS (
    SELECT id
    FROM conversations
    WHERE group_id = @group_id
      AND deleted_at IS NULL
),
insert_members AS (
    INSERT INTO conversation_members (conversation_id, user_id)
    SELECT (SELECT id FROM convo), unnest(@user_ids::bigint[])
    ON CONFLICT (conversation_id, user_id) DO NOTHING
    RETURNING conversation_id
)
SELECT id FROM convo;

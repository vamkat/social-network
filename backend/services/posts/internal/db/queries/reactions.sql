-- user reacted


-- add like
INSERT INTO reactions (content_id, user_id)
VALUES ($1, $2)
ON CONFLICT (user_id, content_id) DO NOTHING
RETURNING *;

-- remove like
DELETE FROM reactions
WHERE content_id = $1 AND user_id = $2
RETURNING *;

-- get who liked entity id
SELECT user_id
FROM reactions
WHERE content_id = $1 AND deleted_at IS NULL;
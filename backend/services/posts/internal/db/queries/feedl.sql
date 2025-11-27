-- get public feed (paginated)
SELECT *
FROM posts
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- get user feed (paginated)
SELECT *
FROM feed_entries
WHERE user_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- add to user feed
INSERT INTO feed_entries (user_id, post_id)
VALUES ($1, $2)
ON CONFLICT (user_id, post_id) DO NOTHING;


-- mark post in feed as seen
UPDATE feed_entries
SET seen = TRUE
WHERE user_id = $1 AND post_id = $2
RETURNING *;

-- delete from user feed
UPDATE feed_entries
SET deleted_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND post_id = $2 AND deleted_at IS NULL
RETURNING *;
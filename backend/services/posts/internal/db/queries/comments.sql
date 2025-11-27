--get comments by parent id
SELECT *
FROM comments
WHERE parent_post_id = $1
  AND deleted_at IS NULL
ORDER BY created_at ASC;

-- create comment
INSERT INTO comments (comment_creator_id, parent_post_id, comment_body)
VALUES ($1, $2, $3)
RETURNING *;

-- edit comment
UPDATE comments
SET comment_body = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

-- delete comment
UPDATE comments
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- get last comment for post id

-- paginated comments by date newest first
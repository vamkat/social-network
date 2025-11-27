-- TODO always return whether user has reacted

--get post by id
SELECT p.*, mi.content_type
FROM posts p
JOIN master_index mi ON mi.id = p.id
WHERE p.id = $1
  AND p.deleted_at IS NULL
  AND mi.deleted_at IS NULL;

--get group post by id?

-- get group posts paginated

-- get user posts paginated

-- get most popular post for group (most liked and most comments)

-- create (group) post
INSERT INTO posts (post_title, post_body, creator_id, group_id, audience)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- edit post 
UPDATE posts
SET post_title = $1,
    post_body  = $2,
    audience   = $3
WHERE id = $4 AND deleted_at IS NULL
RETURNING *;


-- delete post
UPDATE posts
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- change post privacy

-- add user to post audience
INSERT INTO post_audience (post_id, allowed_user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- get post audience
SELECT allowed_user_id
FROM post_audience
WHERE post_id = $1;

-- remove user from audience
DELETE FROM post_audience
WHERE post_id = $1 AND allowed_user_id = $2;

-- clear post audience
DELETE FROM post_audience WHERE post_id = $1;


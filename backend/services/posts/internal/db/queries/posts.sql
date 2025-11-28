-- TODO always return whether user has reacted

---------------------------------------------------------------
-- name: GetPostById :one
------------------------------------------------------------------
SELECT
    p.id,
    p.post_title,
    p.post_body,
    p.creator_id,
    p.group_id,
    p.audience,
    p.comments_count,
    p.reactions_count,
    p.images_count,
    p.last_commented_at,
    p.created_at,
    p.updated_at,

    EXISTS (
        SELECT 1
        FROM reactions r
        WHERE r.content_id = p.id
          AND r.user_id = $2
          AND r.deleted_at IS NULL
    ) AS liked_by_user,


    (SELECT file_name     -- preview = first image by sort_order
     FROM images i
     WHERE i.entity_id = p.id
       AND i.deleted_at IS NULL
     ORDER BY i.sort_order
     LIMIT 1
    ) AS preview_image

FROM posts p
WHERE p.id = $1
  AND p.deleted_at IS NULL

  -- VISIBILITY CHECK
  AND (
        p.creator_id = $2
        OR p.audience = 'everyone' --followers must be checked in users service
        OR (
            p.audience = 'selected'
            AND EXISTS (
                SELECT 1 FROM post_audience pa
                WHERE pa.post_id = p.id
                  AND pa.allowed_user_id = $2
            )
        )
      );


-----------------------------------------------------------------------------
-- name: GetGroupPostsPaginated :many
-------------------------------------------------------------------------------
SELECT
    p.id,
    p.post_title,
    p.post_body,
    p.creator_id,
    p.group_id,
    p.audience,
    p.comments_count,
    p.reactions_count,
    p.images_count,
    p.last_commented_at,
    p.created_at,
    p.updated_at,

    EXISTS (     -- Has the given user liked the post?
        SELECT 1
        FROM reactions r
        WHERE r.content_id = p.id
          AND r.user_id = $2              -- requesting user (check is member from users service)
          AND r.deleted_at IS NULL
    ) AS liked_by_user,
   
    (SELECT file_name     -- preview = first image by sort_order
     FROM images i
     WHERE i.entity_id = p.id
       AND i.deleted_at IS NULL
     ORDER BY i.sort_order
     LIMIT 1
    ) AS preview_image

FROM posts p
LEFT JOIN images i ON i.entity_id = p.id AND i.deleted_at IS NULL
WHERE p.group_id = $1                    -- group id filter
  AND p.deleted_at IS NULL
GROUP BY p.id
ORDER BY p.created_at DESC               -- newest first
LIMIT $3 OFFSET $4;                      -- pagination

------------------------------------------------------------------------------
-- name: GetUserPostsPaginated :many
------------------------------------------------------------------------------
SELECT
    p.id,
    p.post_title,
    p.post_body,
    p.creator_id,
    p.comments_count,
    p.reactions_count,
    p.images_count,
    p.last_commented_at,
    p.created_at,
    p.updated_at,

    EXISTS (    -- Has the requesting user liked the post?
        SELECT 1 FROM reactions r
        WHERE r.content_id = p.id
          AND r.user_id = $2
          AND r.deleted_at IS NULL
    ) AS liked_by_user,

    (SELECT file_name     -- preview = first image by sort_order
     FROM images i
     WHERE i.entity_id = p.id
       AND i.deleted_at IS NULL
     ORDER BY i.sort_order
     LIMIT 1
    ) AS preview_image

FROM posts p
LEFT JOIN images i ON i.entity_id = p.id AND i.deleted_at IS NULL

WHERE p.creator_id = $1                      -- target user we are viewing
  AND p.group_id IS NULL                     -- exclude group posts
  AND p.deleted_at IS NULL

  AND (                    
        p.creator_id = $2    -- If viewer *is* the creator — show all posts                
        OR p.audience = 'everyone' -- followers must be checked in users service
        OR (
            p.audience = 'selected'            -- must be specifically allowed
            AND EXISTS (
                SELECT 1
                FROM post_audience pa
                WHERE pa.post_id = p.id
                  AND pa.allowed_user_id = $2
            )
        )
     )

GROUP BY p.id
ORDER BY p.created_at DESC
LIMIT $3 OFFSET $4;

-----------------------------------------------------------------------
-- name: GetMostPopularPostInGroup :one
-------------------------------------------------------------------------
SELECT
    p.id,
    p.post_title,
    p.post_body,
    p.creator_id,
    p.group_id,
    p.audience,
    p.comments_count,
    p.reactions_count,
    p.images_count,
    p.last_commented_at,
    p.created_at,
    p.updated_at,

    (SELECT file_name     -- preview image (first by sort_order)
     FROM images i
     WHERE i.entity_id = p.id
       AND i.deleted_at IS NULL
     ORDER BY i.sort_order
     LIMIT 1
    ) AS preview_image,


    (p.reactions_count + p.comments_count) AS popularity_score     -- popularity metric (likes + comments)

FROM posts p
WHERE p.group_id = $1
  AND p.deleted_at IS NULL

ORDER BY popularity_score DESC, p.created_at DESC
LIMIT 1;

-----------------------------------------------------
-- create post
-----------------------------------------------------
BEGIN;

-- Insert post
WITH new_post AS (
  INSERT INTO posts (post_title, post_body, creator_id, group_id, audience)
  VALUES ($1, $2, $3, $4, $5)   -- $5 = audience ('everyone', 'selected', etc.)
  RETURNING id
)

-- Insert images
INSERT INTO images (file_name, entity_id, sort_order)
SELECT file_name, id, sort_order
FROM unnest($6_file_names::text[], $6_sort_orders::int[]) AS t(file_name, sort_order),
     new_post;

-- Insert allowed users if audience = 'selected'
INSERT INTO post_audience (post_id, allowed_user_id)
SELECT id, allowed_user_id
FROM unnest($7::bigint[]) AS allowed_user_id,
     new_post
WHERE $5 = 'selected';   -- only insert if audience is 'selected'

COMMIT;

------------------------------------------------------------
-- edit post content
------------------------------------------------------------
UPDATE posts
SET post_body  = $1,
WHERE id = $3 AND deleted_at IS NULL
RETURNING *;

--------------------------------------
-- edit post's intended audience
--------------------------------------
BEGIN;

-- Update audience type
UPDATE posts
SET audience = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- If audience is now 'selected', update allowed users list
-- First, remove previous list
DELETE FROM post_audience WHERE post_id = $1;

-- Then add the new allowed users ONLY if selected
INSERT INTO post_audience (post_id, allowed_user_id)
SELECT $1, unnest_user_id
FROM unnest($3::bigint[]) AS unnest_user_id
WHERE $2 = 'selected';   -- only insert if new audience is selected

COMMIT;

------------------------------------------
-- edit selected audience for post id
--------------------------------------------
BEGIN;

-- Remove current audience list
DELETE FROM post_audience
WHERE post_id = $1;

-- Insert new audience list (full array provided by frontend/backend)
INSERT INTO post_audience (post_id, allowed_user_id)
SELECT $1, unnest_user_id
FROM unnest($2::bigint[]) AS unnest_user_id;

COMMIT;

--------------------------------------------------------
-- delete post
-------------------------------------------------------------
UPDATE posts
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-------------------------------------------------------
-- name: GetPostAudience :many
-----------------------------------------------------------
SELECT allowed_user_id
FROM post_audience
WHERE post_id = $1
ORDER BY allowed_user_id;






---------------------------------------------------
-- name: GetCommentsByPostId :many
---------------------------------------------------
SELECT
    c.id,
    c.comment_creator_id,
    c.comment_body,
    c.reactions_count,
    c.images_count,
    c.created_at,
    c.updated_at,

    -- Whether requesting user liked this comment
    EXISTS (
        SELECT 1
        FROM reactions r
        WHERE r.content_id = c.id
          AND r.user_id = $2
          AND r.deleted_at IS NULL
    ) AS liked_by_user,

    -- Preview image for comment (first image by sort_order)
    (SELECT i.file_name
     FROM images i
     WHERE i.entity_id = c.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ) AS preview_image

FROM comments c
WHERE c.parent_post_id = $1
  AND c.deleted_at IS NULL
ORDER BY c.created_at DESC -- can change to asc, TODO ask front
OFFSET $3
LIMIT $4; 

---------------------------------------------------------
-- create comment
---------------------------------------------------------
BEGIN;

-- Insert the comment
INSERT INTO comments (comment_creator_id, parent_post_id, comment_body)
VALUES ($1, $2, $3)
RETURNING id;

-- Insert associated images (if any)
-- $4 is an array of file names passed from the client
-- sort_order is auto-calculated by the trigger
INSERT INTO images (file_name, entity_id)
SELECT fname, c.id
FROM unnest($4::text[]) AS fname
CROSS JOIN (SELECT id FROM comments WHERE id = c.id) AS c;

COMMIT;

------------------------------------------------------------
-- edit comment
------------------------------------------------------------
UPDATE comments
SET comment_body = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

--------------------------------------------------------
-- delete comment
-------------------------------------------------------
UPDATE comments
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

------------------------------------------------------
-- get latest comment for post id
------------------------------------------------------
SELECT
    c.id,
    c.comment_creator_id,
    c.parent_post_id,
    c.comment_body,
    c.reactions_count,
    c.images_count,
    c.created_at,
    c.updated_at,

    -- liked by requesting user
    EXISTS (
        SELECT 1 FROM reactions r
        WHERE r.content_id = c.id
          AND r.user_id = $2
          AND r.deleted_at IS NULL
    ) AS liked_by_user,

    -- preview image (first by sort_order)
    (SELECT i.file_name
     FROM images i
     WHERE i.entity_id = c.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ) AS preview_image,


FROM comments c
WHERE c.parent_post_id = $1
  AND c.deleted_at IS NULL
ORDER BY c.created_at DESC
LIMIT 1;




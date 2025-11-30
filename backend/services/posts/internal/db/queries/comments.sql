-- name: GetCommentsByPostId :many
SELECT
    c.id,
    c.comment_creator_id,
    c.comment_body,
    c.reactions_count,
    c.images_count,
    c.created_at,
    c.updated_at,

    EXISTS (
        SELECT 1
        FROM reactions r
        WHERE r.content_id = c.id
          AND r.user_id = $2
          AND r.deleted_at IS NULL
    ) AS liked_by_user,

    (SELECT i.file_name
     FROM images i
     WHERE i.entity_id = c.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ) AS preview_image

FROM comments c
WHERE c.parent_id = $1
  AND c.deleted_at IS NULL
ORDER BY c.created_at DESC -- can change to asc, TODO ask front
OFFSET $3
LIMIT $4; 


-- name: CreateComment :one
INSERT INTO comments (comment_creator_id, parent_id, comment_body)
VALUES ($1, $2, $3)
RETURNING id;


-- name: EditComment :one
UPDATE comments
SET comment_body = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteComment :one
UPDATE comments
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: GetLatestCommentforPostId :one
SELECT
    c.id,
    c.comment_creator_id,
    c.parent_id,
    c.comment_body,
    c.reactions_count,
    c.images_count,
    c.created_at,
    c.updated_at,

    EXISTS (
        SELECT 1 FROM reactions r
        WHERE r.content_id = c.id
          AND r.user_id = $2
          AND r.deleted_at IS NULL
    ) AS liked_by_user,

    (SELECT i.file_name
     FROM images i
     WHERE i.entity_id = c.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ) AS preview_image


FROM comments c
WHERE c.parent_id = $1
  AND c.deleted_at IS NULL
ORDER BY c.created_at DESC
LIMIT 1;




-- name: GetPostByID :one
SELECT
    p.id,
    p.post_body,
    p.creator_id,
    p.group_id,
    p.audience,
    p.comments_count,
    p.reactions_count,
    p.last_commented_at,
    p.created_at,
    p.updated_at,

    EXISTS (
        SELECT 1 FROM reactions r
        WHERE r.content_id = p.id
          AND r.user_id = $1
          AND r.deleted_at IS NULL
    ) AS liked_by_user,

COALESCE(
    (SELECT i.id
     FROM images i
     WHERE i.parent_id = p.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ), 0
)::bigint AS image


FROM posts p
WHERE p.id=$2
  AND p.deleted_at IS NULL;

-- name: GetMostPopularPostInGroup :one
SELECT
    p.id,
    p.post_body,
    p.creator_id,
    p.group_id,
    p.audience,
    p.comments_count,
    p.reactions_count,
    p.last_commented_at,
    p.created_at,
    p.updated_at,

    (SELECT i.id    
     FROM images i
     WHERE i.parent_id = p.id
       AND i.deleted_at IS NULL
     ORDER BY i.sort_order
     LIMIT 1
    ) AS image,


    (p.reactions_count + p.comments_count) AS popularity_score     -- popularity metric (likes + comments)

FROM posts p
WHERE p.group_id = $1
  AND p.deleted_at IS NULL

ORDER BY popularity_score DESC, p.created_at DESC
LIMIT 1;

-- name: CreatePost :one
INSERT INTO posts (post_body, creator_id, group_id, audience)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: InsertPostAudience :execrows
INSERT INTO post_audience (post_id, allowed_user_id)
SELECT sqlc.arg(post_id)::bigint,
       allowed_user_id
FROM unnest(sqlc.arg(allowed_user_ids)::bigint[]) AS allowed_user_id;

-- name: EditPostContent :execrows
UPDATE posts
SET post_body  = $1
WHERE id = $2 AND creator_id = $3 AND deleted_at IS NULL;

-- name: UpdatePostAudience :execrows
UPDATE posts
SET audience = $3,
    updated_at = NOW()
WHERE 
    id = $1
    AND creator_id = $2
    AND deleted_at IS NULL
    AND (audience IS DISTINCT FROM $3);

-- name: ClearPostAudience :exec
DELETE FROM post_audience
WHERE post_id = $1;

-- name: DeletePost :execrows
UPDATE posts
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND creator_id=$2 AND deleted_at IS NULL;

-- name: GetPostAudience :many
SELECT allowed_user_id
FROM post_audience
WHERE post_id = $1
ORDER BY allowed_user_id;






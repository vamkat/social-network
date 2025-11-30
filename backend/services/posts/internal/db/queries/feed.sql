-- name: GetPublicFeed :many
SELECT
    p.id,
    p.post_body,
    p.creator_id,
    p.comments_count,
    p.reactions_count,
    p.last_commented_at,
    p.images_count,
    p.created_at,
    p.updated_at,

    EXISTS (
        SELECT 1 FROM reactions r
        WHERE r.content_id = p.id
          AND r.user_id = $1
          AND r.deleted_at IS NULL
    ) AS liked_by_user,

    (SELECT file_name
     FROM images i
     WHERE i.entity_id = p.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ) AS preview_image

FROM posts p
WHERE p.deleted_at IS NULL
  AND p.audience = 'everyone'
ORDER BY p.created_at DESC
OFFSET $2 LIMIT $3;


-- name: GetPersonalizedFeed :many
SELECT
    p.id,
    p.post_body,
    p.creator_id,
    p.comments_count,
    p.reactions_count,
    p.images_count,
    p.last_commented_at,
    p.created_at,
    p.updated_at,

    -- did user like it?
    EXISTS (
        SELECT 1 FROM reactions r
        WHERE r.content_id = p.id
          AND r.user_id = $1
          AND r.deleted_at IS NULL
    ) AS liked_by_user,

    -- first image preview
    (
      SELECT file_name
      FROM images i
      WHERE i.entity_id = p.id
        AND i.deleted_at IS NULL
      ORDER BY i.sort_order
      LIMIT 1
    ) AS preview_image

FROM posts p
WHERE p.deleted_at IS NULL
  AND (
       -- SELECTED audience → only manually approved viewers
       (p.audience = 'selected' AND EXISTS (
           SELECT 1 FROM post_audience pa
           WHERE pa.post_id = p.id AND pa.allowed_user_id = $1
       ))

       -- FOLLOWERS → allowed if creator ∈ list passed in
       OR (p.audience = 'followers' AND p.creator_id = ANY($2))
  )
ORDER BY p.created_at DESC
OFFSET $3 LIMIT $4;



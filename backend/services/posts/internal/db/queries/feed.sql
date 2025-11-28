-------------------------------------------------------
-- name: GetPublicFeed :many
---------------------------------------------------------
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

    -- Whether requesting user has liked it
    EXISTS (
        SELECT 1 FROM reactions r
        WHERE r.content_id = p.id
          AND r.user_id = $1
          AND r.deleted_at IS NULL
    ) AS liked_by_user,

    -- Preview image (first image by sort_order)
    (SELECT i.file_name
     FROM images i
     WHERE i.entity_id = p.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ) AS preview_image

FROM posts p

WHERE p.deleted_at IS NULL
  AND p.audience = 'everyone'
  AND ($2::timestamptz IS NULL OR p.created_at < $2)

ORDER BY p.created_at DESC
LIMIT $3;


-------------------------------------------------
-- name: GetPersonalizedFeed :many
-------------------------------------------------
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

    -- whether requesting user has reacted
    EXISTS (
        SELECT 1 FROM reactions r
        WHERE r.content_id = p.id
        AND r.user_id = $1
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

WHERE p.deleted_at IS NULL

  AND (
        -- selected audience → must be explicitly granted
        (p.audience = 'selected' AND EXISTS (
            SELECT 1 FROM post_audience pa
            WHERE pa.post_id = p.id
              AND pa.allowed_user_id = $1
        ))

        -- followers feed → creator must be in follow list
        OR (p.audience = 'followers' AND p.creator_id = ANY ($3))
      )

  -- cursor pagination: only posts after last seen
  AND ($2 IS NULL OR p.id > $2)

ORDER BY p.created_at DESC
LIMIT $4;

------------------------------------------------------
-- insert or update last seen post id for user
--------------------------------------------------------
INSERT INTO user_feed_state (user_id, last_seen_post_id)
VALUES ($1, $latest_post_id)
ON CONFLICT (user_id) DO UPDATE
SET last_seen_post_id = EXCLUDED.last_seen_post_id,
    updated_at = CURRENT_TIMESTAMP;

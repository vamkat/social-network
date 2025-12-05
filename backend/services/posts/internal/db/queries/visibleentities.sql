-- name: CanUserSeeEntity :one
WITH ent AS (
    -- Post
    SELECT 
        id,
        creator_id,
        audience,
        group_id
    FROM posts
    WHERE id = sqlc.arg(entity_id)::bigint
      AND deleted_at IS NULL

    UNION ALL

    -- Event
    SELECT
        id,
        creator_id,
        NULL AS audience,
        group_id
    FROM events
    WHERE id = sqlc.arg(entity_id)::bigint
      AND deleted_at IS NULL
)
SELECT EXISTS (
    SELECT 1
    FROM ent e
    WHERE
        (
            -- CASE 1: group entity
            e.group_id IS NOT NULL
            AND sqlc.arg(is_member)::bool = TRUE
        )
        OR
        (
            -- CASE 2: post (no group)
            e.group_id IS NULL
            AND (
                e.audience = 'everyone'
                OR (e.audience = 'followers' AND sqlc.arg(is_following)::bool = TRUE)
                OR (
                    e.audience = 'selected'
                    AND EXISTS (
                        SELECT 1 FROM post_audience pa
                        WHERE pa.post_id = e.id
                          AND pa.allowed_user_id = sqlc.arg(user_id)::bigint
                    )
                )
            )
        )
);

-- name: GetEntityCreatorAndGroup :one
SELECT
    mi.content_type,

    -- creator_id: post.creator_id, event.event_creator_id, or parent post creator for comments
    (
        CASE
            WHEN mi.content_type = 'post'
                THEN p.creator_id
            WHEN mi.content_type = 'event'
                THEN e.event_creator_id
            WHEN mi.content_type = 'comment'
                THEN p2.creator_id  -- comment's parent post
        END
    )::BIGINT AS creator_id,

    -- group_id: post.group_id, event.group_id, or parent post group for comments
    (
        CASE
            WHEN mi.content_type = 'post'
                THEN p.group_id
            WHEN mi.content_type = 'event'
                THEN e.group_id
            WHEN mi.content_type = 'comment'
                THEN p2.group_id  -- comment's parent post
        END
    )::BIGINT AS group_id

FROM master_index mi
LEFT JOIN posts p ON p.id = mi.id
LEFT JOIN events e ON e.id = mi.id
LEFT JOIN comments c ON c.id = mi.id
LEFT JOIN posts p2 ON p2.id = c.parent_id  -- parent post for comments
WHERE mi.id = $1
LIMIT 1;

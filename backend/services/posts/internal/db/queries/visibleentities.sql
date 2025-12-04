-- name: CanUserSeeEntity :one
WITH
    ctx_user AS (
        SELECT sqlc.arg(user_id)::bigint AS user_id
    ),
    ctx_following AS (
        SELECT UNNEST(sqlc.arg(following_ids)::bigint[]) AS following_id
    ),
    ctx_groups AS (
        SELECT UNNEST(sqlc.arg(group_ids)::bigint[]) AS group_id
    ),

    -- Entities normalized into one row
    ent AS (
        SELECT 
            id,
            creator_id,
            audience,
            group_id
        FROM posts
        WHERE id = sqlc.arg(entity_id)::bigint
          AND deleted_at IS NULL

        UNION ALL

        SELECT
            id,
            creator_id,
            NULL AS audience,   -- events don't use audience
            group_id
        FROM events
        WHERE id = sqlc.arg(entity_id)::bigint
          AND deleted_at IS NULL
    )

SELECT EXISTS (
    SELECT 1
    FROM ent e
    CROSS JOIN ctx_user u

    -- group membership check
    LEFT JOIN ctx_groups g
        ON g.group_id = e.group_id

    -- following check
    LEFT JOIN ctx_following f
        ON f.following_id = e.creator_id

    WHERE
        (
            -- CASE 1: group entity (group_id != NULL)
            e.group_id IS NOT NULL
            AND g.group_id IS NOT NULL
        )
        OR
        (
            -- CASE 2: post (no group_id), apply audience rules
            e.group_id IS NULL AND (
                e.audience = 'everyone'
                OR (e.audience = 'followers' AND f.following_id IS NOT NULL)
                OR (
                    e.audience = 'selected'
                    AND EXISTS (
                        SELECT 1 FROM post_audience pa
                        WHERE pa.post_id = e.id
                          AND pa.allowed_user_id = u.user_id
                    )
                )
            )
        )
);

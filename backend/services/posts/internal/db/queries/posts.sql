-- name: GetGroupPostsPaginated :many
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

    EXISTS (     -- Has the given user liked the post?
        SELECT 1
        FROM reactions r
        WHERE r.content_id = p.id
          AND r.user_id = $2              -- requesting user (check is member from users service)
          AND r.deleted_at IS NULL
    ) AS liked_by_user,
   
    (SELECT i.id    
     FROM images i
     WHERE i.parent_id = p.id
       AND i.deleted_at IS NULL
     ORDER BY i.sort_order
     LIMIT 1
    ) AS image,

     -- latest comment using LATERAL join
    lc.id AS latest_comment_id,
    lc.comment_creator_id AS latest_comment_creator_id,
    lc.comment_body AS latest_comment_body,
    lc.reactions_count AS latest_comment_reactions_count,
    lc.created_at AS latest_comment_created_at,
    lc.updated_at AS latest_comment_updated_at,
    lc.liked_by_user AS latest_comment_liked_by_user,
    lc.image AS latest_comment_image


FROM posts p

LEFT JOIN LATERAL (
    SELECT
        c.id,
        c.comment_creator_id,
        c.comment_body,
        c.reactions_count,
        c.created_at,
        c.updated_at,
        EXISTS (
            SELECT 1 FROM reactions r
            WHERE r.content_id = c.id
              AND r.user_id = $2
              AND r.deleted_at IS NULL
        ) AS liked_by_user,
        (
            SELECT i.id
            FROM images i
            WHERE i.parent_id = c.id
              AND i.deleted_at IS NULL
            ORDER BY i.sort_order
            LIMIT 1
        ) AS image
    FROM comments c
    WHERE c.parent_id = p.id
      AND c.deleted_at IS NULL
    ORDER BY c.created_at DESC
    LIMIT 1
) lc ON TRUE


WHERE p.group_id = $1                    -- group id filter
  AND p.deleted_at IS NULL
GROUP BY p.id
ORDER BY p.created_at DESC               -- newest first
LIMIT $3 OFFSET $4;                      -- pagination

-- name: GetUserPostsPaginated :many
SELECT
    p.id,
    p.post_body,
    p.creator_id,
    p.comments_count,
    p.reactions_count,
    p.last_commented_at,
    p.created_at,
    p.updated_at,

    EXISTS (    -- Has the requesting user liked the post?
        SELECT 1 FROM reactions r
        WHERE r.content_id = p.id
          AND r.user_id = $2
          AND r.deleted_at IS NULL
    ) AS liked_by_user,

    (SELECT i.id     
     FROM images i
     WHERE i.parent_id = p.id
       AND i.deleted_at IS NULL
     ORDER BY i.sort_order
     LIMIT 1
    ) AS image,

     -- latest comment using LATERAL join
    lc.id AS latest_comment_id,
    lc.comment_creator_id AS latest_comment_creator_id,
    lc.comment_body AS latest_comment_body,
    lc.reactions_count AS latest_comment_reactions_count,
    lc.created_at AS latest_comment_created_at,
    lc.updated_at AS latest_comment_updated_at,
    lc.liked_by_user AS latest_comment_liked_by_user,
    lc.image AS latest_comment_image

FROM posts p

LEFT JOIN LATERAL (
    SELECT
        c.id,
        c.comment_creator_id,
        c.comment_body,
        c.reactions_count,
        c.created_at,
        c.updated_at,
        EXISTS (
            SELECT 1 FROM reactions r
            WHERE r.content_id = c.id
              AND r.user_id = $2
              AND r.deleted_at IS NULL
        ) AS liked_by_user,
        (
            SELECT i.id
            FROM images i
            WHERE i.parent_id = c.id
              AND i.deleted_at IS NULL
            ORDER BY i.sort_order
            LIMIT 1
        ) AS image
    FROM comments c
    WHERE c.parent_id = p.id
      AND c.deleted_at IS NULL
    ORDER BY c.created_at DESC
    LIMIT 1
) lc ON TRUE

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
         OR (
            p.audience = 'followers'
            AND p.creator_id = ANY($3::bigint[])          -- viewer follows creator
        )
     )

GROUP BY p.id
ORDER BY p.created_at DESC
LIMIT $4 OFFSET $5;

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






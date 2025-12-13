-- name: GetPublicFeed :many
SELECT
    p.id,
    p.post_body,
    p.creator_id,
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

    -- image
COALESCE(
    (SELECT i.id
     FROM images i
     WHERE i.parent_id = p.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ), 0
)::bigint AS image   

   FROM posts p



WHERE p.deleted_at IS NULL
  AND (
       -- SELECTED audience → only manually approved viewers
       (p.audience = 'selected' AND EXISTS (
           SELECT 1 FROM post_audience pa
           WHERE pa.post_id = p.id AND pa.allowed_user_id = $1
       ))

       -- FOLLOWERS → allowed if creator ∈ list passed in
       OR (p.audience = 'followers' AND p.creator_id = ANY($2::bigint[]))
  )
ORDER BY p.created_at DESC
OFFSET $3 LIMIT $4;


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
   
COALESCE(
    (SELECT i.id
     FROM images i
     WHERE i.parent_id = p.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ), 0
)::bigint AS image  

  
FROM posts p



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

COALESCE(
    (SELECT i.id
     FROM images i
     WHERE i.parent_id = p.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ), 0
)::bigint AS image   

  
FROM posts p



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
         OR (p.audience = 'followers' AND $3::bool = TRUE)
     )

GROUP BY p.id
ORDER BY p.created_at DESC
LIMIT $4 OFFSET $5;

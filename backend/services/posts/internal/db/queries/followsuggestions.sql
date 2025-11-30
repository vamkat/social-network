-- name: SuggestUsersByPostActivity :many
WITH

-- U1: Users who liked one or more of *your public posts*
u1 AS (
    SELECT 
        r.user_id,
        4 AS score
    FROM reactions r
    JOIN posts p ON p.id = r.content_id
    WHERE p.creator_id = $1
      AND p.audience = 'everyone'
      AND r.user_id <> $1
),

-- U2: Users who commented on your public posts
u2 AS (
    SELECT
        c.comment_creator_id AS user_id,
        4 AS score
    FROM comments c
    JOIN posts p ON p.id = c.parent_id
    WHERE p.creator_id = $1
      AND p.audience = 'everyone'
      AND c.comment_creator_id <> $1
),

-- U3: Users who liked the same posts as you
u3 AS (
    SELECT DISTINCT
        r2.user_id,
        3 AS score
    FROM reactions r1                        -- your likes
    JOIN reactions r2 ON r1.content_id = r2.content_id
    WHERE r1.user_id = $1
      AND r2.user_id <> $1
),

-- U4: Users who commented on the same posts as you
u4 AS (
    SELECT DISTINCT
        c2.comment_creator_id AS user_id,
        2 AS score
    FROM comments c1                         -- your comments
    JOIN comments c2 ON c1.parent_id = c2.parent_id
    WHERE c1.comment_creator_id = $1
      AND c2.comment_creator_id <> $1
),

-- Combine scores
combined AS (
    SELECT user_id, SUM(score) AS total_score
    FROM (
        SELECT * FROM u1
        UNION ALL
        SELECT * FROM u2
        UNION ALL
        SELECT * FROM u3
        UNION ALL
        SELECT * FROM u4
    ) scored
    GROUP BY user_id
)

SELECT user_id
FROM combined
ORDER BY total_score DESC, random()
LIMIT 5;

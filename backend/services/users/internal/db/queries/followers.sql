-- name: FollowUser :one
SELECT follow_user($1, $2);
--1: follower_id
--2: following_id   
-- returns followed or requested depending on target's privacy settings



-- name: UnfollowUser :execrows
WITH deleted_follow AS (
    DELETE FROM follows
    WHERE follower_id = $1
      AND following_id = $2
    RETURNING 1
),
deleted_request AS (
    DELETE FROM follow_requests
    WHERE requester_id = $1
      AND target_id = $2
    RETURNING 1
)
SELECT 1
FROM deleted_follow
UNION ALL
SELECT 1
FROM deleted_request;


-- name: GetFollowers :many
SELECT u.id, u.username, u.avatar_id,u.profile_public, f.created_at AS followed_at
FROM follows f
JOIN users u ON u.id = f.follower_id
WHERE f.following_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetFollowing :many
SELECT u.id, u.username,u.avatar_id,u.profile_public, f.created_at AS followed_at
FROM follows f
JOIN users u ON u.id = f.following_id
WHERE f.follower_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetFollowingIds :many
SELECT following_id
FROM follows 
WHERE follower_id = $1;


-- name: GetFollowerCount :one
SELECT COUNT(*) AS follower_count
FROM follows
WHERE following_id = $1;    


-- name: GetFollowingCount :one
SELECT COUNT(*) AS following_count
FROM follows
WHERE follower_id = $1;


-- name: IsFollowing :one
SELECT EXISTS (
    SELECT 1 FROM follows
    WHERE follower_id =$1 AND following_id = $2
);


-- name: AreFollowingEachOther :one
WITH u1 AS (
  SELECT EXISTS (
    SELECT 1
    FROM follows f
    WHERE f.follower_id = $1 AND f.following_id = $2
  ) AS user1_follows_user2
),
u2 AS (
  SELECT EXISTS (
    SELECT 1
    FROM follows f
    WHERE f.follower_id = $2 AND f.following_id = $1
  ) AS user2_follows_user1
)
SELECT
  u1.user1_follows_user2,
  u2.user2_follows_user1
FROM u1, u2;



-- name: GetMutualFollowers :many
SELECT u.id, u.username
FROM follows f1
JOIN follows f2 ON f1.follower_id = f2.follower_id
JOIN users u ON u.id = f1.follower_id
WHERE f1.following_id = $1
  AND f2.following_id = $2;


-- name: AcceptFollowRequest :exec
WITH updated AS (
    UPDATE follow_requests
    SET status = 'accepted', updated_at = NOW()
    WHERE requester_id = $1
      AND target_id    = $2
      AND deleted_at IS NULL
    RETURNING requester_id, target_id
)
INSERT INTO follows (follower_id, following_id)
SELECT requester_id, target_id
FROM updated
ON CONFLICT DO NOTHING;



-- name: RejectFollowRequest :exec
UPDATE follow_requests
SET status = 'rejected', updated_at = NOW()
WHERE requester_id = $1 AND target_id = $2;

-- name: IsFollowRequestPending :one
SELECT EXISTS(
    SELECT 1 
    FROM follow_requests
    WHERE requester_id = $1
      AND target_id = $2
      AND status = 'pending'
) AS has_pending_request;


-- name: GetFollowSuggestions :many
WITH
-- S1: second-degree follows
s1 AS (
    SELECT 
        f2.following_id AS user_id,
        5 AS score         -- weighted higher
    FROM follows f1 
    JOIN follows f2 ON f1.following_id = f2.follower_id
    WHERE f1.follower_id = $1
      AND f2.following_id <> $1
      AND NOT EXISTS (
          SELECT 1 FROM follows x
          WHERE x.follower_id = $1 AND x.following_id = f2.following_id
      )
),

-- S2: shared groups
s2 AS (
    SELECT
        gm2.user_id,
        3 AS score        -- lighter weight than follows
    FROM group_members gm1
    JOIN group_members gm2 ON gm1.group_id = gm2.group_id
    WHERE gm1.user_id = $1
      AND gm2.user_id <> $1
      AND gm2.deleted_at IS NULL
),

-- Combine & score
combined AS (
    SELECT user_id, SUM(score) AS total_score
    FROM (
        SELECT * FROM s1
        UNION ALL
        SELECT * FROM s2
    ) scored
    GROUP BY user_id
)

SELECT 
    u.id,
    u.username,
    u.avatar_id,
    c.total_score
FROM combined c
JOIN users u ON u.id = c.user_id
ORDER BY c.total_score DESC, random()
LIMIT 5;

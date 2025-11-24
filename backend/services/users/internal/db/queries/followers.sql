-- name: FollowUser :one
SELECT follow_user($1, $2);
--1: follower_id
--2: following_id   
-- returns followed or requested depending on target's privacy settings


-- name: UnfollowUser :exec
DELETE FROM follows
WHERE follower_id = $1 AND following_id = $2;




-- name: GetFollowers :many
SELECT u.id, u.username, u.avatar,u.profile_public, f.created_at AS followed_at
FROM follows f
JOIN users u ON u.id = f.follower_id
WHERE f.following_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetFollowing :many
SELECT u.id, u.username,u.avatar,u.profile_public, f.created_at AS followed_at
FROM follows f
JOIN users u ON u.id = f.following_id
WHERE f.follower_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;


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


-- name: IsFollowingEither :one
SELECT EXISTS (
    SELECT 1
    FROM follows
    WHERE (follower_id = $1 AND following_id = $2)
       OR (follower_id = $2 AND following_id = $1)
) AS is_following_either;



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

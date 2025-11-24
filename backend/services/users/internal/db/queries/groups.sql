-- name: CreateGroup :one
INSERT INTO groups (group_owner, group_title, group_description)
VALUES ($1, $2, $3)
RETURNING id;

-- name: AddGroupOwnerAsMember :exec
INSERT INTO group_members (group_id, user_id, role)
VALUES ($1, $2, 'owner');


-- name: SoftDeleteGroup :exec
UPDATE groups
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1;


-- name: SendGroupJoinRequest :exec
INSERT INTO group_join_requests (group_id, user_id, status)
VALUES ($1, $2, 'pending')
ON CONFLICT (group_id, user_id)
DO UPDATE SET status = 'pending';

-- name: AcceptGroupJoinRequest :exec
UPDATE group_join_requests
SET status = 'accepted'
WHERE group_id = $1
  AND user_id = $2;

-- name: AddUserToGroup :exec
INSERT INTO group_members (group_id, user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;


-- name: RejectGroupJoinRequest :exec
UPDATE group_join_requests
SET status = 'rejected'
WHERE group_id = $1
  AND user_id = $2;

-- name: CancelGroupJoinRequest :exec
DELETE FROM group_join_requests
WHERE group_id = $1
  AND user_id = $2;

-- name: SendGroupInvite :exec
INSERT INTO group_invites (group_id, sender_id, receiver_id, status)
VALUES ($1, $2, $3, 'pending') 
ON CONFLICT (group_id, receiver_id)
DO UPDATE SET status = 'pending';     

-- name: AcceptGroupInvite :exec
UPDATE group_invites
SET status = 'accepted'
WHERE group_id = $1
  AND receiver_id = $2; 


-- name: DeclineGroupInvite :exec
UPDATE group_invites
SET status = 'declined'
WHERE group_id = $1
  AND receiver_id = $2;

-- name: CancelGroupInvite :exec
DELETE FROM group_invites
WHERE group_id = $1
  AND receiver_id = $2
  AND sender_id=$3;

-- name: LeaveGroup :exec
UPDATE group_members
SET deleted_at = CURRENT_TIMESTAMP
WHERE group_id = $1
  AND user_id = $2
  AND role <> 'owner'; -- owners cannot leave the group (transfer ownership logic? TODO)


-- name: TransferOwnership :exec
WITH demote AS (
    UPDATE group_members AS gm_old
    SET role = 'member'
    WHERE gm_old.group_id = $1
      AND gm_old.user_id = $2
      AND gm_old.role = 'owner'
),
promote AS (
    UPDATE group_members AS gm_new
    SET role = 'owner'
    WHERE gm_new.group_id = $1
      AND gm_new.user_id = $3
      AND gm_new.role = 'member'
)
SELECT 1;

-- name: GetAllGroups :many
SELECT
  id,
  group_title,
  group_description,
  members_count
FROM groups
WHERE deleted_at IS NULL
ORDER BY members_count DESC, id ASC
LIMIT $1 OFFSET $2;

-- name: GetUserGroups :many
SELECT DISTINCT
    g.id AS group_id,
    g.group_title,
    g.group_description,
    g.members_count,
    CASE WHEN gm.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS is_member,
    CASE WHEN g.group_owner = $1 THEN TRUE ELSE FALSE END AS is_owner
FROM groups g
LEFT JOIN group_members gm
    ON gm.group_id = g.id
    AND gm.user_id = $1
    AND gm.deleted_at IS NULL
WHERE g.deleted_at IS NULL
  AND (gm.user_id = $1 OR g.group_owner = $1)
ORDER BY COALESCE(gm.joined_at, g.created_at) DESC, g.id DESC
LIMIT $2 OFFSET $3;

      
-- name: GetGroupInfo :one
SELECT
  id,
  group_owner,
  group_title,
  group_description,
  members_count
FROM groups
WHERE id=$1
  AND deleted_at IS NULL;

-- name: GetGroupMembers :many
SELECT
    u.id,
    u.username,
    u.avatar,
    u.profile_public,
    gm.role,
    gm.joined_at
FROM group_members gm
JOIN users u
    ON gm.user_id = u.id
WHERE gm.group_id = $1
  AND gm.deleted_at IS NULL
ORDER BY gm.joined_at DESC, u.id DESC
LIMIT $2 OFFSET $3;


-- name: SearchGroupsFuzzy :many
SELECT
    g.id,
    g.group_title,
    g.group_description,
    g.members_count,
    g.group_owner,
    CASE WHEN gm.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS is_member,
    CASE WHEN g.group_owner = $2 THEN TRUE ELSE FALSE END AS is_owner,
    (
        2 * similarity(g.group_title, $1) +
        1 * similarity(g.group_description, $1)
    ) AS weighted_score
FROM groups g
LEFT JOIN group_members gm
    ON gm.group_id = g.id
    AND gm.user_id = $2
    AND gm.deleted_at IS NULL
WHERE g.deleted_at IS NULL
  AND (
        similarity(g.group_title, $1) > 0.2
        OR similarity(g.group_description, $1) > 0.2
      )
ORDER BY
    -- 1. Weighted fuzzy match score
    weighted_score DESC,
    -- 2. Prioritize groups the user belongs to
    CASE WHEN gm.user_id IS NOT NULL THEN 1 ELSE 0 END DESC,
    -- 3. Prioritize groups with more members
    g.members_count DESC,
    -- 4. Stable pagination
    g.id DESC
LIMIT $3 OFFSET $4;

-- name: GetUserGroupRole :one
SELECT role
FROM group_members
WHERE group_id = $1
  AND user_id = $2
  AND deleted_at IS NULL;

-- name: IsUserGroupMember :one
SELECT EXISTS (
    SELECT 1
    FROM group_members
    WHERE group_id = $1
      AND user_id = $2
      AND deleted_at IS NULL
) AS is_member;

-- name: IsUserGroupOwner :one
SELECT (group_owner = $2) AS is_owner
FROM groups
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UserGroupCountsPerRole :one
SELECT
    COUNT(*) FILTER (WHERE g.group_owner = $1) AS owner_count,
    COUNT(*) FILTER (WHERE gm.role = 'member' AND g.group_owner <> $1) AS member_only_count,
    COUNT(*) AS total_memberships
FROM group_members gm
JOIN groups g ON gm.group_id = g.id
WHERE gm.user_id = $1
  AND gm.deleted_at IS NULL
  AND g.deleted_at IS NULL;
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
WHERE deleted_at IS NULL;

-- name: GetUserGroups :many
SELECT DISTINCT
    g.id AS group_id,
    g.group_title,
    g.group_description,
    g.members_count,
    CASE 
        WHEN g.group_owner = $1 THEN 'owner'
        ELSE 'member'
    END AS role
FROM groups g
LEFT JOIN group_members gm
    ON gm.group_id = g.id
    AND gm.user_id = $1
    AND gm.deleted_at IS NULL
WHERE g.deleted_at IS NULL
  AND (g.group_owner = $1 OR gm.user_id = $1);

      
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
    gm.user_id,
    u.username,
    u.avatar,
    u.profile_public,
    gm.role,
    gm.joined_at
FROM group_members gm
JOIN users u
    ON gm.user_id = u.id
WHERE gm.group_id = $1
  AND gm.deleted_at IS NULL;

-- name: SearchGroupsFuzzy :many
SELECT
    id,
    group_title,
    group_description,
    members_count
FROM groups
WHERE deleted_at IS NULL
  AND (similarity(group_title, $1) > 0.3
       OR similarity(group_description, $1) > 0.3)
ORDER BY GREATEST(similarity(group_title, $1), similarity(group_description, $1)) DESC;

-- name: GetUserGroupRole :one
SELECT role
FROM group_members
WHERE group_id = $1
  AND user_id = $2
  AND deleted_at IS NULL;
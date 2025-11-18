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
  AND receiver_id = $2;

-- name: LeaveGroup :exec
UPDATE group_members
SET deleted_at = CURRENT_TIMESTAMP
WHERE group_id = $1
  AND user_id = $2
  AND role <> 'owner'; -- owners cannot leave the group (transfer ownership logic? TODO)
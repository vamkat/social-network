-- name: GetNotificationsByUserIdAll :many
SELECT
    id,
    user_id,
    notif_type,
    source_service,
    source_entity_id,
    seen,
    needs_action,
    acted,
    payload,
    created_at,
    expires_at
FROM notifications
WHERE user_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;
--1: user_id
--2: limit
--3: offset

-- name: GetNotificationsByUserIdFilteredByStatus :many
SELECT
    id,
    user_id,
    notif_type,
    source_service,
    source_entity_id,
    seen,
    needs_action,
    acted,
    payload,
    created_at,
    expires_at
FROM notifications
WHERE user_id = $1
    AND seen = $2
    AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $3
OFFSET $4;
--1: user_id
--2: seen
--3: limit
--4: offset

-- name: GetNotificationCountAll :one
SELECT COUNT(*)
FROM notifications
WHERE user_id = $1
    AND deleted_at IS NULL;
--1: user_id

-- name: GetNotificationCountFilteredByStatus :one
SELECT COUNT(*)
FROM notifications
WHERE user_id = $1
    AND seen = $2
    AND deleted_at IS NULL;
--1: user_id
--2: seen

-- name: MarkNotificationsAsRead :exec
UPDATE notifications
SET seen = true
WHERE user_id = $1
    AND id = ANY($2::bigint[])
    AND seen = false
    AND deleted_at IS NULL;
--1: user_id
--2: notification_ids

-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    notif_type,
    source_service,
    source_entity_id,
    payload
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING
    id,
    user_id,
    notif_type,
    source_service,
    source_entity_id,
    seen,
    needs_action,
    acted,
    payload,
    created_at,
    expires_at;
--1: user_id
--2: notif_type
--3: source_service
--4: source_entity_id
--5: payload

-- name: MarkNotificationAsActed :exec
UPDATE notifications
SET acted = true
WHERE id = $1
    AND deleted_at IS NULL;
--1: id
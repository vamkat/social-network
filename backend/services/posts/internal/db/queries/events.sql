-- name: GetEventsByGroupId :many
SELECT
    e.id,
    e.event_title,
    e.event_body,
    e.event_creator_id,
    e.group_id,
    e.event_date,
    e.created_at,
    e.updated_at,
    e.going_count,
    e.not_going_count,

    (SELECT i.file_name
     FROM images i
     WHERE i.entity_id = e.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ) AS preview_image,

    (SELECT COUNT(1)
     FROM images i
     WHERE i.entity_id = e.id AND i.deleted_at IS NULL
    ) AS total_images

FROM events e
WHERE e.group_id = $1
  AND e.deleted_at IS NULL
  AND e.event_date >= CURRENT_DATE
ORDER BY e.event_date ASC
OFFSET $2
LIMIT $3;

-- name: CreateEvent :one

INSERT INTO events (
    event_title,
    event_body,
    event_creator_id,
    group_id,
    event_date,
    still_valid
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: UpdateStillValid :exec
-- needs to be run periodically, eg every day
UPDATE events
SET still_valid = FALSE
WHERE event_date < CURRENT_DATE
  AND still_valid = TRUE;

-- name: RespondToEvent :one
INSERT INTO event_responses (event_id, user_id, going)
VALUES ($1, $2, $3)
ON CONFLICT (event_id, user_id)
DO UPDATE
SET going = EXCLUDED.going,
    deleted_at = NULL,           -- restore if it was soft-deleted
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: DeleteEventResponse :one
UPDATE event_responses
SET deleted_at = CURRENT_TIMESTAMP
WHERE event_id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING *;

-- name: EditEvent :one
UPDATE events
SET event_title = $1,
    event_body = $2,
    event_date = $3
WHERE id = $4 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteEvent :one
UPDATE events
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
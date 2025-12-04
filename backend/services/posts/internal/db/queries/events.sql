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

    (SELECT i.id
     FROM images i
     WHERE i.parent_id = e.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ) AS image,

    -- user response (NULL if no response)
    er.going AS user_response

FROM events e
LEFT JOIN event_responses er
    ON er.event_id = e.id
   AND er.user_id = $4
   AND er.deleted_at IS NULL

WHERE e.group_id = $1
  AND e.deleted_at IS NULL
  AND e.event_date >= CURRENT_DATE

ORDER BY e.event_date ASC
OFFSET $2
LIMIT $3;

-- name: CreateEvent :exec

INSERT INTO events (
    event_title,
    event_body,
    event_creator_id,
    group_id,
    event_date
)
VALUES ($1, $2, $3, $4, $5);


-- name: DeleteEventResponse :execrows
UPDATE event_responses
SET deleted_at = CURRENT_TIMESTAMP
WHERE event_id = $1
  AND user_id = $2
  AND deleted_at IS NULL;

-- name: EditEvent :execrows
UPDATE events
SET event_title = $1,
    event_body = $2,
    event_date = $3
WHERE id = $4 AND event_creator_id=$5 AND deleted_at IS NULL;

-- name: DeleteEvent :execrows
UPDATE events
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND event_creator_id=$2 AND deleted_at IS NULL;

-- name: UpsertEventResponse :execrows
INSERT INTO event_responses (event_id, user_id, going)
VALUES ($1, $2, $3)
ON CONFLICT (event_id, user_id)
DO UPDATE
SET
    going = EXCLUDED.going,
    deleted_at = NULL,
    updated_at = CURRENT_TIMESTAMP;
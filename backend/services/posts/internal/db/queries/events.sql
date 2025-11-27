-- get group events paginated

-- EVENTS DON"T HAVE LIKES, COMMENTS but they do have IMAGES

-- get event by id --show created at too
SELECT *
FROM events
WHERE id = $1
  AND deleted_at IS NULL;

  -- create event
INSERT INTO events (event_title, event_body, event_creator_id, group_id, event_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- respond to event
INSERT INTO event_responses (event_id, user_id, going)
VALUES ($1, $2, $3)
ON CONFLICT (event_id, user_id)
DO UPDATE SET going = EXCLUDED.going
RETURNING *;

-- get event responses
SELECT *
FROM event_responses
WHERE event_id = $1 AND deleted_at IS NULL;

-- delete event response
UPDATE event_responses
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;


-- edit event
UPDATE events
SET event_title = $1,
    event_body = $2,
    event_date = $3,
    still_valid = $4
WHERE id = $5 AND deleted_at IS NULL
RETURNING *;

-- delete event
UPDATE events
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
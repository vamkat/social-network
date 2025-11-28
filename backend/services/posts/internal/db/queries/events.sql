-- TODO ask front about preview image vs all images
-- TODO check triggers

-- TODO write queries for events
-- rethink carefully about feed
-- think of cross service needs between users and posts
-- test queries and start with java

-- how to add images count for events? trigger?

--------------------------------------------------------------------
-- name: GetEventsByGroup :many
--------------------------------------------------------------------
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

    -- preview image (first by sort_order)
    (SELECT i.file_name
     FROM images i
     WHERE i.entity_id = e.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ) AS preview_image,

    -- total number of images
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

---------------------------------------------------
  -- create event
  ------------------------------------------------
BEGIN;

-- Insert the event
WITH new_event AS (
    INSERT INTO events (
        event_title,
        event_body,
        event_creator_id,
        group_id,
        event_date,
        still_valid
    )
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *
),

-- Insert associated images (if any)
inserted_images AS (
    INSERT INTO images (file_name, entity_id)
    SELECT fname, e.id
    FROM unnest($7::text[]) AS fname
    CROSS JOIN new_event e
    RETURNING *
)

-- Return event with preview image and total images
SELECT
    e.id,
    e.event_title,
    e.event_body,
    e.event_creator_id,
    e.group_id,
    e.event_date,
    e.still_valid,
    e.going_count,
    e.not_going_count,
    e.created_at,
    e.updated_at,

    -- preview image (first by sort_order)
    (SELECT i.file_name
     FROM images i
     WHERE i.entity_id = e.id AND i.deleted_at IS NULL
     ORDER BY i.sort_order ASC
     LIMIT 1
    ) AS preview_image,

    -- total number of images
    (SELECT COUNT(1)
     FROM images i
     WHERE i.entity_id = e.id AND i.deleted_at IS NULL
    ) AS total_images

FROM new_event e;

COMMIT;

----------------------------------------------------------
-- Change still_valid for past events (needs to be run periodically, eg every day)
---------------------------------------------------------
UPDATE events
SET still_valid = FALSE
WHERE event_date < CURRENT_DATE
  AND still_valid = TRUE;

---------------------------------------------------------------------
-- respond to event
---------------------------------------------------------------------
INSERT INTO event_responses (event_id, user_id, going)
VALUES ($1, $2, $3)
ON CONFLICT (event_id, user_id)
DO UPDATE
SET going = EXCLUDED.going,
    deleted_at = NULL,           -- restore if it was soft-deleted
    updated_at = CURRENT_TIMESTAMP
RETURNING *;



-----------------------------------------------------------
-- delete event response
-------------------------------------------------------------
UPDATE event_responses
SET deleted_at = CURRENT_TIMESTAMP
WHERE event_id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING *;

------------------------------------------------------------
-- edit event
------------------------------------------------------------------
UPDATE events
SET event_title = $1,
    event_body = $2,
    event_date = $3,
WHERE id = $5 AND deleted_at IS NULL
RETURNING *;

------------------------------------------------------------
-- delete event
------------------------------------------------------------
UPDATE events
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
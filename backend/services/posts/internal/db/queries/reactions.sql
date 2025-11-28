----------------------------------------------
-- Get user reaction for entity id
---------------------------------------------
SELECT 1 FROM reactions r
WHERE r.content_id = p.id
AND r.user_id = $1
AND r.deleted_at IS NULL;

-----------------------------------------------
-- Toggle reaction for a user
-----------------------------------------------
WITH updated AS (
    UPDATE reactions
    SET deleted_at = CASE
                         WHEN deleted_at IS NULL THEN CURRENT_TIMESTAMP   -- soft delete
                         ELSE NULL                                       -- restore
                     END,
        updated_at = CURRENT_TIMESTAMP
    WHERE content_id = $1
      AND user_id = $2
    RETURNING *
)
-- If no row existed, insert it
INSERT INTO reactions (content_id, user_id)
SELECT $1, $2
WHERE NOT EXISTS (SELECT 1 FROM updated)
RETURNING *;

----------------------------------------------------
-- get who liked entity id
----------------------------------------------------
SELECT user_id
FROM reactions
WHERE content_id = $1 AND deleted_at IS NULL;
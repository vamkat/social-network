
--------------------------------------------------------------
-- insert images
--------------------------------------------------------------
INSERT INTO images (file_name, entity_id, sort_order)
SELECT file_name, $1::BIGINT, sort_order
FROM unnest($2::text[], $3::int[]) AS t(file_name, sort_order)
RETURNING id, file_name, sort_order, created_at;

------------------------------------------------------------
-- name: GetImages :many
-------------------------------------------------------------
SELECT
    id,
    file_name,
    sort_order,
    created_at
FROM images
WHERE entity_id = $1
  AND deleted_at IS NULL
ORDER BY sort_order;

-----------------------------------------------------------
-- update image (name or sort order)
-------------------------------------------------------------
UPDATE images
SET file_name = $1,
    sort_order = $2
WHERE id = $3 AND deleted_at IS NULL
RETURNING *;

-------------------------------------------------------------
-- delete image
-------------------------------------------------------------
UPDATE images
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

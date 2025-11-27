
-- insert image
INSERT INTO images (file_name, entity_id, sort_order)
VALUES ($1, $2, $3)
RETURNING *;

-- get images for entity id
SELECT *
FROM images
WHERE entity_id = $1 AND deleted_at IS NULL
ORDER BY sort_order ASC;

-- update image (name or sort order)
UPDATE images
SET file_name = $1,
    sort_order = $2
WHERE id = $3 AND deleted_at IS NULL
RETURNING *;

-- delete image
UPDATE images
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

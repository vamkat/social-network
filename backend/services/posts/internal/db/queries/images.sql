-- name: InsertImages :many
INSERT INTO images (file_name, entity_id)
SELECT unnest($2::text[]), $1::BIGINT
RETURNING id, file_name, sort_order, created_at;

-- name: GetImages :many
SELECT
    id,
    file_name,
    sort_order,
    created_at
FROM images
WHERE entity_id = $1
  AND deleted_at IS NULL
ORDER BY sort_order;

-- name: UpdateImage :one
UPDATE images
SET file_name = $1,
    sort_order = $2
WHERE id = $3 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteImage :one
UPDATE images
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

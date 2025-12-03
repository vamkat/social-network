-- name: InsertImage :exec
INSERT INTO images (id, parent_id)
VALUES ($2::BIGINT, $1::BIGINT);

-- name: GetImages :one
SELECT id
FROM images
WHERE parent_id = $1
  AND deleted_at IS NULL
ORDER BY sort_order
  LIMIT 1;

-- name: DeleteImage :exec 
UPDATE images
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;


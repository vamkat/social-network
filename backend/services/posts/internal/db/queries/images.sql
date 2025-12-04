-- name: UpsertImage :exec
INSERT INTO images (id, parent_id)
VALUES (sqlc.arg(id)::BIGINT, sqlc.arg(parent_id)::BIGINT)
ON CONFLICT (id) DO UPDATE
SET parent_id = EXCLUDED.parent_id;

-- name: GetImages :one
SELECT id
FROM images
WHERE parent_id = $1
  AND deleted_at IS NULL
ORDER BY sort_order
  LIMIT 1;

-- name: DeleteImage :execrows 
UPDATE images
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;


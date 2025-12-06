-- name: SaveImageMetadata :many
INSERT INTO images (
    original_name,
    bucket,
    object_key,
    mime_type,
    size_bytes,
    width,
    height,
    checksum
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, created_at;

-- name: GetImageById :many
SELECT *
FROM images
WHERE id = $1;

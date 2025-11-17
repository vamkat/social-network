-- name: GetUserProfile :one
SELECT
    id,
    username,
    first_name,
    last_name,
    date_of_birth
    avatar,
    about_me,
    profile_public
FROM users
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UpdateUserProfile :one
UPDATE users
SET
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    date_of_birth = COALESCE($4, date_of_birth),
    avatar = COALESCE($5, avatar),
    about_me = COALESCE($6, about_me),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;


-- name: UpdateUserPassword :exec
UPDATE auth_user
SET
    password_hash = $2,
    salt = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateUserEmail :exec
UPDATE auth_user
SET
    email = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;
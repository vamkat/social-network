-- name: GetUserProfile :one
SELECT
    public_id,
    username,
    first_name,
    last_name,
    date_of_birth,
    avatar,
    about_me,
    profile_public,
    created_at
FROM users
WHERE id = get_internal_user_id($1)
  AND deleted_at IS NULL;

-- name: GetUserBasic :one
SELECT
    public_id,
    username,
    avatar
FROM users
WHERE id = get_internal_user_id($1);
  
  
-- name: UpdateUserProfile :one
UPDATE users
SET
    username      = $2,
    first_name    = $3,
    last_name     = $4,
    date_of_birth = $5,
    avatar        = $6,
    about_me      = $7,
    updated_at    = CURRENT_TIMESTAMP
WHERE id = get_internal_user_id($1) AND deleted_at IS NULL
RETURNING *;

-- name: UpdateProfilePrivacy :exec
UPDATE users
SET profile_public=$2
WHERE id=get_internal_user_id($1);

-- name: UpdateUserPassword :exec
UPDATE auth_user
SET
    password_hash = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = get_internal_user_id($1);

-- name: UpdateUserEmail :exec
UPDATE auth_user
SET
    email = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = get_internal_user_id($1);

-- name: SearchUsers :many
SELECT
    public_id,
    username,
    avatar,
    profile_public
FROM users
WHERE deleted_at IS NULL
  AND (
        username % $1 OR
        first_name % $1 OR
        last_name % $1
      )
ORDER BY similarity(username, $1) DESC,
         similarity(first_name, $1) DESC,
         similarity(last_name, $1) DESC
LIMIT $2;
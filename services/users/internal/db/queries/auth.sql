-- name: InsertNewUser :one
INSERT INTO users (
    username,
    first_name,
    last_name,
    date_of_birth,
    avatar,
    about_me,
    profile_public
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id;

-- name: InsertNewUserAuth :exec
INSERT INTO auth_user (
    user_id,
    email,
    password_hash,
    salt
) VALUES (
       $1, $2, $3, $4
);



-- name: GetUserForLogin :one
SELECT
    u.id,
    u.username,
    au.password_hash,
    au.salt
FROM users u
JOIN auth_user au ON au.user_id = u.id
WHERE u.username = $1
  AND u.current_status = 'active'
  AND u.deleted_at IS NULL;


-- name: SoftDeleteUser :exec
UPDATE users
SET
    current_status = 'deleted',
    deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND deleted_at IS NULL;


-- name: BanUser :exec
UPDATE users
SET 
    current_status = 'banned',
    ban_ends_at = $2
WHERE id = $1;


-- name: UnbanUser :exec
UPDATE users
SET 
    current_status = 'active',
    ban_ends_at = NULL
WHERE id = $1;


-- name: IncrementFailedLoginAttempts :exec
UPDATE auth_user
SET failed_attempts = failed_attempts + 1
WHERE user_id = $1;


-- name: ResetFailedLoginAttempts :exec
UPDATE auth_user
SET failed_attempts = 0,
    last_login_at = CURRENT_TIMESTAMP
WHERE user_id = $1;



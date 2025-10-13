-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserFromRefreshToken :one
SELECT users.*
FROM users
INNER JOIN refresh_tokens
ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1;

-- name: UpdateUser :one
UPDATE users
SET updated_at = NOW(), email = $2, hashed_password = $3
WHERE id = $1
RETURNING *;

-- name: UpgradeChirpy :exec
UPDATE chirps
SET is_chirpy_red = true
WHERE id = $1;

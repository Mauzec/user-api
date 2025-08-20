-- name: CreateUser :one
INSERT INTO users (
    username,
    full_name,
    sex,
    age,
    email,
    phone,
    hashed_password
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByUsernameForUpdate :one
SELECT * FROM users
WHERE username = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateUser :one
UPDATE users
SET phone = $2,
    full_name = $3,
    sex = $4,
    email = $5
WHERE id = $1
RETURNING *;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1;
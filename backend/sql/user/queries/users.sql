-- name: CreateUser :execresult
INSERT INTO users (email, public_key, nickname, role, created_at)
VALUES (?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT id, email, public_key, nickname, role, created_at
FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, email, public_key, nickname, role, created_at
FROM users
WHERE email = ? LIMIT 1;

-- name: GetUserByPublicKey :one
SELECT id, email, public_key, nickname, role, created_at
FROM users
WHERE public_key = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET email = ?, public_key = ?, nickname = ?, role = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
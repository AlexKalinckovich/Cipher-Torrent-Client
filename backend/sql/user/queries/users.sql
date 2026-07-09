-- name: CreateUser :execresult
INSERT INTO users (email, password_hash, public_key, private_key_enc, nickname, role, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: GetUserByNickname :one
SELECT * FROM users WHERE nickname = ? LIMIT 1;

-- name: GetUserByPublicKey :one
SELECT * FROM users WHERE public_key = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET email = ?, password_hash = ?, nickname = ?, role = ?
WHERE id = ?;

-- name: PatchUser :exec
UPDATE users
SET
    email         = COALESCE(sqlc.narg('email'),    email),
    nickname      = COALESCE(sqlc.narg('nickname'), nickname),
    role          = COALESCE(sqlc.narg('role'),     role),
    password_hash = COALESCE(sqlc.narg('password_hash'), password_hash)
WHERE id = sqlc.arg('id');

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
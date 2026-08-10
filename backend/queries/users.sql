-- name: CreateUser :one
INSERT INTO users (username, password_hash, full_name, role)
VALUES ($1, $2, $3, $4)
RETURNING id, username, password_hash, full_name, role, created_at;

-- name: GetUserByUsername :one
SELECT id, username, password_hash, full_name, role, created_at
FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT id, username, password_hash, full_name, role, created_at
FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT id, username, password_hash, full_name, role, created_at
FROM users ORDER BY full_name;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = $2 WHERE id = $1;

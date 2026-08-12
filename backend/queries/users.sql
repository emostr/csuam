-- name: CreateUser :execlastid
INSERT INTO users (username, password_hash, full_name, role)
VALUES (?, ?, ?, ?);

-- name: GetUserByUsername :one
SELECT id, username, password_hash, full_name, role, created_at
FROM users WHERE username = ?;

-- name: GetUserByID :one
SELECT id, username, password_hash, full_name, role, created_at
FROM users WHERE id = ?;

-- name: ListUsers :many
SELECT id, username, password_hash, full_name, role, created_at
FROM users ORDER BY full_name;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = ? WHERE id = ?;

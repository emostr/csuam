package db

import (
	"context"
)

const createUser = `
INSERT INTO users (username, password_hash, full_name, role)
VALUES ($1, $2, $3, $4)
RETURNING id, username, password_hash, full_name, role, created_at
`

func (q *Queries) CreateUser(ctx context.Context, username, passwordHash, fullName, role string) (User, error) {
	var u User
	err := q.pool.QueryRow(ctx, createUser, username, passwordHash, fullName, role).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt)
	return u, err
}

const getUserByUsername = `
SELECT id, username, password_hash, full_name, role, created_at
FROM users WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	var u User
	err := q.pool.QueryRow(ctx, getUserByUsername, username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt)
	return u, err
}

const getUserByID = `
SELECT id, username, password_hash, full_name, role, created_at
FROM users WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	var u User
	err := q.pool.QueryRow(ctx, getUserByID, id).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt)
	return u, err
}

const listUsers = `
SELECT id, username, password_hash, full_name, role, created_at
FROM users ORDER BY full_name
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.pool.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, u)
	}
	return items, rows.Err()
}

const deleteUser = `DELETE FROM users WHERE id = $1`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.pool.Exec(ctx, deleteUser, id)
	return err
}

const updateUserPassword = `UPDATE users SET password_hash = $2 WHERE id = $1`

func (q *Queries) UpdateUserPassword(ctx context.Context, id int64, passwordHash string) error {
	_, err := q.pool.Exec(ctx, updateUserPassword, id, passwordHash)
	return err
}

const countUsers = `SELECT COUNT(*) FROM users`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	var n int64
	err := q.pool.QueryRow(ctx, countUsers).Scan(&n)
	return n, err
}

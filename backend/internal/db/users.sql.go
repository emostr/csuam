package db

import (
	"context"
)

const createUser = `
INSERT INTO users (username, password_hash, full_name, role)
VALUES (?, ?, ?, ?)
`

func (q *Queries) CreateUser(ctx context.Context, username, passwordHash, fullName, role string) (User, error) {
	res, err := q.db.ExecContext(ctx, createUser, username, passwordHash, fullName, role)
	if err != nil {
		return User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return User{}, err
	}
	return q.GetUserByID(ctx, id)
}

const getUserByUsername = `
SELECT id, username, password_hash, full_name, role, created_at
FROM users WHERE username = ?
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	var u User
	err := q.db.QueryRowContext(ctx, getUserByUsername, username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt)
	return u, err
}

const getUserByID = `
SELECT id, username, password_hash, full_name, role, created_at
FROM users WHERE id = ?
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	var u User
	err := q.db.QueryRowContext(ctx, getUserByID, id).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt)
	return u, err
}

const listUsers = `
SELECT id, username, password_hash, full_name, role, created_at
FROM users ORDER BY full_name
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
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

const deleteUser = `DELETE FROM users WHERE id = ?`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const updateUserPassword = `UPDATE users SET password_hash = ? WHERE id = ?`

func (q *Queries) UpdateUserPassword(ctx context.Context, id int64, passwordHash string) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, passwordHash, id)
	return err
}

const countUsers = `SELECT COUNT(*) FROM users`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	var n int64
	err := q.db.QueryRowContext(ctx, countUsers).Scan(&n)
	return n, err
}

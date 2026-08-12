package db

import (
	"context"
)

const grantPermission = `
INSERT INTO material_permissions (material_id, user_id, granted_by)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE user_id = user_id
`

func (q *Queries) GrantPermission(ctx context.Context, materialID, userID, grantedBy int64) error {
	_, err := q.db.ExecContext(ctx, grantPermission, materialID, userID, grantedBy)
	return err
}

const revokePermission = `
DELETE FROM material_permissions WHERE material_id = ? AND user_id = ?
`

func (q *Queries) RevokePermission(ctx context.Context, materialID, userID int64) error {
	_, err := q.db.ExecContext(ctx, revokePermission, materialID, userID)
	return err
}

const listPermissions = `
SELECT p.id, p.material_id, p.user_id, u.full_name, u.username, u.role, p.created_at
FROM material_permissions p
JOIN users u ON u.id = p.user_id
WHERE p.material_id = ?
ORDER BY p.created_at
`

func (q *Queries) ListPermissions(ctx context.Context, materialID int64) ([]Permission, error) {
	rows, err := q.db.QueryContext(ctx, listPermissions, materialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.MaterialID, &p.UserID, &p.FullName, &p.Username, &p.Role, &p.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

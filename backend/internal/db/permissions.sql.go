package db

import (
	"context"
)

const grantPermission = `
INSERT INTO material_permissions (material_id, user_id, granted_by)
VALUES ($1, $2, $3)
ON CONFLICT (material_id, user_id) DO NOTHING
`

func (q *Queries) GrantPermission(ctx context.Context, materialID, userID, grantedBy int64) error {
	_, err := q.pool.Exec(ctx, grantPermission, materialID, userID, grantedBy)
	return err
}

const revokePermission = `
DELETE FROM material_permissions WHERE material_id = $1 AND user_id = $2
`

func (q *Queries) RevokePermission(ctx context.Context, materialID, userID int64) error {
	_, err := q.pool.Exec(ctx, revokePermission, materialID, userID)
	return err
}

const listPermissions = `
SELECT p.id, p.material_id, p.user_id, u.full_name, u.username, u.role, p.created_at
FROM material_permissions p
JOIN users u ON u.id = p.user_id
WHERE p.material_id = $1
ORDER BY p.created_at
`

func (q *Queries) ListPermissions(ctx context.Context, materialID int64) ([]Permission, error) {
	rows, err := q.pool.Query(ctx, listPermissions, materialID)
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

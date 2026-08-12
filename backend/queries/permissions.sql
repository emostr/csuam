-- name: GrantPermission :exec
INSERT INTO material_permissions (material_id, user_id, granted_by)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE user_id = user_id;

-- name: RevokePermission :exec
DELETE FROM material_permissions WHERE material_id = ? AND user_id = ?;

-- name: ListPermissions :many
SELECT p.id, p.material_id, p.user_id, u.full_name, u.username, u.role, p.created_at
FROM material_permissions p
JOIN users u ON u.id = p.user_id
WHERE p.material_id = ?
ORDER BY p.created_at;

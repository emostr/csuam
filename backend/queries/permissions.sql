-- name: GrantPermission :exec
INSERT INTO material_permissions (material_id, user_id, granted_by)
VALUES ($1, $2, $3)
ON CONFLICT (material_id, user_id) DO NOTHING;

-- name: RevokePermission :exec
DELETE FROM material_permissions WHERE material_id = $1 AND user_id = $2;

-- name: ListPermissions :many
SELECT p.id, p.material_id, p.user_id, u.full_name, u.username, u.role, p.created_at
FROM material_permissions p
JOIN users u ON u.id = p.user_id
WHERE p.material_id = $1
ORDER BY p.created_at;

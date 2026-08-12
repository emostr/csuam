-- name: CreateMaterial :execlastid
INSERT INTO materials (title, description, category, status, file_key, file_name, file_mime, file_size, content, content_format, `condition`, location, origin_date, created_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetMaterial :one
SELECT m.*, u.full_name AS author_name
FROM materials m
LEFT JOIN users u ON u.id = m.created_by
WHERE m.id = ?;

-- name: HasPermission :one
SELECT EXISTS (
  SELECT 1 FROM material_permissions WHERE material_id = ? AND user_id = ?
);

-- name: ListMaterials :many
SELECT m.*, u.full_name AS author_name
FROM materials m
LEFT JOIN users u ON u.id = m.created_by
WHERE (? IS NULL OR m.category = ?)
  AND (? = '' OR m.title LIKE CONCAT('%', ?, '%') OR m.description LIKE CONCAT('%', ?, '%'))
  AND (? IS NULL OR m.origin_date >= ?)
  AND (? IS NULL OR m.origin_date <= ?)
  AND (? IS NULL OR m.status = ?)
  AND (
    ? = 'head_teacher'
    OR m.created_by = ?
    OR (
      (m.category <> 'documents' OR EXISTS (
        SELECT 1 FROM material_permissions p WHERE p.material_id = m.id AND p.user_id = ?
      ))
      AND (m.status = 'approved' OR ? = 'librarian')
    )
  )
ORDER BY m.created_at DESC;

-- name: UpdateMaterial :exec
UPDATE materials
SET title = ?, description = ?, `condition` = ?, location = ?,
    origin_date = ?, content = ?, content_format = ?, updated_at = now()
WHERE id = ?;

-- name: SetMaterialStatus :exec
UPDATE materials SET status = ?, updated_at = now() WHERE id = ?;

-- name: DeleteMaterial :exec
DELETE FROM materials WHERE id = ?;

-- name: CreateMaterial :one
INSERT INTO materials (title, description, category, status, file_key, file_name, file_mime, file_size, content, content_format, condition, location, origin_date, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING id;

-- name: GetMaterial :one
SELECT m.*, u.full_name AS author_name
FROM materials m
LEFT JOIN users u ON u.id = m.created_by
WHERE m.id = $1;

-- name: HasPermission :one
SELECT EXISTS (
  SELECT 1 FROM material_permissions WHERE material_id = $1 AND user_id = $2
);

-- name: ListMaterials :many
SELECT m.*, u.full_name AS author_name
FROM materials m
LEFT JOIN users u ON u.id = m.created_by
WHERE ($1::material_category IS NULL OR m.category = $1)
  AND ($2::text = '' OR m.title ILIKE '%' || $2 || '%' OR m.description ILIKE '%' || $2 || '%')
  AND ($3::date IS NULL OR m.origin_date >= $3)
  AND ($4::date IS NULL OR m.origin_date <= $4)
  AND ($5::material_status IS NULL OR m.status = $5)
  AND (
    $7::user_role = 'head_teacher'
    OR m.created_by = $6
    OR (
      (m.category <> 'documents' OR EXISTS (
        SELECT 1 FROM material_permissions p WHERE p.material_id = m.id AND p.user_id = $6
      ))
      AND (m.status = 'approved' OR $7 = 'librarian')
    )
  )
ORDER BY m.created_at DESC;

-- name: UpdateMaterial :exec
UPDATE materials
SET title = $2, description = $3, condition = $4, location = $5,
    origin_date = $6, content = $7, content_format = $8, updated_at = now()
WHERE id = $1;

-- name: SetMaterialStatus :exec
UPDATE materials SET status = $2, updated_at = now() WHERE id = $1;

-- name: DeleteMaterial :exec
DELETE FROM materials WHERE id = $1;

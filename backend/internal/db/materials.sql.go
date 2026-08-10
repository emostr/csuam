package db

import (
	"context"
	"time"
)

type CreateMaterialParams struct {
	Title         string
	Description   string
	Category      string
	Status        string
	FileKey       *string
	FileName      *string
	FileMime      *string
	FileSize      *int64
	Content       *string
	ContentFormat *string
	Condition     string
	Location      string
	OriginDate    *time.Time
	CreatedBy     *int64
}

const createMaterial = `
INSERT INTO materials (title, description, category, status, file_key, file_name, file_mime, file_size, content, content_format, condition, location, origin_date, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING id
`

func (q *Queries) CreateMaterial(ctx context.Context, p CreateMaterialParams) (int64, error) {
	var id int64
	err := q.pool.QueryRow(ctx, createMaterial,
		p.Title, p.Description, p.Category, p.Status,
		p.FileKey, p.FileName, p.FileMime, p.FileSize,
		p.Content, p.ContentFormat, p.Condition, p.Location,
		p.OriginDate, p.CreatedBy,
	).Scan(&id)
	return id, err
}

const materialColumns = `
m.id, m.title, m.description, m.category, m.status,
m.file_key, m.file_name, m.file_mime, m.file_size,
m.content, m.content_format, m.condition, m.location,
m.origin_date, m.created_by, m.created_at, m.updated_at,
u.full_name AS author_name
`

func scanMaterial(row interface{ Scan(dest ...any) error }) (Material, error) {
	var m Material
	err := row.Scan(
		&m.ID, &m.Title, &m.Description, &m.Category, &m.Status,
		&m.FileKey, &m.FileName, &m.FileMime, &m.FileSize,
		&m.Content, &m.ContentFormat, &m.Condition, &m.Location,
		&m.OriginDate, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
		&m.AuthorName,
	)
	return m, err
}

const getMaterial = `
SELECT ` + materialColumns + `
FROM materials m
LEFT JOIN users u ON u.id = m.created_by
WHERE m.id = $1
`

func (q *Queries) GetMaterial(ctx context.Context, id int64) (Material, error) {
	return scanMaterial(q.pool.QueryRow(ctx, getMaterial, id))
}

const hasPermission = `
SELECT EXISTS (
  SELECT 1 FROM material_permissions WHERE material_id = $1 AND user_id = $2
)
`

func (q *Queries) HasPermission(ctx context.Context, materialID, userID int64) (bool, error) {
	var ok bool
	err := q.pool.QueryRow(ctx, hasPermission, materialID, userID).Scan(&ok)
	return ok, err
}

type ListMaterialsParams struct {
	Category   *string
	Query      string
	FromDate   *time.Time
	ToDate     *time.Time
	Status     *string
	ViewerID   int64
	ViewerRole string
}

const listMaterials = `
SELECT ` + materialColumns + `
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
ORDER BY m.created_at DESC
`

func (q *Queries) ListMaterials(ctx context.Context, p ListMaterialsParams) ([]Material, error) {
	rows, err := q.pool.Query(ctx, listMaterials,
		p.Category, p.Query, p.FromDate, p.ToDate, p.Status, p.ViewerID, p.ViewerRole,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Material
	for rows.Next() {
		m, err := scanMaterial(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

type UpdateMaterialParams struct {
	ID            int64
	Title         string
	Description   string
	Condition     string
	Location      string
	OriginDate    *time.Time
	Content       *string
	ContentFormat *string
}

const updateMaterial = `
UPDATE materials
SET title = $2, description = $3, condition = $4, location = $5,
    origin_date = $6, content = $7, content_format = $8, updated_at = now()
WHERE id = $1
`

func (q *Queries) UpdateMaterial(ctx context.Context, p UpdateMaterialParams) error {
	_, err := q.pool.Exec(ctx, updateMaterial,
		p.ID, p.Title, p.Description, p.Condition, p.Location,
		p.OriginDate, p.Content, p.ContentFormat,
	)
	return err
}

const setMaterialStatus = `
UPDATE materials SET status = $2, updated_at = now() WHERE id = $1
`

func (q *Queries) SetMaterialStatus(ctx context.Context, id int64, status string) error {
	_, err := q.pool.Exec(ctx, setMaterialStatus, id, status)
	return err
}

const deleteMaterial = `DELETE FROM materials WHERE id = $1`

func (q *Queries) DeleteMaterial(ctx context.Context, id int64) error {
	_, err := q.pool.Exec(ctx, deleteMaterial, id)
	return err
}

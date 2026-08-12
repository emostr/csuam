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
INSERT INTO materials (title, description, category, status, file_key, file_name, file_mime, file_size, content, content_format, ` + "`condition`" + `, location, origin_date, created_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func (q *Queries) CreateMaterial(ctx context.Context, p CreateMaterialParams) (int64, error) {
	res, err := q.db.ExecContext(ctx, createMaterial,
		p.Title, p.Description, p.Category, p.Status,
		p.FileKey, p.FileName, p.FileMime, p.FileSize,
		p.Content, p.ContentFormat, p.Condition, p.Location,
		p.OriginDate, p.CreatedBy,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
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
WHERE m.id = ?
`

func (q *Queries) GetMaterial(ctx context.Context, id int64) (Material, error) {
	return scanMaterial(q.db.QueryRowContext(ctx, getMaterial, id))
}

const hasPermission = `
SELECT EXISTS (
  SELECT 1 FROM material_permissions WHERE material_id = ? AND user_id = ?
)
`

func (q *Queries) HasPermission(ctx context.Context, materialID, userID int64) (bool, error) {
	var ok bool
	err := q.db.QueryRowContext(ctx, hasPermission, materialID, userID).Scan(&ok)
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
ORDER BY m.created_at DESC
`

func (q *Queries) ListMaterials(ctx context.Context, p ListMaterialsParams) ([]Material, error) {
	rows, err := q.db.QueryContext(ctx, listMaterials,
		p.Category, p.Category,
		p.Query, p.Query, p.Query,
		p.FromDate, p.FromDate,
		p.ToDate, p.ToDate,
		p.Status, p.Status,
		p.ViewerRole, p.ViewerID, p.ViewerID, p.ViewerRole,
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
SET title = ?, description = ?, ` + "`condition`" + ` = ?, location = ?,
    origin_date = ?, content = ?, content_format = ?, updated_at = now()
WHERE id = ?
`

func (q *Queries) UpdateMaterial(ctx context.Context, p UpdateMaterialParams) error {
	_, err := q.db.ExecContext(ctx, updateMaterial,
		p.Title, p.Description, p.Condition, p.Location,
		p.OriginDate, p.Content, p.ContentFormat, p.ID,
	)
	return err
}

const setMaterialStatus = `
UPDATE materials SET status = ?, updated_at = now() WHERE id = ?
`

func (q *Queries) SetMaterialStatus(ctx context.Context, id int64, status string) error {
	_, err := q.db.ExecContext(ctx, setMaterialStatus, status, id)
	return err
}

const deleteMaterial = `DELETE FROM materials WHERE id = ?`

func (q *Queries) DeleteMaterial(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteMaterial, id)
	return err
}

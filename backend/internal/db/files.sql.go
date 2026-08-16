package db

import (
	"context"
	"strings"
)

type AddMaterialFileParams struct {
	MaterialID int64
	Key        string
	Name       string
	Mime       string
	Size       int64
}

const addMaterialFile = `
INSERT INTO material_files (material_id, file_key, file_name, file_mime, file_size)
VALUES (?, ?, ?, ?, ?)
`

func (q *Queries) AddMaterialFile(ctx context.Context, p AddMaterialFileParams) (int64, error) {
	res, err := q.db.ExecContext(ctx, addMaterialFile, p.MaterialID, p.Key, p.Name, p.Mime, p.Size)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

const materialFileColumns = `id, material_id, file_key, file_name, file_mime, file_size, created_at`

func scanMaterialFile(row interface{ Scan(dest ...any) error }) (MaterialFile, error) {
	var f MaterialFile
	err := row.Scan(&f.ID, &f.MaterialID, &f.Key, &f.Name, &f.Mime, &f.Size, &f.CreatedAt)
	return f, err
}

const getMaterialFile = `
SELECT ` + materialFileColumns + `
FROM material_files
WHERE id = ? AND material_id = ?
`

func (q *Queries) GetMaterialFile(ctx context.Context, materialID, fileID int64) (MaterialFile, error) {
	return scanMaterialFile(q.db.QueryRowContext(ctx, getMaterialFile, fileID, materialID))
}

const listMaterialFiles = `
SELECT ` + materialFileColumns + `
FROM material_files
WHERE material_id = ?
ORDER BY id
`

func (q *Queries) ListMaterialFiles(ctx context.Context, materialID int64) ([]MaterialFile, error) {
	rows, err := q.db.QueryContext(ctx, listMaterialFiles, materialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := []MaterialFile{}
	for rows.Next() {
		f, err := scanMaterialFile(rows)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, rows.Err()
}

func (q *Queries) ListMaterialFilesByMaterials(ctx context.Context, ids []int64) (map[int64][]MaterialFile, error) {
	grouped := make(map[int64][]MaterialFile, len(ids))
	if len(ids) == 0 {
		return grouped, nil
	}
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	query := `SELECT ` + materialFileColumns + `
FROM material_files
WHERE material_id IN (` + strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",") + `)
ORDER BY material_id, id`
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		f, err := scanMaterialFile(rows)
		if err != nil {
			return nil, err
		}
		grouped[f.MaterialID] = append(grouped[f.MaterialID], f)
	}
	return grouped, rows.Err()
}

const deleteMaterialFile = `DELETE FROM material_files WHERE id = ? AND material_id = ?`

func (q *Queries) DeleteMaterialFile(ctx context.Context, materialID, fileID int64) error {
	_, err := q.db.ExecContext(ctx, deleteMaterialFile, fileID, materialID)
	return err
}

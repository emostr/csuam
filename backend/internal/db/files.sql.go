package db

import (
	"context"
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
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

func (q *Queries) AddMaterialFile(ctx context.Context, p AddMaterialFileParams) (int64, error) {
	var id int64
	err := q.pool.QueryRow(ctx, addMaterialFile, p.MaterialID, p.Key, p.Name, p.Mime, p.Size).Scan(&id)
	return id, err
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
WHERE id = $1 AND material_id = $2
`

func (q *Queries) GetMaterialFile(ctx context.Context, materialID, fileID int64) (MaterialFile, error) {
	return scanMaterialFile(q.pool.QueryRow(ctx, getMaterialFile, fileID, materialID))
}

const listMaterialFiles = `
SELECT ` + materialFileColumns + `
FROM material_files
WHERE material_id = $1
ORDER BY id
`

func (q *Queries) ListMaterialFiles(ctx context.Context, materialID int64) ([]MaterialFile, error) {
	rows, err := q.pool.Query(ctx, listMaterialFiles, materialID)
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

const listMaterialFilesByMaterials = `
SELECT ` + materialFileColumns + `
FROM material_files
WHERE material_id = ANY($1::bigint[])
ORDER BY material_id, id
`

func (q *Queries) ListMaterialFilesByMaterials(ctx context.Context, ids []int64) (map[int64][]MaterialFile, error) {
	grouped := make(map[int64][]MaterialFile, len(ids))
	if len(ids) == 0 {
		return grouped, nil
	}
	rows, err := q.pool.Query(ctx, listMaterialFilesByMaterials, ids)
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

const deleteMaterialFile = `DELETE FROM material_files WHERE id = $1 AND material_id = $2`

func (q *Queries) DeleteMaterialFile(ctx context.Context, materialID, fileID int64) error {
	_, err := q.pool.Exec(ctx, deleteMaterialFile, fileID, materialID)
	return err
}

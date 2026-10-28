-- name: AddMaterialFile :one
INSERT INTO material_files (material_id, file_key, file_name, file_mime, file_size)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetMaterialFile :one
SELECT * FROM material_files WHERE id = $1 AND material_id = $2;

-- name: ListMaterialFiles :many
SELECT * FROM material_files WHERE material_id = $1 ORDER BY id;

-- name: ListMaterialFilesByMaterials :many
SELECT * FROM material_files WHERE material_id = ANY($1::bigint[]) ORDER BY material_id, id;

-- name: DeleteMaterialFile :exec
DELETE FROM material_files WHERE id = $1 AND material_id = $2;

-- name: AddMaterialFile :execlastid
INSERT INTO material_files (material_id, file_key, file_name, file_mime, file_size)
VALUES (?, ?, ?, ?, ?);

-- name: GetMaterialFile :one
SELECT * FROM material_files WHERE id = ? AND material_id = ?;

-- name: ListMaterialFiles :many
SELECT * FROM material_files WHERE material_id = ? ORDER BY id;

-- name: ListMaterialFilesByMaterials :many
SELECT * FROM material_files WHERE material_id IN (sqlc.slice(ids)) ORDER BY material_id, id;

-- name: DeleteMaterialFile :exec
DELETE FROM material_files WHERE id = ? AND material_id = ?;

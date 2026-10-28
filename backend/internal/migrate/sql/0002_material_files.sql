CREATE TABLE material_files (
    id BIGSERIAL PRIMARY KEY,
    material_id BIGINT NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    file_key TEXT NOT NULL,
    file_name TEXT NOT NULL,
    file_mime TEXT NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_material_files_material ON material_files(material_id);

INSERT INTO material_files (material_id, file_key, file_name, file_mime, file_size, created_at)
SELECT id, file_key,
       COALESCE(NULLIF(file_name, ''), 'file'),
       COALESCE(NULLIF(file_mime, ''), 'application/octet-stream'),
       COALESCE(file_size, 0),
       created_at
FROM materials
WHERE file_key IS NOT NULL AND file_key <> '';

ALTER TABLE materials
    DROP COLUMN file_key,
    DROP COLUMN file_name,
    DROP COLUMN file_mime,
    DROP COLUMN file_size;

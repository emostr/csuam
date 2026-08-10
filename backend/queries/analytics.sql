-- name: CountMaterials :one
SELECT COUNT(*) FROM materials;

-- name: CountMaterialsByCategory :many
SELECT category, COUNT(*) AS count FROM materials GROUP BY category;

-- name: MonthlyAdditions :many
SELECT to_char(date_trunc('month', created_at), 'YYYY-MM') AS month, COUNT(*) AS count
FROM materials
WHERE created_at >= date_trunc('month', now()) - interval '11 months'
GROUP BY 1
ORDER BY 1;

-- name: CountPendingMaterials :one
SELECT COUNT(*) FROM materials WHERE status = 'pending';

-- name: CountActiveLoans :one
SELECT COUNT(*) FROM loans WHERE returned_at IS NULL;

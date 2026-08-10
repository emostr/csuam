-- name: CreateLoan :one
INSERT INTO loans (material_id, borrower_name, note, due_date)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: ListLoans :many
SELECT l.id, l.material_id, m.title AS material_title, l.borrower_name, l.note,
       l.taken_at, l.due_date, l.returned_at
FROM loans l
JOIN materials m ON m.id = l.material_id
WHERE ($1::boolean AND l.returned_at IS NULL) OR (NOT $1::boolean AND l.returned_at IS NOT NULL)
ORDER BY l.due_date;

-- name: ReturnLoan :exec
UPDATE loans SET returned_at = now() WHERE id = $1 AND returned_at IS NULL;

-- name: CountOverdueLoans :one
SELECT COUNT(*) FROM loans WHERE returned_at IS NULL AND due_date < CURRENT_DATE;

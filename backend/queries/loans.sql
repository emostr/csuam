-- name: CreateLoan :execlastid
INSERT INTO loans (material_id, borrower_name, note, due_date)
VALUES (?, ?, ?, ?);

-- name: ListLoans :many
SELECT l.id, l.material_id, m.title AS material_title, l.borrower_name, l.note,
       l.taken_at, l.due_date, l.returned_at
FROM loans l
JOIN materials m ON m.id = l.material_id
WHERE (l.returned_at IS NULL) = ?
ORDER BY l.due_date;

-- name: ReturnLoan :exec
UPDATE loans SET returned_at = now() WHERE id = ? AND returned_at IS NULL;

-- name: CountOverdueLoans :one
SELECT COUNT(*) FROM loans WHERE returned_at IS NULL AND due_date < CURRENT_DATE;

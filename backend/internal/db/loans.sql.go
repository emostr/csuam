package db

import (
	"context"
	"time"
)

const createLoan = `
INSERT INTO loans (material_id, borrower_name, note, due_date)
VALUES ($1, $2, $3, $4)
RETURNING id
`

func (q *Queries) CreateLoan(ctx context.Context, materialID int64, borrowerName, note string, dueDate time.Time) (int64, error) {
	var id int64
	err := q.pool.QueryRow(ctx, createLoan, materialID, borrowerName, note, dueDate).Scan(&id)
	return id, err
}

const listLoans = `
SELECT l.id, l.material_id, m.title AS material_title, l.borrower_name, l.note,
       l.taken_at, l.due_date, l.returned_at
FROM loans l
JOIN materials m ON m.id = l.material_id
WHERE ($1::boolean AND l.returned_at IS NULL) OR (NOT $1::boolean AND l.returned_at IS NOT NULL)
ORDER BY l.due_date
`

func (q *Queries) ListLoans(ctx context.Context, active bool) ([]Loan, error) {
	rows, err := q.pool.Query(ctx, listLoans, active)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Loan
	for rows.Next() {
		var l Loan
		if err := rows.Scan(&l.ID, &l.MaterialID, &l.MaterialTitle, &l.BorrowerName, &l.Note, &l.TakenAt, &l.DueDate, &l.ReturnedAt); err != nil {
			return nil, err
		}
		items = append(items, l)
	}
	return items, rows.Err()
}

const returnLoan = `
UPDATE loans SET returned_at = now() WHERE id = $1 AND returned_at IS NULL
`

func (q *Queries) ReturnLoan(ctx context.Context, id int64) error {
	_, err := q.pool.Exec(ctx, returnLoan, id)
	return err
}

const countOverdueLoans = `
SELECT COUNT(*) FROM loans WHERE returned_at IS NULL AND due_date < CURRENT_DATE
`

func (q *Queries) CountOverdueLoans(ctx context.Context) (int64, error) {
	var n int64
	err := q.pool.QueryRow(ctx, countOverdueLoans).Scan(&n)
	return n, err
}

package db

import (
	"context"
	"time"
)

const createLoan = `
INSERT INTO loans (material_id, borrower_name, note, due_date)
VALUES (?, ?, ?, ?)
`

func (q *Queries) CreateLoan(ctx context.Context, materialID int64, borrowerName, note string, dueDate time.Time) (int64, error) {
	res, err := q.db.ExecContext(ctx, createLoan, materialID, borrowerName, note, dueDate)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

const listLoans = `
SELECT l.id, l.material_id, m.title AS material_title, l.borrower_name, l.note,
       l.taken_at, l.due_date, l.returned_at
FROM loans l
JOIN materials m ON m.id = l.material_id
WHERE (l.returned_at IS NULL) = ?
ORDER BY l.due_date
`

func (q *Queries) ListLoans(ctx context.Context, active bool) ([]Loan, error) {
	rows, err := q.db.QueryContext(ctx, listLoans, active)
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
UPDATE loans SET returned_at = now() WHERE id = ? AND returned_at IS NULL
`

func (q *Queries) ReturnLoan(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, returnLoan, id)
	return err
}

const countOverdueLoans = `
SELECT COUNT(*) FROM loans WHERE returned_at IS NULL AND due_date < CURRENT_DATE
`

func (q *Queries) CountOverdueLoans(ctx context.Context) (int64, error) {
	var n int64
	err := q.db.QueryRowContext(ctx, countOverdueLoans).Scan(&n)
	return n, err
}

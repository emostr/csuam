package db

import (
	"context"
)

const countMaterials = `SELECT COUNT(*) FROM materials`

func (q *Queries) CountMaterials(ctx context.Context) (int64, error) {
	var n int64
	err := q.db.QueryRowContext(ctx, countMaterials).Scan(&n)
	return n, err
}

const countMaterialsByCategory = `
SELECT category, COUNT(*) AS count FROM materials GROUP BY category
`

func (q *Queries) CountMaterialsByCategory(ctx context.Context) ([]CategoryCount, error) {
	rows, err := q.db.QueryContext(ctx, countMaterialsByCategory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CategoryCount
	for rows.Next() {
		var c CategoryCount
		if err := rows.Scan(&c.Category, &c.Count); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

const monthlyAdditions = `
SELECT DATE_FORMAT(created_at, '%Y-%m') AS month, COUNT(*) AS count
FROM materials
WHERE created_at >= DATE_SUB(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL 11 MONTH)
GROUP BY 1
ORDER BY 1
`

func (q *Queries) MonthlyAdditions(ctx context.Context) ([]MonthCount, error) {
	rows, err := q.db.QueryContext(ctx, monthlyAdditions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MonthCount
	for rows.Next() {
		var m MonthCount
		if err := rows.Scan(&m.Month, &m.Count); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

const countPendingMaterials = `SELECT COUNT(*) FROM materials WHERE status = 'pending'`

func (q *Queries) CountPendingMaterials(ctx context.Context) (int64, error) {
	var n int64
	err := q.db.QueryRowContext(ctx, countPendingMaterials).Scan(&n)
	return n, err
}

const countActiveLoans = `SELECT COUNT(*) FROM loans WHERE returned_at IS NULL`

func (q *Queries) CountActiveLoans(ctx context.Context) (int64, error) {
	var n int64
	err := q.db.QueryRowContext(ctx, countActiveLoans).Scan(&n)
	return n, err
}

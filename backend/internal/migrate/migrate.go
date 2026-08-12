package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sort"
)

//go:embed sql/*.sql
var fs embed.FS

func Run(ctx context.Context, sdb *sql.DB) error {
	if _, err := sdb.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (name VARCHAR(255) PRIMARY KEY, applied_at DATETIME NOT NULL DEFAULT current_timestamp())`); err != nil {
		return err
	}

	entries, err := fs.ReadDir("sql")
	if err != nil {
		return err
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		var exists bool
		if err := sdb.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE name = ?)`, name).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		body, err := fs.ReadFile("sql/" + name)
		if err != nil {
			return err
		}
		if _, err := sdb.ExecContext(ctx, string(body)); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
		if _, err := sdb.ExecContext(ctx, `INSERT INTO schema_migrations (name) VALUES (?)`, name); err != nil {
			return err
		}
	}
	return nil
}

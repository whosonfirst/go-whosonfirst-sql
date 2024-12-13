package database

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
)

func configureSQLiteDatabase(ctx context.Context, db *sql.DB, opts *ConfigureSQLDatabaseOptions) error {

	if opts.CreateTablesIfNecessary {

		table_names := make([]string, 0)

		sql := "SELECT name FROM sqlite_master WHERE type='table'"

		rows, err := db.QueryContext(ctx, sql)

		if err != nil {
			return fmt.Errorf("Failed to query sqlite_master, %w", err)
		}

		defer rows.Close()

		for rows.Next() {

			var name string
			err := rows.Scan(&name)

			if err != nil {
				return fmt.Errorf("Failed scan table name, %w", err)
			}

			table_names = append(table_names, name)
		}

		for _, t := range opts.Tables {

			if slices.Contains(table_names, t.Name()) {
				continue
			}

			schema, err := t.Schema(db)

			if err != nil {
				return fmt.Errorf("Failed to derive schema for table %s, %w", t.Name(), err)
			}
			
			_, err = db.ExecContext(ctx, schema)

			if err != nil {
				return fmt.Errorf("Failed to create %s table, %w", t.Name(), err)
			}
		}
	}

	return nil
}

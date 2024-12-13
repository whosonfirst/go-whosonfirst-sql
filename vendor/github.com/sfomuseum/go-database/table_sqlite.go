package database

// Something something something maybe build tags I am not sure yet...

import (
	"context"
	"database/sql"
	"fmt"
)

func hasSQLiteTable(ctx context.Context, db *sql.DB, table_name string) (bool, error) {

	has_table := false

	// TBD... how to derive database engine...

	sql := "SELECT name FROM sqlite_master WHERE type='table'"

	rows, err := db.QueryContext(ctx, sql)

	if err != nil {
		return false, fmt.Errorf("Failed to query sqlite_master, %w", err)
	}

	defer rows.Close()

	for rows.Next() {

		var name string
		err := rows.Scan(&name)

		if err != nil {
			return false, fmt.Errorf("Failed scan table name, %w", err)
		}

		if name == table_name {
			has_table = true
			break
		}
	}

	return has_table, nil
}

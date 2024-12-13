package database

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
)

var re_mem *regexp.Regexp
var re_vfs *regexp.Regexp
var re_file *regexp.Regexp

func init() {
	re_mem = regexp.MustCompile(`^(file\:)?\:memory\:.*`)
	re_vfs = regexp.MustCompile(`^vfs:\.*`)
	re_file = regexp.MustCompile(`^file\:([^\?]+)(?:\?.*)?$`)
}

type Table interface {
	Name() string
	Schema() string
	InitializeTable(context.Context, *sql.DB) error
	IndexRecord(context.Context, *sql.DB, interface{}) error
}

func HasTable(ctx context.Context, db *sql.DB, table_name string) (bool, error) {

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

func CreateTableIfNecessary(ctx context.Context, db *sql.DB, t Table) error {

	create := false

	has_table, err := HasTable(ctx, db, t.Name())

	if err != nil {
		return err
	}

	if !has_table {
		create = true
	}

	if create {

		sql := t.Schema()

		if err != nil {
			return err
		}

		_, err = db.ExecContext(ctx, sql)

		if err != nil {
			return err
		}

	}

	return nil
}

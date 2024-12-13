package database

import (
	"context"
	"database/sql"
	"fmt"
)

type ConfigureSQLDatabaseOptions struct {
	CreateTablesIfNecessary bool
	Tables                  []Table
	Pragma                  []string
}

func DefaultConfigureSQLDatabaseOptions() *ConfigureSQLDatabaseOptions {
	opts := &ConfigureSQLDatabaseOptions{}
	return opts
}

func ConfigureSQLDatabase(ctx context.Context, db *sql.DB, opts *ConfigureSQLDatabaseOptions) error {

	switch Driver(db) {
	case SQLITE_DRIVER:
		return configureSQLiteDatabase(ctx, db, opts)
	default:
		return fmt.Errorf("Unhandled or unsupported database driver %s", DriverTypeOf(db))
	}
}

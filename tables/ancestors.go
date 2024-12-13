package tables

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sfomuseum/go-database"
	"github.com/whosonfirst/go-whosonfirst-feature/alt"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
)

const ANCESTORS_TABLE_NAME string = "ancestors"

type AncestorsTable struct {
	database.Table
	FeatureTable
	name string
}

type AncestorsRow struct {
	Id                int64
	AncestorID        int64
	AncestorPlacetype string
	LastModified      int64
}

func NewAncestorsTableWithDatabase(ctx context.Context, db *sql.DB) (database.Table, error) {

	t, err := NewAncestorsTable(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to create '%s' table, %w", ANCESTORS_TABLE_NAME, err)
	}

	err = t.InitializeTable(ctx, db)

	if err != nil {
		return nil, database.InitializeTableError(t, err)
	}

	return t, nil
}

func NewAncestorsTable(ctx context.Context) (database.Table, error) {

	t := AncestorsTable{
		name: ANCESTORS_TABLE_NAME,
	}

	return &t, nil
}

func (t *AncestorsTable) Name() string {
	return t.name
}

func (t *AncestorsTable) Schema() string {
	schema, _ := LoadSchema("sqlite", ANCESTORS_TABLE_NAME)
	return schema
}

func (t *AncestorsTable) InitializeTable(ctx context.Context, db *sql.DB) error {
	return database.CreateTableIfNecessary(ctx, db, t)
}

func (t *AncestorsTable) IndexRecord(ctx context.Context, db *sql.DB, i interface{}) error {
	return t.IndexFeature(ctx, db, i.([]byte))
}

func (t *AncestorsTable) IndexFeature(ctx context.Context, db *sql.DB, f []byte) error {

	if alt.IsAlt(f) {
		return nil
	}

	id, err := properties.Id(f)

	if err != nil {
		return database.MissingPropertyError(t, "id", err)
	}

	tx, err := db.Begin()

	if err != nil {
		return database.BeginTransactionError(t, err)
	}

	sql := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return database.PrepareStatementError(t, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return database.ExecuteStatementError(t, err)
	}

	hierarchies := properties.Hierarchies(f)
	lastmod := properties.LastModified(f)

	for _, h := range hierarchies {

		for pt_key, ancestor_id := range h {

			ancestor_placetype := strings.Replace(pt_key, "_id", "", -1)

			sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
				id, ancestor_id, ancestor_placetype, lastmodified
			) VALUES (
			  	 ?, ?, ?, ?
			)`, t.Name())

			stmt, err := tx.Prepare(sql)

			if err != nil {
				return database.PrepareStatementError(t, err)
			}

			defer stmt.Close()

			_, err = stmt.Exec(id, ancestor_id, ancestor_placetype, lastmod)

			if err != nil {
				return database.ExecuteStatementError(t, err)
			}

		}

	}

	err = tx.Commit()

	if err != nil {
		return database.CommitTransactionError(t, err)
	}

	return nil
}

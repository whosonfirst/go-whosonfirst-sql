package tables

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/paulmach/orb/encoding/wkt"
	"github.com/sfomuseum/go-database"
	"github.com/whosonfirst/go-whosonfirst-feature/alt"
	"github.com/whosonfirst/go-whosonfirst-feature/geometry"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
)

const GEOMETRIES_TABLE_NAME string = "geometries"

type GeometriesTableOptions struct {
	IndexAltFiles bool
}

func DefaultGeometriesTableOptions() (*GeometriesTableOptions, error) {

	opts := GeometriesTableOptions{
		IndexAltFiles: false,
	}

	return &opts, nil
}

type GeometriesTable struct {
	database.Table
	FeatureTable
	name    string
	options *GeometriesTableOptions
}

type GeometriesRow struct {
	Id           int64
	Body         string
	LastModified int64
}

func NewGeometriesTable(ctx context.Context) (database.Table, error) {

	opts, err := DefaultGeometriesTableOptions()

	if err != nil {
		return nil, err
	}

	return NewGeometriesTableWithOptions(ctx, opts)
}

func NewGeometriesTableWithOptions(ctx context.Context, opts *GeometriesTableOptions) (database.Table, error) {

	t := GeometriesTable{
		name:    GEOMETRIES_TABLE_NAME,
		options: opts,
	}

	return &t, nil
}

func NewGeometriesTableWithDatabase(ctx context.Context, db *sql.DB) (database.Table, error) {

	opts, err := DefaultGeometriesTableOptions()

	if err != nil {
		return nil, err
	}

	return NewGeometriesTableWithDatabaseAndOptions(ctx, db, opts)
}

func NewGeometriesTableWithDatabaseAndOptions(ctx context.Context, db *sql.DB, opts *GeometriesTableOptions) (database.Table, error) {

	t, err := NewGeometriesTableWithOptions(ctx, opts)

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(ctx, db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *GeometriesTable) Name() string {
	return t.name
}

func (t *GeometriesTable) Schema(db *sql.DB) (string, error) {

	// really this should probably be the SPR table + geom but
	// let's just get this working first and then make it fancy
	// (20180109/thisisaaronland)

	// https://www.gaia-gis.it/spatialite-1.0a/SpatiaLite-tutorial.html
	// http://www.gaia-gis.it/gaia-sins/spatialite-sql-4.3.0.html

	// Note the InitSpatialMetaData() command because this:
	// https://stackoverflow.com/questions/17761089/cannot-create-column-with-spatialite-unexpected-metadata-layout

	return LoadSchema(db, GEOMETRIES_TABLE_NAME)
}

func (t *GeometriesTable) InitializeTable(ctx context.Context, db *sql.DB) error {

	return database.CreateTableIfNecessary(ctx, db, t)
}

func (t *GeometriesTable) IndexRecord(ctx context.Context, db *sql.DB, i interface{}) error {
	return t.IndexFeature(ctx, db, i.([]byte))
}

func (t *GeometriesTable) IndexFeature(ctx context.Context, db *sql.DB, f []byte) error {

	is_alt := alt.IsAlt(f)

	if is_alt && !t.options.IndexAltFiles {
		return nil
	}

	id, err := properties.Id(f)

	if err != nil {
		return MissingPropertyError(t, "id", err)
	}

	alt_label, err := properties.AltLabel(f)

	if err != nil {
		return MissingPropertyError(t, "alt label", err)
	}

	lastmod := properties.LastModified(f)

	tx, err := db.Begin()

	if err != nil {
		return database.BeginTransactionError(t, err)
	}

	geojson_geom, err := geometry.Geometry(f)

	if err != nil {
		return MissingPropertyError(t, "geometry", err)
	}

	orb_geom := geojson_geom.Geometry()

	str_wkt := wkt.MarshalString(orb_geom)

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		id, is_alt, alt_label, type, geom, lastmodified
	) VALUES (
		?, ?, ?, ?, GeomFromText('%s', 4326), ?
	)`, t.Name(), str_wkt)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return database.PrepareStatementError(t, err)
	}

	defer stmt.Close()

	geom_type := "common"

	_, err = stmt.Exec(id, is_alt, alt_label, geom_type, lastmod)

	if err != nil {
		return database.ExecuteStatementError(t, err)
	}

	err = tx.Commit()

	if err != nil {
		return database.CommitTransactionError(t, err)
	}

	return nil
}
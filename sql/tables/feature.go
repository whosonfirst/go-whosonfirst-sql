package tables

import (
	"context"
	"database/sql"

	"github.com/sfomuseum/go-database"
)

type FeatureTable interface {
	database.Table
	IndexFeature(context.Context, *sql.DB, []byte) error
}

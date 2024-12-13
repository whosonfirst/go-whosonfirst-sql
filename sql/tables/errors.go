package tables

import (
	"fmt"

	"github.com/sfomuseum/go-database"
)

// MissingPropertyError returns a new error with a default message for problems deriving a given property ('prop') from a record, wrapping 'err' and prepending with the value of 't's Name() method.
func MissingPropertyError(t database.Table, prop string, err error) error {
	return database.WrapError(t, fmt.Errorf("Failed to determine value for '%s' property, %w", prop, err))
}

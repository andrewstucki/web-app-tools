package migrator

import (
	rice "github.com/GeertJohan/go.rice"
	migrate "github.com/golang-migrate/migrate/v4"

	// import pq
	_ "github.com/lib/pq"
)

// NewBoxMigrator returns an instance of migrate.Migrate wrapping migrations in the given rice.Box
func NewBoxMigrator(migrations *rice.Box, connectionString string) (*migrate.Migrate, error) {
	box := NewBox(migrations)
	if err := box.Initialize(); err != nil {
		return nil, err
	}
	return migrate.NewWithSourceInstance("rice", box, connectionString)
}

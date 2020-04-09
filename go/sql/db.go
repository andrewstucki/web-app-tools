package sql

import (
	"github.com/jmoiron/sqlx"
	// import pq
	_ "github.com/lib/pq"
)

// Connect makes sure we set up a sqlx.DB instance with
// a particular postgres driver
func Connect(connectionString string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", connectionString)
}

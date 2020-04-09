package sql

import (
	"testing"

	managerTest "github.com/andrewstucki/web-app-tools/go/security/testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func database(test func(db *sqlx.DB)) {
	connection := "postgres://postgres:postgres@localhost/security-sql-test?sslmode=disable"
	db := sqlx.MustConnect("postgres", connection)
	defer func() {
		db.MustExec(`DROP TABLE IF EXISTS roles`)
		db.MustExec(`DROP TABLE IF EXISTS memberships`)
		db.Close()
	}()
	db.MustExec(`DROP TABLE IF EXISTS roles`)
	db.MustExec(`DROP TABLE IF EXISTS memberships`)
	db.MustExec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	db.MustExec(`CREATE TABLE roles (
    user_id uuid NOT NULL,
		role varchar(50) NOT NULL,
		PRIMARY KEY (user_id)
	)`)
	db.MustExec(`CREATE TABLE memberships (
		namespace_id uuid NOT NULL,
		user_id uuid NOT NULL,
		role varchar(50) NOT NULL,
		PRIMARY KEY (namespace_id, user_id)
	)`)

	test(db)
}

func TestSQLManager(t *testing.T) {
	database(func(db *sqlx.DB) {
		managerTest.ManagerTest(t, NewNamespaceManager(db))
	})
}

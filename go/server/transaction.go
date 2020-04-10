package server

import (
	"context"

	"github.com/jmoiron/sqlx"

	sqlContext "github.com/andrewstucki/web-app-tools/go/sql/context"
)

// Tx MUST be called only after using the transaction middleware
func Tx(ctx context.Context) *sqlx.Tx {
	tx := sqlContext.FromContext(ctx)
	if tx == nil {
		panic("tx middleware not used")
	}
	return tx
}

package context

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var transactionKey = struct{}{}

// WithTransaction returns a context with the transaction injected
func WithTransaction(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, transactionKey, tx)
}

// FromContext returns a transaction found in the context
func FromContext(ctx context.Context) *sqlx.Tx {
	tx := ctx.Value(transactionKey)
	if tx == nil {
		return nil
	}
	return tx.(*sqlx.Tx)
}

// QueryContext returns something that can query or exec
type QueryContext interface {
	sqlx.QueryerContext
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// GetQueryer returns a transaction or the database passed
func GetQueryer(ctx context.Context, db *sqlx.DB) QueryContext {
	if tx := FromContext(ctx); tx != nil {
		return tx
	}
	return db
}
package context

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var transactionKey = "transaction-context-key"

// WithTransaction returns a context with the transaction injected
func WithTransaction(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, &transactionKey, tx)
}

// FromContext returns a transaction found in the context
func FromContext(ctx context.Context) *sqlx.Tx {
	tx := ctx.Value(&transactionKey)
	if tx == nil {
		return nil
	}
	return tx.(*sqlx.Tx)
}

// StartTx starts a transaction and injects it into the context
func StartTx(ctx context.Context, db *sqlx.DB) (*sqlx.Tx, context.Context, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return tx, WithTransaction(ctx, tx), nil
}

// QueryContext returns something that can query or exec
type QueryContext interface {
	sqlx.QueryerContext
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// GetQueryer returns a transaction or the database passed
func GetQueryer(ctx context.Context, db *sqlx.DB) QueryContext {
	if tx := FromContext(ctx); tx != nil {
		return tx
	}
	return db
}

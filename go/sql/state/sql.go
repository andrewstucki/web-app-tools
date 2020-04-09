package sql

import (
	"context"

	sqlContext "github.com/andrewstucki/web-app-tools/go/sql/context"

	"github.com/jmoiron/sqlx"
)

const (
	findToken = `
		SELECT token FROM tokens WHERE id = $1
	`
	persistToken = `
		INSERT INTO tokens (id, token) VALUES ($1, $2)
	ON CONFLICT (id) DO
		UPDATE SET token = EXCLUDED.token;
	`
)

// TokenManager is a TokenManager
// that writes and retrieves data from
// a SQL database, it expects to have
// a table named "tokens" to read/write from
type TokenManager struct {
	db *sqlx.DB
}

// NewTokenManager creates a new token manager from the given
// database
func NewTokenManager(db *sqlx.DB) *TokenManager {
	return &TokenManager{
		db: db,
	}
}

// Set sets or updates the stored token
func (m *TokenManager) Set(ctx context.Context, subject, value string) error {
	_, err := sqlContext.GetQueryer(ctx, m.db).ExecContext(ctx, persistToken, subject, value)
	return err
}

// Get returns the stored token
func (m *TokenManager) Get(ctx context.Context, subject string) (string, error) {
	var token string
	if err := sqlx.GetContext(ctx, sqlContext.GetQueryer(ctx, m.db), &token, findToken, subject); err != nil {
		return "", err
	}
	return token, nil
}

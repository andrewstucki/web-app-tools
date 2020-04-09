package oauth

import "context"

// TokenManager maintains state
// for storing tokens
type TokenManager interface {
	Set(ctx context.Context, subject, token string) error
	Get(ctx context.Context, subject string) (string, error)
}

package state

import "context"

// HookedTokenManager is just a dummy
// wrapper around a TokenManager where
// a user can provide functions without
// having to wrap everything in a full struct
type HookedTokenManager struct {
	Setter func(ctx context.Context, subject, token string) error
	Getter func(ctx context.Context, subject string) (string, error)
}

// Set invokes the user-defined Setter
func (m *HookedTokenManager) Set(ctx context.Context, subject, token string) error {
	return m.Setter(ctx, subject, token)
}

// Get invokes the user-defined Getter
func (m *HookedTokenManager) Get(ctx context.Context, subject string) (string, error) {
	return m.Getter(ctx, subject)
}

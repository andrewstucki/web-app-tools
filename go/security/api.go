package security

import (
	"context"
	"sync"

	uuid "github.com/satori/go.uuid"
)

var (
	once                   sync.Once
	globalNamespaceManager NamespaceManager
	globalRoleManager      *roleManager
	globalNamespace        = uuid.UUID{}
)

func init() {
	globalRoleManager = newRoleManager()
}

// RegisterManager sets the global namespace manager
func RegisterManager(namespaceManager NamespaceManager) {
	once.Do(func() {
		globalNamespaceManager = namespaceManager
	})
}

// Register adds roles to the internal role manager
func Register(roles ...Role) {
	globalRoleManager.register(roles...)
}

// WithUser initializes a policy evaluation engine for the
// global namespace
func WithUser(user uuid.UUID) *Evaluator {
	return WithNamespaceAndUser(globalNamespace, user)
}

// WithNamespaceAndUser initializes a policy evaluation engine for the
// given namespace
func WithNamespaceAndUser(namespace, user uuid.UUID) *Evaluator {
	return newEvaluator(globalNamespaceManager, namespace, user)
}

// SetRole sets the role of a user in the global namespace
func SetRole(ctx context.Context, role Role, user uuid.UUID) error {
	return AddUserToNamespace(ctx, role, globalNamespace, user)
}

// UnsetRole removes the role of a user in the global namespace
func UnsetRole(ctx context.Context, user uuid.UUID) error {
	return RemoveUserFromNamespace(ctx, globalNamespace, user)
}

// AddUserToNamespace sets the role of a user in the given namespace
func AddUserToNamespace(ctx context.Context, role Role, id, user uuid.UUID) error {
	if globalNamespaceManager != nil {
		return globalNamespaceManager.AddUserToNamespace(ctx, role, id, user)
	}
	return nil
}

// RemoveUserFromNamespace removes the role of a user in the given namespace
func RemoveUserFromNamespace(ctx context.Context, id, user uuid.UUID) error {
	if globalNamespaceManager != nil {
		return globalNamespaceManager.RemoveUserFromNamespace(ctx, id, user)
	}
	return nil
}

package security

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

// NamespaceRole is a storage interface
// that each manager should implement
type NamespaceRole interface {
	// Namespace is the uuid of the namespace the role is associated with
	Namespace() uuid.UUID
	// Name is the name of the role initially registered with the global manager
	Name() string
}

// NamespaceManager is the main storage
// interface for storing roles based off of
// namespaces (including the global namespace)
type NamespaceManager interface {
	// AddUserToNamespace sets the role of a user in the given namespace
	AddUserToNamespace(ctx context.Context, role Role, id, user uuid.UUID) error
	// RemoveUserFromNamespace removes the role of a user in the given namespace
	RemoveUserFromNamespace(ctx context.Context, id, user uuid.UUID) error
	// RolesFor is used in gathering all of the roles for both the global and given namespace for
	// a given user
	RolesFor(ctx context.Context, globalNamespace, namespace, user uuid.UUID) ([]NamespaceRole, error)
}

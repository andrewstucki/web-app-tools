package memory

import (
	"context"
	"sync"

	"github.com/andrewstucki/web-app-tools/go/security"

	uuid "github.com/satori/go.uuid"
)

// NamespaceManager is an abstraction
// that writes and retrieves data from
// memory
type NamespaceManager struct {
	membership sync.Map
}

// NewNamespaceManager creates a new manager that stores
// everything in memory
func NewNamespaceManager() *NamespaceManager {
	return &NamespaceManager{}
}

func compoundKey(namespace, user uuid.UUID) string {
	return namespace.String() + "|" + user.String()
}

type memoryRole struct {
	MemoryName      string
	MemoryNamespace uuid.UUID
}

// Namespace is the uuid of the namespace the role is associated with
func (r *memoryRole) Namespace() uuid.UUID {
	return r.MemoryNamespace
}

// Name is the name of the role initially registered with the global manager
func (r *memoryRole) Name() string {
	return r.MemoryName
}

// AddUserToNamespace sets the role of a user in the given namespace
func (m *NamespaceManager) AddUserToNamespace(ctx context.Context, role security.Role, id, user uuid.UUID) error {
	m.membership.Store(compoundKey(id, user), &memoryRole{role.Name, id})
	return nil
}

// RemoveUserFromNamespace removes the role of a user in the given namespace
func (m *NamespaceManager) RemoveUserFromNamespace(ctx context.Context, id, user uuid.UUID) error {
	m.membership.Delete(compoundKey(id, user))
	return nil
}

// RolesFor is used in gathering all of the roles for both the global and given namespace for
// a given user
func (m *NamespaceManager) RolesFor(ctx context.Context, globalNamespace, namespace, user uuid.UUID) ([]security.NamespaceRole, error) {
	memoryRoles := []*memoryRole{}
	globalRole, ok := m.membership.Load(compoundKey(globalNamespace, user))
	if ok {
		memoryRoles = append(memoryRoles, globalRole.(*memoryRole))
	}
	role, ok := m.membership.Load(compoundKey(namespace, user))
	if ok {
		memoryRoles = append(memoryRoles, role.(*memoryRole))
	}
	roles := make([]security.NamespaceRole, len(memoryRoles))
	for i, memoryRole := range memoryRoles {
		roles[i] = memoryRole
	}
	return roles, nil
}

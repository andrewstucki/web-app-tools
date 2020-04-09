package security

import (
	"context"
	"sync"
	"testing"

	uuid "github.com/satori/go.uuid"
)

type testNamespaceManager struct {
	membership sync.Map
}

func newTestNamespaceManager() *testNamespaceManager {
	return &testNamespaceManager{}
}

func compoundKey(namespace, user uuid.UUID) string {
	return namespace.String() + "|" + user.String()
}

type memoryRole struct {
	MemoryName      string
	MemoryNamespace uuid.UUID
}

func (r *memoryRole) Namespace() uuid.UUID {
	return r.MemoryNamespace
}

func (r *memoryRole) Name() string {
	return r.MemoryName
}

func (m *testNamespaceManager) AddUserToNamespace(ctx context.Context, role Role, id, user uuid.UUID) error {
	m.membership.Store(compoundKey(id, user), &memoryRole{role.Name, id})
	return nil
}

func (m *testNamespaceManager) RemoveUserFromNamespace(ctx context.Context, id, user uuid.UUID) error {
	m.membership.Delete(compoundKey(id, user))
	return nil
}

func (m *testNamespaceManager) RolesFor(ctx context.Context, globalNamespace, namespace, user uuid.UUID) ([]NamespaceRole, error) {
	memoryRoles := []*memoryRole{}
	globalRole, ok := m.membership.Load(compoundKey(globalNamespace, user))
	if ok {
		memoryRoles = append(memoryRoles, globalRole.(*memoryRole))
	}
	role, ok := m.membership.Load(compoundKey(namespace, user))
	if ok {
		memoryRoles = append(memoryRoles, role.(*memoryRole))
	}
	roles := make([]NamespaceRole, len(memoryRoles))
	for i, memoryRole := range memoryRoles {
		roles[i] = memoryRole
	}
	return roles, nil
}

func TestPathMatcher(t *testing.T) {
	type args struct {
		path    string
		pattern string
	}
	tests := []struct {
		args args
		want bool
	}{
		{
			args: args{"/project/1/robot", "/project/1"},
			want: false,
		},
		{
			args: args{"/project/1/robot", "/project/:pid"},
			want: false,
		},
		{
			args: args{"/project/1/robot", "/project/1/*"},
			want: true,
		},
		{
			args: args{"/project/1/robot", "/project/:pid/robot"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run("match-"+tt.args.path+"-"+tt.args.pattern, func(t *testing.T) {
			if have := pathMatch(tt.args.path, tt.args.pattern); have != tt.want {
				t.Errorf("pathMatch() = %v, want %v", have, tt.want)
			}
		})
	}
}

func TestEvaluator(t *testing.T) {
	foo := Resource("foo")
	tests := []struct {
		name     string
		role     Role
		resource Resource
		action   Action
		want     bool
	}{
		{
			name:     "all foo - create",
			role:     Role{"admin", []Policy{{ResourceAll, ActionAll}}},
			resource: foo,
			action:   ActionCreate,
			want:     true,
		},
		{
			name:     "all foo - read",
			role:     Role{"admin", []Policy{{ResourceAll, ActionAll}}},
			resource: foo,
			action:   ActionRead,
			want:     true,
		},
		{
			name:     "all foo - update",
			role:     Role{"admin", []Policy{{ResourceAll, ActionAll}}},
			resource: foo,
			action:   ActionUpdate,
			want:     true,
		},
		{
			name:     "all foo - delete",
			role:     Role{"admin", []Policy{{ResourceAll, ActionAll}}},
			resource: foo,
			action:   ActionDelete,
			want:     true,
		},
		{
			name:     "all foo - list",
			role:     Role{"admin", []Policy{{ResourceAll, ActionAll}}},
			resource: foo,
			action:   ActionList,
			want:     true,
		},
		{
			name:     "create foo - list",
			role:     Role{"admin", []Policy{{ResourceAll, ActionCreate}}},
			resource: foo,
			action:   ActionList,
			want:     false,
		},
		{
			name:     "list foo - list",
			role:     Role{"admin", []Policy{{ResourceAll, ActionList}}},
			resource: foo,
			action:   ActionList,
			want:     true,
		},
		{
			name:     "none - list",
			role:     Role{"admin", []Policy{}},
			resource: foo,
			action:   ActionList,
			want:     false,
		},
		{
			name:     "foo sub resource - list",
			role:     Role{"admin", []Policy{{foo.Sub("bar"), ActionList}}},
			resource: foo.Sub("bar"),
			action:   ActionList,
			want:     true,
		},
		{
			name:     "foo sub resource - list fail",
			role:     Role{"admin", []Policy{{foo.Sub("bar"), ActionList}}},
			resource: foo.Sub("baz"),
			action:   ActionList,
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namespaceManager := newTestNamespaceManager()
			globalRoleManager = newRoleManager()
			globalRoleManager.register(tt.role)
			user := uuid.NewV4()

			namespaceManager.AddUserToNamespace(context.Background(), tt.role, globalNamespace, user)
			have, err := newEvaluator(namespaceManager, globalNamespace, user).Can(context.Background(), tt.action, tt.resource)
			if err != nil {
				t.Error(err)
				return
			}
			if have != tt.want {
				t.Errorf("Can() = %v, want %v", have, tt.want)
			}
		})
	}
}

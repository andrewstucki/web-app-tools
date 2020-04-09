package testing

import (
	"context"
	"testing"

	"github.com/andrewstucki/web-app-tools/go/security"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

// ManagerTest is a simple smoke test to make
// sure that a manager actually works
func ManagerTest(t *testing.T, manager security.NamespaceManager) {
	adminRole := security.Role{"admin", []security.Policy{
		{security.ResourceAll, security.ActionAll},
	}}
	security.Register([]security.Role{
		adminRole,
	}...)

	globalNamespace := uuid.UUID{}
	namespace := uuid.NewV4()
	user := uuid.NewV4()
	ctx := context.Background()

	err := manager.AddUserToNamespace(ctx, adminRole, namespace, user)
	require.NoError(t, err)
	roles, err := manager.RolesFor(ctx, globalNamespace, namespace, user)
	require.NoError(t, err)
	require.Len(t, roles, 1)
	require.Equal(t, "admin", roles[0].Name())
	require.Equal(t, namespace, roles[0].Namespace())

	err = manager.AddUserToNamespace(ctx, adminRole, globalNamespace, user)
	require.NoError(t, err)
	roles, err = manager.RolesFor(ctx, globalNamespace, namespace, user)
	require.NoError(t, err)
	require.Len(t, roles, 2)

	err = manager.RemoveUserFromNamespace(ctx, namespace, user)
	require.NoError(t, err)
	roles, err = manager.RolesFor(ctx, globalNamespace, namespace, user)
	require.NoError(t, err)
	require.Len(t, roles, 1)
	require.Equal(t, "admin", roles[0].Name())
	require.Equal(t, globalNamespace, roles[0].Namespace())

	err = manager.RemoveUserFromNamespace(ctx, globalNamespace, user)
	require.NoError(t, err)
	roles, err = manager.RolesFor(ctx, globalNamespace, namespace, user)
	require.Len(t, roles, 0)
}

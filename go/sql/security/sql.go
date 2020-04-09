package sql

import (
	"context"
	"database/sql"

	"github.com/andrewstucki/web-app-tools/go/security"
	sqlContext "github.com/andrewstucki/web-app-tools/go/sql/context"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

const (
	persistMembership = `
	INSERT INTO memberships (namespace_id, user_id, role)
		VALUES ($1, $2, $3)
	ON CONFLICT (namespace_id, user_id) DO
		UPDATE SET role = EXCLUDED.role;
	`
	deleteMembership = `
	DELETE FROM memberships WHERE
	namespace_id = $1 AND user_id = $2;
	`
	getMembership = `
	SELECT role FROM memberships WHERE
	namespace_id = $1 AND user_id = $2;
	`
	getRolesAndMembership = `
	SELECT role, namespace_id
	FROM memberships WHERE
	(namespace_id = $2 OR namespace_id = $3) AND user_id = $1;	
	`
)

// NamespaceManager is an abstraction
// that writes and retrieves data from
// a SQL database, it expects to have
// "memberships" to read/write from
// and "roles" to read from
type NamespaceManager struct {
	db *sqlx.DB
}

// NewNamespaceManager creates a new manager from the given
// database and driver
func NewNamespaceManager(db *sqlx.DB) *NamespaceManager {
	return &NamespaceManager{
		db: db,
	}
}

// AddUserToNamespace sets the role of a user in the given namespace
func (m *NamespaceManager) AddUserToNamespace(ctx context.Context, role security.Role, id, user uuid.UUID) error {
	_, err := sqlContext.GetQueryer(ctx, m.db).ExecContext(ctx, persistMembership, id, user, role.Name)
	return err
}

// RemoveUserFromNamespace removes the role of a user in the given namespace
func (m *NamespaceManager) RemoveUserFromNamespace(ctx context.Context, id, user uuid.UUID) error {
	_, err := sqlContext.GetQueryer(ctx, m.db).ExecContext(ctx, deleteMembership, id, user)
	return err
}

type role struct {
	DBName      string    `db:"role"`
	DBNamespace uuid.UUID `db:"namespace_id"`
}

// Namespace is the uuid of the namespace the role is associated with
func (r *role) Namespace() uuid.UUID {
	return r.DBNamespace
}

// Name is the name of the role initially registered with the global manager
func (r *role) Name() string {
	return r.DBName
}

// RolesFor is used in gathering all of the roles for both the global and given namespace for
// a given user
func (m *NamespaceManager) RolesFor(ctx context.Context, globalNamespace, namespace, user uuid.UUID) ([]security.NamespaceRole, error) {
	var dbRoles []*role
	if err := sqlx.SelectContext(ctx, sqlContext.GetQueryer(ctx, m.db), &dbRoles, getRolesAndMembership, user, namespace, globalNamespace); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	roles := make([]security.NamespaceRole, len(dbRoles))
	for i, dbRole := range dbRoles {
		roles[i] = dbRole
	}
	return roles, nil
}

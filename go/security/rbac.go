package security

import (
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"
)

// Role is an association of a name and a set of policies
type Role struct {
	Name     string
	Policies []Policy
}

func (r Role) policiesFor(namespace uuid.UUID) []Policy {
	if namespace == globalNamespace {
		return r.Policies
	}
	policies := make([]Policy, len(r.Policies))
	for i, policy := range r.Policies {
		sub := Resource(namespace.String()).Sub(policy.Resource)
		policies[i] = Policy{sub, policy.Action}
	}
	return policies
}

// Policy associates a resource with an action
type Policy struct {
	Resource
	Action
}

// String returns the representation of a policy as a string
func (p Policy) String() string {
	return p.Resource.String() + "|" + p.Action.String()
}

// Resource represents what is trying to be accessed
type Resource string

// String returns the representation of a resource as a string
func (r Resource) String() string {
	return string(r)
}

// Sub return a sub resource of the given resource
func (r Resource) Sub(resources ...Resource) Resource {
	elements := []string{r.String()}

	for _, resource := range resources {
		elements = append(elements, resource.String())
	}

	return Resource(strings.Join(elements, "/"))
}

// Action represents what the user is trying to do
type Action string

// String returns the representation of an action as a string
func (a Action) String() string {
	return string(a)
}

type roleManager struct {
	mutex sync.RWMutex
	roles map[string]Role
}

func newRoleManager() *roleManager {
	return &roleManager{
		roles: make(map[string]Role),
	}
}

func (m *roleManager) register(roles ...Role) {
	m.mutex.Lock()
	for _, role := range roles {
		m.roles[role.Name] = role
	}
	m.mutex.Unlock()
}

func (m *roleManager) getPolicies(roles ...NamespaceRole) []Policy {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	policies := []Policy{}
	for _, namespaceRole := range roles {
		role, ok := m.roles[namespaceRole.Name()]
		if !ok {
			continue
		}
		policies = append(policies, role.policiesFor(namespaceRole.Namespace())...)
	}
	return policies
}

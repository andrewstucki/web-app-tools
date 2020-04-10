package roles

import "github.com/andrewstucki/web-app-tools/go/security"

var (
	// SuperAdminRole can do anything
	SuperAdminRole = security.Role{
		Name: "super_admin",
		Policies: []security.Policy{
			security.Policy{Resource: security.ResourceAll, Action: security.ActionAll},
		},
	}
)

func Register() {
	security.Register(SuperAdminRole)
}

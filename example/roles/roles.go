package roles

import "github.com/andrewstucki/web-app-tools/go/security"

var (
	// SuperAdminRole can do anything
	SuperAdminRole = security.Role{"super_admin", []security.Policy{
		{security.ResourceAll, security.ActionAll},
	}}
)

func Register() {
	security.Register(SuperAdminRole)
}

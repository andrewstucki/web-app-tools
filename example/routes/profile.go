package routes

import (
	"net/http"

	"github.com/andrewstucki/web-app-tools/go/security"

	"example/payload"
	"example/roles"
)

// Me returns the current user and their permission policies
func (h *V1Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r.Context())
	policies, err := security.WithUser(user.ID).Policies(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("error retrieving policies")
		h.InternalError(w)
		return
	}
	h.Render(w, http.StatusOK, &payload.ProfileResponse{
		User:     user,
		Policies: policies,
		IsAdmin:  roles.IsAdmin(policies),
	})
}

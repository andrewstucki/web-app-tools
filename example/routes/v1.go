package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"

	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/andrewstucki/web-app-tools/go/security"

	"example/payload"
	"example/models"
)

type V1Handler struct {
	*baseHandler
}

func NewV1Handler(logger zerolog.Logger, render common.Renderer) *V1Handler {
	return &V1Handler{newBaseHandler(logger, render)}
}

func (h *V1Handler) Register(router chi.Router) {
	router.Get("/me", h.Me)
}

func (h *V1Handler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.MustCurrentUser(ctx, w, func(user *models.User) {
		policies, err := security.WithUser(user.ID).Policies(ctx)
		if err != nil {
			h.logger.Error().Err(err).Msg("error retrieving policies")
			h.InternalError(w)
			return
		}
		h.Render(w, http.StatusOK, &payload.ProfileResponse{
			User: user,
			Policies: policies,
		})
	})
}

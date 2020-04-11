package routes

import (
	"context"
	"net/http"

	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/andrewstucki/web-app-tools/go/security"
	"github.com/andrewstucki/web-app-tools/go/server"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"

	"example/models"
)

var userKey = "user-key"

// V1Handler is a wrapper around v1 api routes
type V1Handler struct {
	common.Renderer
	logger zerolog.Logger
}

// NewV1Handler returns a v1 handler
func NewV1Handler(logger zerolog.Logger, render common.Renderer) *V1Handler {
	return &V1Handler{
		Renderer: render,
		logger:   logger,
	}
}

// Register registers the v1 api handlers
func (h *V1Handler) Register(router chi.Router) {
	router.Get("/me", h.Authenticated(h.Me))
}

// WithCurrentUser retrieves the current logged in user
func (h *V1Handler) WithCurrentUser(ctx context.Context, w http.ResponseWriter, inner func(current *models.User)) {
	getCurrentUser(ctx, h.logger, h.Renderer, w, false, inner)
}

// Authenticated makes sure the user is authenticated
func (h *V1Handler) Authenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getCurrentUser(r.Context(), h.logger, h.Renderer, w, true, func(current *models.User) {
			next(w, r.Clone(context.WithValue(r.Context(), &userKey, current)))
		})
	}
}

// Authorized makes sure that the current request is authorized
func (h *V1Handler) Authorized(action security.Action, resource security.Resource, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		can(r.Context(), h.logger, h.Renderer, w, action, resource, func(current *models.User) {
			next(w, r.Clone(context.WithValue(r.Context(), &userKey, current)))
		})
	}
}

// CurrentUser should only be used when Authenticated or Authorized wraps a handler
func CurrentUser(ctx context.Context) *models.User {
	return ctx.Value(&userKey).(*models.User)
}

func getCurrentUser(ctx context.Context, logger zerolog.Logger, render common.Renderer, w http.ResponseWriter, required bool, inner func(current *models.User)) {
	current, err := server.CurrentUser(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("error getting current user")
		render.InternalError(w)
		return
	}
	if current != nil {
		user := current.(*models.User)
		if user != nil {
			inner(user)
			return
		}
	}
	if required {
		render.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	inner(nil)
}

func can(ctx context.Context, logger zerolog.Logger, render common.Renderer, w http.ResponseWriter, action security.Action, resource security.Resource, inner func(current *models.User)) {
	getCurrentUser(ctx, logger, render, w, true, func(user *models.User) {
		allowed, err := security.WithUser(user.ID).Can(ctx, action, resource)
		if err != nil {
			logger.Error().Err(err).Msg("error getting policies for user")
			render.InternalError(w)
			return
		}
		if !allowed {
			render.Error(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		inner(user)
	})
}

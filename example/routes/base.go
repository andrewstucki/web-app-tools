package routes

import (
	"context"
	"example/models"
	"net/http"

	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"

	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/andrewstucki/web-app-tools/go/security"
	"github.com/andrewstucki/web-app-tools/go/server"
)

type baseHandler struct {
	common.Renderer
	logger zerolog.Logger
}

func newBaseHandler(logger zerolog.Logger, render common.Renderer) *baseHandler {
	return &baseHandler{
		logger:   logger,
		Renderer: render,
	}
}

// WithCurrentUser retrieves the current logged in user
func (h *baseHandler) WithCurrentUser(ctx context.Context, w http.ResponseWriter, inner func(current *models.User)) {
	h.currentUserFromContext(ctx, w, false, inner)
}

// MustCurrentUser retrieves the current logged in user or returns unauthorized
func (h *baseHandler) MustCurrentUser(ctx context.Context, w http.ResponseWriter, inner func(current *models.User)) {
	h.currentUserFromContext(ctx, w, true, inner)
}

// Can enforces permissions
func (h *baseHandler) Can(ctx context.Context, w http.ResponseWriter, action security.Action, resource security.Resource, inner func(current *models.User)) {
	h.MustCurrentUser(ctx, w, func(user *models.User) {
		allowed, err := security.WithUser(user.ID).Can(ctx, action, resource)
		if err != nil {
			h.InternalError(w)
			return
		}
		if !allowed {
			h.Error(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		inner(user)
	})
}

// SuperAdminOnly enforces permissions
func (h *baseHandler) SuperAdminOnly(ctx context.Context, w http.ResponseWriter, inner func(current *models.User)) {
	h.MustCurrentUser(ctx, w, func(user *models.User) {
		allowed, err := security.WithUser(user.ID).Can(ctx, security.ActionAll, security.ResourceAll)
		if err != nil {
			h.InternalError(w)
			return
		}
		if !allowed {
			h.Error(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		inner(user)
	})
}

// CanWithNamespace enforces permissions in a namespace
func (h *baseHandler) CanWithNamespace(ctx context.Context, w http.ResponseWriter, namespace uuid.UUID, action security.Action, resource security.Resource, inner func(current *models.User)) {
	h.MustCurrentUser(ctx, w, func(user *models.User) {
		allowed, err := security.WithNamespaceAndUser(namespace, user.ID).Can(ctx, action, resource)
		if err != nil {
			h.InternalError(w)
			return
		}
		if !allowed {
			h.Error(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		inner(user)
	})
}

func (h *baseHandler) currentUserFromContext(ctx context.Context, w http.ResponseWriter, required bool, inner func(current *models.User)) {
	current, err := server.CurrentUser(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("error getting current user")
		h.InternalError(w)
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
		h.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	inner(nil)
}

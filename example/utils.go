package main

import (
	"context"
	"example/models"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/andrewstucki/web-app-tools/go/server"
)

var (
	logger zerolog.Logger
	render common.Renderer
)

// WithCurrentUser retrieves the current logged in user
func WithCurrentUser(ctx context.Context, w http.ResponseWriter, inner func(current *models.User)) {
	current, err := server.CurrentUser(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("error getting current user")
		render.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if current != nil {
		inner(current.(*models.User))
		return
	}
	inner(nil)
}

// MustCurrentUser retrieves the current logged in user or returns unauthorized
func MustCurrentUser(ctx context.Context, w http.ResponseWriter, inner func(current *models.User)) {
	current, err := server.CurrentUser(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("error getting current user")
		render.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if current == nil {
		render.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	inner(current.(*models.User))
}

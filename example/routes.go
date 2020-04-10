package main

import (
	"example/models"
	"net/http"
)

func me(w http.ResponseWriter, r *http.Request) {
	MustCurrentUser(r.Context(), w, func(currentUser *models.User) {
		render.Render(w, http.StatusOK, currentUser)
	})
}

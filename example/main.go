package main

//go:generate rice embed-go

import (
	"net/http"
	"os"

	rice "github.com/GeertJohan/go.rice"
	"github.com/joho/godotenv"

	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
	"github.com/andrewstucki/web-app-tools/go/server"
)

func init() {
	// ignore the error if no .env file is found
	godotenv.Load()
}

func main() {
	server.RunServer(server.Config{
		Migrations:   rice.MustFindBox("./migrations"),
		Assets:       rice.MustFindBox("./frontend/build"),
		HostPort:     ":3456",
		DatabaseURL:  os.Getenv("POSTGRES_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		BaseURL:      os.Getenv("BASE_URL"),
		SecretKey:    os.Getenv("JWT_SECRET"),
		Domains:      []string{"gpmail.org"},
		Setup: func(config *server.SetupConfig) {
			config.Router.Get("/v1/me", func(w http.ResponseWriter, r *http.Request) {
				claims := config.Handler.MustClaims(r.Context())
				config.Render.Render(w, http.StatusOK, struct {
					Email string `json:"email"`
				}{
					Email: claims.Email,
				})
			})
		},
		OnLogin: func(config *server.SetupConfig, claims *verifier.StandardClaims) error {
			return nil
		},
	})
}

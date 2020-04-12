package routes

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/andrewstucki/web-app-tools/go/security"
	"github.com/andrewstucki/web-app-tools/go/security/memory"
	"github.com/andrewstucki/web-app-tools/go/server"
	sqlContext "github.com/andrewstucki/web-app-tools/go/sql/context"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"example/models"
	"example/roles"
)

const postgresURL = "postgres://postgres:postgres@localhost:5434/example-test?sslmode=disable"

func randomUser() *models.User {
	return &models.User{
		ID:        uuid.NewV4(),
		Email:     randomdata.Email(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func jsonifyUser(user *models.User) string {
	return fmt.Sprintf(`{
		"id": "%s",
		"email": "%s",
		"createdAt": "%s",
		"updatedAt": "%s"
	}`, user.ID.String(), user.Email, user.CreatedAt.Format(time.RFC3339Nano), user.UpdatedAt.Format(time.RFC3339Nano))
}

func setupTest(t *testing.T, user *models.User) (*httptest.Server, func()) {
	router := chi.NewRouter()
	db := sqlx.MustConnect("postgres", postgresURL)
	security.RegisterManager(memory.NewNamespaceManager())
	security.Register(roles.SuperAdminRole)
	router.Route("/api/v1", func(router chi.Router) {
		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r.Clone(server.SetCurrentUserFn(r.Context(), func(_ context.Context) (interface{}, error) {
					return user, nil
				})))
			})
		}, func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tx, ctx, err := sqlContext.StartTx(r.Context(), db)
				require.NoError(t, err)
				next.ServeHTTP(w, r.Clone(ctx))
				require.NoError(t, tx.Rollback())
			})
		})
		NewV1Handler(zerolog.Nop(), common.NewJSONRenderer()).Register(router)
	})
	server := httptest.NewServer(router)
	return server, func() {
		server.Close()
		db.Close()
	}
}

func testRequest(t *testing.T, server *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	request, err := http.NewRequest(method, server.URL+path, body)
	require.NoError(t, err)
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	responseBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	defer response.Body.Close()
	return response, string(responseBody)
}

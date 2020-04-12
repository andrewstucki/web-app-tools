package routes

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/andrewstucki/web-app-tools/go/security"
	"github.com/stretchr/testify/require"

	"example/roles"
)

func TestMe_Admin(t *testing.T) {
	user := randomUser()
	server, cleanup := setupTest(t, user)
	defer cleanup()

	require.NoError(t, security.SetRole(context.Background(), roles.SuperAdminRole, user.ID))
	expected := fmt.Sprintf(`{
		"user": %s,
		"policies": [{"action":"*","resource":"*"}],
		"isAdmin": true
	}`, jsonifyUser(user))

	response, body := testRequest(t, server, "GET", "/api/v1/me", nil)

	require.Equal(t, http.StatusOK, response.StatusCode)
	require.JSONEq(t, expected, body)
}

func TestMe_NoAdmin(t *testing.T) {
	user := randomUser()
	server, cleanup := setupTest(t, user)
	defer cleanup()

	expected := fmt.Sprintf(`{
		"user": %s,
		"policies": [],
		"isAdmin": false
	}`, jsonifyUser(user))

	response, body := testRequest(t, server, "GET", "/api/v1/me", nil)

	require.Equal(t, http.StatusOK, response.StatusCode)
	require.JSONEq(t, expected, body)
}

func TestMe_NoUser(t *testing.T) {
	server, cleanup := setupTest(t, nil)
	defer cleanup()

	response, body := testRequest(t, server, "GET", "/api/v1/me", nil)

	require.Equal(t, http.StatusUnauthorized, response.StatusCode)
	require.JSONEq(t, `{"reason":"Unauthorized"}`, body)
}

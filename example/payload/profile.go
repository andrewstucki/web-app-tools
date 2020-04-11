package payload

import (
	"github.com/andrewstucki/web-app-tools/go/security"

	"example/models"
)

// APIError represents an error returned by the API
type APIError struct {
	Reason string `json:"reason"`
}

// ProfileResponse represents a profile request response
type ProfileResponse struct {
	User     *models.User      `json:"user"`
	Policies []security.Policy `json:"policies"`
	IsAdmin  bool              `json:"isAdmin"`
}

package oauth

import (
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
)

// Callbacks encapsulate the state handling logic when
// the flow endpoints/middleware either encounter an
// error, success, or get a refreshed token
type Callbacks interface {
	// OnError is invoked when any error is encountered in the handlers
	OnError(w http.ResponseWriter, err error)
	// OnSuccess is invoked when an id token is retrieved for the first
	// time at the end of an OAuth flow
	OnSuccess(w http.ResponseWriter, location, raw string, claims *verifier.StandardClaims)
	// OnInvalidToken is invoked when an id token is determined to be invalid
	// based off of the verification configuration passed into the handler
	OnInvalidToken(w http.ResponseWriter, err error)
	// OnRefresh is invoked when an id token is successfully refreshed
	// in middleware
	OnRefresh(w http.ResponseWriter, raw string) error
}

// Config is a configuration object for OAuth handlers.
type Config struct {
	// ClientTimeout is the timeout for doing the OAuth token exchange
	// if none is specified, defaults to 10 seconds
	ClientTimeout time.Duration
	// Verifier specifies the JWT verifier for the id token
	Verifier *verifier.Verifier
	// TokenManager manages token storage
	TokenManager TokenManager
	// Callbacks manage the error/success handling of the endpoint
	Callbacks Callbacks
	// AllowedRedirects whitelists where we can redirect to after getting a token
	AllowedRedirects []string
	// Logger is a zerolog instance used for logging
	Logger *zerolog.Logger

	// All of these must be specified

	// ClientID is the Google Client ID
	ClientID string
	// ClientSecret is the Google Client Secret
	ClientSecret string
	// MountURL is the URL where this handler is mounted
	MountURL string
	// SecretKey is the secret for JWT generation for state management
	SecretKey string

	// not exported
	mountURL *url.URL
}

func (c *Config) validate() error {
	if c.MountURL == "" {
		return ErrNeedMountURL
	}
	if c.ClientID == "" {
		return ErrNeedClientID
	}
	if c.ClientSecret == "" {
		return ErrNeedClientSecret
	}
	if c.SecretKey == "" {
		return ErrNeedSecretKey
	}

	mountURL, err := url.Parse(c.MountURL)
	if err != nil {
		return errors.Wrap(err, MessageMountURLParsingFailed)
	}
	c.mountURL = mountURL

	return nil
}

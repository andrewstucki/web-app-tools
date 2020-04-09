package server

import "errors"

var (
	// ErrNeedMountURL occurs when a mount url is not specified
	ErrNeedMountURL = errors.New("must specify a mount url")
	// ErrNeedClientID occurs when a client id is not specified
	ErrNeedClientID = errors.New("must specify a client id")
	// ErrNeedClientSecret occurs when a client secret is not specified
	ErrNeedClientSecret = errors.New("must specify a client secret")
	// ErrNeedSecretKey occurs when a secret key is not specified
	ErrNeedSecretKey = errors.New("must specify a secret key")
	// ErrInvalidRedirect occurs when we have a non-whitelisted
	// redirect parameter
	ErrInvalidRedirect = errors.New("bad redirect value")
	// ErrInvalidStateValue occurs when we the state returned
	// by the provider fails JWT validation
	ErrInvalidStateValue = errors.New("bad state value")
	// ErrInvalidCodeValue occurs when we the code returned
	// by the provider is blank
	ErrInvalidCodeValue = errors.New("bad code value")
	// ErrInvalidToken occurs when we the token returned after the exchange
	// by the provider is bad
	ErrInvalidToken = errors.New("invalid token")

	// The following values are annotations around the underlying errors

	// MessageMountURLParsingFailed occurs when we can't parse the URL provided
	// by MountURL
	MessageMountURLParsingFailed = "parsing mount url failed"
	// MessageStateCookieRetrieval occurs when we can't retrieve the state cookie after
	// the redirect from the provider
	MessageStateCookieRetrieval = "failed to get oauth state cookie"
	// MessageExchangeFailed occurs when we can't finish the exchange for the longer lived
	// tokens from the provider
	MessageExchangeFailed = "exchange failed"
	// MessageUserFailed occurs when we can't get information about the user from
	// the provider
	MessageUserFailed = "user retrieval failed"
	// MessageStateGenerationFailed occurs when we can't generate the state cookie for some
	// reason
	MessageStateGenerationFailed = "state generation failed"
	// MessageTokenRejected is displayed when a token handed back from Google has been rejected
	// for some reason, often due to an Audience or Domain mismatch
	MessageTokenRejected = "The token received was rejected, make sure you signed in with the right account."
)

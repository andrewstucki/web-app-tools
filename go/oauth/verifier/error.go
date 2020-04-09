package verifier

import (
	"errors"
)

var (
	ErrWrongSignature    = errors.New("token uses the wrong signature algorithm")
	ErrPublicKeyNotFound = errors.New("token references unknown public key")
	ErrInvalidToken      = errors.New("token is invalid")
	ErrIssuedAt          = errors.New("token used before issued")
	ErrExpired           = errors.New("token is expired")
	ErrInvalidIssuer     = errors.New("token has an invalid issuer")
	ErrInvalidAudience   = errors.New("token has an invalid audience")
	ErrInvalidDomain     = errors.New("token has an invalid domain")
)

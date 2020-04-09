package verifier

import (
	"time"
)

var (
	// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
	// You can override it to use another time value.  This is useful for testing or if your
	// server uses a different time zone than your tokens.
	TimeFunc       = time.Now
	allowedIssuers = []string{
		"accounts.google.com",
		"https://accounts.google.com",
	}
)

type GoogleClaims interface {
	Valid() error
	VerifyAudience(audience string) bool
	VerifyDomain(domain string) bool
}

type StandardClaims struct {
	Issuer          string `json:"iss"`
	AuthorizedParty string `json:"azp"`
	Audience        string `json:"aud"`
	Subject         string `json:"sub"`
	Email           string `json:"email"`
	EmailVerified   bool   `json:"email_verified"`
	AtHash          string `json:"at_hash"`
	Name            string `json:"name"`
	Picture         string `json:"picture"`
	GivenName       string `json:"given_name"`
	FamilyName      string `json:"family_name"`
	Locale          string `json:"locale"`
	HD              string `json:"hd"`
	IssuedAt        int64  `json:"iat"`
	ExpiresAt       int64  `json:"exp"`

	skew time.Duration
}

func (c *StandardClaims) SetAllowedSkew(skew time.Duration) {
	c.skew = skew
}

// Validates time based claims "exp, iat, nbf".
// There is no accounting for clock skew.
// As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (c StandardClaims) Valid() error {
	now := TimeFunc()
	expiration := time.Unix(c.ExpiresAt, 0).Add(c.skew)
	issued := time.Unix(c.IssuedAt, 0).Add(-c.skew)

	if issued.After(now) {
		return ErrIssuedAt
	}

	if !now.Before(expiration) {
		return ErrExpired
	}

	if !c.verifyIssuer() {
		return ErrInvalidIssuer
	}

	return nil
}

// Compares the Issuer claim against the valid list of issuers.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) verifyIssuer() bool {
	for _, issuer := range allowedIssuers {
		if c.Issuer == issuer {
			return true
		}
	}
	return false
}

// Compares the Audience claim against audience.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyAudience(audience string) bool {
	return c.Audience == audience
}

// Compares the Domain claim against domain.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyDomain(domain string) bool {
	return c.HD == domain
}

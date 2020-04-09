package verifier

import (
	"time"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

var (
	// account for up to 2 minutes of clock drift
	clockSkew = time.Minute * 2
)

type Verifier struct {
	Audiences *[]string
	Domains   *[]string
}

func NewVerifier() *Verifier {
	return &Verifier{}
}

func (v *Verifier) WithAudiences(audiences ...string) *Verifier {
	v.Audiences = &audiences
	return v
}

func (v *Verifier) WithDomains(domains ...string) *Verifier {
	v.Domains = &domains
	return v
}

func (v *Verifier) VerifyIDToken(token string, claims GoogleClaims) error {
	certs, err := getFederatedSignonCerts()
	if err != nil {
		return err
	}
	return v.verifySignedJWTWithCerts(token, certs, claims)
}

// verifySignedJWTWithCerts verifies the JWT string using the given secret key.
// On success it returns the user ID and the time the token was issued.
func (v *Verifier) verifySignedJWTWithCerts(token string, certs *Certs, claims GoogleClaims) error {
	_, err := jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, ErrWrongSignature
			}
			kid, ok := token.Header["kid"]
			if !ok {
				return nil, ErrInvalidToken
			}
			cert := certs.Key(kid.(string))
			if cert == nil {
				return nil, ErrPublicKeyNotFound
			}
			return cert, nil
		},
	)
	if err != nil {
		return err
	}

	if v.Audiences != nil && len(*v.Audiences) > 0 {
		found := false
		for _, audience := range *v.Audiences {
			if claims.VerifyAudience(audience) {
				found = true
				break
			}
		}
		if !found {
			return ErrInvalidAudience
		}
	}

	if v.Domains != nil && len(*v.Domains) > 0 {
		found := false
		for _, domain := range *v.Domains {
			if claims.VerifyDomain(domain) {
				found = true
				break
			}
		}
		if !found {
			return ErrInvalidDomain
		}
	}

	return nil
}
